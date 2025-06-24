package api

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/magasin/caissier"
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

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
		return c.Render("login", fiber.Map{
			"Title": "Login - Caisse App",
		})
	})

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.Render("product", nil)
	})
	app.Get("/panier", func(c *fiber.Ctx) error {
		return c.Render("checkout", nil)
	})
	app.Get("/transactions", func(c *fiber.Ctx) error {
		return c.Render("transactions", nil)
	})

	port := ":8080"
	caissier.Host = os.Getenv("GATEWAY") + port
	logger.Init("Magasin")
	log.Fatal(app.Listen(port))
}

func authMiddleWare(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	magasin := c.Get("C-Mag")
	caisse := c.Get("C-Caisse")
	err := caissier.CheckLogedIn(authHeader, magasin, caisse)
	if err != nil {
		if strings.HasPrefix(c.Path(), "/api") {
			return GetApiError(c, "this action requires authentification", http.StatusUnauthorized)
		}
		return c.Redirect("/")
	}
	return c.Next()
}
