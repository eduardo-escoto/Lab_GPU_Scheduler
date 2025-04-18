package slackapp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func HandleSlashCommands(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		http.Error(w, "Failed to parse slash command", http.StatusInternalServerError)
		return
	}

	log.Printf("Received request: %s\n", string(s.Command))
	switch s.Command {
	case "/gpu-request":
		// Handle the GPU request command
		response := fmt.Sprintf("Received GPU request: %s", s.Text)
		w.Write([]byte(response))
	case "/schedule":
		response := fmt.Sprintf("Received Schedule Request: %s", s.Text)
		w.Write(([]byte(response)))
	default:
		http.Error(w, "Unknown command", http.StatusBadRequest)
	}
}
