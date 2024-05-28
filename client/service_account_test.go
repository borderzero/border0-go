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
	"github.com/stretchr/testify/require"
)

func Test_APIClient_ServiceAccount(t *testing.T) {
	t.Parallel()

	testServiceAccount := &ServiceAccount{
		Name:        "test-name",
		Description: "Test description",
		Role:        "admin",
		ID:          "test-id",
	}

	tests := []struct {
		name               string
		mockRequester      func(context.Context, *mocks.ClientHTTPRequester)
		givenName          string
		wantServiceAccount *ServiceAccount
		wantErr            error
	}{
		{
			name: "failed to get service account",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, "test-name"), nil, new(ServiceAccount)).
					Return(http.StatusBadRequest, errors.New("failed to get service account"))
			},
			givenName:          "test-name",
			wantServiceAccount: nil,
			wantErr:            errors.New("failed after 1 attempt: failed to get service account"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, "test-name"), nil, new(ServiceAccount)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "service account not found"})
			},
			givenName: "test-name",
			wantErr:   errors.New("service account with name [test-name] not found: failed after 4 attempts: 404: service account not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, "test-name"), nil, new(ServiceAccount)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*ServiceAccount)
						*output = *testServiceAccount
					})
			},
			givenName:          "test-name",
			wantServiceAccount: testServiceAccount,
			wantErr:            nil,
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

			gotServiceAccount, gotErr := api.ServiceAccount(ctx, test.givenName)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantServiceAccount, gotServiceAccount)
		})
	}
}

func Test_APIClient_CreateServiceAccount(t *testing.T) {
	t.Parallel()

	timestamp, err := FlexibleTimeFrom("2020-01-01T00:00:00Z")
	require.NoError(t, err)

	testServiceAccountInput := &ServiceAccount{
		Name:        "test-name",
		Description: "Test description",
		Role:        "admin",
		Active:      true,
	}
	testServiceAccountOutput := &ServiceAccount{
		Name:        "test-name",
		Description: "Test description",
		Role:        "admin",
		Active:      true,
		ID:          "test-id",
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
		LastSeenAt:  timestamp,
	}

	tests := []struct {
		name                string
		mockRequester       func(context.Context, *mocks.ClientHTTPRequester)
		givenServiceAccount *ServiceAccount
		wantServiceAccount  *ServiceAccount
		wantErr             error
	}{
		{
			name: "failed to create serviceAccount",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/service_accounts", defaultBaseURL), testServiceAccountInput, new(ServiceAccount)).
					Return(http.StatusBadRequest, errors.New("failed to create service account"))
			},
			givenServiceAccount: testServiceAccountInput,
			wantServiceAccount:  nil,
			wantErr:             errors.New("failed after 1 attempt: failed to create service account"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/service_accounts", defaultBaseURL), testServiceAccountInput, new(ServiceAccount)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						serviceAccount := output.(*ServiceAccount)
						*serviceAccount = *testServiceAccountOutput
					})
			},
			givenServiceAccount: testServiceAccountInput,
			wantServiceAccount:  testServiceAccountOutput,
			wantErr:             nil,
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

			gotServiceAccount, gotErr := api.CreateServiceAccount(ctx, test.givenServiceAccount)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantServiceAccount, gotServiceAccount)
		})
	}
}

func Test_APIClient_UpdateServiceAccount(t *testing.T) {
	t.Parallel()

	testServiceAccount := &ServiceAccount{
		Name:        "test-name",
		Description: "Test description",
		Role:        "admin",
		ID:          "test-id",
	}

	tests := []struct {
		name                string
		mockRequester       func(context.Context, *mocks.ClientHTTPRequester)
		givenServiceAccount *ServiceAccount
		wantServiceAccount  *ServiceAccount
		wantErr             error
	}{
		{
			name: "failed to update service account",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, testServiceAccount.Name), testServiceAccount, new(ServiceAccount)).
					Return(http.StatusBadRequest, errors.New("failed to update service account"))
			},
			givenServiceAccount: testServiceAccount,
			wantServiceAccount:  nil,
			wantErr:             errors.New("failed after 1 attempt: failed to update service account"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, testServiceAccount.Name), testServiceAccount, new(ServiceAccount)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						serviceAccount := output.(*ServiceAccount)
						*serviceAccount = *testServiceAccount
					})
			},
			givenServiceAccount: testServiceAccount,
			wantServiceAccount:  testServiceAccount,
			wantErr:             nil,
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

			gotServiceAccount, gotErr := api.UpdateServiceAccount(ctx, test.givenServiceAccount)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantServiceAccount, gotServiceAccount)
		})
	}
}

