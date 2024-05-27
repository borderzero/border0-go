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

func Test_APIClient_Group(t *testing.T) {
	t.Parallel()

	testGroup := &Group{
		DisplayName: "Test description",
		ID:          "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantGroup     *Group
		wantErr       error
	}{
		{
			name: "failed to get group",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/groups/%s", defaultBaseURL, "test-id"), nil, new(Group)).
					Return(http.StatusBadRequest, errors.New("failed to get group"))
			},
			givenID:   "test-id",
			wantGroup: nil,
			wantErr:   errors.New("failed after 1 attempt: failed to get group"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/groups/%s", defaultBaseURL, "test-id"), nil, new(Group)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "group not found"})
			},
			givenID: "test-id",
			wantErr: errors.New("group with ID [test-id] not found: failed after 4 attempts: 404: group not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/groups/%s", defaultBaseURL, "test-id"), nil, new(Group)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*Group)
						*output = *testGroup
					})
			},
			givenID:   "test-id",
			wantGroup: testGroup,
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

			gotGroup, gotErr := api.Group(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantGroup, gotGroup)
		})
	}
}

func Test_APIClient_CreateGroup(t *testing.T) {
	t.Parallel()

	testGroupInput := &Group{
		DisplayName: "Test description",
	}

	testGroupOutput := &Group{
		DisplayName: "Test description",
		ID:          "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenGroup    *Group
		wantGroup     *Group
		wantErr       error
	}{
		{
			name: "failed to create connector",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/groups", defaultBaseURL), testGroupInput, new(Group)).
					Return(http.StatusBadRequest, errors.New("failed to create group"))
			},
			givenGroup: testGroupInput,
			wantGroup:  nil,
			wantErr:    errors.New("failed after 1 attempt: failed to create group"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/groups", defaultBaseURL), testGroupInput, new(Group)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						group := output.(*Group)
						*group = *testGroupOutput
					})
			},
			givenGroup: testGroupInput,
			wantGroup:  testGroupOutput,
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

			gotGroup, gotErr := api.CreateGroup(ctx, test.givenGroup)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantGroup, gotGroup)
		})
	}
}

func Test_APIClient_UpdateGroup(t *testing.T) {
	t.Parallel()

	testGroup := &Group{
		DisplayName: "Test description",
		ID:          "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenGroup    *Group
		wantGroup     *Group
		wantErr       error
	}{
		{
			name: "failed to update group",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/organizations/iam/groups", defaultBaseURL), testGroup, new(Group)).
					Return(http.StatusBadRequest, errors.New("failed to update group"))
			},
			givenGroup: testGroup,
			wantGroup:  nil,
			wantErr:    errors.New("failed after 1 attempt: failed to update group"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/organizations/iam/groups", defaultBaseURL), testGroup, new(Group)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						group := output.(*Group)
						*group = *testGroup
					})
			},
			givenGroup: testGroup,
			wantGroup:  testGroup,
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

			gotGroup, gotErr := api.UpdateGroup(ctx, test.givenGroup)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantGroup, gotGroup)
		})
	}
}

func Test_APIClient_DeleteGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantErr       error
	}{
		{
			name: "failed to delete group",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/groups/%s", defaultBaseURL, "test-id"), nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete group"))
			},
			givenID: "test-id",
			wantErr: errors.New("failed after 1 attempt: failed to delete group"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/groups/%s", defaultBaseURL, "test-id"), nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "group not found"})
			},
			givenID: "test-id",
			wantErr: nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/groups/%s", defaultBaseURL, "test-id"), nil, nil).
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

			gotErr := api.DeleteGroup(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
