// Package border0 provides helper functions to:
//
//  1. create a Border API client to manage Border0 resources
//  2. create a Border listener to handle incoming HTTP connections
//
// # About Border0
//
// Border0 enables users to log into various services, including web, SSH, database, and generic TCP, using their existing Single Sign-On (SSO) credentials.
// If you haven't yet registered, [create a new account] and explore our [informative blog posts] and [comprehensive documentation].
//
// [create a new account]: https://portal.border0.com/register
// [informative blog posts]: https://www.border0.com/blog
// [comprehensive documentation]: https://docs.border0.com/docs/quick-start
package border0

import (
	"context"
	"net"

	"github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/listen"
)

// NewAPIClient creates a new Border0 API client with the given options.
// If no options are provided, some default values will be used. See
// the client.Option type for more details.
//
// Explore examples folder for additional examples of how to manage Border0
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
// Explore examples folder for additional examples of how to use this Listen
// function with other HTTP libraries and frameworks.
func Listen(options ...listen.Option) (net.Listener, error) {
	l := listen.New(options...)
	return l, l.Start(context.Background())
}
