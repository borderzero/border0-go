package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/borderzero/border0-go/client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_APIClient_ExchangeWebIdentityToken(t *testing.T) {
	t.Parallel()

	testInput := &WebIdentityTokenExchangeInput{
		OrganizationSubdomain:       "test-org",
		ServiceAccountName:          "test-service-account",
		WebIdentityToken:            "test-web-identity-token",
		Border0TokenDurationSeconds: 3600,
	}

	testOutput := &WebIdentityTokenExchangeOutput{
		Token: "test-border0-token",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenInput    *WebIdentityTokenExchangeInput
		wantOutput    *WebIdentityTokenExchangeOutput
		wantErr       error
	}{
		{
			name: "failed to exchange web identity token",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/auth/web_identity/exchange", defaultBaseURL), testInput, new(WebIdentityTokenExchangeOutput)).
					Return(http.StatusInternalServerError, errors.New("internal server error"))
			},
			givenInput: testInput,
			wantOutput: nil,
			wantErr:    errors.New("failed after 1 attempt: internal server error"),
		},
		{
			name: "bad request error with message extraction",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/auth/web_identity/exchange", defaultBaseURL), testInput, new(WebIdentityTokenExchangeOutput)).
					Return(http.StatusBadRequest, Error{Code: http.StatusBadRequest, Message: "invalid web identity token"})
			},
			givenInput: testInput,
			wantOutput: nil,
			wantErr:    errors.New("invalid web identity token"),
		},
		{
			name: "bad request error with different message",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/auth/web_identity/exchange", defaultBaseURL), testInput, new(WebIdentityTokenExchangeOutput)).
					Return(http.StatusBadRequest, Error{Code: http.StatusBadRequest, Message: "service account not found"})
			},
			givenInput: testInput,
			wantOutput: nil,
			wantErr:    errors.New("service account not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodPost, fmt.Sprintf("%s/auth/web_identity/exchange", defaultBaseURL), testInput, new(WebIdentityTokenExchangeOutput)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*WebIdentityTokenExchangeOutput)
						*output = *testOutput
					})
			},
			givenInput: testInput,
			wantOutput: testOutput,
			wantErr:    nil,
		},
		{
			name: "happy path without optional duration",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				inputWithoutDuration := &WebIdentityTokenExchangeInput{
					OrganizationSubdomain: "test-org",
					ServiceAccountName:    "test-service-account",
					WebIdentityToken:      "test-web-identity-token",
				}
				requester.On("Request", ctx, http.MethodPost, fmt.Sprintf("%s/auth/web_identity/exchange", defaultBaseURL), inputWithoutDuration, new(WebIdentityTokenExchangeOutput)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*WebIdentityTokenExchangeOutput)
						*output = *testOutput
					})
			},
			givenInput: &WebIdentityTokenExchangeInput{
				OrganizationSubdomain: "test-org",
				ServiceAccountName:    "test-service-account",
				WebIdentityToken:      "test-web-identity-token",
			},
			wantOutput: testOutput,
			wantErr:    nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.ClientHTTPRequester)
			test.mockRequester(ctx, requester)

			api := New(
				WithRetryMax(0),
			)
			api.http = requester

			gotOutput, gotErr := api.ExchangeWebIdentityToken(ctx, test.givenInput)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantOutput, gotOutput)
		})
	}
}
