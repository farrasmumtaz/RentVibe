package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/farrasmumtaz/RentVibe/config"
	"github.com/farrasmumtaz/RentVibe/internal/cache"
	"github.com/farrasmumtaz/RentVibe/internal/migration"
	"github.com/farrasmumtaz/RentVibe/internal/routes"
)

func main() {

	config.LoadEnv()

	config.ConnectDatabase()
	if err := migration.Run(); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migration completed")

	redisClient, err := config.ConnectRedis(context.Background())
	var cacheStore cache.Store = cache.NewNopStore()
	if err != nil {
		log.Printf("Redis unavailable, cache disabled: %v", err)
	} else {
		defer redisClient.Close()
		cacheStore = cache.NewRedisStore(redisClient)
		log.Println("Redis connected successfully")
	}

	router, err := routes.SetupRouter(cacheStore)
	if err != nil {
		log.Fatal("Failed to configure router:", err)
	}

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Server running on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}
}
