# Create an authenticated http load balancer using the border0-go SDK

This example demonstrates how to use border0-go to create a net.Listener, and use it as a HTTP load balancer. 
In this example, we'll load balance incomming requests between example.com, example.net and example.org.
We then rewrite the content, and let the user know which origin server served the request. We'll also insert the user's email and full name. This serves as an example for how to access the user's identity.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) If needed create a new socket, or use the existing socket named sdk-socket-http of type http.
2) Create a new net.Listener for the socket using `border0.Listen()`
3) we create a new handler called loadBalancerHandler() and start an http server on the border0 net.Listener.
4) every incoming request is authenticated by the border0 edge servers (done for you, not part of the code). The user's identity is passed in the the X-Auth-Name and X-Auth-Email headers.  
5) The code does round robin load balancing between example.com, example.net and example.org.
5) Finnaly we'll rewrite the content using ModifyResponse() and insert the user's name and email as well as the origin server that served the request.


# Expected output
```
$ go run main.go
2023/08/11 17:30:44 Starting load balancer on Border0
Welcome to Border0.com
sdk-socket-http - https://sdk-socket-http-border0-demo.border0.io

=======================================================
Logs
=======================================================
```

In this case the newly created socket is named sdk-socket-http and is available at https://sdk-socket-http-border0-demo.border0.io. You can now use this socket to authenticate your requests and load balance between your various origin servers.