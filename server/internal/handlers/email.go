// package handlers

// import (
// 	"net/http"

// 	"github.com/eduardo-escoto/gpu_request/server/internal/services"
// )

// func EmailHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Send email
// 	err := services.SendEmail("recipient@example.com", "Subject", "Email body")
// 	if err != nil {
// 		http.Error(w, "Failed to send email", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Email sent successfully"))
// }
