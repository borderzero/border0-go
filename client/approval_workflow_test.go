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

func Test_APIClient_ApprovalWorkflow(t *testing.T) {
	t.Parallel()

	testFlow := &ApprovalWorkflow{
		ID:                "test-id",
		Name:              "test flow",
		Description:       "test flow description",
		SocketIDs:         []string{"socket-1"},
		RequesterGroupIDs: []string{"group-1"},
		ApproverUserIDs:   []string{"user-1"},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantFlow      *ApprovalWorkflow
		wantErr       error
	}{
		{
			name: "failed to get approval workflow",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), nil, new(ApprovalWorkflow)).
					Return(http.StatusBadRequest, errors.New("failed to get approval workflow"))
			},
			givenID:  "test-id",
			wantFlow: nil,
			wantErr:  errors.New("failed after 1 attempt: failed to get approval workflow"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), nil, new(ApprovalWorkflow)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "approval workflow not found"})
			},
			givenID:  "test-id",
			wantFlow: nil,
			wantErr:  errors.New("approval workflow with ID [test-id] not found: failed after 4 attempts: 404: approval workflow not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), nil, new(ApprovalWorkflow)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*ApprovalWorkflow)
						*output = *testFlow
					})
			},
			givenID:  "test-id",
			wantFlow: testFlow,
			wantErr:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.ClientHTTPRequester)
			test.mockRequester(ctx, requester)

			api := New(
				WithRetryMax(0),
			)
			api.http = requester

			gotFlow, gotErr := api.ApprovalWorkflow(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantFlow, gotFlow)
		})
	}
}

func Test_APIClient_ApprovalWorkflows(t *testing.T) {
	t.Parallel()

	testFlows := &ApprovalWorkflows{
		List: []ApprovalWorkflow{
			{ID: "flow-1", Name: "flow-1"},
			{ID: "flow-2", Name: "flow-2"},
		},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		wantFlows     *ApprovalWorkflows
		wantErr       error
	}{
		{
			name: "failed to list approval workflows",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/approval_workflows", defaultBaseURL), nil, new(ApprovalWorkflows)).
					Return(http.StatusBadRequest, errors.New("failed to list approval workflows"))
			},
			wantFlows: nil,
			wantErr:   errors.New("failed after 1 attempt: failed to list approval workflows"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.On("Request", ctx, http.MethodGet, fmt.Sprintf("%s/approval_workflows", defaultBaseURL), nil, new(ApprovalWorkflows)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*ApprovalWorkflows)
						*output = *testFlows
					})
			},
			wantFlows: testFlows,
			wantErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.ClientHTTPRequester)
			test.mockRequester(ctx, requester)

			api := New(
				WithRetryMax(0),
			)
			api.http = requester

			gotFlows, gotErr := api.ApprovalWorkflows(ctx)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantFlows, gotFlows)
		})
	}
}

func Test_APIClient_CreateApprovalWorkflow(t *testing.T) {
	t.Parallel()

	testFlowInput := &ApprovalWorkflow{
		Name:              "test flow",
		Description:       "test flow description",
		SocketTags:        map[string]string{"env": "prod"},
		RequesterGroupIDs: []string{"group-1"},
		ApproverGroupIDs:  []string{"group-2"},
	}

	testFlowOutput := &ApprovalWorkflow{
		ID:                "test-id",
		Name:              "test flow",
		Description:       "test flow description",
		SocketTags:        map[string]string{"env": "prod"},
		RequesterGroupIDs: []string{"group-1"},
		ApproverGroupIDs:  []string{"group-2"},
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenFlow     *ApprovalWorkflow
		wantFlow      *ApprovalWorkflow
		wantErr       error
	}{
		{
			name: "failed to create approval workflow",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/approval_workflows", defaultBaseURL), testFlowInput, new(ApprovalWorkflow)).
					Return(http.StatusBadRequest, errors.New("failed to create approval workflow"))
			},
			givenFlow: testFlowInput,
			wantFlow:  nil,
			wantErr:   errors.New("failed after 1 attempt: failed to create approval workflow"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/approval_workflows", defaultBaseURL), testFlowInput, new(ApprovalWorkflow)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						flow := output.(*ApprovalWorkflow)
						*flow = *testFlowOutput
					})
			},
			givenFlow: testFlowInput,
			wantFlow:  testFlowOutput,
			wantErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.ClientHTTPRequester)
			test.mockRequester(ctx, requester)

			api := New(
				WithRetryMax(0),
			)
			api.http = requester

			gotFlow, gotErr := api.CreateApprovalWorkflow(ctx, test.givenFlow)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantFlow, gotFlow)
		})
	}
}

func Test_APIClient_UpdateApprovalWorkflow(t *testing.T) {
	t.Parallel()

	testFlow := &ApprovalWorkflow{
		ID:                "test-id",
		Name:              "test flow",
		Description:       "test flow description",
		SocketIDs:         []string{"socket-1"},
		RequesterUserIDs:  []string{"user-1"},
		ApproverUserIDs:   []string{"user-2"},
		AllowSelfApproval: true,
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenFlow     *ApprovalWorkflow
		wantFlow      *ApprovalWorkflow
		wantErr       error
	}{
		{
			name: "failed to update approval workflow",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), testFlow, new(ApprovalWorkflow)).
					Return(http.StatusBadRequest, errors.New("failed to update approval workflow"))
			},
			givenFlow: testFlow,
			wantFlow:  nil,
			wantErr:   errors.New("failed after 1 attempt: failed to update approval workflow"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), testFlow, new(ApprovalWorkflow)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "approval workflow not found"})
			},
			givenFlow: testFlow,
			wantFlow:  nil,
			wantErr:   errors.New("approval workflow with ID [test-id] not found: failed after 4 attempts: 404: approval workflow not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), testFlow, new(ApprovalWorkflow)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						flow := output.(*ApprovalWorkflow)
						*flow = *testFlow
					})
			},
			givenFlow: testFlow,
			wantFlow:  testFlow,
			wantErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			requester := new(mocks.ClientHTTPRequester)
			test.mockRequester(ctx, requester)

			api := New(
				WithRetryMax(0),
			)
			api.http = requester

			gotFlow, gotErr := api.UpdateApprovalWorkflow(ctx, test.givenFlow)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantFlow, gotFlow)
		})
	}
}

func Test_APIClient_DeleteApprovalWorkflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantErr       error
	}{
		{
			name: "failed to delete approval workflow",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete approval workflow"))
			},
			givenID: "test-id",
			wantErr: errors.New("failed after 1 attempt: failed to delete approval workflow"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "approval workflow not found"})
			},
			givenID: "test-id",
			wantErr: nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/approval_workflows/%s", defaultBaseURL, "test-id"), nil, nil).
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

			gotErr := api.DeleteApprovalWorkflow(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
