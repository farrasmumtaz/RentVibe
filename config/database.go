package config

import (
	"log"
	"os"

	"github.com/farrasmumtaz/RentVibe/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	DB = database

	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Item{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database connected successfully")
	log.Println("Database migration completed")
}
