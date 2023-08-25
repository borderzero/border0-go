package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/borderzero/border0-go/types/common"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateSshServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *SshServiceConfiguration
		expectedError error
	}{
		// TODO
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_ValidateStandardSshServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *StandardSshServiceConfiguration
		expectedError error
	}{
		// TODO
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_ValidateAwsSsmSshServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *AwsSsmSshServiceConfiguration
		expectedError error
	}{
		// TODO
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_ValidateBuiltInSshServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *BuiltInSshServiceConfiguration
		expectedError error
	}{
		// TODO
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_AwsEc2ICSshServiceConfiguration(t *testing.T) {
	t.Parallel()

	mockInstanceId := "i-123456789012"
	mockInstanceRegion := "us-east-1"

	tests := []struct {
		name          string
		configuration *AwsEc2ICSshServiceConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when all inputs are valid",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "133.33.33.33",
					Port:     22,
				},
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: mockInstanceRegion,
			},
			expectedError: nil,
		},
		{
			name: "Should fail when hostname is missing",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Port: 22,
				},
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: mockInstanceRegion,
			},
			expectedError: errors.New("hostname is a required field"),
		},
		{
			name: "Should fail when port is missing",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "133.33.33.33",
				},
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: mockInstanceRegion,
			},
			expectedError: errors.New("port is a required field"),
		},
		{
			name: "Should fail when ec2 instance id is missing",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "133.33.33.33",
					Port:     22,
				},
				Ec2InstanceRegion: mockInstanceRegion,
			},
			expectedError: errors.New("ec2_instance_id is a required field"),
		},
		{
			name: "Should fail when ec2 instance region is missing",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "133.33.33.33",
					Port:     22,
				},
				Ec2InstanceId: mockInstanceId,
			},
			expectedError: errors.New("ec2_instance_region is a required field"),
		},
		{
			name: "Should fail when ec2 instance region is invalid",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "133.33.33.33",
					Port:     22,
				},
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: "bad region",
			},
			expectedError: fmt.Errorf("invalid ec2_instance_region: region \"%s\" is not a valid aws region", "bad region"),
		},
		{
			name: "Should fail when aws credentials are invalid",
			configuration: &AwsEc2ICSshServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "133.33.33.33",
					Port:     22,
				},
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: mockInstanceRegion,
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId: pointer.To("BAD CREDS"),
				},
			},
			expectedError: fmt.Errorf("invalid aws_credentials: aws_secret_access_key is required when aws_access_key_id is provided"),
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

func Test_ValidateAwsSsmEc2TargetConfiguration(t *testing.T) {
	t.Parallel()

	mockInstanceId := "i-123456789012"
	mockInstanceRegion := "us-east-1"

	tests := []struct {
		name          string
		configuration *AwsSsmEc2TargetConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when all inputs are valid",
			configuration: &AwsSsmEc2TargetConfiguration{
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: mockInstanceRegion,
			},
			expectedError: nil,
		},
		{
			name: "Should fail when ec2 instance id is missing",
			configuration: &AwsSsmEc2TargetConfiguration{
				Ec2InstanceRegion: mockInstanceRegion,
			},
			expectedError: errors.New("ec2_instance_id is a required field"),
		},
		{
			name: "Should fail when ecs cluster region is missing",
			configuration: &AwsSsmEc2TargetConfiguration{
				Ec2InstanceId: mockInstanceId,
			},
			expectedError: errors.New("ec2_instance_region is a required field"),
		},
		{
			name: "Should fail when ecs instance region is invalid",
			configuration: &AwsSsmEc2TargetConfiguration{
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: "bad region",
			},
			expectedError: fmt.Errorf("invalid ec2_instance_region: region \"%s\" is not a valid aws region", "bad region"),
		},
		{
			name: "Should fail when aws credentials are invalid",
			configuration: &AwsSsmEc2TargetConfiguration{
				Ec2InstanceId:     mockInstanceId,
				Ec2InstanceRegion: mockInstanceRegion,
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId: pointer.To("BAD CREDS"),
				},
			},
			expectedError: fmt.Errorf("invalid aws_credentials: aws_secret_access_key is required when aws_access_key_id is provided"),
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

func Test_ValidateAwsSsmEcsTargetConfiguration(t *testing.T) {
	t.Parallel()

	mockClusterName := "mock_cluster"
	mockClusterRegion := "us-east-1"
	mockServiceName := "mock_service"

	tests := []struct {
		name          string
		configuration *AwsSsmEcsTargetConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when all inputs are valid",
			configuration: &AwsSsmEcsTargetConfiguration{
				EcsClusterName:   mockClusterName,
				EcsClusterRegion: mockClusterRegion,
				EcsServiceName:   mockServiceName,
			},
			expectedError: nil,
		},
		{
			name: "Should fail when ecs cluster name is missing",
			configuration: &AwsSsmEcsTargetConfiguration{
				EcsClusterRegion: mockClusterRegion,
				EcsServiceName:   mockServiceName,
			},
			expectedError: errors.New("ecs_cluster_name is a required field"),
		},
		{
			name: "Should fail when ecs cluster region is missing",
			configuration: &AwsSsmEcsTargetConfiguration{
				EcsClusterName: mockClusterName,
				EcsServiceName: mockServiceName,
			},
			expectedError: errors.New("ecs_cluster_region is a required field"),
		},
		{
			name: "Should fail when ecs cluster region is invalid",
			configuration: &AwsSsmEcsTargetConfiguration{
				EcsClusterName:   mockClusterName,
				EcsClusterRegion: "bad region",
				EcsServiceName:   mockServiceName,
			},
			expectedError: fmt.Errorf("invalid ecs_cluster_region: region \"%s\" is not a valid aws region", "bad region"),
		},
		{
			name: "Should fail when ecs service name is missing",
			configuration: &AwsSsmEcsTargetConfiguration{
				EcsClusterName:   mockClusterName,
				EcsClusterRegion: mockClusterRegion,
			},
			expectedError: errors.New("ecs_service_name is a required field"),
		},
		{
			name: "Should fail when aws credentials are invalid",
			configuration: &AwsSsmEcsTargetConfiguration{
				EcsClusterName:   mockClusterName,
				EcsClusterRegion: mockClusterRegion,
				EcsServiceName:   mockServiceName,
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId: pointer.To("BAD CREDS"),
				},
			},
			expectedError: fmt.Errorf("invalid aws_credentials: aws_secret_access_key is required when aws_access_key_id is provided"),
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

func Test_ValidateUsernameAndPasswordAuthConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *UsernameAndPasswordAuthConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when username and password are valid",
			configuration: &UsernameAndPasswordAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Username:         "username",
				Password:         "password",
			},
			expectedError: nil,
		},
		{
			name: "Should fail when username type is 'defined' and username is missing",
			configuration: &UsernameAndPasswordAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Password:         "password",
			},
			expectedError: fmt.Errorf("username must be provided when username_provider is \"%s\"", UsernameProviderDefined),
		},
		{
			name: "Should fail when username type is 'prompt_client' and username is present",
			configuration: &UsernameAndPasswordAuthConfiguration{
				UsernameProvider: UsernameProviderPromptClient,
				Username:         "username",
				Password:         "password",
			},
			expectedError: fmt.Errorf("username must be empty when username_provider is %s", UsernameProviderPromptClient),
		},
		{
			name: "Should fail when username provider is invalid",
			configuration: &UsernameAndPasswordAuthConfiguration{
				UsernameProvider: "not valid",
				Username:         "username",
				Password:         "password",
			},
			expectedError: fmt.Errorf("username_provider %s is not valid", "not valid"),
		},
		{
			name: "Should fail when password is missing",
			configuration: &UsernameAndPasswordAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Username:         "username",
			},
			expectedError: errors.New("password is a required field"),
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

