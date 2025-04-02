package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/saurify/saurify-backend/internal/handlers"
	sqldb "github.com/saurify/saurify-backend/internal/postgres"
	redisdb "github.com/saurify/saurify-backend/internal/redis"
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

	sqldb.InitDB(dsn)
	redis_db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisdb.InitRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWD"), redis_db)

	//initialize fiber
	app := fiber.New()

	//routes
	app.Post("/shorten", handlers.ShortenedURL)
	app.Get("/:shortCode", handlers.ResolveURL)

	//start server
	log.Fatal(app.Listen(os.Getenv("PORT")))
}
