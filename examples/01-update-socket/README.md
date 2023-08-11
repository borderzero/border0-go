# Create a socket using the border0-go SDK 

This example demonstrates how to use the border0-go SDK to interact with the Border0 API. This particualr example covers fetching an existing socket by name and updating its description.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) Fetch the socket named sdk-socket-http and print its details to the console.
2) we then update the description of the socket and set it to "updated description"
3) We then post the updated socket object back to the api using `UpdateSocket()` and print the updated socket details to the console.

# Expected output
```
$ go run main.go
2023/08/11 13:49:45 ✅ socket from previous example [00-create-socket] = {
  "name": "sdk-socket-http",
  "socket_id": "1a1c11bd-b06f-4924-81b7-cafde4fd6595",
  "socket_type": "http",
  "upstream_type": "http",
  "recording_enabled": false,
  "connector_authentication_enabled": false
}
2023/08/11 13:49:45 ✅ socket from previous example [00-create-socket] = {
  "name": "sdk-socket-http",
  "socket_id": "1a1c11bd-b06f-4924-81b7-cafde4fd6595",
  "socket_type": "http",
  "description": "updated description",
  "upstream_type": "http",
  "recording_enabled": false,
  "connector_authentication_enabled": false
}
```
