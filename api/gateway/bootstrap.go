package gateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func createAWSConfig(ctx context.Context, awsEndpoint string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: awsEndpoint}, nil
			})),
	)
	return cfg, err
}

func Bootstrap(config Config) (*echo.Echo, error) {
	ctx := context.Background()

	cfg, err := createAWSConfig(ctx, config.AWSEndpoint)
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
