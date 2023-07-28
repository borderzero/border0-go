package ccs

import (
	"context"

	"google.golang.org/grpc/credentials"
)

const (
	// ConnectorControlStreamMetadataKeyToken is the GRPC
	// stream metadata key for the authorization token.
	ConnectorControlStreamMetadataKeyToken = "token"

	// ConnectorControlStreamMetadataKeyConnectorId is
	// the GRPC stream metadata key for the connector id.
	ConnectorControlStreamMetadataKeyConnectorId = "connector_id"
)

// ConnectorControlStreamCredentials represents the authentication mechanism
// against the Border0 API's connector-control-plain (GRPC) server.
type ConnectorControlStreamCredentials struct {
	token             string
	connectorId       string
	insecureTransport bool
}

// ensures border0GrpcTunnelCredentials implements credentials.PerRPCCredentials
// (the generic authentication interface for GRPC) at compile-time.
var _ credentials.PerRPCCredentials = (*ConnectorControlStreamCredentials)(nil)

// CredentialOption is the constructor option type for ConnectorControlStreamCredentials.
type CredentialOption func(*ConnectorControlStreamCredentials)

// WithToken is the CredentialOption to set the token.
func WithToken(token string) CredentialOption {
	return func(c *ConnectorControlStreamCredentials) { c.token = token }
}

// WithConnectorId is the CredentialOption to set the connector id.
func WithConnectorId(connectorId string) CredentialOption {
	return func(c *ConnectorControlStreamCredentials) { c.connectorId = connectorId }
}

// WithInsecureTransport is the CredentialOption to toggle insecure transport.
func WithInsecureTransport(insecureTransport bool) CredentialOption {
	return func(c *ConnectorControlStreamCredentials) { c.insecureTransport = insecureTransport }
}

// NewConnectorControlStreamCredentials returns a new ConnectorControlStreamCredentials
// object initialized with the given options.
func NewConnectorControlStreamCredentials(opts ...CredentialOption) *ConnectorControlStreamCredentials {
	creds := &ConnectorControlStreamCredentials{
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
func (c *ConnectorControlStreamCredentials) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	md := map[string]string{}
	if c.token != "" {
		md[ConnectorControlStreamMetadataKeyToken] = c.token
	}
	if c.connectorId != "" {
		md[ConnectorControlStreamMetadataKeyConnectorId] = c.connectorId
	}
	return md, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
//
// ^ copied straight from the interface defintion.
func (c *ConnectorControlStreamCredentials) RequireTransportSecurity() bool {
	return !c.insecureTransport
}
