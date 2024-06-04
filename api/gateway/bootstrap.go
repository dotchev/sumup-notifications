package gateway

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"sumup-notifications/pkg/awsclient"
	"sumup-notifications/pkg/storage"
)

func Start(config Config) error {
	ctx := context.Background()

	err := storage.MigrateDB(ctx, config.PostgresURL)
	if err != nil {
		return err
	}

	awsConfig, err := awsclient.LoadConfig(ctx, config.AWSEndpoint)
	if err != nil {
		return err
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
		return err
	}

	nh := NotificationsHandler{
		SNSClient:             snsClient,
		NotificationsTopicARN: config.NotificationsTopicARN,
		db:                    dbPool,
	}
	nh.Mount(e)

	rh := RecipientsHandler{
		db: dbPool,
	}
	rh.Mount(e)

	return e.Start(":" + strconv.Itoa(config.Port))
}
