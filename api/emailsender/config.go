package emailsender

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AWSEndpoint             string `envconfig:"AWS_ENDPOINT" required:"true"`
	EmailNotificationsQueue string `envconfig:"EMAIL_NOTIFICATIONS_QUEUE" default:"email_notifications"`
	EmailSender             string `envconfig:"EMAIL_SENDER" default:"sender@notifications.com"`
}

func LoadConfig() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
