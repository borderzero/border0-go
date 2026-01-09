package listen

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/listen/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Listener_ensureSocketCreated(t *testing.T) {
	t.Parallel()

	socket := &client.Socket{
		Name:       "test-socket",
		SocketType: "http",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.APIClientRequester)
		wantSocket    *client.Socket
		wantErr       error
	}{
		{
			name: "failed to fetch socket and it's not 404 found error",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					Socket(ctx, "test-socket").
					Return(nil, errors.New("failed to fetch socket"))
			},
			wantErr: errors.New("failed to fetch socket"),
		},
		{
			name: "socket not found, create a new one, but encountered an error",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					Socket(ctx, "test-socket").
					Return(nil, client.Error{Code: http.StatusNotFound, Message: "socket not found"})
				requester.EXPECT().
					CreateSocket(ctx, socket).
					Return(nil, errors.New("failed to create socket"))
			},
			wantSocket: nil,
			wantErr:    errors.New("failed to create socket"),
		},
		{
			name: "socket not found, create a new one, and succeeded",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					Socket(ctx, "test-socket").
					Return(nil, client.Error{Code: http.StatusNotFound, Message: "socket not found"})
				requester.EXPECT().
					CreateSocket(ctx, socket).
					Return(socket, nil)
			},
			wantSocket: socket,
			wantErr:    nil,
		},
		{
			name: "found socket and return it",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					Socket(ctx, "test-socket").
					Return(socket, nil)
			},
			wantSocket: socket,
			wantErr:    nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.APIClientRequester)
			test.mockRequester(ctx, requester)

			l := New(
				WithAPIClient(requester),
				WithSocketName("test-socket"),
			)
			gotSocket, gotErr := l.ensureSocketCreated(context.Background())

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantSocket, gotSocket)
		})
	}
}

func Test_Listener_ensurePoliciesAttached(t *testing.T) {
	t.Parallel()

	socketWithPolicies := &client.Socket{
		SocketID:   "test-socket-id-2",
		Name:       "test-socket-2",
		SocketType: "http",
		Policies: []client.Policy{
			{ID: "test-policy-id-1", Name: "test-policy-1"},
			{ID: "test-policy-id-2", Name: "test-policy-2"},
		},
	}

	policiesToAttach := []client.Policy{
		{ID: "test-policy-id-3", Name: "test-policy-3"},
		{ID: "test-policy-id-4", Name: "test-policy-4"},
	}

	tests := []struct {
		name              string
		mockRequester     func(context.Context, *mocks.APIClientRequester)
		givenWithPolicies Option
		givenSocket       *client.Socket
		wantErr           error
	}{
		{
			name:              "no-op, no policy names given",
			mockRequester:     func(ctx context.Context, requester *mocks.APIClientRequester) {},
			givenWithPolicies: nil,
			givenSocket:       socketWithPolicies,
			wantErr:           nil,
		},
		{
			name:              "no-op, policy names given, but no changes",
			mockRequester:     func(ctx context.Context, requester *mocks.APIClientRequester) {},
			givenWithPolicies: WithPolicies([]string{"test-policy-1", "test-policy-2"}),
			givenSocket:       socketWithPolicies,
			wantErr:           nil,
		},
		{
			name: "attach new policies, but new policies are not found",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					PoliciesByNames(ctx, []string{"test-policy-3"}).
					Return(nil, errors.New("policy [test-policy-3] does not exist, please create the policy first"))
			},
			givenWithPolicies: WithPolicies([]string{"test-policy-1", "test-policy-2", "test-policy-3"}),
			givenSocket:       socketWithPolicies,
			wantErr:           errors.New("policy [test-policy-3] does not exist, please create the policy first"),
		},
		{
			name: "attach new policies, found new policies, but failed to attach",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					PoliciesByNames(ctx, []string{"test-policy-3", "test-policy-4"}).
					Return(policiesToAttach, nil)
				requester.EXPECT().
					AttachPoliciesToSocket(ctx, []string{"test-policy-id-3", "test-policy-id-4"}, "test-socket-id-2").
					Return(errors.New("failed to attach policies to socket"))
			},
			givenWithPolicies: WithPolicies([]string{"test-policy-1", "test-policy-2", "test-policy-3", "test-policy-4"}),
			givenSocket:       socketWithPolicies,
			wantErr:           errors.New("failed to attach policies to socket"),
		},
		{
			name: "happy path - attach new policies, found new policies, and attached successfully",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					PoliciesByNames(ctx, []string{"test-policy-3", "test-policy-4"}).
					Return(policiesToAttach, nil)
				requester.EXPECT().
					AttachPoliciesToSocket(ctx, []string{"test-policy-id-3", "test-policy-id-4"}, "test-socket-id-2").
					Return(nil)
			},
			givenWithPolicies: WithPolicies([]string{"test-policy-1", "test-policy-2", "test-policy-3", "test-policy-4"}),
			givenSocket:       socketWithPolicies,
			wantErr:           nil,
		},
		{
			name: "detach policies, but failed to detach",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					RemovePoliciesFromSocket(ctx, []string{"test-policy-id-1", "test-policy-id-2"}, "test-socket-id-2").
					Return(errors.New("failed to detach policies from socket"))
			},
			givenWithPolicies: WithPolicies([]string{}),
			givenSocket:       socketWithPolicies,
			wantErr:           errors.New("failed to detach policies from socket"),
		},
		{
			name: "happy path - detach policies, and succeeded",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					RemovePoliciesFromSocket(ctx, []string{"test-policy-id-1", "test-policy-id-2"}, "test-socket-id-2").
					Return(nil)
			},
			givenWithPolicies: WithPolicies([]string{}),
			givenSocket:       socketWithPolicies,
			wantErr:           nil,
		},
		{
			name: "happy path - attach 2 new policies and detach 1 policy",
			mockRequester: func(ctx context.Context, requester *mocks.APIClientRequester) {
				requester.EXPECT().
					PoliciesByNames(ctx, []string{"test-policy-3", "test-policy-4"}).
					Return(policiesToAttach, nil)
				requester.EXPECT().
					AttachPoliciesToSocket(ctx, []string{"test-policy-id-3", "test-policy-id-4"}, "test-socket-id-2").
					Return(nil)
				requester.EXPECT().
					RemovePoliciesFromSocket(ctx, []string{"test-policy-id-1"}, "test-socket-id-2").
					Return(nil)
			},
			givenWithPolicies: WithPolicies([]string{"test-policy-2", "test-policy-3", "test-policy-4"}),
			givenSocket:       socketWithPolicies,
			wantErr:           nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.APIClientRequester)
			test.mockRequester(ctx, requester)

			options := []Option{
				WithAPIClient(requester),
			}
			if test.givenWithPolicies != nil {
				options = append(options, test.givenWithPolicies)
			}
			l := New(options...)
			gotErr := l.ensurePoliciesAttached(context.Background(), test.givenSocket)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
