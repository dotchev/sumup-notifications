package smssender

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"

	"sumup-notifications/pkg/model"
)

type NotificationHandler struct {
	SNSClient *sns.Client
	DB        *pgxpool.Pool
}

func (h NotificationHandler) Handle(ctx context.Context, notification model.NotificationMessage) error {
	if notification.PhoneNumber == "" {
		log.Println("No phone number, ignore notification")
		return nil
	}

	log.Printf("Sending SMS for %s to %s", notification.Recipient, notification.PhoneNumber)
	_, err := h.SNSClient.Publish(ctx, &sns.PublishInput{
		Message:     &notification.Message,
		PhoneNumber: &notification.PhoneNumber,
	})
	if err != nil {
		return fmt.Errorf("could not send SMS: %w", err)
	}
	return nil
}
