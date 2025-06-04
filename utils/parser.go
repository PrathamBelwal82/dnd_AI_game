package utils

import (
	"regexp"
	"fmt"
)

// Extracts damage if mentioned in AI response like "You lose 10 HP"
func ExtractHPChange(response string) int {
	re := regexp.MustCompile(`(?i)(lose|lost|take|took) (\d+) HP`)
	matches := re.FindStringSubmatch(response)
	if len(matches) == 3 {
		var dmg int
		fmt.Sscanf(matches[2], "%d", &dmg)
		return -dmg
	}
	return 0
}