package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateVpnServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *VpnServiceConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when vpn subnet and routes are valid",
			configuration: &VpnServiceConfiguration{
				DHCPPoolSubnet:   "10.0.0.0/24",
				AdvertisedRoutes: []string{"10.0.0.0/8"},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when vpn subnet is not valid cidr",
			configuration: &VpnServiceConfiguration{
				DHCPPoolSubnet:   "10.0.0.0",
				AdvertisedRoutes: []string{"10.0.0.0/8"},
			},
			expectedError: fmt.Errorf("dhcp_pool_subnet \"%s\" is not valid IPv4 CIDR notation", "10.0.0.0"),
		},
		{
			name: "Should fail when routes contains invalid cidr in index 0",
			configuration: &VpnServiceConfiguration{
				DHCPPoolSubnet:   "10.0.0.0/24",
				AdvertisedRoutes: []string{"10.0.0.0"},
			},
			expectedError: fmt.Errorf("advertised_routes[%d] (\"%s\") is not valid IPv4 CIDR notation", 0, "10.0.0.0"),
		},
		{
			name: "Should fail when routes contains invalid cidr in index non-zero",
			configuration: &VpnServiceConfiguration{
				DHCPPoolSubnet:   "10.0.0.0/24",
				AdvertisedRoutes: []string{"192.168.0.0/24", "10.0.0.0"},
			},
			expectedError: fmt.Errorf("advertised_routes[%d] (\"%s\") is not valid IPv4 CIDR notation", 1, "10.0.0.0"),
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
