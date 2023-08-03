package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/borderzero/border0-go/service/connector/types"
)

// SocketService is an interface for API client methods that interact with Border0 API to manage sockets.
type SocketService interface {
	Socket(ctx context.Context, idOrName string) (out *Socket, err error)
	Sockets(ctx context.Context) (out []Socket, err error)
	CreateSocket(ctx context.Context, in *Socket) (out *Socket, err error)
	UpdateSocket(ctx context.Context, idOrName string, in *Socket) (out *Socket, err error)
	DeleteSocket(ctx context.Context, idOrName string) (err error)
	SignSocketKey(ctx context.Context, idOrName string, in *SocketKeyToSign) (out *SignedSocketKey, err error)
}

// Socket fetches a socket by socket UUID or name. Socket UUID is globally unique and socket name is unique within an
// organization.
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

// Sockets fetches all sockets in your organization.
func (api *APIClient) Sockets(ctx context.Context) (out []Socket, err error) {
	_, err = api.request(ctx, http.MethodGet, "/socket", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateSocket creates a new socket in your organization. Socket name must be unique within your organization,
// otherwise, an error will be returned. Socket type is required and must be one of the following: "http", "ssh",
// "tls" or "database".
func (api *APIClient) CreateSocket(ctx context.Context, in *Socket) (out *Socket, err error) {
	out = new(Socket)
	_, err = api.request(ctx, http.MethodPost, "/socket", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateSocket updates a socket in your organization.
func (api *APIClient) UpdateSocket(ctx context.Context, idOrName string, in *Socket) (out *Socket, err error) {
	out = new(Socket)
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/socket/%s", idOrName), in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteSocket deletes a socket in your organization. If the socket does not exist, no error will be returned.
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

// Socket represents a socket in Border0 API.
type Socket struct {
	Name                 string            `json:"name"`
	SocketID             string            `json:"socket_id"`
	SocketType           string            `json:"socket_type"`
	Description          string            `json:"description,omitempty"`
	UpstreamType         string            `json:"upstream_type,omitempty"`
	UpstreamHTTPHostname string            `json:"upstream_http_hostname,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty"`

	RecordingEnabled               bool `json:"recording_enabled"`
	ConnectorAuthenticationEnabled bool `json:"connector_authentication_enabled"`

	ConnectorData *SocketConnectorData `json:"connector_data,omitempty"`
}

type SocketConnectorData struct {
	ConnectorID string                                `json:"connector_id,omitempty"`
	Config      *types.ConnectorServiceUpstreamConfig `json:"config,omitempty"`
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
