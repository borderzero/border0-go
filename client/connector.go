package client

import (
	"context"
	"fmt"
	"net/http"
)

type ConnectorService interface {
	Connector(ctx context.Context, id string) (out *Connector, err error)
	Connectors(ctx context.Context) (out []Connector, err error)
	CreateConnector(ctx context.Context, in *Connector) (out *Connector, err error)
	UpdateConnector(ctx context.Context, in *Connector) (out *Connector, err error)
	DeleteConnector(ctx context.Context, id string) (err error)
	ConnectorTokens(ctx context.Context, connectorID string) (out *ConnectorTokens, err error)
	CreateConnectorToken(ctx context.Context, in *ConnectorToken) (out *ConnectorToken, err error)
	DeleteConnectorToken(ctx context.Context, connectorID, tokenID string) (err error)
}

func (api *APIClient) Connector(ctx context.Context, id string) (out *Connector, err error) {
	out = new(Connector)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/connector/%s", id), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("connector [%s] not found: %w", id, err)
		}
		return nil, err
	}
	return out, nil
}

func (api *APIClient) Connectors(ctx context.Context) (out []Connector, err error) {
	_, err = api.request(ctx, http.MethodGet, "/connectors", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (api *APIClient) CreateConnector(ctx context.Context, in *Connector) (out *Connector, err error) {
	out = new(Connector)
	_, err = api.request(ctx, http.MethodPost, "/connector", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (api *APIClient) UpdateConnector(ctx context.Context, in *Connector) (out *Connector, err error) {
	out = new(Connector)
	_, err = api.request(ctx, http.MethodPut, "/connector", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (api *APIClient) DeleteConnector(ctx context.Context, id string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/connector/%s", id), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func (api *APIClient) ConnectorTokens(ctx context.Context, connectorID string) (out *ConnectorTokens, err error) {
	out = new(ConnectorTokens)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/connector/%s/tokens", connectorID), nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (api *APIClient) CreateConnectorToken(ctx context.Context, in *ConnectorToken) (out *ConnectorToken, err error) {
	out = new(ConnectorToken)
	_, err = api.request(ctx, http.MethodPost, "/connector/token", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (api *APIClient) DeleteConnectorToken(ctx context.Context, connectorID, tokenID string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/connector/%s/token/%s", connectorID, tokenID), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

type Connector struct {
	// input and output fields
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// output field
	ConnectorID string `json:"connector_id"`
}

type ConnectorToken struct {
	// input and output fields
	ConnectorID string       `json:"connector_id"`
	Name        string       `json:"name"`
	ExpiresAt   FlexibleTime `json:"expires_at,omitempty"`

	// additional output fields
	ID        string       `json:"id"`
	Token     string       `json:"token"`
	CreatedBy string       `json:"created_by"`
	CreatedAt FlexibleTime `json:"created_at"`
}

type ConnectorTokens struct {
	List      []ConnectorToken `json:"list"`
	Connector Connector        `json:"connector"`
}
