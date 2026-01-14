package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/borderzero/border0-go/client/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_APIClient_IsAuthenticated(t *testing.T) {
	t.Parallel()

	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	expiredToken := func() string {
		claims := jwt.MapClaims{
			"exp": float64(time.Now().Add(-1 * time.Hour).Unix()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, _ := token.SignedString([]byte("secret"))
		return s
	}()
	validTokenWithExp := func() string {
		claims := jwt.MapClaims{
			"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, _ := token.SignedString([]byte("secret"))
		return s
	}()

	type mockConfig struct {
		method string
		path   string
		output interface{}
		code   int
		err    error
	}

	tests := []struct {
		name       string
		authToken  string
		mockConfig *mockConfig
		want       bool
		wantErr    bool
	}{
		{
			name:      "empty token returns false",
			authToken: "",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "invalid token parsing returns false",
			authToken: "invalid-token",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "expired token returns false locally",
			authToken: expiredToken,
			want:      false,
			wantErr:   false,
		},
		{
			name:      "valid local token, server 200 OK returns true",
			authToken: validTokenWithExp,
			mockConfig: &mockConfig{
				method: http.MethodGet,
				path:   "/iam/whoami",
				code:   http.StatusOK, 
				//we don't care about the body content for IsAuthenticated, just the status code
				output: &paginatedResponse[Socket]{}, 
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "valid local token (no exp), server 200 OK returns true",
			authToken: validToken, 
			mockConfig: &mockConfig{
				method: http.MethodGet,
				path:   "/iam/whoami",
				code:   http.StatusOK,
				output: &paginatedResponse[Socket]{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "server 401 returns false",
			authToken: validToken,
			mockConfig: &mockConfig{
				method: http.MethodGet,
				path:   "/iam/whoami",
				code:   http.StatusUnauthorized,
				err:    Error{Code: http.StatusUnauthorized},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:      "server 403 returns true",
			authToken: validToken,
			mockConfig: &mockConfig{
				method: http.MethodGet,
				path:   "/iam/whoami",
				code:   http.StatusForbidden,
				err:    Error{Code: http.StatusForbidden},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "server error (network) returns false and error",
			authToken: validToken,
			mockConfig: &mockConfig{
				method: http.MethodGet,
				path:   "/iam/whoami",
				code:   0,
				err:    errors.New("network error"),
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(WithAuthToken(tt.authToken))
			
			if tt.mockConfig != nil {
				requesterMock := new(mocks.ClientHTTPRequester)
				// mock
				fullURL := "https://api.border0.com/api/v1" + tt.mockConfig.path
				call := requesterMock.EXPECT().
					Request(context.Background(), tt.mockConfig.method, fullURL, nil, mock.Anything)
				
				call.Return(tt.mockConfig.code, tt.mockConfig.err).Once()
				
				api.http = requesterMock
			} else {
				requesterMock := new(mocks.ClientHTTPRequester)
				api.http = requesterMock
			}

			got, err := api.IsAuthenticated(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
