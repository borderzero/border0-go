# `border0-go` examples

## Prerequisite

Go to [Border0 Admin Portal](https://portal.border0.com) -> Organization Settings -> Access Tokens, create a token in `Member` permission group.

And then either export the token as an environment variable:

```shell
export BORDER0_AUTH_TOKEN=_your_access_token_
go run main.go
```

Or prefix `go run main.go` command with the token as environment variable:

```shell
BORDER0_AUTH_TOKEN=_your_access_token_ go run main.go
```

## List of examples

- [00-create-socket](./00-create-socket): Create an HTTP socket with Border0 API client
- [01-delete-socket](./01-delete-socket): Delete the socket that was created from previous example [00-create-socket].
- [02-echo-listener](./02-echo-listener): Start an echo HTTP server with Border0 tunnel listener.
- [02-fiber-listener](./02-fiber-listener): Start a fiber server with the Border0 listner.
- [02-gin-listener](./02-gin-listener): Using the Border0 listener and serve HTTP requests with gin.
- [02-http-listener](./02-http-listener): Another listener example, and this time we use the http package standard library.
- [03-reverse-proxy](./03-reverse-proxy): Make a reverse proxy in Go, and rewrite the HTTP response body, then serve the requests with Border0 listener.
- [04-list-sockets](./04-list-sockets): List all sockets from your Border0 organization using Border0 API client.
