package gateway

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                  int    `default:"8080"`
	AWSEndpoint           string `envconfig:"AWS_ENDPOINT" required:"true"`
	NotificationsTopicARN string `envconfig:"NOTIFICATIONS_TOPIC_ARN" required:"true"`
	PostgresURL           string `envconfig:"POSTGRES_URL" required:"true"`
}

func LoadConfig() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
