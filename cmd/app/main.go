package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidcm146/bus-booking-be/internal/app"
	"github.com/joho/godotenv"
)

// @title       Bus Booking API
// @version     1.0
// @description Backend API for the Bus Booking system
// @host        localhost:8080
// @BasePath    /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Shutdown()

	// Graceful shutdown on OS signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Use goRoutine to run the server to avoid blocking the main thread
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
}
