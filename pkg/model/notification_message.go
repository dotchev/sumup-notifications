package model

type NotificationMessage struct {
	Notification     `json:"notification"`
	RecipientContact `json:"recipient_contact"`
}
