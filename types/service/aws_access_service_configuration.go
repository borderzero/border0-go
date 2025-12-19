package service

import (
	"fmt"
	"regexp"

	"github.com/borderzero/border0-go/lib/types/set"
	"github.com/borderzero/border0-go/types/common"
)

const (
	maxAwsAccessAccountsEntries = 500
	maxAwsRoleNameLen           = 100
	maxAwsRoleDescriptionLen    = 250
)

var (
	awsAccountIDRegex = regexp.MustCompile(`^\d{12}$`)

	// awsArnRegex validates ARN format: arn:partition:service:region:account-id:resource
	// We use a regex pattern instead of aws-sdk-go-v2/aws/arn to avoid adding a large
	// dependency for simple format validation. The regex covers the general ARN structure
	// which is sufficient for our validation purposes.
	awsArnRegex = regexp.MustCompile(`^arn:[\w-]+:[\w-]+:[\w-]*:\d*:.+$`)
)

// AwsAccessServiceConfiguration represents service
// configuration for aws access services (fka sockets).
type AwsAccessServiceConfiguration struct {
	AwsCredentials *common.AwsCredentials `json:"aws_credentials,omitempty"`
	AwsAccounts    []AwsAccount           `json:"aws_accounts"`
}

type AwsAccount struct {
	AwsAccountAlias string       `json:"aws_account_alias,omitempty"`
	AwsAccountID    string       `json:"aws_account_id"`
	AwsIamRoles     []AwsIamRole `json:"aws_iam_roles"`
}

type AwsIamRole struct {
	ARN         string `json:"arn"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validate validates the AwsAccessServiceConfiguration.
func (c *AwsAccessServiceConfiguration) Validate() error {
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid AWS credentials: %w", err)
		}
	}

	if len(c.AwsAccounts) > 0 {
		if len(c.AwsAccounts) > maxAwsAccessAccountsEntries {
			return fmt.Errorf("the aws accounts list cannot contain more than %d entries", maxAwsAccessAccountsEntries)
		}
		accountIDs := set.New[string]()
		for i, account := range c.AwsAccounts {
			// validate account id
			if account.AwsAccountID != "LOCAL" {
				if !awsAccountIDRegex.MatchString(account.AwsAccountID) {
					return fmt.Errorf("the aws account entry in index %d (\"%s\") has invalid characters", i, account.AwsAccountID)
				}
			}
			if accountIDs.Has(account.AwsAccountID) {
				return fmt.Errorf("the aws account entry in index %d (\"%s\") is repeated", i, account.AwsAccountID)
			}
			accountIDs.Add(account.AwsAccountID)

			// validate roles
			if len(account.AwsIamRoles) == 0 {
				return fmt.Errorf("the aws account entry in index %d (\"%s\") has no roles configured", i, account.AwsAccountID)
			} else {
				for j, role := range account.AwsIamRoles {
					if role.Name == "" {
						return fmt.Errorf("the aws iam role entry in index %d of account entry in index %d (\"%s\") has an empty name", j, i, account.AwsAccountID)
					}
					if len(role.Name) > maxAwsRoleNameLen {
						return fmt.Errorf("the aws iam role name for entry in index %d (\"%s\") of account entry in index %d (\"%s\") exceeds maximum length of %d characters", j, role.Name, i, account.AwsAccountID, maxAwsRoleNameLen)
					}
					if len(role.Description) > maxAwsRoleDescriptionLen {
						return fmt.Errorf("the aws iam role description for entry in index %d (\"%s\") of account entry in index %d (\"%s\") exceeds maximum length of %d characters", j, role.Name, i, account.AwsAccountID, maxAwsRoleDescriptionLen)
					}
					if !awsArnRegex.MatchString(role.ARN) {
						return fmt.Errorf("the aws iam role arn for entry in index %d (\"%s\") of account entry in index %d (\"%s\") is not a valid ARN", j, role.ARN, i, account.AwsAccountID)
					}
				}
			}
		}
	}
	return nil
}
