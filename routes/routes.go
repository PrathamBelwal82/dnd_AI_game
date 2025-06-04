package routes

import (
	"github.com/gofiber/fiber/v2"
	"dnd_rpg/handlers"
)

// SetupRoutes registers API endpoints
func SetupRoutes(app *fiber.App) {
	app.Get("/start-game", handlers.StartGame)
	app.Post("/player-action", handlers.PlayerAction)
	app.Post("/combat/turn", handlers.CombatTurnHandler)
	app.Post("/roll-dice",handlers.RollDice)

}
