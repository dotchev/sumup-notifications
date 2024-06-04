package itest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/assert"
)

type SMS struct {
	PhoneNumber string
	Message     string
}

func resetSMS() error {
	del, err := http.NewRequest(http.MethodDelete, localstackURL+"/_aws/sns/sms-messages", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(del)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete SMS messages: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func listSMS() ([]SMS, error) {
	resp, err := http.Get(localstackURL + "/_aws/sns/sms-messages")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("could not list SMS messages: unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Messages map[string][]SMS `json:"sms_messages"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var list []SMS
	for _, messages := range result.Messages {
		list = append(list, messages...)
	}
	return list, nil
}

func waitForSMS(t *testing.T, expected []SMS) {
	var list []SMS
	Retry(t, func() error {
		var err error
		list, err = listSMS()
		if err != nil {
			return retry.Unrecoverable(err)
		}
		if len(list) < len(expected) {
			return errors.New("not enough SMS messages")
		}
		return nil
	})
	assert.Equal(t, expected, list)
}
