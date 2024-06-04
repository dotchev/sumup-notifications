package awsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"sumup-notifications/pkg/model"
)

type NotificationConsumer struct {
	SQSClient           *sqs.Client
	QueueName           string
	NotificationHandler func(ctx context.Context, notification model.NotificationMessage) error
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	getQueueUrlOutput, err := c.SQSClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &c.QueueName,
	})
	if err != nil {
		return fmt.Errorf("could not get URL for SQS queue %s: %w", c.QueueName, err)
	}
	queueUrl := getQueueUrlOutput.QueueUrl

	log.Printf("listening for messages on %s", *queueUrl)

	for {
		// Terminate the loop if the context is canceled
		if err := ctx.Err(); err != nil {
			return err
		}

		output, err := c.SQSClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            queueUrl,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})
		if err != nil {
			log.Println("Error receiving SQS message:", err)
			continue
		}
		log.Printf("%d messages received from SQS queue", len(output.Messages))

		for _, message := range output.Messages {
			notification, err := parseSNSMessage(*message.Body)
			if err != nil {
				log.Println("Error parsing SNS message:", err)
				continue
			}

			err = c.NotificationHandler(ctx, notification)
			if err != nil {
				log.Println("Error processing notification:", err)
				// SQS will resend the message after the visibility timeout, so we will retry it
				continue
			}

			_, err = c.SQSClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      queueUrl,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Println("Error deleting SQS message:", err)
			}
		}
	}
}

func parseSNSMessage(body string) (model.NotificationMessage, error) {
	var notification model.NotificationMessage

	var bodyData struct {
		Message string `json:"Message"`
	}
	err := json.Unmarshal([]byte(body), &bodyData)
	if err != nil {
		return notification, err
	}

	err = json.Unmarshal([]byte(bodyData.Message), &notification)
	if err != nil {
		return notification, err
	}

	return notification, nil
}
