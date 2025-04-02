package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/saurify/saurify-backend/internal/db"
	"github.com/saurify/saurify-backend/internal/handlers"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error: where is env file!?!?!")
	}

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("no db to connect to :/")
	}

	db.InitDB(dsn)

	//initialize fiber
	app := fiber.New()

	//routes
	app.Post("/shorten", handlers.ShortenedURL)
	app.Get("/:shortCode", handlers.ResolveURL)

	//start server
	log.Fatal(app.Listen(":8080"))
}
