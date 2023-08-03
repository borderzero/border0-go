package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPClient_Request(t *testing.T) {
	t.Parallel()

	type (
		mockInput struct {
			Ping string `json:"ping"`
		}
		mockOutput struct {
			Pong string `json:"pong"`
		}
	)

	testToken := "test_token"

	tests := []struct {
		name           string
		mockServerCode int
		mockServerResp string
		givenMethod    string
		givenPath      string
		givenInput     any
		givenOutput    any
		wantErrMsg     string
		wantOutput     any
	}{
		{
			name:        "failed to create new request",
			givenMethod: http.MethodGet,
			givenPath:   `%%`, // invalid url path
			wantErrMsg:  `failed to create request`,
		},
		{
			name:           "failed to decode a successful response",
			mockServerCode: http.StatusOK,
			mockServerResp: `%%`,
			givenMethod:    http.MethodGet,
			givenPath:      "/failed_to_decode_a_successful_response",
			givenOutput:    new(mockOutput),
			wantErrMsg:     "failed to decode response from JSON: invalid character '%' looking for beginning of value",
		},
		{
			name:           "bad request returned from the server",
			mockServerCode: http.StatusBadRequest,
			mockServerResp: `{"status_code": 400, "error_message":"bad request"}`,
			givenMethod:    http.MethodPost,
			givenPath:      "/bad_request",
			givenInput:     &mockInput{Ping: "bad request"},
			givenOutput:    new(mockOutput),
			wantErrMsg:     "400: bad request",
		},
		{
			name:           "not an error response and not a successful response",
			mockServerCode: http.StatusMovedPermanently,
			mockServerResp: ``, // empty response
			givenMethod:    http.MethodGet,
			givenPath:      "/not_an_error_response_and_not_a_successful_response",
		},
		{
			name:           "successful response from get request without input",
			mockServerCode: http.StatusOK,
			mockServerResp: `{"pong":"successful response from get request without input"}`,
			givenMethod:    http.MethodGet,
			givenPath:      "/successful_response_from_get_request_without_input",
			givenOutput:    new(mockOutput),
		},
		{
			name:           "successful response from post request with input",
			mockServerCode: http.StatusOK,
			mockServerResp: `{"pong":"successful response from post request with input"}`,
			givenMethod:    http.MethodPost,
			givenPath:      "/successful_response_from_post_request_with_input",
			givenInput:     &mockInput{Ping: "successful response from post request with input"},
			givenOutput:    new(mockOutput),
			wantOutput:     &mockOutput{Pong: "successful response from post request with input"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.givenMethod, r.Method)
				assert.Equal(t, test.givenPath, r.URL.Path)
				assert.Equal(t, testToken, r.Header.Get(headerAccessToken))
				if test.givenInput == nil {
					assert.Equal(t, applicationJSON, r.Header.Get(headerAccept))
				} else {
					assert.Equal(t, applicationJSON, r.Header.Get(headerContentType))
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					input, err := json.Marshal(test.givenInput)
					assert.NoError(t, err)
					assert.JSONEq(t, string(input), string(body))
				}
				w.WriteHeader(test.mockServerCode)
				w.Write([]byte(test.mockServerResp))
			}))
			defer ts.Close()

			requester := new(HTTPClient)
			requester.token = testToken
			requester.client = ts.Client()
			defer requester.Close()

			gotCode, gotErr := requester.Request(
				context.Background(),
				test.givenMethod,
				ts.URL+test.givenPath,
				test.givenInput,
				test.givenOutput,
			)

			// check error
			if test.wantErrMsg == "" {
				assert.NoError(t, gotErr)
			} else {
				assert.Error(t, gotErr)
				assert.Contains(t, gotErr.Error(), test.wantErrMsg)
			}
			// check status code
			assert.Equal(t, test.mockServerCode, gotCode)
			// check output if it is expected
			if test.wantOutput != nil {
				assert.Equal(t, test.wantOutput, test.givenOutput)
			}
		})
	}
}

func Test_APIErrorFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		givenResp *http.Response
		wantErr   error
	}{
		{
			name: "cannot decode error response so just return whatever responded",
			givenResp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("not a json error response")),
			},
			wantErr: errors.New("400: not a json error response"),
		},
		{
			name: "error response can be decoded, but error_message field is empty, so fallback on message field",
			givenResp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"message":"bad request"}`)),
			},
			wantErr: errors.New("400: bad request"),
		},
		{
			name: "error response can be decoded, but both error_message and message fields are empty",
			givenResp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			wantErr: errors.New("400: unexpected status code"),
		},
		{
			name: "error response can be decoded",
			givenResp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"status_code": 400, "error_message":"bad request"}`)),
			},
			wantErr: errors.New("400: bad request"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotErr := APIErrorFrom(test.givenResp)
			assert.EqualError(t, gotErr, test.wantErr.Error())
		})
	}
}

func Test_NotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		givenErr error
		want     bool
	}{
		{
			name:     "no - not an error",
			givenErr: nil,
			want:     false,
		},
		{
			name:     "no - not an Error typed error",
			givenErr: errors.New("not an Error typed error"),
			want:     false,
		},
		{
			name: "no - it's a bad request error",
			givenErr: Error{
				Code:    http.StatusBadRequest,
				Message: "bad request",
			},
			want: false,
		},
		{
			name: "yes - not found error",
			givenErr: Error{
				Code:    http.StatusNotFound,
				Message: "not found",
			},
			want: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NotFound(test.givenErr)
			assert.Equal(t, test.want, got)
		})
	}
}
