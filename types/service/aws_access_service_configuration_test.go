package service

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/borderzero/border0-go/types/common"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateAwsAccessServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *AwsAccessServiceConfiguration
		expectedError error
	}{
		{
			name:          "Should succeed when config is valid (empty accounts)",
			configuration: &AwsAccessServiceConfiguration{AwsAccounts: []AwsAccount{}},
			expectedError: nil,
		},
		{
			name: "Should succeed when config is valid (single account with single role)",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when config is valid (multiple accounts with multiple roles)",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/Role1",
								Name:        "Role1",
								Description: "First role",
							},
							{
								ARN:         "arn:aws:iam::123456789012:role/Role2",
								Name:        "Role2",
								Description: "Second role",
							},
						},
					},
					{
						AwsAccountID: "987654321098",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::987654321098:role/Role3",
								Name:        "Role3",
								Description: "Third role",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when config is valid (LOCAL account ID)",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "LOCAL",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/LocalRole",
								Name:        "LocalRole",
								Description: "Local test role",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when config is valid (role with empty description)",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: "",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when config is valid with AWS credentials",
			configuration: &AwsAccessServiceConfiguration{
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId:     pointer.To("AKIAIOSFODNN7EXAMPLE"),
					AwsSecretAccessKey: pointer.To("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
				},
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when AWS credentials are invalid",
			configuration: &AwsAccessServiceConfiguration{
				AwsCredentials: &common.AwsCredentials{
					AwsAccessKeyId: pointer.To("AKIAIOSFODNN7EXAMPLE"),
				},
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: fmt.Errorf("invalid AWS credentials: %w", errors.New("aws_secret_access_key is required when aws_access_key_id is provided")),
		},
		{
			name: "Should fail when accounts list exceeds maximum entries",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: func() []AwsAccount {
					accounts := make([]AwsAccount, maxAwsAccessAccountsEntries+1)
					for i := range accounts {
						accounts[i] = AwsAccount{
							AwsAccountID: fmt.Sprintf("%012d", i),
							AwsIamRoles: []AwsIamRole{
								{
									ARN:         fmt.Sprintf("arn:aws:iam::%012d:role/Role", i),
									Name:        fmt.Sprintf("Role%d", i),
									Description: "Test role",
								},
							},
						}
					}
					return accounts
				}(),
			},
			expectedError: fmt.Errorf("the aws accounts list cannot contain more than %d entries", maxAwsAccessAccountsEntries),
		},
		{
			name: "Should fail when account ID has invalid characters",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789abc",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws account entry in index 0 (\"123456789abc\") has invalid characters"),
		},
		{
			name: "Should fail when account ID is too short",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "12345678901",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws account entry in index 0 (\"12345678901\") has invalid characters"),
		},
		{
			name: "Should fail when account ID is repeated",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/Role1",
								Name:        "Role1",
								Description: "First role",
							},
						},
					},
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/Role2",
								Name:        "Role2",
								Description: "Second role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws account entry in index 1 (\"123456789012\") is repeated"),
		},
		{
			name: "Should fail when account has no roles",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles:  []AwsIamRole{},
					},
				},
			},
			expectedError: errors.New("the aws account entry in index 0 (\"123456789012\") has no roles configured"),
		},
		{
			name: "Should fail when role name is empty",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws iam role entry in index 0 of account entry in index 0 (\"123456789012\") has an empty name"),
		},
		{
			name: "Should fail when role name exceeds maximum length",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/LongRoleName",
								Name:        strings.Repeat("a", maxAwsRoleNameLen+1),
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: fmt.Errorf("the aws iam role name for entry in index 0 (\"%s\") of account entry in index 0 (\"123456789012\") exceeds maximum length of %d characters", strings.Repeat("a", maxAwsRoleNameLen+1), maxAwsRoleNameLen),
		},
		{
			name: "Should fail when role description exceeds maximum length",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: strings.Repeat("a", maxAwsRoleDescriptionLen+1),
							},
						},
					},
				},
			},
			expectedError: fmt.Errorf("the aws iam role description for entry in index 0 (\"MyRole\") of account entry in index 0 (\"123456789012\") exceeds maximum length of %d characters", maxAwsRoleDescriptionLen),
		},
		{
			name: "Should succeed when role name is at maximum length",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/LongRoleName",
								Name:        strings.Repeat("a", maxAwsRoleNameLen),
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should succeed when role description is at maximum length",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/MyRole",
								Name:        "MyRole",
								Description: strings.Repeat("a", maxAwsRoleDescriptionLen),
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when role ARN is invalid",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "not-a-valid-arn",
								Name:        "MyRole",
								Description: "Test role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws iam role arn for entry in index 0 (\"not-a-valid-arn\") of account entry in index 0 (\"123456789012\") is not a valid ARN"),
		},
		{
			name: "Should fail when second role in first account has empty name",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/Role1",
								Name:        "Role1",
								Description: "First role",
							},
							{
								ARN:         "arn:aws:iam::123456789012:role/Role2",
								Name:        "",
								Description: "Second role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws iam role entry in index 1 of account entry in index 0 (\"123456789012\") has an empty name"),
		},
		{
			name: "Should fail when role in second account has invalid ARN",
			configuration: &AwsAccessServiceConfiguration{
				AwsAccounts: []AwsAccount{
					{
						AwsAccountID: "123456789012",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "arn:aws:iam::123456789012:role/Role1",
								Name:        "Role1",
								Description: "First role",
							},
						},
					},
					{
						AwsAccountID: "987654321098",
						AwsIamRoles: []AwsIamRole{
							{
								ARN:         "invalid-arn",
								Name:        "Role2",
								Description: "Second role",
							},
						},
					},
				},
			},
			expectedError: errors.New("the aws iam role arn for entry in index 0 (\"invalid-arn\") of account entry in index 1 (\"987654321098\") is not a valid ARN"),
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
