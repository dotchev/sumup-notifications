package model

import (
	"encoding/json"
)

type Notification struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
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
