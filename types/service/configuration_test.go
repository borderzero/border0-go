package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Configuration_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *Configuration
		expectedError error
	}{
		{
			name: "Should succeed for valid standard ssh socket",
			configuration: &Configuration{
				ServiceType: ServiceTypeSsh,
				SshServiceConfiguration: &SshServiceConfiguration{
					SshServiceType: SshServiceTypeStandard,
					StandardSshServiceConfiguration: &StandardSshServiceConfiguration{
						HostnameAndPort: HostnameAndPort{
							Hostname: "hello.com",
							Port:     443,
						},
						SshAuthenticationType: StandardSshServiceAuthenticationTypeUsernameAndPassword,
						UsernameAndPasswordAuthConfiguration: &UsernameAndPasswordAuthConfiguration{
							Username: "root",
							Password: "mypassword",
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed for valid built-in ssh socket",
			configuration: &Configuration{
				ServiceType: ServiceTypeSsh,
				SshServiceConfiguration: &SshServiceConfiguration{
					SshServiceType: SshServiceTypeConnectorBuiltIn,
					BuiltInSshServiceConfiguration: &BuiltInSshServiceConfiguration{
						UsernameProvider: UsernameProviderUseConnectorUser,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed for valid tls socket",
			configuration: &Configuration{
				ServiceType: ServiceTypeTls,
				TlsServiceConfiguration: &TlsServiceConfiguration{
					TlsServiceType: TlsServiceTypeStandard,
					StandardTlsServiceConfiguration: &StandardTlsServiceConfiguration{
						HostnameAndPort: HostnameAndPort{
							Hostname: "192.0.2.2",
							Port:     5900,
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed for valid DEPRECATED subnet routes socket",
			configuration: &Configuration{
				ServiceType: DEPRECATED_ServiceTypeSubnetRoutes,
				DEPRECATED_SubnetRoutesServiceConfiguration: &SubnetRouterServiceConfiguration{
					IPv4CIDRRanges: []string{"66.66.66.66/24", "55.55.55.55/16"},
					IPv6CIDRRanges: []string{"2001:0db8:85a3::/64", "2001:0db8:85a4::/64"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed for valid DEPRECATED subnet routes socket",
			configuration: &Configuration{
				ServiceType: ServiceTypeSubnetRouter,
				SubnetRouterServiceConfiguration: &SubnetRouterServiceConfiguration{
					IPv4CIDRRanges: []string{"66.66.66.66/24", "55.55.55.55/16"},
					IPv6CIDRRanges: []string{"2001:0db8:85a3::/64", "2001:0db8:85a4::/64"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when multiple socket type configurations are present",
			configuration: &Configuration{
				ServiceType: ServiceTypeSsh,
				SshServiceConfiguration: &SshServiceConfiguration{
					SshServiceType: SshServiceTypeStandard,
					StandardSshServiceConfiguration: &StandardSshServiceConfiguration{
						HostnameAndPort: HostnameAndPort{
							Hostname: "hello.com",
							Port:     443,
						},
						SshAuthenticationType: StandardSshServiceAuthenticationTypeUsernameAndPassword,
						UsernameAndPasswordAuthConfiguration: &UsernameAndPasswordAuthConfiguration{
							Username: "root",
							Password: "mypassword",
						},
					},
				},
				HttpServiceConfiguration: &HttpServiceConfiguration{
					HttpServiceType: HttpServiceTypeStandard,
				},
			},
			expectedError: errors.New(`service configuration for service type "ssh" can only have ssh service configuration defined`),
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

func Test_ConnectorSocketConfiguration_Validate(t *testing.T) {
	t.Parallel()

	badUpstreamConfig := Configuration{
		ServiceType: ServiceTypeSsh,
		// SSH configuration is missing
	}
	goodUpstreamConfig := Configuration{
		ServiceType: ServiceTypeSsh,
		SshServiceConfiguration: &SshServiceConfiguration{
			SshServiceType: SshServiceTypeConnectorBuiltIn,
			BuiltInSshServiceConfiguration: &BuiltInSshServiceConfiguration{
				UsernameProvider: UsernameProviderUseConnectorUser,
			},
		},
	}

	tests := []struct {
		name          string
		configuration *ConnectorServiceConfiguration
		expectedError error
	}{
		{
			name: "failed the upstream config validation",
			configuration: &ConnectorServiceConfiguration{
				ConnectorAuthenticationEnabled: true,
				Upstream:                       badUpstreamConfig,
			},
			expectedError: errors.New(`invalid upstream configuration: service configuration for service type "ssh" must have ssh service configuration defined`),
		},
		{
			name: "happy path",
			configuration: &ConnectorServiceConfiguration{
				ConnectorAuthenticationEnabled: true,
				Upstream:                       goodUpstreamConfig,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.configuration.Validate()
			if test.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, test.expectedError, err.Error())
			}
		})
	}
}
