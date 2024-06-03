package main

import (
	"log"

	"sumup-notifications/api/smssender"
)

func main() {
	appConfig, err := smssender.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Fatal(smssender.Start(appConfig))
}
