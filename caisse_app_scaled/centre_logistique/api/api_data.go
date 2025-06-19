package api

import (
	"caisse-app-scaled/caisse_app_scaled/centre_logistique/logistics"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func newDataApi() *fiber.App {
	// api mount
	api := fiber.New(fiber.Config{})
	api.Get("/notify", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API is running",
		})
	})
	api.Get("/commands", authMiddleWare, func(c *fiber.Ctx) error {
		return c.JSON(logistics.GetAllCommands())
	})
	api.Post("/commande/:magasin/:id", func(c *fiber.Ctx) error {
		mag, _ := url.QueryUnescape(c.Params("magasin"))

		var body struct {
			Host string `json:"host"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		host := body.Host

		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec",
			})
		}
		logistics.AjouterUneCommande(id, mag, host)
		return c.JSON(fiber.Map{
			"message": "Commande reçu",
		})
	})
	api.Put("/commande/:id", authMiddleWare, func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec",
			})
		}
		logistics.AccepterUneCommande(id)
		return c.JSON(fiber.Map{
			"message": "Commande acceptée",
		})
	})
	api.Delete("/commande/:id", authMiddleWare, func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec",
			})
		}
		logistics.RefuserUneCommande(id)
		return c.JSON(fiber.Map{
			"message": "Commande refusé avec succes",
		})
	})

	api.Get("/produits/:nom", func(c *fiber.Ctx) error {
		nom, _ := url.QueryUnescape(c.Params("nom"))

		produits, err := logistics.TrouverProduit(nom)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search products",
			})
		}
		return c.JSON(produits)
	})

	api.Get("/produits/id/:id", func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec",
			})
		}
		prod, err := logistics.TrouverProduitParID(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to find product",
			})
		}
		return c.JSON(prod)
	})
	// PUT /produit/:id body : {"nom":nom,"prix":prix,"description":description}
	api.Put("/produit/:id", func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec de id",
			})
		}
		var body struct {
			Nom         string  `json:"nom"`
			Prix        float64 `json:"prix"`
			Description string  `json:"description"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		nom := body.Nom
		prix := body.Prix
		description := body.Description
		err := logistics.MiseAJourProduit(id, nom, prix, description)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Produit non mise a jour",
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "success",
		})
	})
	return api
}
