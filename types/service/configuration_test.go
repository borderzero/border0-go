package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *Configuration
		expectedError error
	}{
		{
			name: "Should succeed for valid ssh socket",
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
