package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/listen"
)

func main() {
	if err := ensurePolicyCreated("sdk-created-policy"); err != nil {
		log.Fatalln("failed to create policy:", err)
	}

	listener, err := border0.Listen(
		// http socket name the listener will be bound to, socket will be created if not exists
		listen.WithSocketName("sdk-socket-http"),

		// optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")),

		// Let's attach a policy; make sure this policy exist
		listen.WithPolicies([]string{"sdk-created-policy"}),
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// border0 will set this header along with a few other identity related headers
		name := r.Header.Get("X-Auth-Name")
		email := r.Header.Get("X-Auth-Email")
		fmt.Fprintf(w, "Hello, %s %s! This is Border0-go + standard library http.", name, email)
	})

	log.Fatalln(http.Serve(listener, handler))
}

func ensurePolicyCreated(name string) error {
	ctx := context.Background()
	api := border0.NewAPIClient(
		// optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
		client.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")),

		// no retry, just 1 attempt
		client.WithRetryMax(0),
	)

	policy := client.Policy{
		Name: name,
		PolicyData: client.PolicyData{
			Action: []string{"http"},
			Condition: client.PolicyCondition{
				Who: client.PolicyWho{
					Domain: []string{"gmail.com"}, // any email address with gmail.com domain will be allowed
				},
				When: client.PolicyWhen{
					After:  time.Now().Format(time.RFC3339),
					Before: time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339),
				},
			},
		},
	}

	if _, err := api.CreatePolicy(ctx, &policy); err != nil {
		if strings.Contains(err.Error(), "policy already exists") {
			// we can ignore the error if the policy already exists
			return nil
		}
		return err
	}

	return nil
}
