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
- [05-load-balancer](./05-load-balancer): a http load balancer example with Border0 listener. The load balancer will forward the requests to the upstream servers in round-robin fashion and rewrites the response body so that it contains the server's hostname and user's full name.
- [06-http-listener](./06-http-listener): Border0 tunnel listener example, and this time we use the http package standard library.
- [06-gin-listener](./06-gin-listener): Border0 tunnel listenerand serve HTTP requests with gin.
- [06-echo-listener](./06-echo-listener): Start an echo HTTP server with Border0 tunnel listener.
- [06-fiber-listener](./06-fiber-listener): Start a fiber server with the Border0 listner.
- [07-create-connector-and-socket](./07-create-connector-and-socket): Create a connector, a connector token and an SSH socket that's linked to the connector with upstream config.
- [08-delete-connector-and-socket](./08-delete-connector-and-socket): Delete the resources created from the previous example (a connector, a connector token and an SSH socket).
- [09-create-policy-and-attach](./09-create-policy-and-attach): Create a policy and a socket, and then attach the policy to the socket.
- [10-detach-policy-and-delete](./10-detach-policy-and-delete): Detach the policy and delete the policy and socket created from the previous example.
