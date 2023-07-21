package listen

import (
	"github.com/borderzero/border0-go/client"
)

// Option is a function that configures a Listener.
type Option func(*Listener)

// WithAPIClient sets the API client to use for the Listener.
func WithAPIClient(api client.Requester) Option {
	return func(l *Listener) {
		l.apiClient = api
	}
}

// WithAuthToken sets the auth token to use with the API client.
func WithAuthToken(token string) Option {
	return func(l *Listener) {
		l.authToken = token
	}
}

// WithSocketName sets the socket name to use for the HTTP socket that the Listener will create.
func WithSocketName(name string) Option {
	return func(l *Listener) {
		l.socketName = name
	}
}

// WithTunnelServer sets the tunnel server address for the Listener.
func WithTunnelServer(server string) Option {
	return func(l *Listener) {
		l.tunnelServer = server
	}
}
