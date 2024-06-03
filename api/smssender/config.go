package smssender

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AWSEndpoint           string `envconfig:"AWS_ENDPOINT" required:"true"`
	NotificationsTopicARN string `envconfig:"NOTIFICATIONS_TOPIC_ARN" required:"true"`
	SMSNotificationsQueue string `envconfig:"SMS_NOTIFICATIONS_QUEUE" default:"sms_notifications"`
}

func LoadConfig() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
