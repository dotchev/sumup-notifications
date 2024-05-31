package gateway

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Notification struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type NotificationHandler struct {
}

func (handler NotificationHandler) Register(router *mux.Router) {
	router.HandleFunc("", handler.Post).Methods(http.MethodPost)
}

func (handler NotificationHandler) Post(w http.ResponseWriter, r *http.Request) {
	var n Notification
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Publish the notification to the SNS topic "notifications"
	publishNotification(n)

	w.WriteHeader(http.StatusCreated)
	log.Printf(`Notification sent to %s with message: "%s"`, n.Recipient, n.Message)
}

func publishNotification(n Notification) {
	// TODO: Implement the logic to publish the notification to the SNS topic "notifications"
	// You can use the AWS SDK or any other library to interact with SNS
	// Example:
	// snsClient := createSNSClient()
	// snsClient.Publish(&sns.PublishInput{
	//     TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:notifications"),
	//     Message:  aws.String(fmt.Sprintf(`{"recipient": "%s", "message": "%s"}`, n.Recipient, n.Message)),
	// })
}
