package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/borderzero/border0-go/client/mocks"
	"github.com/borderzero/border0-go/client/reqedit"
	"github.com/borderzero/border0-go/types/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_APIClient_Socket(t *testing.T) {
	t.Parallel()

	testSocket := &Socket{
		Name:       "test-name",
		SocketID:   "test-id",
		SocketType: "test-type",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenIDOrName string
		wantSocket    *Socket
		wantErr       error
	}{
		{
			name: "failed to get socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/socket/test-name", nil, new(Socket)).
					Return(http.StatusBadRequest, errors.New("failed to get socket"))
			},
			givenIDOrName: "test-name",
			wantSocket:    nil,
			wantErr:       errors.New("failed after 1 attempt: failed to get socket"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/socket/test-name", nil, new(Socket)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "socket not found"})
			},
			givenIDOrName: "test-name",
			wantErr:       fmt.Errorf("socket [test-name] not found: failed after %d %s: 404: socket not found", notFoundRetryMax+1, attemptOrAttempts(notFoundRetryMax+1)),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/socket/test-name", nil, new(Socket)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*Socket)
						*output = *testSocket
					})
			},
			givenIDOrName: "test-name",
			wantSocket:    testSocket,
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

			gotSocket, gotErr := api.Socket(ctx, test.givenIDOrName)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantSocket, gotSocket)
		})
	}
}

func Test_APIClient_Sockets(t *testing.T) {
	t.Parallel()

	testSockets := []Socket{
		{Name: "test-name-1", SocketID: "test-id-1", SocketType: "http"},
		{Name: "test-name-2", SocketID: "test-id-2", SocketType: "ssh"},
		{Name: "test-name-3", SocketID: "test-id-3", SocketType: "database"},
		{Name: "test-name-4", SocketID: "test-id-4", SocketType: "tls"},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		wantSockets   []Socket
		wantErr       error
	}{
		{
			name: "failed to get sockets",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/socket", nil, new([]Socket)).
					Return(http.StatusInternalServerError, errors.New("failed to get sockets"))
			},
			wantSockets: nil,
			wantErr:     errors.New("failed after 1 attempt: failed to get sockets"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/socket", nil, new([]Socket)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*[]Socket)
						*output = testSockets
					})
			},
			wantSockets: testSockets,
			wantErr:     nil,
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

			gotSockets, gotErr := api.Sockets(ctx)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantSockets, gotSockets)
		})
	}
}

func Test_APIClient_CreateSocket(t *testing.T) {
	t.Parallel()

	testSocket := &Socket{
		Name:       "test-name",
		SocketID:   "test-id",
		SocketType: "test-type",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenSocket   *Socket
		wantSocket    *Socket
		wantErr       error
	}{
		{
			name: "failed to create socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/socket", testSocket, new(Socket)).
					Return(http.StatusBadRequest, errors.New("failed to create socket"))
			},
			givenSocket: testSocket,
			wantSocket:  nil,
			wantErr:     errors.New("failed after 1 attempt: failed to create socket"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/socket", testSocket, new(Socket)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any, _ ...reqedit.EditRequestFunc) {
						socket := output.(*Socket)
						*socket = *testSocket
					})
			},
			givenSocket: testSocket,
			wantSocket:  testSocket,
			wantErr:     nil,
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

			gotSocket, gotErr := api.CreateSocket(ctx, test.givenSocket)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantSocket, gotSocket)
		})
	}
}

func Test_APIClient_UpdateSocket(t *testing.T) {
	t.Parallel()

	testSocket := &Socket{
		Name:       "test-name",
		SocketID:   "test-id",
		SocketType: "test-type",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenIDOrName string
		givenSocket   *Socket
		wantSocket    *Socket
		wantErr       error
	}{
		{
			name: "failed to update socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/socket/test-name", testSocket, new(Socket)).
					Return(http.StatusBadRequest, errors.New("failed to update socket"))
			},
			givenIDOrName: "test-name",
			givenSocket:   testSocket,
			wantSocket:    nil,
			wantErr:       errors.New("failed after 1 attempt: failed to update socket"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/socket/test-name", testSocket, new(Socket)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any, _ ...reqedit.EditRequestFunc) {
						socket := output.(*Socket)
						*socket = *testSocket
					})
			},
			givenIDOrName: "test-name",
			givenSocket:   testSocket,
			wantSocket:    testSocket,
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

			gotSocket, gotErr := api.UpdateSocket(ctx, test.givenIDOrName, test.givenSocket)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantSocket, gotSocket)
		})
	}
}

