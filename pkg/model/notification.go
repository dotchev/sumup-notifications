package model

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Notification struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

func (n Notification) Validate() error {
	if n.Recipient == "" {
		return errors.New("missing recipient")
	}
	if n.Message == "" {
		return errors.New("missing message")
	}
	return nil
}

func ParseSNSMessage(body string) (Notification, error) {
	var notification Notification

	var bodyData struct {
		Message string `json:"Message"`
	}
	err := json.Unmarshal([]byte(body), &bodyData)
	if err != nil {
		return notification, err
	}

	err = json.Unmarshal([]byte(bodyData.Message), &notification)
	if err != nil {
		return notification, err
	}

	return notification, nil
}
