package main

import (
	"log"
	"strconv"

	"sumup-notifications/api/gateway"
)

func main() {
	appConfig, err := gateway.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	e, err := gateway.Bootstrap(appConfig)
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(appConfig.Port)))
}
