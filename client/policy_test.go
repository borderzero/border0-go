package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/borderzero/border0-go/client/mocks"
	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testPolicyData = PolicyData{
	Version: "v1",
	Action:  []string{"database", "ssh", "http", "tls"},
	Condition: PolicyCondition{
		Who: PolicyWho{
			Email:  []string{"johndoe@example.com"},
			Domain: []string{"example.com"},
			Group:  []string{"707a6c88-5675-4c5f-9db2-13b91c1a43e8"},
		},
		Where: PolicyWhere{
			AllowedIP:  []string{"0.0.0.0/0", "::/0"},
			Country:    []string{"NL", "CA", "US", "BR", "FR"},
			CountryNot: []string{"BE"},
		},
		When: PolicyWhen{
			After:           "2022-10-13T05:12:27Z",
			Before:          "",
			TimeOfDayAfter:  "00:00 UTC",
			TimeOfDayBefore: "23:59 UTC",
		},
	},
}

var maxDuration = 3600
var testPolicyDataV2 = PolicyDataV2{
	Permissions: PolicyPermissions{
		Database: &DatabasePermissions{
			AllowedDatabases: &[]DatabasePermission{
				{
					Database:          "testdb",
					AllowedQueryTypes: &[]string{"ReadOnly"},
				},
			},
			MaxSessionDurationSeconds: &maxDuration,
		},
		SSH: &SSHPermissions{
			Shell: &SSHShellPermission{},
			Exec: &SSHExecPermission{
				Commands: &[]string{"ls", "pwd"},
			},
			SFTP: &SSHSFTPPermission{},
			TCPForwarding: &SSHTCPForwardingPermission{
				AllowedConnections: &[]SSHTcpForwardingConnection{
					{
						DestinationAddress: pointer.To("*"),
						DestinationPort:    pointer.To("443"),
					},
				},
			},
			KubectlExec: &SSHKubectlExecPermission{
				AllowedNamespaces: &[]KubectlExecNamespace{
					{
						Namespace: "test-namespace",
						PodSelector: &map[string]string{
							"app": "test-app",
						},
					},
				},
			},
			DockerExec: &SSHDockerExecPermission{
				AllowedContainers: &[]string{"test-container"},
			},
			MaxSessionDurationSeconds: &maxDuration,
			AllowedUsernames:          &[]string{"root"},
		},
		HTTP:    &HTTPPermissions{},
		TLS:     &TLSPermissions{},
		VNC:     &VNCPermissions{},
		RDP:     &RDPPermissions{},
		VPN:     &VPNPermissions{},
		Network: &NetworkPermissions{},
	},
	Condition: PolicyConditionV2{
		Who: PolicyWhoV2{
			Email:          []string{"johndoe@example.com"},
			Group:          []string{"707a6c88-5675-4c5f-9db2-13b91c1a43e8"},
			ServiceAccount: []string{"test-service-account"},
		},
		Where: PolicyWhere{
			AllowedIP:  []string{"0.0.0.0/0", "::/0"},
			Country:    []string{"NL", "CA", "US", "BR", "FR"},
			CountryNot: []string{"BE"},
		},
		When: PolicyWhen{
			After:           "2022-10-13T05:12:27Z",
			Before:          "",
			TimeOfDayAfter:  "00:00 UTC",
			TimeOfDayBefore: "23:59 UTC",
		},
	},
}

func Test_APIClient_Policy(t *testing.T) {
	t.Parallel()

	testPolicy := &Policy{
		Name:        "test-name",
		Version:     "v1",
		Description: "Test description",
		OrgWide:     true,
		PolicyData:  testPolicyData,
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantPolicy    *Policy
		wantErr       error
	}{
		{
			name: "failed to get policy",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/policy/test-id", nil, new(Policy)).
					Return(http.StatusBadRequest, errors.New("failed to get policy"))
			},
			givenID:    "test-id",
			wantPolicy: nil,
			wantErr:    errors.New("failed after 1 attempt: failed to get policy"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/policy/test-id", nil, new(Policy)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "policy not found"})
			},
			givenID: "test-id",
			wantErr: fmt.Errorf("policy [test-id] not found: failed after %d %s: 404: policy not found", notFoundRetryMax+1, attemptOrAttempts(notFoundRetryMax+1)),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/policy/test-id", nil, new(Policy)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*Policy)
						*output = *testPolicy
					})
			},
			givenID:    "test-id",
			wantPolicy: testPolicy,
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

			gotPolicy, gotErr := api.Policy(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantPolicy, gotPolicy)
		})
	}
}

