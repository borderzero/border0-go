// Package client provides API client methods that interact with our API to manage Border0 resources.
//
// Example to create a new client:
//
//	api := client.New(
//		client.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
//		client.WithRetryMax(2),                                // 1 initial + 2 retries = 3 attempts
//	)
//
// See [Option] for more configurable options.
package client
