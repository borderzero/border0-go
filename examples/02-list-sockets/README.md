# Modify a socket using the border0-go SDK 

This example demonstrates how to use the border0-go SDK to interact with the Border0 API. This particualr example covers fetching all sockets from the API and printing the details of each socket to the console.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) Fetch all the sockets in your organization using the `Sockets()` method. 
2) we then loop over each socket and print its details to the console.


# Expected output
```
$ go run main.go
✅ found 1 sockets in this Border0 organization

✅ socket details for sdk-socket-http
{
  "name": "sdk-socket-http",
  "socket_id": "87e78f05-e9f7-477c-b264-1745f6288ad5",
  "socket_type": "http",
  "upstream_type": "http",
  "recording_enabled": false,
  "connector_authentication_enabled": false
}
```
