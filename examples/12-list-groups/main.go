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

// 2 groups per page. This number is purposefully low
// so that the example fetches a few pages. A reasonable
// number to use it 100.
const pageSize = 2

func main() {
	api := border0.NewAPIClient(
		client.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
		client.WithRetryMax(2),                                // 1 initial + 2 retries = 3 attempts
	)
	ctx := context.Background()

	// Use the paginator to iterate groups page by page.
	paginator := api.GroupsPaginator(ctx, pageSize)
	total := 0
	for paginator.HasNext() {
		page, err := paginator.Next(ctx)
		if err != nil {
			log.Fatalln("❌ failed to get groups using paginator:", err)
		}
		fmt.Printf("✅ page returned %d groups\n", len(page))
		for _, group := range page {
			output, _ := json.MarshalIndent(group, "", "  ")
			fmt.Printf("✅ group details for %s\n", group.DisplayName)
			fmt.Println(string(output))
		}
		total += len(page)
	}
	fmt.Printf("✅ found %d groups in this Border0 organization\n", total)
}
