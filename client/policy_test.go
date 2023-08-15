package client

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/borderzero/border0-go/mocks"
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
			wantErr: errors.New("policy [test-id] not found: failed after 1 attempt: 404: policy not found"),
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
		{Name: "test-name-1", Description: "Test description 1", PolicyData: testPolicyData},
		{Name: "test-name-2", Description: "Test description 2", PolicyData: testPolicyData},
		{Name: "test-name-3", Description: "Test description 3", PolicyData: testPolicyData},
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

func Test_APIClient_CreatePolicy(t *testing.T) {
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
			name: "failed to attach policy socket",
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
			name: "failed to remove policy socket",
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
