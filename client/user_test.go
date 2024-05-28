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

func Test_APIClient_User(t *testing.T) {
	t.Parallel()

	testUser := &User{
		Email:       "test-name",
		DisplayName: "Test description",
		Role:        "admin",
		ID:          "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantUser      *User
		wantErr       error
	}{
		{
			name: "failed to get user",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/users/%s", defaultBaseURL, "test-id"), nil, new(User)).
					Return(http.StatusBadRequest, errors.New("failed to get user"))
			},
			givenID:  "test-id",
			wantUser: nil,
			wantErr:  errors.New("failed after 1 attempt: failed to get user"),
		},
		{
			name: "404 not found error returned, let's make sure we wrap the error with more info",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/users/%s", defaultBaseURL, "test-id"), nil, new(User)).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "user not found"})
			},
			givenID: "test-id",
			wantErr: errors.New("user with ID [test-id] not found: failed after 4 attempts: 404: user not found"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				// have to use On() instead of EXPECT() because we need to set the output
				// and the Run() function would raise nil pointer panic if we use it with
				// EXPECT()
				requester.On("Request", ctx, http.MethodGet, fmt.Sprintf("%s/organizations/iam/users/%s", defaultBaseURL, "test-id"), nil, new(User)).
					Return(http.StatusOK, nil).
					Run(func(args mock.Arguments) {
						output := args.Get(4).(*User)
						*output = *testUser
					})
			},
			givenID:  "test-id",
			wantUser: testUser,
			wantErr:  nil,
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

			gotUser, gotErr := api.User(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantUser, gotUser)
		})
	}
}

func Test_APIClient_CreateUser(t *testing.T) {
	t.Parallel()

	testUserInput := &User{
		Email:       "testemail@gmail.com",
		DisplayName: "Test description",
		Role:        "admin",
	}

	testUserOutput := &User{
		Email:       "testemail@gmail.com",
		DisplayName: "Test description",
		Role:        "admin",
		ID:          "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenUser     *User
		wantUser      *User
		wantErr       error
	}{
		{
			name: "failed to create user",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/users", defaultBaseURL), testUserInput, new(User)).
					Return(http.StatusBadRequest, errors.New("failed to create user"))
			},
			givenUser: testUserInput,
			wantUser:  nil,
			wantErr:   errors.New("failed after 1 attempt: failed to create user"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPost, fmt.Sprintf("%s/organizations/iam/users", defaultBaseURL), testUserInput, new(User)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						user := output.(*User)
						*user = *testUserOutput
					})
			},
			givenUser: testUserInput,
			wantUser:  testUserOutput,
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

			gotUser, gotErr := api.CreateUser(ctx, test.givenUser)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantUser, gotUser)
		})
	}
}

func Test_APIClient_UpdateUser(t *testing.T) {
	t.Parallel()

	testUser := &User{
		Email:       "testemail@gmail.com",
		DisplayName: "Test description",
		Role:        "admin",
		ID:          "test-id",
	}

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenUser     *User
		wantUser      *User
		wantErr       error
	}{
		{
			name: "failed to update user",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/organizations/iam/users", defaultBaseURL), testUser, new(User)).
					Return(http.StatusBadRequest, errors.New("failed to update user"))
			},
			givenUser: testUser,
			wantUser:  nil,
			wantErr:   errors.New("failed after 1 attempt: failed to update user"),
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodPut, fmt.Sprintf("%s/organizations/iam/users", defaultBaseURL), testUser, new(User)).
					Return(http.StatusOK, nil).
					Run(func(_ context.Context, _, _ string, _, output any) {
						user := output.(*User)
						*user = *testUser
					})
			},
			givenUser: testUser,
			wantUser:  testUser,
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

			gotUser, gotErr := api.UpdateUser(ctx, test.givenUser)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantUser, gotUser)
		})
	}
}

func Test_APIClient_DeleteUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockRequester func(context.Context, *mocks.ClientHTTPRequester)
		givenID       string
		wantErr       error
	}{
		{
			name: "failed to delete user",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/users/%s", defaultBaseURL, "test-id"), nil, nil).
					Return(http.StatusBadRequest, errors.New("failed to delete user"))
			},
			givenID: "test-id",
			wantErr: errors.New("failed after 1 attempt: failed to delete user"),
		},
		{
			name: "404 not found error returned, but we will ignore it and return nil",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/users/%s", defaultBaseURL, "test-id"), nil, nil).
					Return(http.StatusNotFound, Error{Code: http.StatusNotFound, Message: "user not found"})
			},
			givenID: "test-id",
			wantErr: nil,
		},
		{
			name: "happy path",
			mockRequester: func(ctx context.Context, requester *mocks.ClientHTTPRequester) {
				requester.EXPECT().
					Request(ctx, http.MethodDelete, fmt.Sprintf("%s/organizations/iam/users/%s", defaultBaseURL, "test-id"), nil, nil).
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

			gotErr := api.DeleteUser(ctx, test.givenID)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
		})
	}
}
