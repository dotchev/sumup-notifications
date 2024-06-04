package itest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sumup-notifications/pkg/model"
)

func postNotification(t *testing.T, notification model.Notification) (*http.Response, string) {
	body, err := json.Marshal(notification)
	require.NoError(t, err)
	resp, err := http.Post(gatewayURL+"/notifications", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	if http.StatusCreated == resp.StatusCode {
		var respNotification model.Notification
		err = json.Unmarshal(bodyBytes, &respNotification)
		require.NoError(t, err)
		assert.Equal(t, notification, respNotification)
	}

	return resp, bodyString
}

func TestNotifications_Validation(t *testing.T) {
	tests := []struct {
		name         string
		notification model.Notification
		expectedCode int
		expectedBody string
	}{
		{
			name:         "missing recipient",
			notification: model.Notification{},
			expectedCode: http.StatusBadRequest,
			expectedBody: "missing recipient",
		},
		{
			name:         "missing message",
			notification: model.Notification{Recipient: "john"},
			expectedCode: http.StatusBadRequest,
			expectedBody: "missing message",
		},
		{
			name:         "unknown recipient",
			notification: model.Notification{Recipient: "none", Message: "hello"},
			expectedCode: http.StatusBadRequest,
			expectedBody: "unknown recipient",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := postNotification(t, tt.notification)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			assert.Contains(t, body, tt.expectedBody)
		})
	}
}

func TestNotifications_SMS(t *testing.T) {
	t.Run("single one", func(t *testing.T) {
		require.NoError(t, resetSMS())

		recipient := "john"
		sms := SMS{PhoneNumber: "+1234567890", Message: "hello"}

		resp, _ := putRecipient(t, recipient, model.RecipientContact{PhoneNumber: sms.PhoneNumber})
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		notification := model.Notification{Recipient: recipient, Message: sms.Message}
		resp, _ = postNotification(t, notification)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		waitForSMS(t, []SMS{sms})
	})

	t.Run("multiple recipients", func(t *testing.T) {
		require.NoError(t, resetSMS())

		john := "john"
		smsJohn := SMS{PhoneNumber: "+1234567890", Message: "hello"}

		resp, _ := putRecipient(t, john, model.RecipientContact{PhoneNumber: smsJohn.PhoneNumber})
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		jane := "jane"
		smsJane := SMS{PhoneNumber: "+9876543210", Message: "world"}

		resp, _ = putRecipient(t, jane, model.RecipientContact{PhoneNumber: smsJane.PhoneNumber})
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		notification := model.Notification{Recipient: john, Message: smsJohn.Message}
		resp, _ = postNotification(t, notification)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		notification = model.Notification{Recipient: jane, Message: smsJane.Message}
		resp, _ = postNotification(t, notification)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		waitForSMS(t, []SMS{smsJohn, smsJane})
	})
}
