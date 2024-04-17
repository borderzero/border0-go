# List Connectors Using the border0-go SDK 

This example demonstrates how to use the border0-go SDK to interact with the Border0 API. This particular, this example covers fetching all connectors from the API and printing the details of each connector to the console.

Make sure you've exported your token like `export BORDER0_AUTH_TOKEN=your_auth_token_here` before running the example.

# Running the example
```
go run main.go
```

This will:

1) Fetch all the connectors in your organization using the `Connectors()` method. 
2) we then loop over each connector and print its details to the console.

# Expected output
```
$ go run main.go
✅ found 2 connectors in this Border0 organization
✅ connector details for border-clab-3430
{
  "name": "border-clab-3430",
  "description": "ContainerLab Connector",
  "connector_id": "77ce3654-0bf4-4a94-9dc0-3de38ff9a112"
}
✅ connector details for border-clab-6821
{
  "name": "border-clab-6821",
  "description": "ContainerLab Connector",
  "connector_id": "d56c2df9-f86c-4abd-a584-de3e0a213ee9"
}
```
