package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	nh := NotificationHandler{}
	nh.Register(r.PathPrefix("/notifications").Subrouter())

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func createSNSClient(ctx context.Context) (*sns.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "https://mock.amazonaws.com"}, nil
			})),
	)
	if err != nil {
		return nil, err
	}

	return sns.NewFromConfig(cfg), nil
}
