package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/labstack/echo/v4"
)

type Notification struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type NotificationHandler struct {
	SNSClient             *sns.Client
	NotificationsTopicARN string
}

func (handler NotificationHandler) Mount(e *echo.Echo) {
	e.POST("/notifications", echo.HandlerFunc(handler.Post))
}

func (handler NotificationHandler) Post(c echo.Context) error {
	ctx := c.Request().Context()

	var n Notification
	err := c.Bind(&n)
	if err != nil {
		return err
	}

	err = handler.publishNotification(ctx, n) // Check the error return value
	if err != nil {
		return fmt.Errorf("failed to publish notification to SNS: %w", err)
	}

	c.Logger().Infof(`Notification sent to %s with message: "%s"`, n.Recipient, n.Message)
	return c.JSON(http.StatusCreated, n)
}

func (handler NotificationHandler) publishNotification(ctx context.Context, n Notification) error {
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
