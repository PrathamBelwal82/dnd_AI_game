// File: handlers/dice.go
package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// DiceRollRequest - Struct for incoming JSON body
type DiceRollRequest struct {
	Player string `json:"player"`
	Die    string `json:"die"` // e.g., "d6", "d20"
}

func RollDice(c *fiber.Ctx) error {
	var req DiceRollRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	dieMap := map[string]int{
		"d4":  4,
		"d6":  6,
		"d8":  8,
		"d10": 10,
		"d12": 12,
		"d20": 20,
	}

	sides, ok := dieMap[req.Die]
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported die type"})
	}

	roll := rand.Intn(sides) + 1 // Roll between 1 and sides

	resultText := "You rolled a " + strconv.Itoa(roll) + " on a " + req.Die + "."

	// Optional: Save result to player's context in DB if needed

	return c.JSON(fiber.Map{
		"player": req.Player,
		"die":    req.Die,
		"roll":   roll,
		"result": resultText,
	})
}
