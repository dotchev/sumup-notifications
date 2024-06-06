package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"sumup-notifications/pkg/awsclient"
	"sumup-notifications/pkg/storage"
)

func Start(config Config) error {
	ctx := context.Background()

	err := storage.MigrateDB(ctx, config.PostgresURL)
	if err != nil {
		return err
	}

	awsConfig, err := awsclient.LoadConfig(ctx, config.AWSEndpoint)
	if err != nil {
		return err
	}

	snsClient := sns.NewFromConfig(awsConfig)

	e := echo.New()
	e.Logger.SetHeader("${time_rfc3339} ${level}")
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${error}\n",
	}))

	dbPool, err := pgxpool.New(ctx, config.PostgresURL)
	if err != nil {
		return err
	}

	nh := NotificationsHandler{
		SNSClient:             snsClient,
		NotificationsTopicARN: config.NotificationsTopicARN,
		db:                    dbPool,
	}
	nh.Mount(e)

	rh := RecipientsHandler{
		db: dbPool,
	}
	rh.Mount(e)

	tlsConfig, err := loadMTLSConfig(config)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:      ":" + strconv.Itoa(config.Port),
		Handler:   e,
		TLSConfig: tlsConfig,
	}
	e.Logger.Infof("Starting server on port %d", config.Port)
	return server.ListenAndServeTLS("", "")
}

func loadMTLSConfig(appConfig Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(appConfig.ServerCertFile, appConfig.ServerKeyFile)
	if err != nil {
		return nil, err
	}

	clientCAPool := x509.NewCertPool()
	caCert, err := os.ReadFile(appConfig.CACertFile)
	if err != nil {
		return nil, err
	}
	clientCAPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    clientCAPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil

}
