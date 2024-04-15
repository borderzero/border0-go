package client

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/borderzero/border0-go/client/mocks"
	"github.com/borderzero/border0-go/client/reqedit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_APIClient_Connector(t *testing.T) {
	t.Parallel()

	testConnector := &Connector{
		Name:        "test-name",
		Description: "Test description",
		ConnectorID: "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantConnector *Connector
		wantErr       error
	}{
		{
			name: "failed to get connector",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/connector/test-id", nil, new(Connector)).
					Return(http.StatusBadRequest, errors.New("failed to get connector"))
			},
			givenID:       "test-id",
			wantConnector: nil,
			wantErr:       errors.New("failed after 1 attempt: failed to get connector"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/connector/test-id", nil, new(Connector)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "connector not found"})
			},
			givenID: "test-id",
			wantErr: errors.New("connector [test-id] not found: failed after 4 attempts: 404: connector not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/connector/test-id", nil, new(Connector)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*Connector)
						*output = *testConnector
					})
			},
			givenID:       "test-id",
			wantConnector: testConnector,
			wantErr:       nil,
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

			gotConnector, gotErr := api.Connector(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnector, gotConnector)
		})
	}
}

func Test_APIClient_Connectors(t *testing.T) {
	t.Parallel()

	testConnectors := []Connector{
		{Name: "test-name-1", ConnectorID: "test-id-1", Description: "Test description 1"},
		{Name: "test-name-2", ConnectorID: "test-id-2", Description: "Test description 2"},
		{Name: "test-name-3", ConnectorID: "test-id-3", Description: "Test description 3"},
		{Name: "test-name-4", ConnectorID: "test-id-4", Description: "Test description 4"},
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		wantConnectors []Connector
		wantErr        error
	}{
		{
			name: "failed to get connectors",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/connectors", nil, new([]Connector)).
					Return(http.StatusInternalServerError, errors.New("failed to get connectors"))
			},
			wantConnectors: nil,
			wantErr:        errors.New("failed after 1 attempt: failed to get connectors"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/connectors", nil, new([]Connector)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*[]Connector)
						*output = testConnectors
					})
			},
			wantConnectors: testConnectors,
			wantErr:        nil,
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

			gotConnectors, gotErr := api.Connectors(ctx)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnectors, gotConnectors)
		})
	}
}

func Test_APIClient_CreateConnector(t *testing.T) {
	t.Parallel()

	testConnectorInput := &Connector{
		Name:        "test-name",
		Description: "Test description",
	}
	testConnectorOutput := &Connector{
		Name:        "test-name",
		Description: "Test description",
		ConnectorID: "test-id",
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		givenConnector *Connector
		wantConnector  *Connector
		wantErr        error
	}{
		{
			name: "failed to create connector",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/connector", testConnectorInput, new(Connector)).
					Return(http.StatusBadRequest, errors.New("failed to create connector"))
			},
			givenConnector: testConnectorInput,
			wantConnector:  nil,
			wantErr:        errors.New("failed after 1 attempt: failed to create connector"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/connector", testConnectorInput, new(Connector)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any, _ ...reqedit.EditRequestFunc) {
						connector := output.(*Connector)
						*connector = *testConnectorOutput
					})
			},
			givenConnector: testConnectorInput,
			wantConnector:  testConnectorOutput,
			wantErr:        nil,
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

			gotConnector, gotErr := api.CreateConnector(ctx, test.givenConnector)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnector, gotConnector)
		})
	}
}

