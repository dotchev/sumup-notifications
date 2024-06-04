package emailsender

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"sumup-notifications/pkg/awsclient"
)

func Start(config Config) error {
	ctx := context.Background()

	awsConfig, err := awsclient.LoadConfig(ctx, config.AWSEndpoint)
	if err != nil {
		return fmt.Errorf("error loading AWS config: %w", err)
	}
	sqsClient := sqs.NewFromConfig(awsConfig)
	sesClient := ses.NewFromConfig(awsConfig)

	_, err = sesClient.VerifyEmailIdentity(ctx, &ses.VerifyEmailIdentityInput{
		EmailAddress: &config.EmailSender,
	})
	if err != nil {
		return fmt.Errorf("error verifying SES email identity: %w", err)
	}

	handler := NotificationHandler{
		SESClient:   sesClient,
		EmailSender: config.EmailSender,
	}

	consumer := awsclient.NotificationConsumer{
		SQSClient:           sqsClient,
		QueueName:           config.EmailNotificationsQueue,
		NotificationHandler: handler.Handle,
	}
	return consumer.Start(ctx)
}
