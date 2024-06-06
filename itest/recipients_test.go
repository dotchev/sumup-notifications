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

func putRecipient(t *testing.T, recipient string, contact model.RecipientContact) (*http.Response, string) {
	body, err := json.Marshal(contact)
	require.NoError(t, err)
	put, err := http.NewRequest(http.MethodPut, gatewayURL+"/recipients/"+recipient, bytes.NewReader(body))
	require.NoError(t, err)
	put.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(put)
	require.NoError(t, err)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	if http.StatusNoContent == resp.StatusCode {
		t.Cleanup(func() {
			deleteRecipient(t, recipient)
		})
	}

	return resp, bodyString
}

func deleteRecipient(t *testing.T, recipient string) *http.Response {
	del, err := http.NewRequest(http.MethodDelete, gatewayURL+"/recipients/"+recipient, nil)
	require.NoError(t, err)
	resp, err := httpClient.Do(del)
	require.NoError(t, err)
	return resp
}

func getRecipient(t *testing.T, recipient string) (*http.Response, model.RecipientContact) {
	resp, err := httpClient.Get(gatewayURL + "/recipients/" + recipient)
	require.NoError(t, err)
	defer resp.Body.Close()

	var contact model.RecipientContact
	if http.StatusOK == resp.StatusCode {
		err = json.NewDecoder(resp.Body).Decode(&contact)
		require.NoError(t, err)
	}
	return resp, contact
}

func TestRecipients_Get(t *testing.T) {
	t.Run("unknown recipient", func(t *testing.T) {
		resp, _ := getRecipient(t, "none")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("recipient exists", func(t *testing.T) {
		contact := model.RecipientContact{PhoneNumber: "+1234567890"}
		resp, _ := putRecipient(t, "john", contact)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, respContact := getRecipient(t, "john")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, contact, respContact)
	})
}

func TestRecipients_Put(t *testing.T) {
	t.Run("invalid recipient contact", func(t *testing.T) {
		contact := model.RecipientContact{}
		resp, body := putRecipient(t, "john", contact)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, body, "at least one contact method (phone number, email or Slack ID) must be provided")
	})

	t.Run("update recipient", func(t *testing.T) {
		contact := model.RecipientContact{PhoneNumber: "+1234567890"}
		resp, _ := putRecipient(t, "john", contact)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, respContact := getRecipient(t, "john")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, contact, respContact)

		contact.PhoneNumber = "+0987654321"
		resp, _ = putRecipient(t, "john", contact)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, respContact = getRecipient(t, "john")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, contact, respContact)
	})
}

func TestRecipients_Delete(t *testing.T) {
	t.Run("delete unknown recipient", func(t *testing.T) {
		resp, _ := getRecipient(t, "john")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		resp = deleteRecipient(t, "john")
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("delete recipient", func(t *testing.T) {
		contact := model.RecipientContact{PhoneNumber: "+1234567890"}
		resp, _ := putRecipient(t, "john", contact)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, respContact := getRecipient(t, "john")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, contact, respContact)

		resp = deleteRecipient(t, "john")
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, _ = getRecipient(t, "john")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
