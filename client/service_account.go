package client

import (
	"context"
	"fmt"
	"net/http"
)

// ServiceAccountService is an interface for API client methods that interact with Border0 API to manage service accounts.
type ServiceAccountService interface {
	ServiceAccount(ctx context.Context, name string) (out *ServiceAccount, err error)
	CreateServiceAccount(ctx context.Context, in *ServiceAccount) (out *ServiceAccount, err error)
	UpdateServiceAccount(ctx context.Context, in *ServiceAccount) (out *ServiceAccount, err error)
	DeleteServiceAccount(ctx context.Context, name string) (err error)
	CreateServiceAccountToken(ctx context.Context, serviceAccountName string, in *ServiceAccountToken) (out *ServiceAccountToken, err error)
	DeleteServiceAccountToken(ctx context.Context, serviceAccountName, tokenID string) (err error)
}

// ServiceAccount fetches a service account from your Border0 organization
// by name. Service Account name must be unique and immutable.
func (api *APIClient) ServiceAccount(ctx context.Context, name string) (out *ServiceAccount, err error) {
	out = new(ServiceAccount)
	_, err = api.request(ctx, http.MethodGet, fmt.Sprintf("/organizations/iam/service_accounts/%s", name), nil, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("service account with name [%s] not found: %w", name, err)
		}
		return nil, err
	}
	return out, nil
}

// CreateServiceAccount creates a new service account in your Border0 organization. Service
// Account name must be in slug format (alphanumeric and dashes) and be unique in the organization.
func (api *APIClient) CreateServiceAccount(ctx context.Context, in *ServiceAccount) (out *ServiceAccount, err error) {
	out = new(ServiceAccount)
	_, err = api.request(ctx, http.MethodPost, "/organizations/iam/service_accounts", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateServiceAccount updates an existing service account in your Border0 organization.
func (api *APIClient) UpdateServiceAccount(ctx context.Context, in *ServiceAccount) (out *ServiceAccount, err error) {
	out = new(ServiceAccount)
	_, err = api.request(ctx, http.MethodPut, fmt.Sprintf("/organizations/iam/service_accounts/%s", in.Name), in, out)
	if err != nil {
		if NotFound(err) {
			return nil, fmt.Errorf("service account with name [%s] not found: %w", in.Name, err)
		}
		return nil, err
	}
	return out, nil
}

// DeleteServiceAccount deletes an existing service account from your Border0 organization.
func (api *APIClient) DeleteServiceAccount(ctx context.Context, name string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/organizations/iam/service_accounts/%s", name), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// CreateServiceAccountToken creates a new token for a service account. The token is used to authenticate connector with the
// Border0 API. The token can be created with or without a expiration date. If ExpiresAt field is not set, token will not expire.
func (api *APIClient) CreateServiceAccountToken(ctx context.Context, serviceAccountName string, in *ServiceAccountToken) (out *ServiceAccountToken, err error) {
	out = new(ServiceAccountToken)
	_, err = api.request(ctx, http.MethodPost, fmt.Sprintf("/organizations/iam/service_accounts/%s/tokens", serviceAccountName), in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteServiceAccountToken deletes a token for a service account by service account name and token UUID.
func (api *APIClient) DeleteServiceAccountToken(ctx context.Context, serviceAccountName, tokenID string) (err error) {
	_, err = api.request(ctx, http.MethodDelete, fmt.Sprintf("/organizations/iam/service_accounts/%s/tokens/%s", serviceAccountName, tokenID), nil, nil)
	if err != nil {
		if NotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// ServiceAccount represents a service account in your Border0 organization.
type ServiceAccount struct {
	// input fields
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Active      bool   `json:"active"`

	// output fields
	ID         string       `json:"service_account_id"`
	CreatedAt  FlexibleTime `json:"created_at"`
	UpdatedAt  FlexibleTime `json:"updated_at"`
	LastSeenAt FlexibleTime `json:"last_seen_at,omitempty"`
}

// ServiceAccountToken represents a service account token in your Border0 organization.
type ServiceAccountToken struct {
	// input fields
	Name      string       `json:"name"`
	ExpiresAt FlexibleTime `json:"expires_at,omitempty"`

	// output fields
	ID        string       `json:"id"`
	Token     string       `json:"token"`
	CreatedAt FlexibleTime `json:"created_at"`
}
