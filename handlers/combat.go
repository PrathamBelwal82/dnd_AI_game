package handlers

import (
	"dnd_rpg/models"
)

func ProcessTurn(action string, player *models.Character, npc *models.Character) (string, error) {
	// 1. Get AI response to the player's action
	aiReply := " 💀 You were defeated!"

	return aiReply, nil
}

func CalculateDamage(attacker, defender *models.Character) int {
	damage := attacker.Attack - defender.Defense/2
	if damage < 1 {
		damage = 1
	}
	defender.HP -= damage
	return damage
}
