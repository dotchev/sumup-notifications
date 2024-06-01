package main

import (
	"log"
	"strconv"

	"sumup-notifications/api/gateway"
)

func main() {
	c, err := gateway.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	e, err := gateway.Bootstrap(c)

	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(c.Port)))
}
