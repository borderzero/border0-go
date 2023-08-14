package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/borderzero/border0-go/service/connector/types"
)

// SocketService is an interface for API client methods that interact with Border0 API to manage sockets.
type SocketService interface {
	Socket(ctx context.Context, idOrName string) (out *Socket, err error)
	Sockets(ctx context.Context) (out []Socket, err error)
	CreateSocket(ctx context.Context, in *Socket) (out *Socket, err error)
	UpdateSocket(ctx context.Context, idOrName string, in *Socket) (out *Socket, err error)
	DeleteSocket(ctx context.Context, idOrName string) (err error)
	SocketConnectors(ctx context.Context, idOrName string) (out *SocketConnectors, err error)
	SocketUpstreamConfigs(ctx context.Context, idOrName string) (out *SocketUpstreamConfigs, err error)
	SignSocketKey(ctx context.Context, idOrName string, in *SocketKeyToSign) (out *SignedSocketKey, err error)
}

// Socket fetches a socket by socket UUID or name. Socket UUID is globally unique and socket name is unique within
// a Border0 organization.
func (api *APIClient) Socket(ctx context.Context, idOrName string) (out *Socket, err error) {
	out = new(Socket)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/socket/%s", idOrName), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("socket [%s] not found: %w", idOrName, err)
		}
		return nil, err
	}
	return out, nil
}

// Sockets fetches all sockets in your Border0 organization.
func (api *APIClient) Sockets(ctx context.Context) (out []Socket, err error) {
	_, err = api.request(ctx, http.MethodGet, "/socket", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateSocket creates a new socket in your Border0 organization. Socket name must be unique within your organization,
// otherwise, an error will be returned. Socket type is required and must be one of the following: "http", "ssh",
// "tls" or "database". Socket name name must contain only lowercase letters, numbers and dashes.
func (api *APIClient) CreateSocket(ctx context.Context, in *Socket) (out *Socket, err error) {
	out = new(Socket)
	_, err = api.request(ctx, http.MethodPost, "/socket", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateSocket updates an existing socket in your Border0 organization.
func (api *APIClient) UpdateSocket(ctx context.Context, idOrName string, in *Socket) (out *Socket, err error) {
	out = new(Socket)
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/socket/%s", idOrName), in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteSocket deletes a socket in your Border0 organization. If the socket does not exist, no error will be returned.
func (api *APIClient) DeleteSocket(ctx context.Context, idOrName string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/socket/%s", idOrName), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// SocketConnectors fetches all connectors that are linked to a socket.
func (api *APIClient) SocketConnectors(ctx context.Context, idOrName string) (out *SocketConnectors, err error) {
	out = new(SocketConnectors)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/socket/%s/connectors", idOrName), nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SocketUpstreamConfigs fetches all upstream configurations for a socket.
func (api *APIClient) SocketUpstreamConfigs(ctx context.Context, idOrName string) (out *SocketUpstreamConfigs, err error) {
	out = new(SocketUpstreamConfigs)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/socket/%s/upstream_configurations", idOrName), nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SignSocketKey generates a signed SSH certificate for a socket. The SSH public key must be in OpenSSH format.
// The SSH certificate will be valid for 5 minutes. The host key is the public key Border0 server. It can be used
// to verify the SSH certificate.
func (api *APIClient) SignSocketKey(ctx context.Context, idOrName string, in *SocketKeyToSign) (out *SignedSocketKey, err error) {
	out = new(SignedSocketKey)
	_, err = api.request(ctx, http.MethodPost, fmt.Sprintf("/socket/%s/signkey", idOrName), in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Socket represents a socket in Border0 API. A socket can be linked to a connector with upstream configuration.
// Use `ConnectorID` to link a socket to a connector, and use `UpstreamConfig` to configure upstream for a socket.
type Socket struct {
	Name                           string            `json:"name"`
	SocketID                       string            `json:"socket_id"`
	SocketType                     string            `json:"socket_type"`
	Description                    string            `json:"description,omitempty"`
	UpstreamType                   string            `json:"upstream_type,omitempty"`
	UpstreamHTTPHostname           string            `json:"upstream_http_hostname,omitempty"`
	RecordingEnabled               bool              `json:"recording_enabled"`
	ConnectorAuthenticationEnabled bool              `json:"connector_authentication_enabled"`
	Tags                           map[string]string `json:"tags,omitempty"`

	// link to a connector with upstream config
	ConnectorID    string                                `json:"connector_id,omitempty"`
	UpstreamConfig *types.ConnectorServiceUpstreamConfig `json:"upstream_configuration,omitempty"`

	// associated policies
	Policies []Policy `json:"policies,omitempty"`
}

// SocketConnectors represents a list of connectors that are linked to a socket.
type SocketConnectors struct {
	List []SocketConnector `json:"list"`
}

// SocketConnector represents a connector that is linked to a socket.
type SocketConnector struct {
	ID            uint64 `json:"id"`
	ConnectorID   string `json:"connector_id"`
	ConnectorName string `json:"connector_name"`
	SocketID      string `json:"socket_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// SocketUpstreamConfigs represents a list of upstream configurations for a socket.
type SocketUpstreamConfigs struct {
	List []SocketUpstreamConfig `json:"list"`
}

// SocketUpstreamConfig represents an upstream configuration for a socket.
type SocketUpstreamConfig struct {
	Config    types.ConnectorServiceUpstreamConfig `json:"config"`
	CreatedAt time.Time                            `json:"created_at"`
	UpdatedAt time.Time                            `json:"updated_at"`
}

// SocketKeyToSign represents a SSH public key to sign.
type SocketKeyToSign struct {
	SSHPublicKey string `json:"ssh_public_key"`
}

// SignedSocketKey represents a signed SSH certificate and the host key.
type SignedSocketKey struct {
	SignedSSHCert string `json:"signed_ssh_cert"`
	HostKey       string `json:"host_key"`
}
