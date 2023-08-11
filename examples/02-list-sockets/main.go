package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	fmt.Printf("✅ found %d sockets in this Border0 organization", len(sockets))

	//  print the details of each socket
	for _, socket := range sockets {
		output, _ := json.MarshalIndent(socket, "", "  ")
		fmt.Printf("✅ socket details for %s\n", socket.Name)
		fmt.Println(string(output))
	}

}
