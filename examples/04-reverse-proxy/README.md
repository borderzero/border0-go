# Create an authenticated Reverse proxy using the border0-go SDK

This example demonstrates how to use border0-go to create a net.Listener, and use it as a reverse proxy to an upstream server. We'll use the `http` socket type for this example. We then start a reverse proxy server on the listerner, and forward all incoming requests to the upstream server, in this case BBC.com. We'll also rewrite the content. Finally, since all requests will be authenticated, we'll log to the console the name and email of the authenticated user.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) If needed create a new socket, or use the existing socket named sdk-socket-http of type http.
2) Create a new net.Listener for the socket using `border0.Listen()`
3) Create a new reverse proxy NewSingleHostReverseProxy() and start an http server on the border0 net.Listener.
4) every incoming request will be authenticated we'll print the name and email of the authenticated user to the console. This information is in the X-Auth-Name X-Auth-email headers.
5) Finnaly we'll reverse the requests to bbc.com and rewrite the content.


# Expected output
```
$ go run main.go
Welcome to Border0.com
sdk-socket-http - https://sdk-socket-http-border0-demo.border0.io

=======================================================
Logs
=======================================================
2023/08/11 14:26:22 serving request from Andree Toonk andree@border0.com
```

In this case the newly created socket is named sdk-socket-http and is available at https://sdk-socket-http-border0-demo.border0.io. You can now use this socket to authenticate your requests to bbc.com. 