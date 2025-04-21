package slackapp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
)

func HandleInteractions(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Failed to parse interaction payload", http.StatusInternalServerError)
		return
	}

	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		fmt.Printf("Button clicked: %v\n", payload.ActionCallback.BlockActions)
		// Handle button clicks
	case slack.InteractionTypeViewSubmission:
		fmt.Printf("Modal submitted: %v\n", payload.View)
		// Handle modal submissions
	}
}
