package main

import (
	"log"

	"sumup-notifications/api/slacksender"
)

func main() {
	appConfig, err := slacksender.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Fatal(slacksender.Start(appConfig))
}
