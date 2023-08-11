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
	// STEP 0: get the socket (sdk-socket-ssh) created from the previous example (05-create-connector-and-socket)
	//
	previousExample := "05-create-connector-and-socket"
	socketFromPreviousExample := "sdk-socket-ssh"
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

	//
	// STEP 1: get socket's linked connectors (should be 1 connector)
	//
	connectors, err := api.SocketConnectors(ctx, socket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to get socket connectors using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(connectors, "", "  ")
	log.Println("✅ socket connectors =", string(output))

	//
	// STEP 2: delete the socket (sdk-socket-ssh)
	//
	err = api.DeleteSocket(ctx, socket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to delete socket using border0 api client sdk:", err)
	}

	log.Printf("❎ socket [%s] deleted", socket.Name)

	//
	// STEP 3: delete the connector, associated connector tokens will be deleted automatically
	//
	for _, connector := range connectors.List {
		err = api.DeleteConnector(ctx, connector.ConnectorID)
		if err != nil {
			log.Fatalln("failed to delete connector using border0 api client sdk:", err)
		}

		log.Printf("❎ connector [%s] deleted", connector.ConnectorName)
	}
}
