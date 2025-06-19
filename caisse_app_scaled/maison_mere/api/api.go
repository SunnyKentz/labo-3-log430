package api

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/db"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/mere"
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
	app.Mount("/api", newDataApi())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", nil)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		employe := c.FormValue("username")
		role := c.FormValue("role")
		if !mere.Login(employe, role) {
			return c.Status(400).Render("login", nil)
		}
		return c.Redirect("/home", 302)
	})

	app.Get("/home", authMiddleWare, func(c *fiber.Ctx) error {
		return c.Render("analytics", nil)
	})

	app.Get("/rapport", authMiddleWare, func(c *fiber.Ctx) error {
		return c.Render("reports", nil)
	})

	app.Get("/produits", authMiddleWare, func(c *fiber.Ctx) error {
		return c.Render("products", nil)
	})

	db.Init()
	logger.Init("Mere")
	db.SetupLog()
	log.Fatal(app.Listen(":8090"))
}

func authMiddleWare(c *fiber.Ctx) error {
	if _, err := mere.Nom(); err != nil && c.Path() != "/login" {
		return c.Redirect("/")
	}
	return c.Next()
}