func Test_APIClient_DeleteSocket(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenIDOrName string
		wantErr       error
	}{
		{
			name: "failed to delete socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/socket/test-name", nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete socket"))
			},
			givenIDOrName: "test-name",
			wantErr:       errors.New("failed after 1 attempt: failed to delete socket"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/socket/test-name", nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "socket not found"})
			},
			givenIDOrName: "test-name",
			wantErr:       nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/socket/test-name", nil, nil).
					Return(http.StatusOK, nil)
			},
			givenIDOrName: "test-name",
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

			gotErr := api.DeleteSocket(ctx, test.givenIDOrName)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_SocketUpstreamConfigs(t *testing.T) {
	t.Parallel()

	// SSH upstream config with username and password
	testConfig := service.Configuration{
		ServiceType: service.ServiceTypeSsh,
		SshServiceConfiguration: &service.SshServiceConfiguration{
			SshServiceType: service.SshServiceTypeStandard,
			StandardSshServiceConfiguration: &service.StandardSshServiceConfiguration{
				SshAuthenticationType: service.StandardSshServiceAuthenticationTypeUsernameAndPassword,
				UsernameAndPasswordAuthConfiguration: &service.UsernameAndPasswordAuthConfiguration{
					Username: "test-username",
					Password: "test-password",
				},
				HostnameAndPort: service.HostnameAndPort{
					Hostname: "test-hostname",
					Port:     22,
				},
			},
		},
	}

	testConfigs := &SocketUpstreamConfigs{
		List: []SocketUpstreamConfig{
			{
				Config: testConfig,
			},
		},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenIDOrName string
		wantConfigs   *SocketUpstreamConfigs
		wantErr       error
	}{
		{
			name: "failed to get socket upstream configs",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/socket/test-name/upstream_configurations", nil, new(SocketUpstreamConfigs)).
					Return(http.StatusInternalServerError, errors.New("failed to get socket upstream configs"))
			},
			givenIDOrName: "test-name",
			wantConfigs:   nil,
			wantErr:       errors.New("failed after 1 attempt: failed to get socket upstream configs"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/socket/test-name/upstream_configurations", nil, new(SocketUpstreamConfigs)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*SocketUpstreamConfigs)
						*output = *testConfigs
					})
			},
			givenIDOrName: "test-name",
			wantConfigs:   testConfigs,
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

			gotConfigs, gotErr := api.SocketUpstreamConfigs(ctx, test.givenIDOrName)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConfigs, gotConfigs)
		})
	}
}

func Test_APIClient_SocketConnectors(t *testing.T) {
	t.Parallel()

	testConnectors := &SocketConnectors{
		List: []SocketConnector{
			{ID: 1, ConnectorID: "test-id-1", ConnectorName: "test-name-1"},
			{ID: 2, ConnectorID: "test-id-2", ConnectorName: "test-name-2"},
			{ID: 3, ConnectorID: "test-id-3", ConnectorName: "test-name-3"},
		},
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		givenIDOrName  string
		wantConnectors *SocketConnectors
		wantErr        error
	}{
		{
			name: "failed to get socket upstream configs",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/socket/test-name/connectors", nil, new(SocketConnectors)).
					Return(http.StatusInternalServerError, errors.New("failed to get socket upstream configs"))
			},
			givenIDOrName:  "test-name",
			wantConnectors: nil,
			wantErr:        errors.New("failed after 1 attempt: failed to get socket upstream configs"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/socket/test-name/connectors", nil, new(SocketConnectors)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*SocketConnectors)
						*output = *testConnectors
					})
			},
			givenIDOrName:  "test-name",
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

			gotConnectors, gotErr := api.SocketConnectors(ctx, test.givenIDOrName)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantConnectors, gotConnectors)
		})
	}
}

func Test_APIClient_SignSocketKey(t *testing.T) {
	t.Parallel()

	testKeyToSign := &SocketKeyToSign{
		SSHPublicKey: "test-key",
	}
	testSignedKey := &SignedSocketKey{
		SignedSSHCert: "test-cert",
		HostKey:       "test-host-key",
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		givenIDOrName  string
		givenKeyToSign *SocketKeyToSign
		wantSignedKey  *SignedSocketKey
		wantErr        error
	}{
		{
			name: "failed to sign socket key",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/socket/test-name/signkey", testKeyToSign, new(SignedSocketKey)).
					Return(http.StatusBadRequest, errors.New("failed to sign socket key"))
			},
			givenIDOrName:  "test-name",
			givenKeyToSign: testKeyToSign,
			wantSignedKey:  nil,
			wantErr:        errors.New("failed after 1 attempt: failed to sign socket key"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/socket/test-name/signkey", testKeyToSign, new(SignedSocketKey)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any, _ ...reqedit.EditRequestFunc) {
						signedKey := output.(*SignedSocketKey)
						*signedKey = *testSignedKey
					})
			},
			givenIDOrName:  "test-name",
			givenKeyToSign: testKeyToSign,
			wantSignedKey:  testSignedKey,
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

			gotSignedKey, gotErr := api.SignSocketKey(ctx, test.givenIDOrName, testKeyToSign)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantSignedKey, gotSignedKey)
		})
	}
}
