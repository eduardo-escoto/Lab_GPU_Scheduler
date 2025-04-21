package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/users", UsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", UserHandler).Methods("GET")
	r.HandleFunc("/users", CreateUserHandler).Methods("POST")
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/slack/events", SlackEventsHandler)
	mux.HandleFunc("/slack/commands", SlackCommandsHandler)
	mux.HandleFunc("/slack/interactions", SlackInteractionsHandler)
	mux.HandleFunc("/send-email", EmailHandler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Go Webserver!"))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Logic to retrieve and return users
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Logic to retrieve and return a specific user
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Logic to create a new user
}
