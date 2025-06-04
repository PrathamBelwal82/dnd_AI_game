package main

import (
	"dnd_rpg/config"
	"dnd_rpg/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)
import "github.com/gofiber/fiber/v2/middleware/cors"

func main() {
	// Initialize Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // or your frontend origin
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	// Connect to Database
	config.ConnectDB()

	// Register Routes
	routes.SetupRoutes(app)

	// Start Server
	log.Fatal(app.Listen(":3001"))
}
