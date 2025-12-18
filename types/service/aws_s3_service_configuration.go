package service

import (
	"fmt"
	"regexp"

	"github.com/borderzero/border0-go/lib/types/set"
	"github.com/borderzero/border0-go/types/common"
)

const (
	// maxAwsS3BucketAllowlistEntries is the maximum number of entries allowed in the bucket allowlist
	maxAwsS3BucketAllowlistEntries = 250
)

// AwsS3ServiceConfiguration represents service
// configuration for aws s3 services (fka sockets).
type AwsS3ServiceConfiguration struct {
	AwsCredentials  *common.AwsCredentials `json:"aws_credentials,omitempty"`
	BucketAllowlist []string               `json:"bucket_allowlist,omitempty"`
}

// Validate validates the AwsS3ServiceConfiguration.
func (c *AwsS3ServiceConfiguration) Validate() error {
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid AWS credentials: %w", err)
		}
	}
	if len(c.BucketAllowlist) > 0 {
		if len(c.BucketAllowlist) > maxAwsS3BucketAllowlistEntries {
			return fmt.Errorf("the bucket allowlist cannot contain more than %d entries", maxAwsS3BucketAllowlistEntries)
		}
		// regex matches bucket names that start with alphanumeric or wildcard,
		// followed by any combination of alphanumeric, wildcard, dot, underscore, or hyphen
		regex := regexp.MustCompile(`^[a-zA-Z0-9*][a-zA-Z0-9*._\-]*$`)
		entries := set.New[string]()
		for i, name := range c.BucketAllowlist {
			// reject empty string
			if name == "" {
				return fmt.Errorf("the bucket allowlist entry in index %d is an empty string", i)
			}
			// make sure its valid
			if !regex.MatchString(name) {
				return fmt.Errorf("the bucket allowlist entry in index %d (\"%s\") has invalid characters", i, name)
			}
			// make sure its not repeated
			if entries.Has(name) {
				return fmt.Errorf("the bucket allowlist entry in index %d (\"%s\") is repeated", i, name)
			}
			entries.Add(name)
		}
	}
	return nil
}
