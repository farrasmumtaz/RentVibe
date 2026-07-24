package main

import (
	"log"
	"os"

	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/routes"
)

func main() {

	config.LoadEnv()

	config.ConnectDatabase()

	router, err := routes.SetupRouter()
	if err != nil {
		log.Fatal("Failed to configure router:", err)
	}

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
