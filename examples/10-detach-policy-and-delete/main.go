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
	// STEP 0: get the socket (another-sdk-socket-http) created from the previous example (09-create-policy-and-attach)
	//
	previousExample := "09-create-policy-and-attach"
	socketFromPreviousExample := "another-sdk-socket-http"
	socket, err := api.Socket(ctx, socketFromPreviousExample)
	if err != nil {
		if client.NotFound(err) {
			log.Fatalf("⚠️  socket [%s] not found, run example [%s] first", socketFromPreviousExample, previousExample)
		} else {
			log.Fatalln("❌ failed to get socket using border0 api client sdk:", err)
		}
	}

	output, _ := json.MarshalIndent(socket, "", "  ")
	log.Printf("✅ socket from previous example [%s] = %s", previousExample, string(output))

	for _, policy := range socket.Policies {
		//
		// STEP 1: print the socket's associated policies
		//
		output, _ = json.MarshalIndent(policy, "", "  ")
		log.Printf("✅ socket associated policy [%s] = %s", policy.Name, string(output))

		//
		// STEP 2: detach the policy from the socket
		//
		err = api.RemovePolicyFromSocket(ctx, policy.ID, socket.SocketID)
		if err != nil {
			log.Fatalln("❌ failed to detach policy from socket using border0 api client sdk:", err)
		}

		log.Printf("✅ policy [%s] detached from socket [%s]", policy.Name, socket.Name)

		//
		// STEP 3: delete the policy
		//
		err = api.DeletePolicy(ctx, policy.ID)
		if err != nil {
			log.Fatalln("❌ failed to delete policy using border0 api client sdk:", err)
		}

		log.Printf("✅ policy [%s] deleted", policy.Name)
	}

	//
	// STEP 4: delete the socket
	//
	err = api.DeleteSocket(ctx, socket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to delete socket using border0 api client sdk:", err)
	}

	log.Printf("✅ socket [%s] deleted", socket.Name)
}
