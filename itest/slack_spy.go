package itest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"testing"

	"github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/assert"
)

const fakeSlackURL = "http://localhost:7000/__admin/requests"

type SlackMessage struct {
	Channel string
	Text    string
}

func resetSlackMessages() error {
	del, err := http.NewRequest(http.MethodDelete, fakeSlackURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(del)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not delete WireMock requests: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func listSlackMessages() ([]SlackMessage, error) {
	resp, err := http.Get(fakeSlackURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("could not list WireMock requests: unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Requests []struct {
			Request struct {
				Body string `json:"body"`
				Time int    `json:"loggedDate"`
			} `json:"request"`
		} `json:"requests"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	sort.Slice(result.Requests, func(i, j int) bool {
		return result.Requests[i].Request.Time < result.Requests[j].Request.Time
	})

	var list []SlackMessage
	for _, req := range result.Requests {
		values, err := url.ParseQuery(req.Request.Body)
		if err != nil {
			return nil, fmt.Errorf(`could not parse request body "%s": %w`, req.Request.Body, err)
		}
		list = append(list, SlackMessage{
			Channel: values.Get("channel"),
			Text:    values.Get("text"),
		})
	}

	return list, nil
}

func waitForSlackMessages(t *testing.T, expected []SlackMessage) {
	var list []SlackMessage
	Retry(t, func() error {
		var err error
		list, err = listSlackMessages()
		if err != nil {
			return retry.Unrecoverable(err)
		}
		if len(list) < len(expected) {
			return errors.New("not enough Slack messages")
		}
		return nil
	})
	assert.Equal(t, expected, list)
}
