package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateRdpServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *RdpServiceConfiguration
		expectedError error
	}{
		// Happy cases
		{
			name: "Happy case for rdp service without credentials",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com", Port: 3389},
			},
			expectedError: nil,
		},
		{
			name: "Happy case for rdp service with username and password",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret123"},
			},
			expectedError: nil,
		},
		{
			name: "Happy case for rdp service with username, password and domain",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret123"},
				Domain:              "CORP",
			},
			expectedError: nil,
		},
		{
			name: "Happy case for rdp service with UPN format username",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin@corp.example.com", Password: "secret123"},
			},
			expectedError: nil,
		},
		{
			name: "Happy case for rdp service with FQDN domain",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret123"},
				Domain:              "corp.example.com",
			},
			expectedError: nil,
		},
		{
			name: "Happy case for rdp service with special characters in password",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "P@ssw0rd!#$%^&*()"},
			},
			expectedError: nil,
		},
		// Hostname/port validation errors
		{
			name: "Should fail for rdp service with missing port",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com"},
			},
			expectedError: fmt.Errorf("invalid hostname or port: port is a required field"),
		},
		{
			name: "Should fail for rdp service with missing hostname",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Port: 3389},
			},
			expectedError: fmt.Errorf("invalid hostname or port: hostname is a required field"),
		},
		// Username/password mutual dependency errors
		{
			name: "Should fail when username is provided without password",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin"},
			},
			expectedError: fmt.Errorf("password is required when username is provided"),
		},
		{
			name: "Should fail when password is provided without username",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Password: "secret123"},
			},
			expectedError: fmt.Errorf("username is required when password is provided"),
		},
		// Domain dependency errors
		{
			name: "Should fail when domain is provided without username",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{Hostname: "hello.com", Port: 3389},
				Domain:          "CORP",
			},
			expectedError: fmt.Errorf("username is required when domain is provided"),
		},
		// Length validation errors
		{
			name: "Should fail when username exceeds max length",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: strings.Repeat("a", 257), Password: "secret123"},
			},
			expectedError: fmt.Errorf("username exceeds maximum length of %d characters", maxRdpUsernameLength),
		},
		{
			name: "Should fail when password exceeds max length",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: strings.Repeat("a", 257)},
			},
			expectedError: fmt.Errorf("password exceeds maximum length of %d characters", maxRdpPasswordLength),
		},
		{
			name: "Should fail when domain exceeds max length",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret123"},
				Domain:              strings.Repeat("a", 254),
			},
			expectedError: fmt.Errorf("domain exceeds maximum length of %d characters", maxRdpDomainLength),
		},
		// Control character validation errors
		{
			name: "Should fail when username contains control characters",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin\x00", Password: "secret123"},
			},
			expectedError: fmt.Errorf("username contains invalid control characters"),
		},
		{
			name: "Should fail when password contains control characters",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret\x00123"},
			},
			expectedError: fmt.Errorf("password contains invalid control characters"),
		},
		{
			name: "Should fail when domain contains control characters",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret123"},
				Domain:              "CORP\x00",
			},
			expectedError: fmt.Errorf("domain contains invalid control characters"),
		},
		// Edge cases - max length values should pass
		{
			name: "Should pass with username at max length",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: strings.Repeat("a", 256), Password: "secret123"},
			},
			expectedError: nil,
		},
		{
			name: "Should pass with password at max length",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: strings.Repeat("a", 256)},
			},
			expectedError: nil,
		},
		{
			name: "Should pass with domain at max length",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort:     HostnameAndPort{Hostname: "hello.com", Port: 3389},
				UsernameAndPassword: UsernameAndPassword{Username: "admin", Password: "secret123"},
				Domain:              strings.Repeat("a", 253),
			},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_containsControlCharacters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Normal string",
			input:    "hello world",
			expected: false,
		},
		{
			name:     "String with null character",
			input:    "hello\x00world",
			expected: true,
		},
		{
			name:     "String with bell character",
			input:    "hello\x07world",
			expected: true,
		},
		{
			name:     "String with escape character",
			input:    "hello\x1bworld",
			expected: true,
		},
		{
			name:     "String with tab (allowed)",
			input:    "hello\tworld",
			expected: false,
		},
		{
			name:     "String with newline (allowed)",
			input:    "hello\nworld",
			expected: false,
		},
		{
			name:     "String with carriage return (allowed)",
			input:    "hello\rworld",
			expected: false,
		},
		{
			name:     "String with special characters (not control)",
			input:    "P@ssw0rd!#$%^&*()",
			expected: false,
		},
		{
			name:     "String with unicode",
			input:    "héllo wörld",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, containsControlCharacters(test.input))
		})
	}
}
