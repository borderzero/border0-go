package border0

import (
	"net"

	"github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/listen"
)

// NewAPIClient creates a new Border0 API client with the given options.
// If no options are provided, some default values will be used. See
// the client.Option type for more details.
//
// Also see examples folder for a basic example of how to manage Border0
// resources using this API client.
func NewAPIClient(options ...client.Option) client.Requester {
	return client.New(options...)
}

// Listen creates a new Border0 listener with the given options. It returns
// a net.Listener that can be used to accept incoming connections. When the
// listener is passed to http.Serve, the server will accept HTTP requests
// sent by Border0 and forward them to the handler. The handler's response
// will be sent back to Border0. If no options are provided, some default
// values will be used. See the listen.Option type for more details.
//
// Also see examples folder for a simple and advanced examples of how to
// use this Listen function.
func Listen(options ...listen.Option) (net.Listener, error) {
	l := listen.New(options...)
	return l, l.Start()
}
