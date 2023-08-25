package common

import (
	"errors"
	"testing"

	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateSshServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		creds         *AwsCredentials
		expectedError error
	}{
		{
			name:          "Should succeed when nothing is set",
			creds:         &AwsCredentials{},
			expectedError: nil,
		},
		{
			name: "Should succeed when aws profile is set",
			creds: &AwsCredentials{
				AwsProfile: pointer.To("profile"),
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when both aws access key id and aws secret access key are set",
			creds: &AwsCredentials{
				AwsAccessKeyId:     pointer.To("AKIAIOSFODNN7EXAMPLE"),
				AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when all of aws profile, aws access key id, and aws secret access key are set",
			creds: &AwsCredentials{
				AwsProfile:         pointer.To("profile"),
				AwsAccessKeyId:     pointer.To("AKIAIOSFODNN7EXAMPLE"),
				AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when all of aws profile, aws access key id, aws secret access key, and aws session token are set",
			creds: &AwsCredentials{
				AwsProfile:         pointer.To("profile"),
				AwsAccessKeyId:     pointer.To("AKIAIOSFODNN7EXAMPLE"),
				AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
				AwsSessionToken:    pointer.To("FQoGZXIvYXdzEFYaDm8MhGSN/fLqNq3KUyK3AzH/2TP8xTBGtvKbTAYpaC8TMSG1kTzO/m3E9dLdcpJPLHXqBgj3TjvkgDCIq1Tm3fLmF+pJEJaCWZ2RcDXLxIyPZxFHCG9z3h8IZZC7ODM6YzPm3DzK1ZfBSU/EXAMPLETOKEN=="),
			},
			expectedError: nil,
		},
		{
			name: "Should fail when aws access key id is invalid",
			creds: &AwsCredentials{
				AwsProfile:         pointer.To("profile"),
				AwsAccessKeyId:     pointer.To("badkey"),
				AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			},
			expectedError: errors.New("invalid aws_access_key_id"),
		},
		{
			name: "Should fail when aws secret access key is invalid",
			creds: &AwsCredentials{
				AwsProfile:         pointer.To("profile"),
				AwsAccessKeyId:     pointer.To("AKIAIOSFODNN7EXAMPLE"),
				AwsSecretAccessKey: pointer.To("bad"),
			},
			expectedError: errors.New("invalid aws_secret_access_key"),
		},
		{
			name: "Should fail when only aws access key id is set",
			creds: &AwsCredentials{
				AwsProfile:     pointer.To("profile"),
				AwsAccessKeyId: pointer.To("AKIAIOSFODNN7EXAMPLE"),
			},
			expectedError: errors.New("aws_secret_access_key is required when aws_access_key_id is provided"),
		},
		{
			name: "Should fail when only aws secret access key is set",
			creds: &AwsCredentials{
				AwsProfile:         pointer.To("profile"),
				AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
			},
			expectedError: errors.New("aws_access_key_id is required when aws_secret_access_key is provided"),
		},
		{
			name: "Should fail when only aws session token is set",
			creds: &AwsCredentials{
				AwsProfile:      pointer.To("profile"),
				AwsSessionToken: pointer.To("FQoGZXIvYXdzEFYaDm8MhGSN/fLqNq3KUyK3AzH/2TP8xTBGtvKbTAYpaC8TMSG1kTzO/m3E9dLdcpJPLHXqBgj3TjvkgDCIq1Tm3fLmF+pJEJaCWZ2RcDXLxIyPZxFHCG9z3h8IZZC7ODM6YzPm3DzK1ZfBSU/EXAMPLETOKEN=="),
			},
			expectedError: errors.New("both aws_access_key_id and aws_secret_access_key are required when aws_session_token is provided"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.creds.Validate())
		})
	}
}
