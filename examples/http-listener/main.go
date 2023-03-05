package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
)

func main() {
	listener, err := border0.Listen(
		listen.WithSocketName("border0-go-http-listener"),     // http socket name the listener will be bound to
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello world! This is border0-go.")
		},
	)

	log.Fatal(http.Serve(listener, handler))
}
