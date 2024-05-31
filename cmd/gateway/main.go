package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gorilla/mux"
)

type Notification struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

func main() {

	r := SetupRouter()

	// TODO graceful shutdown, see https://pkg.go.dev/github.com/gorilla/mux@v1.8.1#readme-graceful-shutdown
	log.Fatal(http.ListenAndServe(":8080", r))
}

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/notifications", PostNotification).Methods(http.MethodPost)
	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func PostNotification(w http.ResponseWriter, r *http.Request) {
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

func createSNSClient() *sns.SNS {

	// TODO: Implement the logic to create an SNS client
	// You need to provide your AWS credentials and configure the region
	// Example:
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-west-2"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sns.New(sess)

}
