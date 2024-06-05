package slacksender

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/slack-go/slack"

	"sumup-notifications/pkg/awsclient"
)

func Start(config Config) error {
	ctx := context.Background()

	awsConfig, err := awsclient.LoadConfig(ctx, config.AWSEndpoint)
	if err != nil {
		return fmt.Errorf("error loading AWS config: %w", err)
	}
	sqsClient := sqs.NewFromConfig(awsConfig)

	slackClient := slack.New(
		config.SlackToken,
		slack.OptionAPIURL(config.SlackAPIURL),
		slack.OptionLog(logger{}),
	)

	handler := NotificationHandler{SlackClient: slackClient}

	consumer := awsclient.NotificationConsumer{
		SQSClient:           sqsClient,
		QueueName:           config.SlackNotificationsQueue,
		NotificationHandler: handler.Handle,
	}
	return consumer.Start(ctx)
}

func Log(int, msg string) error {
	log.Println(msg)
	return nil
}

type logger struct{}

func (l logger) Output(n int, msg string) error {
	log.Println(n, msg)
	return nil
}
