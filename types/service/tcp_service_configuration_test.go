package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateTcpServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *TcpServiceConfiguration
		expectedError error
	}{
		{
			name: "Happy case for tcp service",
			configuration: &TcpServiceConfiguration{
				HostnameAndPort{Hostname: "hello.com", Port: 443},
			},
			expectedError: nil,
		},
		{
			name: "Should fail for tcp service with invalid hostname/port (missing port)",
			configuration: &TcpServiceConfiguration{
				HostnameAndPort{Hostname: "hello.com"},
			},
			expectedError: fmt.Errorf("invalid hostname or port: port is a required field"),
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
