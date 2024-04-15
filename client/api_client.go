package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/borderzero/border0-go/client/reqedit"
	"github.com/golang-jwt/jwt"
)

// APIClient is the client for the Border0 API.
type APIClient struct {
	http          HTTPRequester
	timeout       time.Duration
	authToken     string
	baseURL       string
	portalBaseURL string
	retryWaitMin  time.Duration // minimum time to wait
	retryWaitMax  time.Duration // maximum time to wait
	retryMax      int           // maximum number of retries
	backoff       Backoff
}

// Requester is the interface for the Border0 API client.
type Requester interface {
	TokenClaims() (jwt.MapClaims, error)
	AuthenticationService
	SocketService
	ConnectorService
	PolicyService
}

const (
	notFoundRetryMax     = 3
	notFoundRetryWaitMin = time.Millisecond * 500
	notFoundRetryWaitMax = time.Millisecond * 1000
)

// default config setting values
const (
	defaultTimeout       = 10 * time.Second                 // default timeout for requests
	defaultBaseURL       = "https://api.border0.com/api/v1" // default base URL for Border0 API
	defaultPortalBaseURL = "https://portal.border0.com"     // default base URL for Border0 (Admin) Portal
	defaultRetryWaitMin  = 1 * time.Second                  // default minimum time to wait between retries
	defaultRetryWaitMax  = 30 * time.Second                 // default maximum time to wait between retries
	defaultRetryMax      = 4                                // default maximum number of retries
)

// New creates a new Border0 API client.
func New(options ...Option) *APIClient {
	api := &APIClient{
		timeout:       defaultTimeout,
		authToken:     os.Getenv("BORDER0_AUTH_TOKEN"),
		baseURL:       os.Getenv("BORDER0_BASE_URL"),
		portalBaseURL: os.Getenv("BORDER0_PORTAL_BASE_URL"),
		retryWaitMin:  defaultRetryWaitMin,
		retryWaitMax:  defaultRetryWaitMax,
		retryMax:      defaultRetryMax,
		backoff:       ExponentialBackoff,
	}
	if api.baseURL == "" {
		api.baseURL = defaultBaseURL
	}
	if api.portalBaseURL == "" {
		api.portalBaseURL = defaultPortalBaseURL
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

func (api *APIClient) request(ctx context.Context, method, path string, input, output any, edits ...reqedit.EditRequestFunc) (code int, err error) {

	shouldRetry := false
	retryMax := api.retryMax
	waitMin := api.retryWaitMin
	waitMax := api.retryWaitMax
	retryCount := 0

	for ; ; retryCount++ {

		shouldRetry = false
		code, err = api.http.Request(ctx, method, api.baseURL+path, input, output, edits...)
		if err != nil {
			if code == http.StatusNotFound {
				shouldRetry = true
				retryMax = notFoundRetryMax
				waitMin = notFoundRetryWaitMin
				waitMax = notFoundRetryWaitMax
			}
			if code >= http.StatusInternalServerError {
				shouldRetry = true
				retryMax = api.retryMax
				waitMin = api.retryWaitMin
				waitMax = api.retryWaitMax
			}
		}

		if !shouldRetry {
			break
		}

		remain := retryMax - retryCount
		if remain <= 0 {
			break
		}

		wait := api.backoff(waitMin, waitMax, retryCount)
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

	return code, fmt.Errorf("failed after %d %s: %w", retryCount+1, attemptOrAttempts(retryCount+1), err)
}

func attemptOrAttempts(attempt int) string {
	if attempt == 1 {
		return "attempt"
	}
	return "attempts"
}
