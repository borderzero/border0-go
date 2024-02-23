package service

import (
	"fmt"
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
		{
			name: "Happy case for vnc service",
			configuration: &RdpServiceConfiguration{
				HostnameAndPort{Hostname: "hello.com", Port: 443},
			},
			expectedError: nil,
		},
		{
			name: "Should fail for vnc service with invalid hostname/port (missing port)",
			configuration: &RdpServiceConfiguration{
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
