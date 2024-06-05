package slacksender

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AWSEndpoint             string `envconfig:"AWS_ENDPOINT" required:"true"`
	SlackAPIURL             string `envconfig:"SLACK_API_URL" required:"true"`
	SlackToken              string `envconfig:"SLACK_TOKEN" required:"true"`
	SlackNotificationsQueue string `envconfig:"SLACK_NOTIFICATIONS_QUEUE" default:"slack_notifications"`
}

func LoadConfig() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
