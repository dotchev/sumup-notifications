package smssender

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"sumup-notifications/pkg/awsconfig"
)

func Start(appConfig Config) error {
	ctx := context.Background()

	cfg, err := awsconfig.Load(ctx, appConfig.AWSEndpoint)
	if err != nil {
		return fmt.Errorf("Error loading AWS config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)

	getQueueUrlOutput, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &appConfig.SMSNotificationsQueue,
	})
	if err != nil {
		return fmt.Errorf("Could not get URL for SQS queue %s: %w", appConfig.SMSNotificationsQueue, err)
	}
	queueUrl := getQueueUrlOutput.QueueUrl

	log.Printf("SMS Sender started, listening for messages on %s", *queueUrl)

	return processSQSMessages(ctx, sqsClient, snsClient, queueUrl)
}
