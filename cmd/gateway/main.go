package main

import (
	"log"

	"sumup-notifications/api/gateway"
)

func main() {
	appConfig, err := gateway.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Fatal(gateway.Start(appConfig))
}
