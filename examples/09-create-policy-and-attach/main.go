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
	// STEP 0: create a new policy
	//
	policy := client.Policy{
		Name: "sdk-policy",
		PolicyData: client.PolicyData{
			Version: "v1",
			Action:  []string{"database", "ssh", "http", "tls"},
			Condition: client.PolicyCondition{
				Who: client.PolicyWho{
					Email:  []string{"johndoe@example.com"},
					Domain: []string{"example.com"},
				},
				Where: client.PolicyWhere{
					AllowedIP:  []string{"0.0.0.0/0", "::/0"},
					Country:    []string{"NL", "CA", "US", "BR", "FR"},
					CountryNot: []string{"BE"},
				},
				When: client.PolicyWhen{
					After:           "2022-10-13T05:12:27Z",
					Before:          "",
					TimeOfDayAfter:  "00:00 UTC",
					TimeOfDayBefore: "23:59 UTC",
				},
			},
		},
	}
	createdPolicy, err := api.CreatePolicy(ctx, &policy)
	if err != nil {
		log.Fatalln("❌ failed to create policy using border0 api client sdk:", err)
	}

	output, _ := json.MarshalIndent(createdPolicy, "", "  ")
	log.Println("✅ created policy =", string(output))

	//
	// STEP 1: create a new socket
	//
	socket := client.Socket{
		Name:       "another-sdk-socket-http",
		SocketType: "http",
	}
	createdSocket, err := api.CreateSocket(ctx, &socket)
	if err != nil {
		log.Fatalln("❌ failed to create socket using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(createdSocket, "", "  ")
	log.Println("✅ created socket =", string(output))

	//
	// STEP 2: attach the policy to the socket
	//
	err = api.AttachPolicyToSocket(ctx, createdPolicy.ID, createdSocket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to attach policy to socket using border0 api client sdk:", err)
	}

	log.Println("✅ attached policy to socket")

	//
	// STEP 3: get the created socket and listed associated policies
	//
	fetchedSocket, err := api.Socket(ctx, createdSocket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to get socket using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(fetchedSocket, "", "  ")
	log.Println("✅ socket with associated policies =", string(output))
}
