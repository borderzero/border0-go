package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/client"
)

func main() {
	api := border0.NewAPIClient(
		client.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
		client.WithRetryMax(2),                                // 1 initial + 2 retries = 3 attempts
	)
	ctx := context.Background()

	// get all sockets from your organization
	sockets, err := api.Sockets(ctx)
	if err != nil {
		log.Fatalln("❌ failed to get sockets using border0 api client sdk:", err)
	}

	output, _ := json.MarshalIndent(sockets, "", "  ")
	log.Printf("✅ found %d sockets from my Border0 organization = %s", len(sockets), string(output))
}
