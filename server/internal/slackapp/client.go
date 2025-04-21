package slackapp

import (
	"os"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	Client        *slack.Client
	SigningSecret string
}

func NewSlackClient() *SlackClient {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

	client := slack.New(botToken)
	return &SlackClient{
		Client:        client,
		SigningSecret: signingSecret,
	}
}
