package smssender

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"sumup-notifications/pkg/awsclient"
)

func Start(config Config) error {
	ctx := context.Background()

	awsConfig, err := awsclient.LoadConfig(ctx, config.AWSEndpoint)
	if err != nil {
		return fmt.Errorf("Error loading AWS config: %w", err)
	}
	sqsClient := sqs.NewFromConfig(awsConfig)
	snsClient := sns.NewFromConfig(awsConfig)

	handler := NotificationHandler{
		SNSClient: snsClient,
	}

	consumer := awsclient.NotificationConsumer{
		SQSClient:           sqsClient,
		QueueName:           config.SMSNotificationsQueue,
		NotificationHandler: handler.Handle,
	}
	return consumer.Start(ctx)
}
