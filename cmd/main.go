package main

import (
    "log"
    "net/http"
    "go-webserver-project/internal/handlers"
    "go-webserver-project/config"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize routes
    mux := http.NewServeMux()
    handlers.RegisterRoutes(mux)

    // Start the server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}