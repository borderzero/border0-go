package connector

import (
	"context"

	"google.golang.org/grpc/credentials"
)

const (
	// ControlStreamMetadataKeyToken is the GRPC
	// stream metadata key for the authorization token.
	ControlStreamMetadataKeyToken = "token"

	// ControlStreamMetadataKeyConnectorId is
	// the GRPC stream metadata key for the connector id.
	ControlStreamMetadataKeyConnectorId = "connector_id"
)

// ControlStreamCredentials represents the authentication mechanism
// against the Border0 API's connector-control-plain (GRPC) server.
type ControlStreamCredentials struct {
	token             string
	connectorId       string
	insecureTransport bool
}

// ensures ControlStreamCredentials implements credentials.PerRPCCredentials
// (the generic authentication interface for GRPC) at compile-time.
var _ credentials.PerRPCCredentials = (*ControlStreamCredentials)(nil)

// CredentialOption is the constructor option type for ControlStreamCredentials.
type CredentialOption func(*ControlStreamCredentials)

// WithToken is the CredentialOption to set the token.
func WithToken(token string) CredentialOption {
	return func(c *ControlStreamCredentials) { c.token = token }
}

// WithConnectorId is the CredentialOption to set the connector id.
func WithConnectorId(connectorId string) CredentialOption {
	return func(c *ControlStreamCredentials) { c.connectorId = connectorId }
}

// WithInsecureTransport is the CredentialOption to toggle insecure transport.
func WithInsecureTransport(insecureTransport bool) CredentialOption {
	return func(c *ControlStreamCredentials) { c.insecureTransport = insecureTransport }
}

// NewControlStreamCredentials returns a new ControlStreamCredentials
// object initialized with the given options.
func NewControlStreamCredentials(opts ...CredentialOption) *ControlStreamCredentials {
	creds := &ControlStreamCredentials{
		insecureTransport: false,
	}
	for _, opt := range opts {
		opt(creds)
	}
	return creds
}

// GetRequestMetadata gets the current request metadata, refreshing tokens
// if required. This should be called by the transport layer on each
// request, and the data should be populated in headers or other
// context. If a status code is returned, it will be used as the status for
// the RPC (restricted to an allowable set of codes as defined by gRFC
// A54). uri is the URI of the entry point for the request.  When supported
// by the underlying implementation, ctx can be used for timeout and
// cancellation. Additionally, RequestInfo data will be available via ctx
// to this call.
//
// ^ copied straight from the interface defintion.
func (c *ControlStreamCredentials) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	md := map[string]string{}
	if c.token != "" {
		md[ControlStreamMetadataKeyToken] = c.token
	}
	if c.connectorId != "" {
		md[ControlStreamMetadataKeyConnectorId] = c.connectorId
	}
	return md, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
//
// ^ copied straight from the interface defintion.
func (c *ControlStreamCredentials) RequireTransportSecurity() bool {
	return !c.insecureTransport
}
