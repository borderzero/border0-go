package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient is a wrapper around http.Client that handles authentication,
// request/response encoding/decoding, and error handling.
type HTTPClient struct {
	client *http.Client
	token  string
}

// HTTPRequester is an interface for HTTPClient.
type HTTPRequester interface {
	Request(ctx context.Context, method, path string, input, output any) (int, error)
	Close()
}

const (
	// HTTP header names
	headerAccessToken = "x-access-token"
	headerAccept      = "Accept"
	headerContentType = "Content-Type"

	// HTTP header values
	applicationJSON = "application/json"
)

// Request sends an HTTP request to the API server.
func (h *HTTPClient) Request(ctx context.Context, method, path string, input, output any) (int, error) {
	// create request
	var buf bytes.Buffer
	if input != nil {
		if err := json.NewEncoder(&buf).Encode(input); err != nil {
			return 0, fmt.Errorf("failed to encode input into JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, path, &buf)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(headerAccessToken, h.token)
	if input == nil {
		req.Header.Set(headerAccept, applicationJSON)
	} else {
		req.Header.Set(headerContentType, applicationJSON)
	}

	// send request
	resp, err := h.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// handle successful response (2xx)
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		if output != nil {
			if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
				return resp.StatusCode, fmt.Errorf("failed to decode response from JSON: %w", err)
			}
		}
		return resp.StatusCode, nil
	}

	// handle error response (4xx, 5xx)
	if resp.StatusCode >= http.StatusBadRequest {
		return resp.StatusCode, APIErrorFrom(resp)
	}

	// handle unexpected response (1xx, 3xx)
	return resp.StatusCode, nil
}

// Close closes idle connections in the underlying HTTP client.
func (h *HTTPClient) Close() {
	h.client.CloseIdleConnections()
}

// APIErrorFrom creates an Error from an HTTP response.
func APIErrorFrom(resp *http.Response) Error {
	apiErr := Error{
		Code:    resp.StatusCode,
		Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
	}

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	if err := json.NewDecoder(tee).Decode(&apiErr); err != nil {
		peeked, _ := bufio.NewReader(&buf).Peek(1024)
		if len(peeked) > 0 {
			apiErr.Message = string(peeked)
		}
	}

	return apiErr
}

// Error is an error returned by the API server.
type Error struct {
	Code    int    `json:"status_code"`
	Message string `json:"error_message"`
}

// Error returns string representation of an Error.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
