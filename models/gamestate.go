package models
import "github.com/lib/pq"

type GameState struct {
    ID        uint   `gorm:"primaryKey"`
    Player    string `gorm:"unique"`
    Story     string
    PlayerHP  int      `gorm:"default:100"`
    Inventory pq.StringArray `gorm:"type:text[]"`
}
