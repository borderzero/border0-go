package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/borderzero/border0-go/lib/types/set"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateHostnameAndPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *HostnameAndPort
		expectedError error
	}{
		{
			name: "Should succeed when host and post are defined",
			configuration: &HostnameAndPort{
				Hostname: "hello.com",
				Port:     443,
			},
			expectedError: nil,
		},
		{
			name: "Should fail when hostname is not defined",
			configuration: &HostnameAndPort{
				Port: 443,
			},
			expectedError: errors.New("hostname is a required field"),
		},
		{
			name: "Should fail when port is not defined",
			configuration: &HostnameAndPort{
				Hostname: "hello.com",
			},
			expectedError: errors.New("port is a required field"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_validateUsernameWithProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		usernameProvider          string
		username                  string
		allowedEmptyUserProviders set.Set[string]
		expectedError             error
	}{
		{
			name:                      "Should succeed when username provider is UsernameProviderDefined and username is non empty",
			usernameProvider:          UsernameProviderDefined,
			username:                  "user123",
			allowedEmptyUserProviders: set.New[string](),
			expectedError:             nil,
		},
		{
			name:                      "Should fail when username provider is UsernameProviderDefined and username is empty",
			usernameProvider:          UsernameProviderDefined,
			username:                  "",
			allowedEmptyUserProviders: set.New[string](),
			expectedError:             fmt.Errorf("username must be provided when username_provider is \"%s\"", UsernameProviderDefined),
		},
		{
			name:                      "Should succeed when username provider is in the empty-username-provider-allowlist and username is empty",
			usernameProvider:          UsernameProviderUseConnectorUser,
			username:                  "",
			allowedEmptyUserProviders: set.New[string](UsernameProviderUseConnectorUser),
			expectedError:             nil,
		},
		{
			name:                      "Should fail when username provider is NOT the empty-username-provider-allowlist and username is empty",
			usernameProvider:          UsernameProviderUseConnectorUser,
			username:                  "",
			allowedEmptyUserProviders: set.New[string](),
			expectedError:             fmt.Errorf("username_provider %s is not valid", UsernameProviderUseConnectorUser),
		},
		{
			name:                      "Should succeed when username provider is in the empty-username-provider-allowlist and username is NOT empty",
			usernameProvider:          UsernameProviderUseConnectorUser,
			username:                  "user123",
			allowedEmptyUserProviders: set.New[string](UsernameProviderUseConnectorUser),
			expectedError:             fmt.Errorf("username must be empty when username_provider is %s", UsernameProviderUseConnectorUser),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				test.expectedError,
				validateUsernameWithProvider(
					test.usernameProvider,
					test.username,
					test.allowedEmptyUserProviders,
				),
			)
		})
	}
}
