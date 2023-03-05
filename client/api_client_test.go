package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/borderzero/border0-go/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func Test_APIClient_TokenClaims(t *testing.T) {
	t.Parallel()

	badJWT := "bad.jwt.token"
	testJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	claimsFromTestJWT := jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  float64(1516239022),
	}

	tests := []struct {
		name       string
		givenToken string
		wantClaims jwt.MapClaims
		wantErr    error
	}{
		{
			name:       "cannot parse invalid token",
			givenToken: badJWT,
			wantErr:    errors.New("failed to parse token"),
		},
		{
			name:       "happy path",
			givenToken: testJWT,
			wantClaims: claimsFromTestJWT,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			api := New(
				WithAuthToken(test.givenToken),
			)

			gotClaims, gotErr := api.TokenClaims()
			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantClaims, gotClaims)
		})
	}
}

func Test_APIClient_request(t *testing.T) {
	t.Parallel()

	type (
		mockInput struct {
			Ping string `json:"ping"`
		}
		mockOutput struct {
			Pong string `json:"pong"`
		}
	)

	mockBackoff := func(min, max time.Duration, attemptNum int) time.Duration {
		return 0 // no backoff in tests
	}

	testMethod := http.MethodPost
	testPath := "/api/v1/test"
	testBaseURL := "http://test.base.url"
	testURL := testBaseURL + testPath
	testInput := &mockInput{Ping: "ping"}
	testOutput := new(mockOutput)

	errUnitTest := errors.New("expected unit test error")

	tests := []struct {
		name          string
		mockRequester func(*mocks.ClientHTTPRequester, context.Context, context.CancelFunc)
		givenOptions  []Option
		wantCode      int
		wantErr       error
	}{
		{
			name: "successful request on first try",
			mockRequester: func(requester *mocks.ClientHTTPRequester, ctx context.Context, _ context.CancelFunc) {
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusOK, nil).
					Once()
			},
			givenOptions: []Option{
				WithBaseURL(testBaseURL),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "failed request with no retry configured",
			mockRequester: func(requester *mocks.ClientHTTPRequester, ctx context.Context, _ context.CancelFunc) {
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusInternalServerError, errUnitTest).
					Once()
			},
			givenOptions: []Option{
				WithRetryMax(0),
				WithBaseURL(testBaseURL),
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  fmt.Errorf("failed after 1 attempt: %w", errUnitTest),
		},
		{
			name: "unauthorized request and should not retry",
			mockRequester: func(requester *mocks.ClientHTTPRequester, ctx context.Context, _ context.CancelFunc) {
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusUnauthorized, errUnitTest).
					Once()
			},
			givenOptions: []Option{
				WithRetryMax(3),
				WithBackoff(mockBackoff),
				WithBaseURL(testBaseURL),
			},
			wantCode: http.StatusUnauthorized,
			wantErr:  fmt.Errorf("failed after 1 attempt: %w", errUnitTest),
		},
		{
			name: "failed request with 3 max retries configured",
			mockRequester: func(requester *mocks.ClientHTTPRequester, ctx context.Context, _ context.CancelFunc) {
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusInternalServerError, errUnitTest).
					Times(4) // 1 initial + 3 retries
			},
			givenOptions: []Option{
				WithRetryMax(3),
				WithBackoff(mockBackoff),
				WithBaseURL(testBaseURL),
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  fmt.Errorf("failed after 4 attempts: %w", errUnitTest),
		},
		{
			name: "failed on first try, successful on second try",
			mockRequester: func(requester *mocks.ClientHTTPRequester, ctx context.Context, _ context.CancelFunc) {
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusInternalServerError, errUnitTest).
					Once()
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusOK, nil).
					Once()
			},
			givenOptions: []Option{
				WithRetryMax(3),
				WithRetryWaitMin(500 * time.Millisecond),
				WithRetryWaitMax(10 * time.Second),
				WithBackoff(mockBackoff),
				WithTimeout(30 * time.Second),
				WithBaseURL(testBaseURL),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "failed requests with canclled context",
			mockRequester: func(requester *mocks.ClientHTTPRequester, ctx context.Context, cancel context.CancelFunc) {
				requester.EXPECT().
					Request(ctx, testMethod, testURL, testInput, testOutput).
					Return(http.StatusInternalServerError, errUnitTest).
					Run(func(_ context.Context, _, _ string, _, _ any) {
						cancel()
					}).
					Once()
				requester.EXPECT().Close().Once()
			},
			givenOptions: []Option{
				WithRetryMax(3),
				WithBackoff(mockBackoff),
				WithBaseURL(testBaseURL),
			},
			wantCode: 0,
			wantErr:  context.Canceled,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			requesterMock := new(mocks.ClientHTTPRequester)
			test.mockRequester(requesterMock, ctx, cancel)

			api := New(test.givenOptions...)
			api.http = requesterMock

			gotCode, gotErr := api.request(ctx, testMethod, testPath, testInput, testOutput)

			if test.wantErr == nil {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, test.wantErr.Error())
			}
			assert.Equal(t, test.wantCode, gotCode)
		})
	}
}
