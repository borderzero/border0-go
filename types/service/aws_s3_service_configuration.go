package service

import (
	"fmt"

	"github.com/borderzero/border0-go/types/common"
)

// AwsS3ServiceConfiguration represents service
// configuration for aws s3 services (fka sockets).
type AwsS3ServiceConfiguration struct {
	AwsCredentials *common.AwsCredentials `json:"aws_credentials,omitempty"`
}

// Validate validates the AwsS3ServiceConfiguration.
func (c *AwsS3ServiceConfiguration) Validate() error {
	if c.AwsCredentials != nil {
		if err := c.AwsCredentials.Validate(); err != nil {
			return fmt.Errorf("invalid AWS credentials: %w", err)
		}
	}
	return nil
}
