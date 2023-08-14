package client

import (
	"context"
	"fmt"
	"net/http"
)

// ConnectorService is an interface for API client methods that interact with Border0 API to manage connectors and connector tokens.
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

// Connector fetches a connector from your Border0 organization by UUID. Connector UUID is globally unique and immutable.
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

// Connectors fetches all connectors in your Border0 organization.
func (api *APIClient) Connectors(ctx context.Context) (out []Connector, err error) {
	_, err = api.request(ctx, http.MethodGet, "/connectors", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateConnector creates a new connector in your Border0 organization. Connector name must be unique within your organization,
// otherwise API will return an error. Connector name must contain only lowercase letters, numbers and dashes.
func (api *APIClient) CreateConnector(ctx context.Context, in *Connector) (out *Connector, err error) {
	out = new(Connector)
	_, err = api.request(ctx, http.MethodPost, "/connector", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateConnector updates an existing connector in your Border0 organization.
func (api *APIClient) UpdateConnector(ctx context.Context, in *Connector) (out *Connector, err error) {
	out = new(Connector)
	_, err = api.request(ctx, http.MethodPut, "/connector", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteConnector deletes a connector from your Border0 organization by connector's UUID.
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

// ConnectorTokens fetches all tokens for a connector by connector's UUID.
func (api *APIClient) ConnectorTokens(ctx context.Context, connectorID string) (out *ConnectorTokens, err error) {
	out = new(ConnectorTokens)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/connector/%s/tokens", connectorID), nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateConnectorToken creates a new token for a connector. Token is used to authenticate connector with Border0 API.
// Token can be created with or without a expiration date. If ExpiresAt field is not set, token will not expire.
func (api *APIClient) CreateConnectorToken(ctx context.Context, in *ConnectorToken) (out *ConnectorToken, err error) {
	out = new(ConnectorToken)
	_, err = api.request(ctx, http.MethodPost, "/connector/token", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteConnectorToken deletes a token for a connector by connector's UUID and token's UUID.
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

// Connector represents a connector in your Border0 organization.
type Connector struct {
	// input and output fields
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// output field
	ConnectorID string `json:"connector_id"`
}

// ConnectorToken represents a token for a connector.
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

// ConnectorTokens represents a list of tokens for a connector.
type ConnectorTokens struct {
	List      []ConnectorToken `json:"list"`
	Connector Connector        `json:"connector"`
}