func Test_APIClient_UpdateConnector(t *testing.T) {
	t.Parallel()

	testConnector := &Connector{
		Name:                     "test-name",
		Description:              "Test description",
		BuiltInSshServiceEnabled: true,
		ConnectorID:              "test-id",
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		givenConnector *Connector
		wantConnector  *Connector
		wantErr        error
	}{
		{
			name: "failed to update connector",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/connector", testConnector, new(Connector)).
					Return(http.StatusBadRequest, errors.New("failed to update connector"))
			},
			givenConnector: testConnector,
			wantConnector:  nil,
			wantErr:        errors.New("failed after 1 attempt: failed to update connector"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/connector", testConnector, new(Connector)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any, _ ...reqedit.EditRequestFunc) {
						connector := output.(*Connector)
						*connector = *testConnector
					})
			},
			givenConnector: testConnector,
			wantConnector:  testConnector,
			wantErr:        nil,
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

			gotConnector, gotErr := api.UpdateConnector(ctx, test.givenConnector)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnector, gotConnector)
		})
	}
}

func Test_APIClient_DeleteConnector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantErr       error
	}{
		{
			name: "failed to delete connector",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/connector/test-id", nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete connector"))
			},
			givenID: "test-id",
			wantErr: errors.New("failed after 1 attempt: failed to delete connector"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/connector/test-id", nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "connector not found"})
			},
			givenID: "test-id",
			wantErr: nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/connector/test-id", nil, nil).
					Return(http.StatusOK, nil)
			},
			givenID: "test-id",
			wantErr: nil,
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

			gotErr := api.DeleteConnector(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_ConnectorTokens(t *testing.T) {
	t.Parallel()

	expiresAt, err := FlexibleTimeFrom("2030-01-01T00:00:00Z")
	require.NoError(t, err)

	createdAt, err := FlexibleTimeFrom("2020-01-01T00:00:00Z")
	require.NoError(t, err)

	testConnectorTokens := &ConnectorTokens{
		List: []ConnectorToken{
			{ID: "test-id-1", Name: "test-name-1", ExpiresAt: expiresAt, CreatedBy: "test-created-by", CreatedAt: createdAt},
			{ID: "test-id-2", Name: "test-name-2", ExpiresAt: expiresAt, CreatedBy: "test-created-by", CreatedAt: createdAt},
			{ID: "test-id-3", Name: "test-name-3", ExpiresAt: expiresAt, CreatedBy: "test-created-by", CreatedAt: createdAt},
		},
		Connector: Connector{
			Name:        "test-name",
			Description: "Test description",
			ConnectorID: "test-id",
		},
	}

	tests := []struct {
		name                string
		mockRequester       func(context.Context, *mocks.ClientHTTPRequester)
		givenConnectorID    string
		wantConnectorTokens *ConnectorTokens
		wantErr             error
	}{
		{
			name: "failed to get connector tokens",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/connector/test-id/tokens", nil, new(ConnectorTokens)).
					Return(http.StatusInternalServerError, errors.New("failed to get connector tokens"))
			},
			givenConnectorID:    "test-id",
			wantConnectorTokens: nil,
			wantErr:             errors.New("failed after 1 attempt: failed to get connector tokens"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/connector/test-id/tokens", nil, new(ConnectorTokens)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*ConnectorTokens)
						*output = *testConnectorTokens
					})
			},
			givenConnectorID:    "test-id",
			wantConnectorTokens: testConnectorTokens,
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

			gotConnectorTokens, gotErr := api.ConnectorTokens(ctx, test.givenConnectorID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnectorTokens, gotConnectorTokens)
		})
	}
}

