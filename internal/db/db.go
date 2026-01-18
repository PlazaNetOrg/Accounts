package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"plazanet-accounts/internal/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("accounts.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.User{})
	log.Println("Database connected and migrated")
}