package main

import (
	"log"
	"net/http"
	"os"

	"github.com/eduardo-escoto/gpu_request/server/internal/database"
	"github.com/eduardo-escoto/gpu_request/server/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	if os.Getenv("GPU_SCHED_ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	log.Printf("Loaded DSN: %+v", os.Getenv("DATABASE_DSN"))

	db, err := database.Connect(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	// Initialize routes
	mux := http.NewServeMux()
	handlers.RegisterRoutesWithDB(mux, db)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
