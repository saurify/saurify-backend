package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/saurify/saurify-backend/internal/db"
	"github.com/saurify/saurify-backend/internal/services"
)

func ShortenedURL(c *fiber.Ctx) error {
	var request struct {
		URL         string `json:"url"`
		IsTemporary bool   `json:"is_temporary"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	shortCode := services.GenerateShortCode(request.URL)

	if request.IsTemporary {
		//saving in redis
		if err := db.SaveURLToCache(shortCode, request.URL); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store in cache"})
		}
	} else {
		if err := db.SaveURLToSQL(shortCode, request.URL, false); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store in sql db"})
		}
		db.SaveURLToCache(shortCode, request.URL)
	}
	return c.JSON(fiber.Map{"short_url": shortCode})
}

func ResolveURL(c *fiber.Ctx) error {
	var request struct {
		ShortCode string `json:"short_code"`
	}

	// Try to parse JSON body (for POST requests)
	if err := c.BodyParser(&request); err != nil {
		request.ShortCode = c.Params("shortCode") // Fallback to URL param
	}

	shortCode := request.ShortCode
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Short code is required"})
	}

	// Check Redis cache first
	originalURL, err := db.GetURLToCache(shortCode)
	if err == nil && originalURL != "" {
		return c.Redirect(originalURL, http.StatusFound)
	}

	// Check SQL database
	originalURL, isTemporary, err := db.GetURLFromSQL(shortCode)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}

	// Save to cache if it's permanent
	if !isTemporary {
		go db.SaveURLToCache(shortCode, originalURL)
	}

	// If temporary, delete from SQL after use
	if isTemporary {
		go db.DeleteFromSQL(shortCode)
	}

	return c.Redirect(originalURL, http.StatusFound)
}
