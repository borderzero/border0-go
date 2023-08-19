package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// APIClient is the client for the Border0 API.
type APIClient struct {
	http         HTTPRequester
	timeout      time.Duration
	authToken    string
	baseURL      string
	retryWaitMin time.Duration // minimum time to wait
	retryWaitMax time.Duration // maximum time to wait
	retryMax     int           // maximum number of retries
	backoff      Backoff
}

// Requester is the interface for the Border0 API client.
type Requester interface {
	TokenClaims() (jwt.MapClaims, error)
	SocketService
	ConnectorService
	PolicyService
}

const (
	defaultTimeout      = 10 * time.Second                 // default timeout for requests
	defaultBaseURL      = "https://api.border0.com/api/v1" // default base URL for Border0 API
	defaultRetryWaitMin = 1 * time.Second                  // default minimum time to wait between retries
	defaultRetryWaitMax = 30 * time.Second                 // default maximum time to wait between retries
	defaultRetryMax     = 4                                // default maximum number of retries
)

// New creates a new Border0 API client.
func New(options ...Option) *APIClient {
	api := &APIClient{
		timeout:      defaultTimeout,
		authToken:    os.Getenv("BORDER0_AUTH_TOKEN"),
		baseURL:      os.Getenv("BORDER0_BASE_URL"),
		retryWaitMin: defaultRetryWaitMin,
		retryWaitMax: defaultRetryWaitMax,
		retryMax:     defaultRetryMax,
		backoff:      ExponentialBackoff,
	}
	if api.baseURL == "" {
		api.baseURL = defaultBaseURL
	}
	for _, option := range options {
		option(api)
	}
	api.http = &HTTPClient{
		client: &http.Client{
			Timeout: api.timeout,
		},
		token: api.authToken,
	}
	return api
}

// TokenClaims returns the claims of the JWT token.
func (api *APIClient) TokenClaims() (jwt.MapClaims, error) {
	parsedJWT, _ := jwt.Parse(api.authToken, nil)
	if parsedJWT == nil || parsedJWT.Claims == nil {
		return nil, fmt.Errorf("failed to parse token")
	}

	return parsedJWT.Claims.(jwt.MapClaims), nil
}

func (api *APIClient) request(ctx context.Context, method, path string, input, output any) (code int, err error) {
	var (
		attemptedCount int
		shouldRetry    bool
	)

	for i := 0; ; i++ {
		attemptedCount++

		shouldRetry = false

		code, err = api.http.Request(ctx, method, api.baseURL+path, input, output)
		if err != nil {
			if code == http.StatusUnauthorized || code == http.StatusNotFound {
				shouldRetry = false
			} else {
				shouldRetry = true
			}
		}

		if !shouldRetry {
			break
		}

		remain := api.retryMax - i
		if remain <= 0 {
			break
		}

		wait := api.backoff(api.retryWaitMin, api.retryWaitMax, i)
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			api.http.Close()
			return 0, ctx.Err()
		case <-timer.C:
		}
	}

	// request was successful or we should not retry, return the result.
	if err == nil && !shouldRetry {
		return code, nil
	}

	return code, fmt.Errorf("failed after %d %s: %w", attemptedCount, attemptOrAttempts(attemptedCount), err)
}

func attemptOrAttempts(attempt int) string {
	if attempt == 1 {
		return "attempt"
	}
	return "attempts"
}
