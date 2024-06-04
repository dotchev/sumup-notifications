package main

import (
	"log"

	"sumup-notifications/api/emailsender"
)

func main() {
	appConfig, err := emailsender.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Fatal(emailsender.Start(appConfig))
}
