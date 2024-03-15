package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateVncServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *VncServiceConfiguration
		expectedError error
	}{
		{
			name: "Happy case for vnc service",
			configuration: &VncServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com", Port: 443},
				Password:        "password",
			},
			expectedError: nil,
		},
		{
			name: "Happy case for vnc service with empty password",
			configuration: &VncServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com", Port: 443},
			},
			expectedError: nil,
		},
		{
			name: "Should fail for vnc service with invalid hostname/port (missing port)",
			configuration: &VncServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com"},
				Password:        "password",
			},
			expectedError: fmt.Errorf("invalid hostname or port: port is a required field"),
		},
		{
			name: "Should fail for vnc service with invalid password (too long)",
			configuration: &VncServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com", Port: 5900},
				Password:        "longer-password",
			},
			expectedError: fmt.Errorf("password can be at most %d characters long", maxVncPasswordLength),
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
