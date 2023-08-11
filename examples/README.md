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
- [01-update-socket](./01-update-socket): Update the descriton of the socket that was created in the previous example [00-create-socket].
- [02-list-socket](./02-list-socket): List all your sockets and print the details of each socket to the console.
- [03-delete-socket](./03-delete-socket): Delete the socket that was created from previous example [00-create-socket].
- [04-reverse-proxy](./04-reverse-proxy): Make a reverse proxy in Go, and rewrite the HTTP response body, then serve the requests with Border0 listener.
- [05-http-listener](./05-http-listener): Border0 tunnel listener example, and this time we use the http package standard library.
- [05-gin-listener](./05-gin-listener): Border0 tunnel listenerand serve HTTP requests with gin.
- [05-echo-listener](./05-echo-listener): Start an echo HTTP server with Border0 tunnel listener.
- [05-fiber-listener](./05-fiber-listener): Start a fiber server with the Border0 listner.
- [06-create-connector-and-socket](./05-create-connector-and-socket): Create a connector, a connector token and an SSH socket that's linked to the connector with upstream config.
- [07-delete-connector-and-socket](./06-delete-connector-and-socket): Delete the resources created from the previous example (a connector, a connector token and an SSH socket).
