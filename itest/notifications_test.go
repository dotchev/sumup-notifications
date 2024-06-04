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

func TestNotifications(t *testing.T) {
	t.Run("send notification", func(t *testing.T) {
		recipient := "john"
		resp, _ := putRecipient(t, recipient, model.RecipientContact{PhoneNumber: "+1234567890"})
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		notification := model.Notification{Recipient: recipient, Message: "hello"}
		resp, _ = postNotification(t, notification)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}