func Test_APIClient_DeleteServiceAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenName     string
		wantErr       error
	}{
		{
			name: "failed to delete service account",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, "test-name"), nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete service account"))
			},
			givenName: "test-name",
			wantErr:   errors.New("failed after 1 attempt: failed to delete service account"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, "test-name"), nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "service account not found"})
			},
			givenName: "test-name",
			wantErr:   nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/service_accounts/%s", defaultBaseURL, "test-name"), nil, nil).
					Return(http.StatusOK, nil)
			},
			givenName: "test-name",
			wantErr:   nil,
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

			gotErr := api.DeleteServiceAccount(ctx, test.givenName)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_CreateServiceAccountToken(t *testing.T) {
	t.Parallel()

	expiresAt, err := FlexibleTimeFrom("2030-01-01T00:00:00Z")
	require.NoError(t, err)

	createdAt, err := FlexibleTimeFrom("2020-01-01T00:00:00Z")
	require.NoError(t, err)

	testServiceAccountName := "test-name"

	testServiceAccountTokenInput := &ServiceAccountToken{
		Name:      "test-name",
		ExpiresAt: expiresAt,
	}
	testServiceAccountTokenOutput := &ServiceAccountToken{
		Token:     "test-token",
		ID:        "test-id",
		Name:      "test-name",
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	tests := []struct {
		name                     string
		mockRequester            func(context.Context, *mocks.ClientHTTPRequester)
		givenServiceAccountToken *ServiceAccountToken
		wantServiceAccountToken  *ServiceAccountToken
		wantErr                  error
	}{
		{
			name: "failed to create service account token",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/service_accounts/%s/tokens", defaultBaseURL, testServiceAccountName), testServiceAccountTokenInput, new(ServiceAccountToken)).
					Return(http.StatusBadRequest, errors.New("failed to create service account token"))
			},
			givenServiceAccountToken: testServiceAccountTokenInput,
			wantServiceAccountToken:  nil,
			wantErr:                  errors.New("failed after 1 attempt: failed to create service account token"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/service_accounts/%s/tokens", defaultBaseURL, testServiceAccountName), testServiceAccountTokenInput, new(ServiceAccountToken)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						serviceAccountToken := output.(*ServiceAccountToken)
						*serviceAccountToken = *testServiceAccountTokenOutput
					})
			},
			givenServiceAccountToken: testServiceAccountTokenInput,
			wantServiceAccountToken:  testServiceAccountTokenOutput,
			wantErr:                  nil,
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

			gotServiceAccountToken, gotErr := api.CreateServiceAccountToken(ctx, testServiceAccountName, test.givenServiceAccountToken)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantServiceAccountToken, gotServiceAccountToken)
		})
	}
}

func Test_APIClient_DeleteServiceAccountToken(t *testing.T) {
	t.Parallel()

	testServiceAccountName := "test-service-account-name"

	tests := []struct {
		name                    string
		mockRequester           func(context.Context, *mocks.ClientHTTPRequester)
		givenServiceAccountName string
		givenTokenID            string
		wantErr                 error
	}{
		{
			name: "failed to delete connector token",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/service_accounts/%s/tokens/%s", defaultBaseURL, testServiceAccountName, "test-token-id"), nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete service account token"))
			},
			givenServiceAccountName: testServiceAccountName,
			givenTokenID:            "test-token-id",
			wantErr:                 errors.New("failed after 1 attempt: failed to delete service account token"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/service_accounts/%s/tokens/%s", defaultBaseURL, testServiceAccountName, "test-token-id"), nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "service account token not found"})
			},
			givenServiceAccountName: testServiceAccountName,
			givenTokenID:            "test-token-id",
			wantErr:                 nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/service_accounts/%s/tokens/%s", defaultBaseURL, testServiceAccountName, "test-token-id"), nil, nil).
					Return(http.StatusOK, nil)
			},
			givenServiceAccountName: testServiceAccountName,
			givenTokenID:            "test-token-id",
			wantErr:                 nil,
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

			gotErr := api.DeleteServiceAccountToken(ctx, test.givenServiceAccountName, test.givenTokenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
