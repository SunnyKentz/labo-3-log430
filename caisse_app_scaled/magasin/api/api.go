package api

import (
	"caisse-app-scaled/caisse_app_scaled/magasin/caissier"
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func NewApp() {

	engine := html.New("./view", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./view")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "Login - Caisse App",
		})
	})
	app.Mount("/api", newDataApi())

	app.Post("/login", func(c *fiber.Ctx) error {
		employe := c.FormValue("username")
		caisse := c.FormValue("caisse")
		magasin := c.FormValue("magasin")
		if magasin != "" {
			caissier.Magasin = magasin
		}
		if !caissier.InitialiserPOS(employe, caisse) {
			return c.SendStatus(400)
		}
		return c.Redirect("/home", 302)
	})

	app.Get("/home", authMiddleWare, func(c *fiber.Ctx) error {
		return c.Render("product", nil)
	})
	app.Get("/panier", authMiddleWare, func(c *fiber.Ctx) error {
		return c.Render("checkout", nil)
	})

	app.Get("/transactions", authMiddleWare, func(c *fiber.Ctx) error {

		return c.Render("transactions", nil)
	})
	port := ":8080"
	caissier.Host = GATEWAY + port
	log.Fatal(app.Listen(port))
}

func authMiddleWare(c *fiber.Ctx) error {
	if _, err := caissier.Nom(); err != nil && c.Path() != "/login" {
		return c.Redirect("/")
	}
	return c.Next()
}
