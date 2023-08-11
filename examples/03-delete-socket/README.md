# Delete a socket using the border0-go SDK 

This example demonstrates how to use the border0-go SDK to interact with the Border0 API. This particualr example covers fetching an existing socket by name and deleting it.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) fetch the socket named sdk-socket-http using `	fetched, err := api.Socket(ctx, socketFromPreviousExample)
` and print its details to the console. 
2) We then delete the socket using `DeleteSocket()` and print the result to the console.

# Expected output
```
$ go run main.go
2023/08/11 13:50:18 ✅ socket from previous example [00-create-socket] = {
  "name": "sdk-socket-http",
  "socket_id": "1a1c11bd-b06f-4924-81b7-cafde4fd6595",
  "socket_type": "http",
  "description": "updated description",
  "upstream_type": "http",
  "recording_enabled": false,
  "connector_authentication_enabled": false
}
2023/08/11 13:50:18 ❎ socket [sdk-socket-http] deleted
```