func Test_APIClient_Policies(t *testing.T) {
	t.Parallel()

	testPolicies := []Policy{
		{Name: "test-name-1", Version: "v1", Description: "Test description 1", PolicyData: testPolicyData},
		{Name: "test-name-2", Version: "v1", Description: "Test description 2", PolicyData: testPolicyData},
		{Name: "test-name-3", Version: "v2", Description: "Test description 3", PolicyData: testPolicyDataV2},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		wantPolicies  []Policy
		wantErr       error
	}{
		{
			name: "failed to get policies",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/policies", nil, new([]Policy)).
					Return(http.StatusInternalServerError, errors.New("failed to get policies"))
			},
			wantPolicies: nil,
			wantErr:      errors.New("failed after 1 attempt: failed to get policies"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/policies", nil, new([]Policy)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*[]Policy)
						*output = testPolicies
					})
			},
			wantPolicies: testPolicies,
			wantErr:      nil,
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

			gotPolicies, gotErr := api.Policies(ctx)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantPolicies, gotPolicies)
		})
	}
}

func Test_APIClient_PoliciesByNames(t *testing.T) {
	t.Parallel()

	onlyOnePolicy := &Policy{Name: "test-name-1", Description: "Test description 1", PolicyData: testPolicyData}

	multiplePolicies := []Policy{
		{Name: "test-name-1", Version: "v1", Description: "Test description 1", PolicyData: testPolicyData},
		{Name: "test-name-2", Version: "v1", Description: "Test description 2", PolicyData: testPolicyData},
		{Name: "test-name-3", Version: "v2", Description: "Test description 3", PolicyData: testPolicyDataV2},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenNames    []string
		wantPolicies  []Policy
		wantErr       error
	}{
		{
			name: "no policy names provided",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
			},
			givenNames:   nil,
			wantPolicies: nil,
			wantErr:      errors.New("no policy names provided"),
		},
		{
			name: "1 name provided but not found",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/policies/find?name=test-name-1", nil, new(Policy)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "policy not found"})
			},
			givenNames:   []string{"test-name-1"},
			wantPolicies: nil,
			wantErr:      errors.New("policy [test-name-1] does not exist, please create the policy first"),
		},
		{
			name: "1 name provided but failed to find policy",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/policies/find?name=test-name-1", nil, new(Policy)).
					Return(http.StatusInternalServerError, errors.New("failed to find policy"))
			},
			givenNames:   []string{"test-name-1"},
			wantPolicies: nil,
			wantErr:      errors.New("failed after 1 attempt: failed to find policy"),
		},
		{
			name: "happy path - 1 name provided and found policy",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/policies/find?name=test-name-1", nil, new(Policy)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*Policy)
						*output = *onlyOnePolicy
					})
			},
			givenNames:   []string{"test-name-1"},
			wantPolicies: []Policy{*onlyOnePolicy},
			wantErr:      nil,
		},
		{
			name: "multiple names provided but failed to fetch policies",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, defaultBaseURL+"/policies", nil, new([]Policy)).
					Return(http.StatusInternalServerError, errors.New("failed to fetch policies"))
			},
			givenNames:   []string{"test-name-1", "test-name-2", "test-name-3"},
			wantPolicies: nil,
			wantErr:      errors.New("failed after 1 attempt: failed to fetch policies"),
		},
		{
			name: "multiple names provided but one of them is not found",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/policies", nil, new([]Policy)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*[]Policy)
						*output = multiplePolicies
					})
			},
			givenNames:   []string{"not-exists", "test-name-1"},
			wantPolicies: nil,
			wantErr:      errors.New("policy [not-exists] does not exist, please create the policy first"),
		},
		{
			name: "happy path - multiple names provided and found policies",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, defaultBaseURL+"/policies", nil, new([]Policy)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*[]Policy)
						*output = multiplePolicies
					})
			},
			givenNames:   []string{"test-name-1", "test-name-2", "test-name-3"},
			wantPolicies: multiplePolicies,
			wantErr:      nil,
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

			gotPolicies, gotErr := api.PoliciesByNames(ctx, test.givenNames...)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantPolicies, gotPolicies)
		})
	}
}

