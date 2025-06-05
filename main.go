package main

import (
	"dnd_rpg/config"
	"dnd_rpg/routes"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)
import "github.com/gofiber/fiber/v2/middleware/cors"

func main() {
	// Initialize Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", 
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	// Connect to Database
	config.ConnectDB()

	// Register Routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Fallback for local dev
	}

	log.Fatal(app.Listen(":" + port))
}
