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

	//
	// STEP 0: get the socket created from the previous example
	//
	previousExample := "00-create-socket"
	socketFromPreviousExample := "sdk-socket-http"
	fetched, err := api.Socket(ctx, socketFromPreviousExample)
	if err != nil {
		if client.NotFound(err) {
			log.Fatalf("⚠️  socket [%s] not found, run example [%s] first", socketFromPreviousExample, previousExample)
		} else {
			log.Fatalln("❌ failed to get socket using border0 api client sdk:", err)
		}
	}

	output, _ := json.MarshalIndent(fetched, "", "  ")
	log.Printf("✅ socket from previous example [%s] = %s", previousExample, string(output))

	// Update the description of the socket
	fetched.Description = "updated description"

	s, err := api.UpdateSocket(ctx, socketFromPreviousExample, fetched)
	if err != nil {
		log.Fatalln("❌ failed to delete socket using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(s, "", "  ")
	log.Printf("✅ socket from previous example [%s] = %s", previousExample, string(output))

}
