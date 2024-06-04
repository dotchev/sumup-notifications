package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"sumup-notifications/pkg/model"
	"sumup-notifications/pkg/storage"
)

type NotificationsHandler struct {
	SNSClient             *sns.Client
	NotificationsTopicARN string
	db                    *pgxpool.Pool
}

func (handler NotificationsHandler) Mount(e *echo.Echo) {
	e.POST("/notifications", echo.HandlerFunc(handler.Post))
}

func (handler NotificationsHandler) Post(c echo.Context) error {
	ctx := c.Request().Context()

	var notification model.Notification
	err := c.Bind(&notification)
	if err != nil {
		return err
	}

	err = notification.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	recipients := storage.Recipients{DB: handler.db}
	_, err = recipients.Load(ctx, notification.Recipient)
	if err != nil {
		if err == storage.ErrNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, "Unknown recipient")
		}
		return fmt.Errorf("failed to load recipient: %w", err)
	}

	err = handler.publishNotification(ctx, notification) // Check the error return value
	if err != nil {
		return fmt.Errorf("failed to publish notification to SNS: %w", err)
	}

	c.Logger().Infof(`Notification sent to %s with message: "%s"`, notification.Recipient, notification.Message)
	return c.JSON(http.StatusCreated, notification)
}

func (handler NotificationsHandler) publishNotification(ctx context.Context, n model.Notification) error {
	jsonData, err := json.Marshal(n)
	if err != nil {
		return err
	}

	_, err = handler.SNSClient.Publish(ctx, &sns.PublishInput{
		TopicArn: &handler.NotificationsTopicARN,
		Message:  aws.String(string(jsonData)),
	})
	return err
}
