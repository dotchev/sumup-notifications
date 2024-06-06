package gateway

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                  int    `default:"8443"`
	AWSEndpoint           string `envconfig:"AWS_ENDPOINT" required:"true"`
	NotificationsTopicARN string `envconfig:"NOTIFICATIONS_TOPIC_ARN" required:"true"`
	PostgresURL           string `envconfig:"POSTGRES_URL" required:"true"`
	ServerCertFile        string `envconfig:"SERVER_CERT_FILE" default:"gateway.crt.pem"`
	ServerKeyFile         string `envconfig:"SERVER_KEY_FILE" default:"gateway.key.pem"`
	CACertFile            string `envconfig:"CA_CERT_FILE" default:"ca.crt.pem"`
}

func LoadConfig() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return c, err
}
