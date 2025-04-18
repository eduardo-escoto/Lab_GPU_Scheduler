package handlers

import (
	"go-webserver-project/internal/slackapp"
	"net/http"
)

func SlackEventsHandler(w http.ResponseWriter, r *http.Request) {
	slackapp.HandleSlackEvents(w, r, slackapp.NewSlackClient().SigningSecret)
}

func SlackCommandsHandler(w http.ResponseWriter, r *http.Request) {
	slackapp.HandleSlashCommands(w, r)
}

func SlackInteractionsHandler(w http.ResponseWriter, r *http.Request) {
	slackapp.HandleInteractions(w, r)
}
