package model

import (
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
