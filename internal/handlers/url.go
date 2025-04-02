package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/saurify/saurify-backend/internal/db"
	"github.com/saurify/saurify-backend/internal/services"
)

func ShortenedURL(c *fiber.Ctx) error {
	var request struct {
		URL string `json:"url"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	shortCode := services.GenerateShortCode(request.URL)

	_, err := db.DB.Exec(context.Background(), "INSERT INTO shortlinks (short_code, original_url) VALUES ($1, $2)", shortCode, request.URL)
	if err != nil {
		log.Println("X Db insert error: ", err)
		return c.Status(500).JSON(fiber.Map{"error": "Could not store URL"})
	}

	return c.JSON(fiber.Map{"short_url": "http://localhost:8080/" + shortCode})
}

func ResolveURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	var originalURL string

	err := db.DB.QueryRow(context.Background(), "SELECT origin_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalURL)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
	}

	return c.Redirect(originalURL, 301)
}
