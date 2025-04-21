package slackapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

func HandleSlackEvents(w http.ResponseWriter, r *http.Request, signingSecret string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	log.Printf("Received request: %s\n", string(body))

	// Parse the event
	var event slackevents.EventsAPIEvent
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "Failed to parse Slack event", http.StatusInternalServerError)
		return
	}

	// Handle URL Verification Challenge
	if event.Type == slackevents.URLVerification {
		var challengeResponse struct {
			Challenge string `json:"challenge"`
		}
		err := json.Unmarshal(body, &challengeResponse)
		if err != nil {
			http.Error(w, "Failed to parse challenge response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(challengeResponse.Challenge))
		return
	}

	// Handle other event types
	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			fmt.Printf("App was mentioned: %s\n", ev.Text)
			// Respond to the mention
		}
	}
}
