package client

import (
	"context"
	"fmt"
	"net/http"
)

// SocketService is an interface for API client methods that interact with Border0 API to manage sockets.
type SocketService interface {
	Socket(ctx context.Context, idOrName string) (out *Socket, err error)
	CreateSocket(ctx context.Context, in *Socket) (out *Socket, err error)
	SignSocketKey(ctx context.Context, idOrName string, in *SocketKeyToSign) (out *SignedSocketKey, err error)
}

// Socket fetches a socket by socket UUID or name. Socket UUID is globally unique and socket name is unique within an
// organization.
func (api *APIClient) Socket(ctx context.Context, idOrName string) (out *Socket, err error) {
	out = new(Socket)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/socket/%s", idOrName), nil, out)
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
	Name       string `json:"name"`
	SocketID   string `json:"socket_id"`
	SocketType string `json:"socket_type"`
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
