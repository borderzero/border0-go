package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/service/connector/types"
)

func main() {
	api := border0.NewAPIClient(
		client.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
		client.WithRetryMax(2),                                // 1 initial + 2 retries = 3 attempts
	)
	ctx := context.Background()

	//
	// STEP 0: create a new connector
	//
	connector := client.Connector{
		Name:        "sdk-connector",
		Description: "Connector created by Border0 SDK",
	}
	createdConnector, err := api.CreateConnector(ctx, &connector)
	if err != nil {
		log.Fatalln("❌ failed to create connector using border0 api client sdk:", err)
	}

	output, _ := json.MarshalIndent(createdConnector, "", "  ")
	log.Println("✅ created connector =", string(output))

	//
	// STEP 1: create a connector token for the connector you just created
	//
	connectorToken := client.ConnectorToken{
		ConnectorID: createdConnector.ConnectorID,
		Name:        "sdk-connector-token",
	}
	createdConnectorToken, err := api.CreateConnectorToken(ctx, &connectorToken)
	if err != nil {
		log.Fatalln("❌ failed to create connector token using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(createdConnectorToken, "", "  ")
	log.Println("✅ created connector token =", string(output))

	//
	// STEP 2: create a socket that's linked to the connector you just created
	//
	socket := client.Socket{
		Name:        "sdk-socket-ssh",
		SocketType:  "ssh",
		ConnectorID: createdConnector.ConnectorID,
		UpstreamConfig: &types.ConnectorServiceUpstreamConfig{ // SSH upstream config with username and password
			UpstreamConnectionType: types.UpstreamConnectionTypeSSH,
			BaseUpstreamDetails: types.BaseUpstreamDetails{
				Hostname: "127.0.0.1",
				Port:     22,
			},
			SSHConfiguration: &types.SSHConfiguration{
				UpstreamAuthenticationType: types.UpstreamAuthenticationTypeUsernamePassword,
				BasicCredentials: &types.BasicCredentials{
					Username: "test-username",
					Password: "test-password",
				},
			},
		},
	}
	createdSocket, err := api.CreateSocket(ctx, &socket)
	if err != nil {
		log.Fatalln("❌ failed to create socket using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(createdSocket, "", "  ")
	log.Println("✅ created socket =", string(output))

	//
	// STEP 3: fetch the socket's linked connectors
	//
	fetchedConnectors, err := api.SocketConnectors(ctx, createdSocket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to fetch socket connectors using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(fetchedConnectors, "", "  ")
	log.Println("✅ fetched socket connectors =", string(output))

	//
	// STEP 4: fetch the socket's linked upstream configs
	//
	fetchedUpstreamConfigs, err := api.SocketUpstreamConfigs(ctx, createdSocket.SocketID)
	if err != nil {
		log.Fatalln("❌ failed to fetch socket upstream configs using border0 api client sdk:", err)
	}

	output, _ = json.MarshalIndent(fetchedUpstreamConfigs, "", "  ")
	log.Println("✅ fetched socket upstream configs =", string(output))
}
