package services

import (
    "encoding/json"
    "io"
    "log"
)

func HandleSlackEvent(body io.ReadCloser) error {
    defer body.Close()

    var event map[string]interface{}
    if err := json.NewDecoder(body).Decode(&event); err != nil {
        log.Printf("Failed to decode Slack event: %v", err)
        return err
    }

    log.Printf("Received Slack event: %v", event)
    // Process the event here
    return nil
}