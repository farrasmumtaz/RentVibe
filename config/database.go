package config

import (
	"log"
	"net"
	"net/url"
	"os"

	"github.com/farrasmumtaz/RentVibe/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := databaseDSN()

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

func databaseDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return os.Getenv("DATABASE_URL")
	}

	connectionURL := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			envOrDefault("DB_USER", "postgres"),
			os.Getenv("DB_PASSWORD"),
		),
		Host: net.JoinHostPort(host, envOrDefault("DB_PORT", "5432")),
		Path: envOrDefault("DB_NAME", "postgres"),
	}

	query := connectionURL.Query()
	query.Set("sslmode", envOrDefault("DB_SSLMODE", "disable"))
	connectionURL.RawQuery = query.Encode()

	return connectionURL.String()
}
