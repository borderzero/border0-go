# Create a socket using the border0-go SDK 

This example demonstrates how to use the border0-go SDK to interact with the Border0 API. It covers creating a new socket and fetching its details.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) Create a new socket named sdk-socket-http of type http and print the details to the console.
2) Fetch the details of the socket you just created and print the details to the console.

# Expected output
```
$ go run main.go
2023/08/11 13:48:48 ✅ created socket = {
  "name": "sdk-socket-http",
  "socket_id": "1a1c11bd-b06f-4924-81b7-cafde4fd6595",
  "socket_type": "http",
  "upstream_type": "http",
  "recording_enabled": false,
  "connector_authentication_enabled": false
}
2023/08/11 13:48:48 ✅ fetched socket = {
  "name": "sdk-socket-http",
  "socket_id": "1a1c11bd-b06f-4924-81b7-cafde4fd6595",
  "socket_type": "http",
  "upstream_type": "http",
  "recording_enabled": false,
  "connector_authentication_enabled": false
}
```
