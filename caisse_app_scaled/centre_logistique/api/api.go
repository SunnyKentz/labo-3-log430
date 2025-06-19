package api

import (
	"caisse-app-scaled/caisse_app_scaled/centre_logistique/db"
	"caisse-app-scaled/caisse_app_scaled/centre_logistique/logistics"
	"caisse-app-scaled/caisse_app_scaled/logger"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func NewApp() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}
	log.Println("Working directory:", workingDir)
	engine := html.New("./view", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./view")
	app.Mount("/api", newDataApi())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "Login - Caisse App",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		employe := c.FormValue("username")
		if !logistics.Login(employe) {
			return c.Status(400).Render("login", nil)
		}
		return c.Redirect("/home", 302)
	})

	app.Get("/home", authMiddleWare, func(c *fiber.Ctx) error {
		return c.Render("commande", nil)
	})
	db.Init()
	logger.Init("Logistique")
	db.SetupLog()
	log.Fatal(app.Listen(":8091"))
}

func authMiddleWare(c *fiber.Ctx) error {
	if _, err := logistics.Nom(); err != nil && c.Path() != "/login" {
		return c.Redirect("/")
	}
	return c.Next()
}
