package config

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"dnd_rpg/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=Pratham#123 dbname=dnd-rpg port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate models
	DB.AutoMigrate(&models.GameState{})
	log.Println("Database connected")
}
