package client

import (
	"context"
	"net/http"
)

type ServerInfoService interface {
	ServerInfo(context.Context) (*ServerInfo, error)
}

func (api *APIClient) ServerInfo(ctx context.Context) (*ServerInfo, error) {
	info := new(ServerInfo)
	_, err := api.request(ctx, http.MethodGet, "/serverinfo", nil, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

type ServerInfo struct {
	Primary *bool `json:"primary,omitempty"`
}
