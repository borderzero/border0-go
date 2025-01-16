package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateSubnetRouterServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *SubnetRouterServiceConfiguration
		expectError   bool
		expectedError string
	}{
		{
			name:          "Happy case for subnet routes with no cidrs",
			configuration: &SubnetRouterServiceConfiguration{},
			expectError:   false,
		},
		{
			name: "Happy case for subnet routes with only one v4 cidrs, no v6 routes",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"66.66.66.66/24"},
			},
			expectError: false,
		},
		{
			name: "Happy case for subnet routes with only multiple v4 cidrs, no v6 routes",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"66.66.66.66/24", "55.55.55.55/16"},
			},
			expectError: false,
		},
		{
			name: "Happy case for subnet routes with only one v6 cidrs, no v4 routes",
			configuration: &SubnetRouterServiceConfiguration{
				IPv6CIDRRanges: []string{"2001:0db8:85a3::/64"},
			},
			expectError: false,
		},
		{
			name: "Happy case for subnet routes with only multiple v6 cidrs, no v4 routes",
			configuration: &SubnetRouterServiceConfiguration{
				IPv6CIDRRanges: []string{"2001:0db8:85a3::/64", "2001:0db8:85a4::/64"},
			},
			expectError: false,
		},
		{
			name: "Happy case for subnet routes with multiple v4 and v6 cidrs",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"66.66.66.66/24", "55.55.55.55/16"},
				IPv6CIDRRanges: []string{"2001:0db8:85a3::/64", "2001:0db8:85a4::/64"},
			},
			expectError: false,
		},
		{
			name: "Should fail with invalid v4 cidr - not a cidr",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"not a cidr"},
			},
			expectError:   true,
			expectedError: fmt.Sprintf("the cidr \"%s\" is not a valid IPv4 CIDR", "not a cidr"),
		},
		{
			name: "Should fail with invalid v4 cidr - is v6 cidr",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"2001:0db8:85a3::/64"},
			},
			expectError:   true,
			expectedError: fmt.Sprintf("the cidr \"%s\" is not a valid IPv4 CIDR", "2001:0db8:85a3::/64"),
		},
		{
			name: "Should fail with invalid v4 cidr - is duplicate",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"66.66.66.66/24", "55.55.55.55/16", "66.66.66.66/24"},
			},
			expectError:   true,
			expectedError: fmt.Sprintf("duplicate IPv4 CIDR \"%s\"", "66.66.66.66/24"),
		},
		{
			name: "Should fail with invalid v6 cidr - not a cidr",
			configuration: &SubnetRouterServiceConfiguration{
				IPv6CIDRRanges: []string{"not a cidr"},
			},
			expectError:   true,
			expectedError: fmt.Sprintf("the cidr \"%s\" is not a valid IPv6 CIDR", "not a cidr"),
		},
		{
			name: "Should fail with invalid v6 cidr - is v4 cidr",
			configuration: &SubnetRouterServiceConfiguration{
				IPv6CIDRRanges: []string{"66.66.66.66/24"},
			},
			expectError:   true,
			expectedError: fmt.Sprintf("the cidr \"%s\" is not a valid IPv6 CIDR", "66.66.66.66/24"),
		},
		{
			name: "Should fail with invalid v6 cidr - is duplicate",
			configuration: &SubnetRouterServiceConfiguration{
				IPv6CIDRRanges: []string{"2001:0db8:85a3::/64", "2001:0db8:85a4::/64", "2001:0db8:85a3::/64"},
			},
			expectError:   true,
			expectedError: fmt.Sprintf("duplicate IPv6 CIDR \"%s\"", "2001:0db8:85a3::/64"),
		},
		{
			name: "Should NOT accept IPv4-mapped IPv6 CIDRs as IPv4 CIDRs",
			configuration: &SubnetRouterServiceConfiguration{
				IPv4CIDRRanges: []string{"::ffff:c0a8:0100/24"}, // ::ffff:c0a8:0100 == 192.168.1.0
			},
			expectError:   true,
			expectedError: fmt.Sprintf("the cidr \"%s\" is an IPv4-mapped IPv6 CIDR and must be passed as an IPv6 CIDR", "::ffff:c0a8:0100/24"),
		},
		{
			name: "Should accept IPv4-mapped IPv6 CIDRs as IPv6 CIDRs",
			configuration: &SubnetRouterServiceConfiguration{
				IPv6CIDRRanges: []string{"::ffff:c0a8:0100/24"}, // ::ffff:c0a8:0100 == 192.168.1.0
			},
			expectError: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.configuration.Validate()
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
