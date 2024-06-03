package gateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"sumup-notifications/pkg/awsconfig"
	"sumup-notifications/pkg/storage"
)

func Bootstrap(config Config) (*echo.Echo, error) {
	ctx := context.Background()

	err := storage.MigrateDB(ctx, config.PostgresURL)
	if err != nil {
		return nil, err
	}

	awsConfig, err := awsconfig.Load(ctx, config.AWSEndpoint)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(awsConfig)

	e := echo.New()
	e.Logger.SetHeader("${time_rfc3339} ${level}")
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${error}\n",
	}))

	dbPool, err := pgxpool.New(ctx, config.PostgresURL)
	if err != nil {
		return nil, err
	}

	nh := NotificationHandler{
		SNSClient:             snsClient,
		NotificationsTopicARN: config.NotificationsTopicARN,
		dbPool:                dbPool,
	}
	nh.Mount(e)

	return e, nil
}
