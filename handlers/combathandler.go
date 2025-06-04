package handlers

import (
	"github.com/gofiber/fiber/v2"
	"dnd_rpg/models"
	
)

func CombatTurnHandler(c *fiber.Ctx) error {
	// Parse JSON body into struct
	var payload struct {
		Action string `json:"action"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	// Initialize characters (example)
	player := &models.Character{
		Name:    "Hero",
		HP:      100,
		Attack:  20,
		Defense: 10,
	}
	npc := &models.Character{
		Name:    "Goblin",
		HP:      60,
		Attack:  15,
		Defense: 5,
	}

	// Process turn
	result, err := ProcessTurn(payload.Action, player, npc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error processing turn")
	}

	// Return JSON response
	return c.JSON(fiber.Map{
		"response": result,
	})
}