func Test_APIClient_CreatePolicy(t *testing.T) {
	t.Parallel()

	testPolicyInput := &Policy{
		Name:        "test-name",
		Description: "Test description",
		OrgWide:     true,
		Version:     "v1",
		PolicyData:  testPolicyData,
	}
	testPolicyOutput := &Policy{
		ID:          "test-id",
		Name:        "test-name",
		Version:     "v1",
		Description: "Test description",
		OrgWide:     true,
		PolicyData:  testPolicyData,
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenPolicy   *Policy
		wantPolicy    *Policy
		wantErr       error
	}{
		{
			name: "failed to create policy",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/policies", testPolicyInput, new(Policy)).
					Return(http.StatusBadRequest, errors.New("failed to create policy"))
			},
			givenPolicy: testPolicyInput,
			wantPolicy:  nil,
			wantErr:     errors.New("failed after 1 attempt: failed to create policy"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, defaultBaseURL+"/policies", testPolicyInput, new(Policy)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						policy := output.(*Policy)
						*policy = *testPolicyOutput
					})
			},
			givenPolicy: testPolicyInput,
			wantPolicy:  testPolicyOutput,
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

			gotPolicy, gotErr := api.CreatePolicy(ctx, test.givenPolicy)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantPolicy, gotPolicy)
		})
	}
}

func Test_APIClient_UpdatePolicy(t *testing.T) {
	t.Parallel()

	testPolicyInput := &Policy{
		Name:        "test-name",
		Description: "Test description",
		OrgWide:     true,
		PolicyData:  testPolicyData,
	}
	testPolicyOutput := &Policy{
		ID:          "test-id",
		Name:        "test-name",
		Description: "Test description",
		OrgWide:     true,
		PolicyData:  testPolicyData,
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		givenPolicy   *Policy
		wantPolicy    *Policy
		wantErr       error
	}{
		{
			name: "failed to update policy",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/policy/test-id", testPolicyInput, new(Policy)).
					Return(http.StatusBadRequest, errors.New("failed to update policy"))
			},
			givenID:     "test-id",
			givenPolicy: testPolicyInput,
			wantPolicy:  nil,
			wantErr:     errors.New("failed after 1 attempt: failed to update policy"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/policy/test-id", testPolicyInput, new(Policy)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						policy := output.(*Policy)
						*policy = *testPolicyOutput
					})
			},
			givenID:     "test-id",
			givenPolicy: testPolicyInput,
			wantPolicy:  testPolicyOutput,
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

			gotPolicy, gotErr := api.UpdatePolicy(ctx, test.givenID, test.givenPolicy)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantPolicy, gotPolicy)
		})
	}
}

