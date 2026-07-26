package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/borderzero/border0-go/client"
)

func main() {
	authToken := os.Getenv("BORDER0_AUTH_TOKEN")
	if authToken == "" {
		fmt.Println("Warning: BORDER0_AUTH_TOKEN is not set. IsAuthenticated should return false.")
	}

	api := client.New()
	ctx := context.Background()

	isAuthenticated, err := api.IsAuthenticated(ctx)
	if err != nil {
		log.Fatalf("IsAuthenticated returned error: %v", err)
	}

	if isAuthenticated {
		fmt.Println("User IS authenticated.")
	} else {
		fmt.Println("User is NOT authenticated.")
	}
}
