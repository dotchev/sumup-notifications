package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadConfig(ctx context.Context) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint != "" {
		opts = append(opts,
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: awsEndpoint}, nil
				})),
		)
	}

	return config.LoadDefaultConfig(ctx, opts...)
}
