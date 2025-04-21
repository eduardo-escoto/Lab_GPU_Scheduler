package main

import (
	"go-webserver-project/config"
	"go-webserver-project/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Log the loaded configuration (for debugging purposes)
	log.Printf("Loaded configuration: %+v\n", cfg)

	// Initialize routes
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
