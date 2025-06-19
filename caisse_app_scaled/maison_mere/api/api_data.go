package api

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/mere"
	"caisse-app-scaled/caisse_app_scaled/models"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Body struct {
	Message string `json:"message"`
	Host    string `json:"host"`
}

func newDataApi() *fiber.App {
	// api mount
	api := fiber.New(fiber.Config{})

	api.Post("/login", func(c *fiber.Ctx) error {
		employe := c.FormValue("username")
		role := c.FormValue("role")
		if !mere.Login(employe, role) {
			return c.Status(400).Render("login", nil)
		}
		return c.SendStatus(200)
	})

	api.Post("/notify", func(c *fiber.Ctx) error {
		var b Body
		if err := c.BodyParser(&b); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid notification format",
			})
		}
		println(b.Message)
		mere.Notifications = append(mere.Notifications, b.Message)
		return c.Status(200).JSON(fiber.Map{
			"message": "ok",
		})
	})

	api.Post("/subscribe", func(c *fiber.Ctx) error {
		var b Body
		if err := c.BodyParser(&b); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid notification format",
			})
		}
		if !slices.Contains(mere.Magasins, b.Host) {
			mere.Magasins = append(mere.Magasins, b.Host)
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "ok",
		})
	})

	api.Get("/alerts", authMiddleWare, func(c *fiber.Ctx) error {

		var notifs []Body = []Body{}
		for _, v := range mere.Notifications {
			notifs = append(notifs, Body{Message: v})
		}
		return c.Status(200).JSON(notifs)
	})

	api.Get("/transactions", func(c *fiber.Ctx) error {
		if transactions := mere.AfficherTransactions(); transactions != nil {
			return c.JSON(transactions)
		}
		logger.Error("transactions is nil")
		return c.Status(400).JSON(fiber.Map{
			"error": "transaction invalide",
		})
	})

	api.Get("/transactions/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid transaction ID",
			})
		}

		transaction, err := mere.AfficherUneTransactions(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}

		return c.JSON(transaction)
	})
	api.Post("/transactions", func(c *fiber.Ctx) error {
		var transaction models.Transaction
		if err := c.BodyParser(&transaction); err == nil {
			if err = mere.FaireUneVente(transaction); err == nil {
				return c.Status(200).JSON(transaction)
			}
			return c.Status(400).JSON(fiber.Map{
				"error": "Transaction non efectuer",
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"error": "transaction invalide",
		})
	})

	api.Delete("/transactions/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid transaction ID",
			})
		}

		if err := mere.FaireUnRetour(id); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message": "API is running",
		})
	})

	api.Get("/magasins", authMiddleWare, func(c *fiber.Ctx) error {

		magasins := mere.AfficherTousLesMagasins()
		return c.JSON(magasins)
	})

	api.Get("/analytics/:mag", authMiddleWare, func(c *fiber.Ctx) error {
		var data struct {
			Magasin string `json:"magasin"`
			Vente   struct {
				Total float64     `json:"total"`
				Dates []time.Time `json:"dates"`
				Sales []float64   `json:"sales"`
			} `json:"vente"`
		}
		mag, _ := url.QueryUnescape(c.Params("mag"))
		if mag == "tout" {
			total, date, sales := mere.AnalyticsVentetout()
			data.Vente.Total = total
			data.Vente.Dates = date
			data.Vente.Sales = sales
			data.Magasin = "tout"
		} else {
			total, date, sales := mere.AnalyticsVenteMagasin(mag)
			data.Vente.Total = total
			data.Vente.Dates = date
			data.Vente.Sales = sales
			data.Magasin = mag
		}
		return c.JSON(data)
	})

	api.Get("/raport", authMiddleWare, func(c *fiber.Ctx) error {
		type data struct {
			Magasin string         `json:"magasin"`
			Total   float64        `json:"total"`
			Best5   []string       `json:"best5"`
			Stock5  map[string]int `json:"stock5"`
		}
		var datas []data = []data{}
		mags := mere.AfficherTousLesMagasins()
		for _, mag := range mags {
			total, best5, stock5 := mere.GetRaportMagasin(mag)
			datas = append(datas, data{
				Magasin: mag,
				Total:   total,
				Best5:   best5,
				Stock5:  stock5,
			})
		}
		return c.JSON(datas)
	})

	api.Get("/produits/:nom", authMiddleWare, func(c *fiber.Ctx) error {
		nom, _ := url.QueryUnescape(c.Params("nom"))

		produits, err := mere.TrouverProduit(nom)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search products",
			})
		}
		return c.JSON(produits)
	})
	api.Put("/produit", authMiddleWare, func(c *fiber.Ctx) error {
		var data struct {
			ID          int     `json:"productId"`
			Nom         string  `json:"nom"`
			Prix        float64 `json:"prix"`
			Description string  `json:"description"`
		}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		err := mere.MiseAJourProduit(data.ID, data.Nom, data.Prix, data.Description)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search products",
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "success",
		})
	})

	return api
}