func Test_APIClient_DeletePolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantErr       error
	}{
		{
			name: "failed to delete policy",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/policy/test-id", nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete policy"))
			},
			givenID: "test-id",
			wantErr: errors.New("failed after 1 attempt: failed to delete policy"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/policy/test-id", nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "policy not found"})
			},
			givenID: "test-id",
			wantErr: nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, defaultBaseURL+"/policy/test-id", nil, nil).
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

			gotErr := api.DeletePolicy(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_AttachPolicySocket(t *testing.T) {
	t.Parallel()

	testPolicySocketInput := &PolicySocketAttachments{
		Actions: []PolicySocketAttachment{
			{Action: "add", ID: "test-socket-id"},
		},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenPolicyID string
		givenSocketID string
		wantErr       error
	}{
		{
			name: "failed to attach policy to socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/policy/test-policy-id/socket", testPolicySocketInput, nil).
					Return(http.StatusBadRequest, errors.New("failed to attach policy to socket"))
			},
			givenPolicyID: "test-policy-id",
			givenSocketID: "test-socket-id",
			wantErr:       errors.New("failed after 1 attempt: failed to attach policy to socket"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/policy/test-policy-id/socket", testPolicySocketInput, nil).
					Return(http.StatusOK, nil)
			},
			givenPolicyID: "test-policy-id",
			givenSocketID: "test-socket-id",
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

			gotErr := api.AttachPolicyToSocket(ctx, test.givenPolicyID, test.givenSocketID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_RemovePolicySocket(t *testing.T) {
	t.Parallel()

	testPolicySocketInput := &PolicySocketAttachments{
		Actions: []PolicySocketAttachment{
			{Action: "remove", ID: "test-socket-id"},
		},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenPolicyID string
		givenSocketID string
		wantErr       error
	}{
		{
			name: "failed to remove policy from socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/policy/test-policy-id/socket", testPolicySocketInput, nil).
					Return(http.StatusBadRequest, errors.New("failed to remove policy socket attachment"))
			},
			givenPolicyID: "test-policy-id",
			givenSocketID: "test-socket-id",
			wantErr:       errors.New("failed after 1 attempt: failed to remove policy socket attachment"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/policy/test-policy-id/socket", testPolicySocketInput, nil).
					Return(http.StatusOK, nil)
			},
			givenPolicyID: "test-policy-id",
			givenSocketID: "test-socket-id",
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

			gotErr := api.RemovePolicyFromSocket(ctx, test.givenPolicyID, test.givenSocketID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_AttachPoliciesSocket(t *testing.T) {
	t.Parallel()

	testPolicySocketInput := &PolicySocketAttachments{
		Actions: []PolicySocketAttachment{
			{Action: "add", ID: "test-policy-id-1"},
			{Action: "add", ID: "test-policy-id-2"},
			{Action: "add", ID: "test-policy-id-3"},
		},
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		givenPolicyIDs []string
		givenSocketID  string
		wantErr        error
	}{
		{
			name: "failed to attach policies to socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/socket/test-socket-id/policy", testPolicySocketInput, nil).
					Return(http.StatusBadRequest, errors.New("failed to attach policies to socket"))
			},
			givenPolicyIDs: []string{"test-policy-id-1", "test-policy-id-2", "test-policy-id-3"},
			givenSocketID:  "test-socket-id",
			wantErr:        errors.New("failed after 1 attempt: failed to attach policies to socket"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/socket/test-socket-id/policy", testPolicySocketInput, nil).
					Return(http.StatusOK, nil)
			},
			givenPolicyIDs: []string{"test-policy-id-1", "test-policy-id-2", "test-policy-id-3"},
			givenSocketID:  "test-socket-id",
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

			gotErr := api.AttachPoliciesToSocket(ctx, test.givenPolicyIDs, test.givenSocketID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}

func Test_APIClient_RemovePoliciesSocket(t *testing.T) {
	t.Parallel()

	testPolicySocketInput := &PolicySocketAttachments{
		Actions: []PolicySocketAttachment{
			{Action: "remove", ID: "test-policy-id-1"},
			{Action: "remove", ID: "test-policy-id-2"},
			{Action: "remove", ID: "test-policy-id-3"},
		},
	}

	tests := []struct {
		name           string
		mockRequester  func(context.Context, *mocks.ClientHTTPRequester)
		givenPolicyIDs []string
		givenSocketID  string
		wantErr        error
	}{
		{
			name: "failed to remove policy from socket",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/socket/test-socket-id/policy", testPolicySocketInput, nil).
					Return(http.StatusBadRequest, errors.New("failed to remove policies from socket"))
			},
			givenPolicyIDs: []string{"test-policy-id-1", "test-policy-id-2", "test-policy-id-3"},
			givenSocketID:  "test-socket-id",
			wantErr:        errors.New("failed after 1 attempt: failed to remove policies from socket"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, defaultBaseURL+"/socket/test-socket-id/policy", testPolicySocketInput, nil).
					Return(http.StatusOK, nil)
			},
			givenPolicyIDs: []string{"test-policy-id-1", "test-policy-id-2", "test-policy-id-3"},
			givenSocketID:  "test-socket-id",
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

			gotErr := api.RemovePoliciesFromSocket(ctx, test.givenPolicyIDs, test.givenSocketID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