func Test_APIClient_ConnectorToken(t *testing.T) {
	t.Parallel()

	expiresAt, err := FlexibleTimeFrom("2030-01-01T00:00:00Z")
	require.NoError(t, err)

	createdAt, err := FlexibleTimeFrom("2020-01-01T00:00:00Z")
	require.NoError(t, err)

	testConnectorToken := &ConnectorToken{
		ID:        "test-id-1",
		Name:      "test-name-1",
		ExpiresAt: expiresAt,
		CreatedBy: "test-created-by",
		CreatedAt: createdAt,
	}

	tests := []struct {
		name                  string
		mockRequester         func(context.Context, *mocks.ClientHTTPRequester)
		givenConnectorID      string
		givenConnectorTokenID string
		wantConnectorToken    *ConnectorToken
		wantErr               error
	}{
		{
			name: "failed to get connector token",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/connector/test-id/token/"+testConnectorToken.ID, nil, new(ConnectorToken)).
					Return(http.StatusInternalServerError, errors.New("failed to get connector tokens"))
			},
			givenConnectorID:      "test-id",
			givenConnectorTokenID: testConnectorToken.ID,
			wantConnectorToken:    nil,
			wantErr:               errors.New("failed after 1 attempt: failed to get connector tokens"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/connector/test-id/token/"+testConnectorToken.ID, nil, new(ConnectorToken)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*ConnectorToken)
						*output = *testConnectorToken
					})
			},
			givenConnectorID:      "test-id",
			givenConnectorTokenID: testConnectorToken.ID,
			wantConnectorToken:    testConnectorToken,
			wantErr:               nil,
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

			gotConnectorToken, gotErr := api.ConnectorToken(ctx, test.givenConnectorID, test.givenConnectorTokenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnectorToken, gotConnectorToken)
		})
	}
}

func Test_APIClient_CreateConnectorToken(t *testing.T) {
	t.Parallel()

	expiresAt, err := FlexibleTimeFrom("2030-01-01T00:00:00Z")
	require.NoError(t, err)

	createdAt, err := FlexibleTimeFrom("2020-01-01T00:00:00Z")
	require.NoError(t, err)

	testConnectorTokenInput := &ConnectorToken{
		ConnectorID: "test-connector-id",
		Name:        "test-name",
		ExpiresAt:   expiresAt,
	}
	testConnectorTokenOutput := &ConnectorToken{
		Token:     "test-token",
		ID:        "test-id",
		Name:      "test-name",
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	tests := []struct {
		name                string
		mockRequester       func(context.Context, *mocks.ClientHTTPRequester)
		givenConnectorToken *ConnectorToken
		wantConnectorToken  *ConnectorToken
		wantErr             error
	}{
		{
			name: "failed to create connector token",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/connector/token", testConnectorTokenInput, new(ConnectorToken)).
					Return(http.StatusBadRequest, errors.New("failed to create connector token"))
			},
			givenConnectorToken: testConnectorTokenInput,
			wantConnectorToken:  nil,
			wantErr:             errors.New("failed after 1 attempt: failed to create connector token"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/connector/token", testConnectorTokenInput, new(ConnectorToken)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any, _ ...reqedit.EditRequestFunc) {
						connectorToken := output.(*ConnectorToken)
						*connectorToken = *testConnectorTokenOutput
					})
			},
			givenConnectorToken: testConnectorTokenInput,
			wantConnectorToken:  testConnectorTokenOutput,
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

			gotConnectorToken, gotErr := api.CreateConnectorToken(ctx, test.givenConnectorToken)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnectorToken, gotConnectorToken)
		})
	}
}

func Test_APIClient_DeleteConnectorToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		mockRequester    func(context.Context, *mocks.ClientHTTPRequester)
		givenConnectorID string
		givenTokenID     string
		wantErr          error
	}{
		{
			name: "failed to delete connector token",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/connector/test-connector-id/token/test-token-id", nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete connector token"))
			},
			givenConnectorID: "test-connector-id",
			givenTokenID:     "test-token-id",
			wantErr:          errors.New("failed after 1 attempt: failed to delete connector token"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/connector/test-connector-id/token/test-token-id", nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "connector token not found"})
			},
			givenConnectorID: "test-connector-id",
			givenTokenID:     "test-token-id",
			wantErr:          nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/connector/test-connector-id/token/test-token-id", nil, nil).
					Return(http.StatusOK, nil)
			},
			givenConnectorID: "test-connector-id",
			givenTokenID:     "test-token-id",
			wantErr:          nil,
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

			gotErr := api.DeleteConnectorToken(ctx, test.givenConnectorID, test.givenTokenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
