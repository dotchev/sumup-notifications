package slacksender

import (
	"context"
	"fmt"
	"log"

	"github.com/slack-go/slack"

	"sumup-notifications/pkg/model"
)

type NotificationHandler struct {
	SlackClient *slack.Client
}

func (h NotificationHandler) Handle(ctx context.Context, notification model.NotificationMessage) error {
	if notification.SlackID == "" {
		log.Println("No Slack ID, ignore notification")
		return nil
	}

	log.Printf("Sending Slack message for %s to %s", notification.Recipient, notification.SlackID)
	_, _, err := h.SlackClient.PostMessageContext(
		ctx,
		notification.SlackID,
		slack.MsgOptionText(notification.Message, true), // display text verbatim, no fancy formatting
	)
	if err != nil {
		return fmt.Errorf("could not send Slack message: %w", err)
	}
	return nil
}
