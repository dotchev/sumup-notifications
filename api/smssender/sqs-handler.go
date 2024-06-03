package smssender

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"sumup-notifications/pkg/model"
)

func processSQSMessages(ctx context.Context, sqsClient *sqs.Client, snsClient *sns.Client, queueUrl *string) error {
	for {
		// Terminate the loop if the context is canceled
		if err := ctx.Err(); err != nil {
			return err
		}

		output, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
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
			log.Println("Message body:", *message.Body)

			notification, err := model.ParseSNSMessage(*message.Body)
			if err != nil {
				log.Println("Error parsing SNS message:", err)
				continue
			}

			phoneNumber := lookupPhoneNumber(notification.Recipient)
			_, err = snsClient.Publish(ctx, &sns.PublishInput{
				Message:     message.Body,
				PhoneNumber: &phoneNumber,
			})
			if err != nil {
				log.Println("Error sending SMS:", err)
				// SQS will resend the message after the visibility timeout, so we will retry it
				continue
			}

			_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      queueUrl,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Println("Error deleting SQS message:", err)
			}
		}
	}
}

func lookupPhoneNumber(recipient string) string {
	// This is a placeholder for a real lookup implementation
	return "+1234567890"
}
