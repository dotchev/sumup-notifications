package main

import (
	"log"
	"net/http"

	"sumup-notifications/api/gateway"
)

func main() {

	r := gateway.SetupRouter()

	// TODO graceful shutdown, see https://pkg.go.dev/github.com/gorilla/mux@v1.8.1#readme-graceful-shutdown
	log.Fatal(http.ListenAndServe(":8080", r))
}