func Test_ValidatePrivateKeyAuthConfiguration(t *testing.T) {
	t.Parallel()

	mockPrivateKey := `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`

	tests := []struct {
		name          string
		configuration *PrivateKeyAuthConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when username and password are valid",
			configuration: &PrivateKeyAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Username:         "username",
				PrivateKey:       mockPrivateKey,
			},
			expectedError: nil,
		},
		{
			name: "Should fail when username type is 'defined' and username is missing",
			configuration: &PrivateKeyAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				PrivateKey:       mockPrivateKey,
			},
			expectedError: fmt.Errorf("username must be provided when username_provider is \"%s\"", UsernameProviderDefined),
		},
		{
			name: "Should fail when username type is 'prompt_client' and username is present",
			configuration: &PrivateKeyAuthConfiguration{
				UsernameProvider: UsernameProviderPromptClient,
				Username:         "username",
				PrivateKey:       mockPrivateKey,
			},
			expectedError: fmt.Errorf("username must be empty when username_provider is %s", UsernameProviderPromptClient),
		},
		{
			name: "Should fail when username provider is invalid",
			configuration: &PrivateKeyAuthConfiguration{
				UsernameProvider: "not valid",
				Username:         "username",
				PrivateKey:       mockPrivateKey,
			},
			expectedError: fmt.Errorf("username_provider %s is not valid", "not valid"),
		},
		{
			name: "Should fail when private key is missing",
			configuration: &PrivateKeyAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Username:         "username",
			},
			expectedError: errors.New("private_key is a required field"),
		},
		{
			name: "Should fail when private key is not a valid private key",
			configuration: &PrivateKeyAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Username:         "username",
				PrivateKey:       "bad priv",
			},
			expectedError: errors.New("private_key is not a valid PEM or DER encoded private key"),
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

func Test_ValidateBorder0CertificateAuthConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *Border0CertificateAuthConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when username is valid",
			configuration: &Border0CertificateAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
				Username:         "username",
			},
			expectedError: nil,
		},
		{
			name: "Should fail when username type is 'defined' and username is missing",
			configuration: &Border0CertificateAuthConfiguration{
				UsernameProvider: UsernameProviderDefined,
			},
			expectedError: fmt.Errorf("username must be provided when username_provider is \"%s\"", UsernameProviderDefined),
		},
		{
			name: "Should fail when username type is 'prompt_client' and username is present",
			configuration: &Border0CertificateAuthConfiguration{
				UsernameProvider: UsernameProviderPromptClient,
				Username:         "username",
			},
			expectedError: fmt.Errorf("username must be empty when username_provider is %s", UsernameProviderPromptClient),
		},
		{
			name: "Should fail when username provider is invalid",
			configuration: &Border0CertificateAuthConfiguration{
				UsernameProvider: "not valid",
				Username:         "username",
			},
			expectedError: fmt.Errorf("username_provider %s is not valid", "not valid"),
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
