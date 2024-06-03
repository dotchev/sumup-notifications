package gateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"sumup-notifications/pkg/awsconfig"
)

func Bootstrap(config Config) (*echo.Echo, error) {
	ctx := context.Background()

	cfg, err := awsconfig.Load(ctx, config.AWSEndpoint)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(cfg)

	e := echo.New()
	e.Logger.SetHeader("${time_rfc3339} ${level}")
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${error}\n",
	}))
	e.Use(middleware.Recover())

	nh := NotificationHandler{
		SNSClient:             snsClient,
		NotificationsTopicARN: config.NotificationsTopicARN,
	}
	nh.Mount(e)

	return e, nil
}
