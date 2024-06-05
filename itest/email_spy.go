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

const localstackEmailURL = localstackURL + "/_aws/ses"

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

func resetEmails() error {
	del, err := http.NewRequest(http.MethodDelete, localstackEmailURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(del)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete emails: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func listEmails() ([]Email, error) {
	resp, err := http.Get(localstackEmailURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("could not list emails: unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Messages []struct {
			Destination struct {
				ToAddresses []string `json:"ToAddresses"`
			} `json:"Destination"`
			Source  string `json:"Source"`
			Subject string `json:"Subject"`
			Body    struct {
				TextPart string `json:"text_part"`
			} `json:"Body"`
		} `json:"messages"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var list []Email
	for _, msg := range result.Messages {
		email := Email{
			From:    msg.Source,
			To:      msg.Destination.ToAddresses[0],
			Subject: msg.Subject,
			Body:    msg.Body.TextPart,
		}
		list = append(list, email)
	}
	return list, nil
}

func waitForEmails(t *testing.T, expected []Email) {
	var list []Email
	Retry(t, func() error {
		var err error
		list, err = listEmails()
		if err != nil {
			return retry.Unrecoverable(err)
		}
		if len(list) < len(expected) {
			return errors.New("not enough emails")
		}
		return nil
	})
	assert.Equal(t, expected, list)
}
