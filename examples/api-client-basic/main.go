package main

import (
	"context"
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

	// create socket
	socket := client.Socket{
		Name:       "sdk-socket-http",
		SocketType: "http",
	}
	created, err := api.CreateSocket(ctx, &socket)
	if err != nil {
		log.Fatalln("failed to create socket using border0 api client sdk:", err)
	}

	log.Printf("created socket: %+v", created)

	// get socket
	fetched, err := api.Socket(ctx, created.Name)
	if err != nil {
		log.Fatalln("failed to get socket using border0 api client sdk:", err)
	}

	log.Printf("fetched socket: %+v", fetched)
}
