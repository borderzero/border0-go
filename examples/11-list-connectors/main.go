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

	// get all connectors from your organization
	connectors, err := api.Connectors(ctx)
	if err != nil {
		log.Fatalln("❌ failed to get connectors using border0 api client sdk:", err)
	}
	fmt.Printf("✅ found %d connectors in this Border0 organization\n", len(connectors.List))

	//  print the details of each socket
	for _, connector := range connectors.List {
		output, _ := json.MarshalIndent(connector, "", "  ")
		fmt.Printf("✅ connector details for %s\n", connector.Name)
		fmt.Println(string(output))
	}
}
