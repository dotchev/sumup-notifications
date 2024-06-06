package itest

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

const (
	gatewayURL    = "https://localhost:8443"
	localstackURL = "http://localhost:4566"
)

var httpClient = createHTTPClient()

func createHTTPClient() *http.Client {
	caCert, err := os.ReadFile("../mtls/ca.crt.pem")
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair("../mtls/client.crt.pem", "../mtls/client.key.pem")
	if err != nil {
		panic(err)
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
			},
		},
	}
}
