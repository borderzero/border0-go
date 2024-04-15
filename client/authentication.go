package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/borderzero/border0-go/client/auth"
	"github.com/cenkalti/backoff"
	"github.com/golang-jwt/jwt"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/term"
)

// AuthenticationService is an interface for API client methods that interact with Border0 API to manage authentication.
type AuthenticationService interface {
	// TODO: IsAuthenticated(ctx context.Context) (bool, error)
	Authenticate(ctx context.Context, opts ...auth.Option) error
}

// Authenticate authenticates the client.
func (api *APIClient) Authenticate(ctx context.Context, opts ...auth.Option) error {
	config, err := auth.GetConfig(opts...)
	if err != nil {
		return fmt.Errorf("failed to initialize authentication configuration: %v", err)
	}

	token, err := doAuthFlow(ctx, api, config)
	if err != nil {
		return err
	}

	if config.ShouldWriteTokensToDisk() {
		tokenStorageFilePath := config.GetTokenStorageFilePath()
		if err = os.MkdirAll(tokenStorageFilePath, 0750); err != nil {
			return fmt.Errorf("failed to ensure directories for Border0 token file: %v", err)
		}
		if err = os.WriteFile(tokenStorageFilePath, []byte(token), 0600); err != nil {
			return fmt.Errorf("failed to write Border0 token: %v", err)
		}
	}

	api.authToken = token // set token as the authentication token to use in api client
	return nil
}

type legacyLoginResponse struct {
	Token string `json:"token"`
}

type legacyLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type deviceAuthorization struct {
	Token string `json:"token,omitempty"`
}

type deviceAuthorizationStatus struct {
	Token string `json:"token,omitempty"`
	State string `json:"state,omitempty"`
}

func doAuthFlow(ctx context.Context, api *APIClient, config *auth.Config) (string, error) {
	if config.ShouldUseLegacyAuthentication() {
		email := config.GetEmail()
		password := config.GetPassword()

		// if email not set read from terminal
		if email == "" {
			var err error
			email, err = readFromTerminal("email")
			if err != nil {
				return "", err
			}
		}

		// if password not set read from terminal
		if password == "" {
			var err error
			password, err = readFromTerminal("password")
			if err != nil {
				return "", err
			}
		}

		loginReq := &legacyLoginRequest{Email: email, Password: password}

		var loginResp legacyLoginResponse
		_, err := api.request(ctx, http.MethodPost, "login", loginReq, &loginResp)
		if err != nil {
			return "", err
		}

		return loginResp.Token, nil
	}

	// execute client device authorization
	deviceAuthToken, err := createDeviceAuthorization(ctx, api)
	if err != nil {
		return "", fmt.Errorf("failed to initiate Border0 device authorization flow: %v", err)
	}
	token, err := handleDeviceAuthorization(ctx, api, deviceAuthToken, config.ShouldTryOpeningBrowser())
	if err != nil {
		return "", fmt.Errorf("failed to authenticate you against Border0: %v", err)
	}
	return token, err
}

func readFromTerminal(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(pass), nil
}

func createDeviceAuthorization(ctx context.Context, api *APIClient) (string, error) {
	var out deviceAuthorization
	if _, err := api.request(ctx, http.MethodPost, "device_authorizations", nil, &out); err != nil {
		return "", fmt.Errorf("failed to request Border0 device authorization: %v", err)
	}
	return out.Token, nil
}

func getDeviceAuthorizationStatus(ctx context.Context, api *APIClient, deviceAuthToken string) (*deviceAuthorizationStatus, error) {
	setDeviceAuthToken := func(r *http.Request) { r.Header.Set(headerDeviceAuthToken, deviceAuthToken) }

	var out deviceAuthorizationStatus
	if _, err := api.request(ctx, http.MethodGet, "device_authorizations", nil, &out, setDeviceAuthToken); err != nil {
		return nil, fmt.Errorf("failed to retrieve Border0 device authorization status: %v", err)
	}
	return &out, nil
}

func handleDeviceAuthorization(ctx context.Context, api *APIClient, deviceAuthToken string, tryOpenBrowser bool) (string, error) {
	deviceAuthJWT, _ := jwt.Parse(deviceAuthToken, nil)
	if deviceAuthJWT == nil {
		return "", fmt.Errorf("failed to decode Border0 device authorization token")
	}
	claims := deviceAuthJWT.Claims.(jwt.MapClaims)
	deviceIdentifier := fmt.Sprint(claims["identifier"])

	// Try opening the system's browser automatically. The error is ignored because the desired behavior of the
	// handler is the same regardless of whether opening the browser fails or succeeds -- we still print the URL.
	// This is desirable because in the event opening the browser succeeds, the customer may still accidentally
	// close the new tab / browser session, or may want to authenticate in a different browser / session. In the
	// event that opening the browser fails, the customer may still complete authenticating by navigating to the
	// URL in a different device.

	url := fmt.Sprintf("%s/login?device_identifier=%v", api.portalBaseURL, url.QueryEscape(deviceIdentifier))

	fmt.Printf("Please navigate to the URL below in order to complete the login process:\n%s\n", url)

	if tryOpenBrowser {
		// check if we're on darwin (MacOS) and if we're running as sudo, if so, make sure we open the browser as the user
		// this prevents folks from not having access to credentials , sessions, etc
		sudoUsername := os.Getenv("SUDO_USER")
		sudoAttempt := false
		if runtime.GOOS == "darwin" && sudoUsername != "" {
			err := exec.Command("sudo", "-u", sudoUsername, "open", url).Run()
			if err == nil {
				// If for some reason this failed, we'll try again to standard way
				sudoAttempt = true
			}
		}
		if !sudoAttempt {
			_ = open.Run(url)
		}
	}

	receivedToken, err := pollForToken(ctx, api, deviceAuthToken)
	if err != nil {
		return "", err
	}

	return receivedToken, nil
}

// pollForToken will poll until the device is authorized.
func pollForToken(ctx context.Context, api *APIClient, deviceAuthorizationToken string) (string, error) {
	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.InitialInterval = 1 * time.Second
	exponentialBackoff.MaxInterval = 5 * time.Second
	exponentialBackoff.Multiplier = 1.3
	exponentialBackoff.MaxElapsedTime = 3 * time.Minute

	var token string

	retryFn := func() error {
		tk, err := getDeviceAuthorizationStatus(ctx, api, deviceAuthorizationToken)
		if err != nil {
			return err
		}
		token = tk.Token
		return err
	}

	err := backoff.Retry(retryFn, exponentialBackoff)
	if err != nil {
		if errors.Is(err, errors.New("unauthorized")) {
			fmt.Printf("We couldn't log you in, your session is expired or you are not authorized to perform this action: %v\n", err)
		}
		fmt.Printf("We couldn't log you in, make sure that you are properly logged in using the link above: %v\n", err)
		return "", err
	}

	fmt.Println("Login successful")

	return token, nil
}
