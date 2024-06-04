package emailsender

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"sumup-notifications/pkg/model"
)

type NotificationHandler struct {
	SESClient   *ses.Client
	EmailSender string
}

func (h NotificationHandler) Handle(ctx context.Context, notification model.NotificationMessage) error {
	if notification.Email == "" {
		log.Println("No email, ignore notification")
		return nil
	}

	log.Printf("Sending email for %s to %s", notification.Recipient, notification.Email)
	_, err := h.SESClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{notification.Email},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: &notification.Message,
				},
			},
			Subject: &types.Content{
				Data: aws.String("Notification"),
			},
		},
		Source: &h.EmailSender,
	})
	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}
	return nil
}
