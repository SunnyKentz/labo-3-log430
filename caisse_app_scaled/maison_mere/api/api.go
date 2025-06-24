package api

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/db"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/mere"
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

// @title Maison Mere API
// @version 1.0
// @description This is the API for the Maison Mere service.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8090
// @BasePath /api/v1
func NewApp() {

	engine := html.New("./view", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return false
		},
	}))

	app.Static("/static", "./view")
	app.Static("/js", "./commonjs")
	app.Mount("/api/v1", newDataApi())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", nil)
	})

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.Render("analytics", nil)
	})

	app.Get("/rapport", func(c *fiber.Ctx) error {
		return c.Render("reports", nil)
	})

	app.Get("/produits", func(c *fiber.Ctx) error {
		return c.Render("products", nil)
	})

	db.Init()
	logger.Init("Mere")
	db.SetupLog()
	log.Fatal(app.Listen(":8090"))
}

func authMiddleWare(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	err := mere.CheckLogedIn(authHeader)
	if err != nil {
		if c.Path() == "api/v1/login" {
			return c.Next()
		}
		if strings.HasPrefix(c.Path(), "/api") {
			return GetApiError(c, "this action requires authentification", http.StatusUnauthorized)
		}
		return c.Redirect("/")
	}
	return c.Next()
}
