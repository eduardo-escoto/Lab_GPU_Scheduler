package handlers

import (
    "net/http"
    "go-webserver-project/internal/services"
)

func SlackHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Handle Slack events
    err := services.HandleSlackEvent(r.Body)
    if err != nil {
        http.Error(w, "Failed to process Slack event", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}