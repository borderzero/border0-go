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
	DataConsistency *DataConsistency `json:"data_consistency,omitempty"`
}

type DataConsistency struct {
	RxAfterTxDelayMS int64 `json:"rx_after_tx_delay_ms,omitempty"`
}
