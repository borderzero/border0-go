package border0_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/listen"
)

func ExampleNewAPIClient() {
	api := border0.NewAPIClient(
		client.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
		client.WithRetryMax(2),                                // 1 initial + 2 retries = 3 attempts
	)
	ctx := context.Background()

	// create socket
	socket := client.Socket{
		Name:       "sdk-socket-http",
		SocketType: "http",
	}
	created, err := api.CreateSocket(ctx, &socket)
	if err != nil {
		log.Fatalln("failed to create socket using border0 api client sdk:", err)
	}

	log.Printf("created socket: %+v", created)
}

func ExampleListen() {
	listener, err := border0.Listen(
		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("X-Auth-Name") // border0 will set this header along with a few other identity related headers
		fmt.Fprintf(w, "Hello, %s! This is border0-go + standard library http.", name)
	})

	log.Fatalln(http.Serve(listener, handler))
}
