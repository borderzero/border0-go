# Border0 Go SDK

[![Run tests](https://github.com/borderzero/border0-go/actions/workflows/run_tests.yml/badge.svg)](https://github.com/borderzero/border0-go/actions/workflows/run_tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/borderzero/border0-go.svg)](https://pkg.go.dev/github.com/borderzero/border0-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/borderzero/border0-go)](https://goreportcard.com/report/github.com/borderzero/border0-go)
[![license](https://img.shields.io/github/license/borderzero/border0-go)](https://github.com/borderzero/border0-go/blob/master/LICENSE)

Border0 enables users to log into various services, including web, SSH, database, and generic TCP, using their existing Single Sign-On (SSO) credentials.
If you haven't yet registered, [create a new account](https://portal.border0.com/register) and explore our [informative blog posts](https://www.border0.com/blog)
and [comprehensive documentation](https://docs.border0.com/docs/quick-start).

This SDK contains 2 major components:

- [Border0 API Client](./client): provides API client methods that interact with our API to manage Border0 resources. See [examples](./examples)
  folder for a basic example of how to manage Border0 resources using this API client.
- [Border0 Listen](./listen): `border0.Listen` creates a Go `net.listener`, that can be used to accept incoming connections. When the
  listener is passed to http.Serve, the server will accept HTTP requests sent by Border0 and forward them to an HTTP handler. The handler's
  response will be sent back to Border0. See [examples](./examples) folder for a simple and advanced examples of how to use the `border0.Listen`
  function.

## Installation

```shell
go get github.com/borderzero/border0-go
```

## Quickstart

Explore the [examples](./examples) folder for additional use cases and examples. To run these examples, you'll need a Border0 access token.
You can generate one by going to [Border0 Admin Portal](https://portal.border0.com) -> Organization Settings -> Access Tokens, create a token in `Member` or `Admin` permission groups.

Once you have the token, you can proceed to run the example code with:

```shell
BORDER0_AUTH_TOKEN=_your_access_token_ go run main.go
```

### Border API Client

Create an HTTP socket using `border0-go`:

```go
package main

import (
	"context"
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
```

### `border0.Listen`

The following example will:

- Automatically create an HTTP socket with name `sdk-socket-http`
- Connect to Border0 and return a Border0 `net.Listener`
- Serve HTTP requests that are sent by Border0 right from your local machine

```go
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
```

Also see `border0.Listen` examples for [gin](./examples/05-gin-listener/main.go), [echo](./examples/05-echo-listener/main.go) and [fiber](./examples/05-fiber-listener/main.go).
