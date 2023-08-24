package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateTlsServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *TlsServiceConfiguration
		expectedError error
	}{
		{
			name: "Happy case for tls service type standard",
			configuration: &TlsServiceConfiguration{
				TlsServiceType: TlsServiceTypeStandard,
				StandardTlsServiceConfiguration: &StandardTlsServiceConfiguration{
					HostnameAndPort{Hostname: "hello.com", Port: 443},
				},
			},
			expectedError: nil,
		},
		{
			name: "Happy case for tls service type vpn",
			configuration: &TlsServiceConfiguration{
				TlsServiceType: TlsServiceTypeVpn,
				VpnTlsServiceConfiguration: &VpnTlsServiceConfiguration{
					VpnSubnet: "10.0.0.0/24",
					Routes:    []string{"10.0.0.0/8"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Happy case for tls service type http proxy",
			configuration: &TlsServiceConfiguration{
				TlsServiceType:                   TlsServiceTypeHttpProxy,
				HttpProxyTlsServiceConfiguration: &HttpProxyTlsServiceConfiguration{ /* NO-OP */ },
			},
			expectedError: nil,
		},
		{
			name:          "Should fail for tls service type standard with missing config",
			configuration: &TlsServiceConfiguration{TlsServiceType: TlsServiceTypeStandard},
			expectedError: fmt.Errorf("tls service configuration for tls service type \"%s\" must have standard tls service configuration defined", TlsServiceTypeStandard),
		},
		{
			name:          "Should fail for tls service type vpn with missing config",
			configuration: &TlsServiceConfiguration{TlsServiceType: TlsServiceTypeVpn},
			expectedError: fmt.Errorf("tls service configuration for tls service type \"%s\" must have vpn tls service configuration defined", TlsServiceTypeVpn),
		},
		{
			name:          "Should fail for tls service type http proxy with missing config",
			configuration: &TlsServiceConfiguration{TlsServiceType: TlsServiceTypeHttpProxy},
			expectedError: fmt.Errorf("tls service configuration for tls service type \"%s\" must have http proxy tls service configuration defined", TlsServiceTypeHttpProxy),
		},
		{
			name: "Should fail for tls service type standard with extraneous config",
			configuration: &TlsServiceConfiguration{
				TlsServiceType: TlsServiceTypeStandard,
				StandardTlsServiceConfiguration: &StandardTlsServiceConfiguration{
					HostnameAndPort{Hostname: "hello.com", Port: 443},
				},
				HttpProxyTlsServiceConfiguration: &HttpProxyTlsServiceConfiguration{}, // extra
			},
			expectedError: fmt.Errorf("tls service type \"%s\" can only have standard tls service configuration defined", TlsServiceTypeStandard),
		},
		{
			name: "Should fail for tls service type vpn with extraneous config",
			configuration: &TlsServiceConfiguration{
				TlsServiceType: TlsServiceTypeVpn,
				VpnTlsServiceConfiguration: &VpnTlsServiceConfiguration{
					VpnSubnet: "10.0.0.0/24",
					Routes:    []string{"10.0.0.0/8"},
				},
				StandardTlsServiceConfiguration: &StandardTlsServiceConfiguration{}, // extra
			},
			expectedError: fmt.Errorf("tls service type \"%s\" can only have vpn tls service configuration defined", TlsServiceTypeVpn),
		},
		{
			name: "Should fail for tls service type http proxy with extraneous config",
			configuration: &TlsServiceConfiguration{
				TlsServiceType:                   TlsServiceTypeHttpProxy,
				HttpProxyTlsServiceConfiguration: &HttpProxyTlsServiceConfiguration{ /* NO-OP */ },
				StandardTlsServiceConfiguration:  &StandardTlsServiceConfiguration{}, // extra
			},
			expectedError: fmt.Errorf("tls service type \"%s\" can only have http proxy tls service configuration defined", TlsServiceTypeHttpProxy),
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

func Test_ValidateStandardTlsServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *StandardTlsServiceConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when hostname and port are valid",
			configuration: &StandardTlsServiceConfiguration{
				HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when hostname-and-port is invalid",
			configuration: &StandardTlsServiceConfiguration{
				HostnameAndPort{
					Port: 443,
				},
			},
			expectedError: errors.New("hostname is a required field"),
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

func Test_ValidateVpnTlsServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *VpnTlsServiceConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when vpn subnet and routes are valid",
			configuration: &VpnTlsServiceConfiguration{
				VpnSubnet: "10.0.0.0/24",
				Routes:    []string{"10.0.0.0/8"},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when vpn subnet is not valid cidr",
			configuration: &VpnTlsServiceConfiguration{
				VpnSubnet: "10.0.0.0",
				Routes:    []string{"10.0.0.0/8"},
			},
			expectedError: fmt.Errorf("vpn_subnet \"%s\" is not valid IPv4 CIDR notation", "10.0.0.0"),
		},
		{
			name: "Should fail when routes contains invalid cidr in index 0",
			configuration: &VpnTlsServiceConfiguration{
				VpnSubnet: "10.0.0.0/24",
				Routes:    []string{"10.0.0.0"},
			},
			expectedError: fmt.Errorf("routes[%d] (\"%s\") is not valid IPv4 CIDR notation", 0, "10.0.0.0"),
		},
		{
			name: "Should fail when routes contains invalid cidr in index non-zero",
			configuration: &VpnTlsServiceConfiguration{
				VpnSubnet: "10.0.0.0/24",
				Routes:    []string{"192.168.0.0/24", "10.0.0.0"},
			},
			expectedError: fmt.Errorf("routes[%d] (\"%s\") is not valid IPv4 CIDR notation", 1, "10.0.0.0"),
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

func Test_ValidateHttpProxyTlsServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *HttpProxyTlsServiceConfiguration
		expectedError error
	}{
		{
			name:          "Should succeed",
			configuration: &HttpProxyTlsServiceConfiguration{ /* NO-OP */ },
			expectedError: nil,
		},
		// nothing to test
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}
