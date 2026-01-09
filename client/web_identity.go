package client

import (
	"context"
	"errors"
	"net/http"
)

// WebIdentityTokenExchangeInput is the input for the web identity exchange endpoint.
type WebIdentityTokenExchangeInput struct {
	OrganizationSubdomain       string `json:"organization_subdomain"`
	ServiceAccountName          string `json:"service_account_name"`
	WebIdentityToken            string `json:"web_identity_token"`
	Border0TokenDurationSeconds uint32 `json:"border0_token_duration_seconds,omitempty"`
}

// WebIdentityTokenExchangeOutput is the output for the web identity exchange endpoint.
type WebIdentityTokenExchangeOutput struct {
	Token string `json:"token"`
}

// WebIdentityService is an interface for API client methods that interact with
// Border0 API to exchange web identity tokens for Border0 service account tokens.
type WebIdentityService interface {
	ExchangeWebIdentityToken(ctx context.Context, input *WebIdentityTokenExchangeInput) (*WebIdentityTokenExchangeOutput, error)
}

// ExchangeWebIdentityToken retrieves a Border0 service account token given a web identity token for a trusted issuer and subject.
func (api *APIClient) ExchangeWebIdentityToken(ctx context.Context, input *WebIdentityTokenExchangeInput) (*WebIdentityTokenExchangeOutput, error) {
	output := new(WebIdentityTokenExchangeOutput)
	_, err := api.request(ctx, http.MethodPost, "/auth/web_identity/exchange", input, output)
	if err != nil {
		if apiErr, ok := BadRequest(err); ok {
			return nil, errors.New(apiErr.Message)
		}
		return nil, err
	}
	return output, nil
}
