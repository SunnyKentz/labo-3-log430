package api

import (
	"caisse-app-scaled/caisse_app_scaled/magasin/caissier"
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
	api.Get("/produits", authMiddleWare, func(c *fiber.Ctx) error {
		produits, err := caissier.AfficherProduits()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to fetch products",
			})
		}
		return c.JSON(produits)
	})

	api.Get("/produits/:nom", authMiddleWare, func(c *fiber.Ctx) error {
		nom, _ := url.QueryUnescape(c.Params("nom"))

		produits, err := caissier.TrouverProduit(nom)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search products",
			})
		}
		return c.JSON(produits)
	})

	api.Post("/cart/:id", authMiddleWare, func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to add product to cart",
			})
		}

		err = caissier.AjouterALaCart(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "Product added to cart successfully",
		})
	})
	api.Get("/cart", authMiddleWare, func(c *fiber.Ctx) error {
		items, err := caissier.GetCartItems()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to get cart items",
			})
		}
		return c.JSON(items)
	})

	api.Delete("/cart/:id", authMiddleWare, func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to get cart items",
			})
		}
		caissier.RetirerDeLaCart(id)
		return c.Status(200).JSON(fiber.Map{
			"message": "Product deleted from cart successfully",
		})
	})

	api.Post("/vendre", authMiddleWare, func(c *fiber.Ctx) error {
		err := caissier.FaireUneVente()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to complete sale",
			})
		}
		caissier.ViderLaCart()
		return c.Status(200).JSON(fiber.Map{
			"message": "Sale completed successfully",
			"success": true,
		})
	})
	api.Get("/transactions", authMiddleWare, func(c *fiber.Ctx) error {
		transactions := caissier.AfficherTransactions()
		if transactions == nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to get transactions",
			})
		}
		return c.JSON(transactions)
	})
	api.Post("/rembourser/:id", authMiddleWare, func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid transaction ID",
			})
		}

		err = caissier.FaireUnRetour(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to process refund",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Refund processed successfully",
			"success": true,
		})
	})

	api.Post("/produit/:id", authMiddleWare, func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to Faire une commande",
			})
		}

		caissier.DemmandeReapprovisionner(id)
		return c.Status(200).JSON(fiber.Map{
			"message": "Commande passé successfully",
		})
	})
	api.Put("/produit/:id", func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec de réaprovisionement",
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
		err := caissier.MiseAJourProduit(id, nom, prix, description)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Produit non mise a jour",
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "success",
		})
	})

	api.Put("/produit/:id/:qt", func(c *fiber.Ctx) error {
		qt, err2 := strconv.Atoi(c.Params("qt"))
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil || err2 != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec de réaprovisionement",
			})
		}

		err := caissier.Reapprovisionner(id, qt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "echec de réaprovisionement",
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "réaprovisionement avec succes",
		})
	})
	return api
}
