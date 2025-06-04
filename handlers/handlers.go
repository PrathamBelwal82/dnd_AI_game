package handlers

import (
	"bytes"
	"dnd_rpg/config"
	"dnd_rpg/models"
	"dnd_rpg/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
 
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func StartGame(c *fiber.Ctx) error {
	player := c.Query("player")
	log.Println("Starting game for player:", player)

	// Delete old game state if exists
	config.DB.Where("player = ?", player).Delete(&models.GameState{})

	// Create fresh initial story
	initialStory := "🌲 You wake up in a misty forest. The scent of damp moss fills your lungs as strange whispers echo from beyond the trees."
	gameState := models.GameState{
		Player:    player,
		Story:     initialStory,
		PlayerHP:  100,
		Inventory: []string{"Old Cloak", "Rusty Dagger"},
	}

	config.DB.Create(&gameState)

	return c.JSON(fiber.Map{
		"message":   "Game started",
		"story":     initialStory,
		"hp":        gameState.PlayerHP,
		"inventory": gameState.Inventory,
	})

}

func PlayerAction(c *fiber.Ctx) error {
	var req struct {
		Player string `json:"player"`
		Action string `json:"action"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var gameState models.GameState
	result := config.DB.Where("player = ?", req.Player).First(&gameState)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Game not found"})
	}

	// Append action to story
	currentContext := gameState.Story + "\n🧍 " + req.Action + "\n"

	// Generate AI response
	aiResp, err := getAIResponse(req.Action, gameState.Story, gameState.PlayerHP, gameState.Inventory)
	hpChange := utils.ExtractHPChange(aiResp)
	gameState.PlayerHP += hpChange
	if gameState.PlayerHP < 0 {
		gameState.PlayerHP = 0
	} else if gameState.PlayerHP > 100 {
		gameState.PlayerHP = 100
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "AI failed"})
	}

	lowerResp := strings.ToLower(aiResp)
	// Remove any "Player's action:" and everything after it in AI response
	if idx := strings.Index(lowerResp, "Player's action"); idx != -1 {
		aiResp = aiResp[:idx]
		aiResp = strings.TrimSpace(aiResp)
	}

	// Update story
	gameState.Story = currentContext + "🧙 " + aiResp
	config.DB.Save(&gameState)

	return c.JSON(fiber.Map{"response": aiResp})
}

func getAIResponse(action, storySoFar string, playerHP int, inventory pq.StringArray) (string, error) {
	apiURL := "http://localhost:11434/api/generate"

	gameStateSummary :=
		"Current game state:\n" +
			"- Player HP: " + fmt.Sprintf("%d/100", playerHP) + "\n" +
			"- Inventory: " + strings.Join(inventory, ", ") + "\n"

	fullPrompt := `
You are a Dungeon Master guiding a player through a fantasy world based on dark and dangerous forest theme.  Your job is to build suspense first, and only introduce challenges or combat occasionally.

Respond based on the player's latest action in one of two phases:
1. Scene-building (describe forest/mystery)
2. Combat (only if triggered, and clearly state HP loss and item gains)

Rules:
- Be vivid, but do NOT reveal hidden things unless the player earns it.
- If there's combat, you MUST say things like "You lose 10 HP" or "You gain a Sword".
` + gameStateSummary + `
Use vivid language, mystery, and tension.

Story so far:
` + storySoFar + `

Player's action: ` + action + `

Dungeon Master:
`

	payload, _ := json.Marshal(map[string]interface{}{
		"model":  "phi",
		"prompt": fullPrompt,
		"stream": false,
	})

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Ollama request failed:", err)
		return "⚠️ Dungeon Master is silent...", nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Ollama Raw Response: " + string(body))

	var result struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Response parsing failed:", err)
		return "⚠️ DM response error...", nil
	}

	return strings.TrimSpace(result.Response), nil
}
