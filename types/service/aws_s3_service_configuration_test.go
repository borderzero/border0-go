package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/borderzero/border0-go/types/common"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateAwsS3ServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *AwsS3ServiceConfiguration
		expectedError error
	}{
		{
			name:          "Should succeed when config is valid (nil bucket allowlist)",
			configuration: &AwsS3ServiceConfiguration{},
			expectedError: nil,
		},
		{
			name:          "Should succeed when config is valid (empty bucket allowlist)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{}},
			expectedError: nil,
		},
		{
			name:          "Should succeed when config is valid (valid bucket allowlist just wildcard)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{"*"}},
			expectedError: nil,
		},
		{
			name:          "Should succeed when config is valid (valid bucket allowlist no wildcards)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{"my-bucket", "another-bucket", "test-bucket"}},
			expectedError: nil,
		},
		{
			name:          "Should succeed when config is valid (valid bucket allowlist with wildcards)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{"prod-*", "dev-*", "test-bucket"}},
			expectedError: nil,
		},
		{
			name:          "Should succeed when config is valid with AWS credentials",
			configuration: &AwsS3ServiceConfiguration{
				BucketAllowlist: []string{"my-bucket"},
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId:     pointer.To("AKIAIOSFODNN7EXAMPLE"),
					AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
				},
			},
			expectedError: nil,
		},
		{
			name:          "Should fail when config is not valid (empty string entries in bucket allowlist)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{"my-bucket", ""}},
			expectedError: errors.New("the bucket allowlist entry in index 1 is an empty string"),
		},
		{
			name:          "Should fail when config is not valid (repeated entries in bucket allowlist)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{"my-bucket", "my-bucket"}},
			expectedError: errors.New("the bucket allowlist entry in index 1 (\"my-bucket\") is repeated"),
		},
		{
			name:          "Should fail when config is not valid (bad pattern in bucket allowlist)",
			configuration: &AwsS3ServiceConfiguration{BucketAllowlist: []string{"my bucket"}}, // space is not valid
			expectedError: errors.New("the bucket allowlist entry in index 0 (\"my bucket\") has invalid characters"),
		},
		{
			name: "Should fail when AWS credentials are invalid",
			configuration: &AwsS3ServiceConfiguration{
				BucketAllowlist: []string{"my-bucket"},
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId: pointer.To("AKIAIOSFODNN7EXAMPLE"),
				},
			},
			expectedError: fmt.Errorf("invalid AWS credentials: %w", errors.New("aws_secret_access_key is required when aws_access_key_id is provided")),
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
