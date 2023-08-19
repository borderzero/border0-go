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

// WithPolicies sets the policy names to use for the Listener's underlaying HTTP socket. Policies with the
// given names will be attached to the socket. If the policy names list is empty, then no policies will be
// attached. If there are any changes made to the policy names list between listener startups, they will get
// properly handled. The listener will check the socket's already attached policies, and compare them with
// the given policy names. Removed policies will be detached from the socket, and new policies will be checked
// and made sure they exist, and then they will be attached to the listener's socket.
func WithPolicies(names []string) Option {
	return func(l *Listener) {
		l.policyNames = names
	}
}

// WithTunnelServer sets the tunnel server address for the Listener.
func WithTunnelServer(server string) Option {
	return func(l *Listener) {
		l.tunnelServer = server
	}
}
