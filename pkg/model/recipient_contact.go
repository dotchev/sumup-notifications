package model

import "github.com/pkg/errors"

type RecipientContact struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	SlackID     string `json:"slack_id,omitempty"`
}

func (rc *RecipientContact) Validate() error {
	if rc.PhoneNumber == "" && rc.Email == "" && rc.SlackID == "" {
		return errors.New("at least one contact method (phone number, email or Slack ID) must be provided")
	}
	return nil
}
