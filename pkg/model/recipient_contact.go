package model

import "github.com/pkg/errors"

type RecipientContact struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (rc *RecipientContact) Validate() error {
	if rc.PhoneNumber == "" && rc.Email == "" {
		return errors.New("at least one contact method (phone number or email) must be provided")
	}
	return nil
}
