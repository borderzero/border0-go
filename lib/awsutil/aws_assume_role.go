package awsutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/borderzero/border0-go/lib/types/pointer"
)

// AssumedRoleFn represents a function to be performed with an assumed role.
type AssumedRoleFn func(ctx context.Context, cfg aws.Config) error

// WithAssumedRole assumes an AWS IAM role and performs a function as the assumed role.
func WithAssumedRole(
	ctx context.Context,
	cfg aws.Config,
	ari *sts.AssumeRoleInput,
	assumeRoleTimeout time.Duration,
	fn AssumedRoleFn,
) error {
	// new AWS STS service client
	stsService := sts.NewFromConfig(cfg)

	// new context object with a timeout
	assumeRoleContext, assumeRoleContextCancel := context.WithTimeout(ctx, assumeRoleTimeout)
	defer assumeRoleContextCancel()

	// perform STS AssumeRole API call to get temporary credentials for AWS IAM role
	assumeRoleOutput, err := stsService.AssumeRole(assumeRoleContext, ari)
	if err != nil {
		return fmt.Errorf("failed to assume role \"%s\": %s", pointer.ValueOrZero(ari.RoleArn), err)
	}

	// initialize new AWS configuration using temporary credentials from assumed role
	assumedRoleCfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			pointer.ValueOrZero(assumeRoleOutput.Credentials.AccessKeyId),
			pointer.ValueOrZero(assumeRoleOutput.Credentials.SecretAccessKey),
			pointer.ValueOrZero(assumeRoleOutput.Credentials.SessionToken),
		),
	}

	// run the passed function with the assumed role session
	return fn(ctx, assumedRoleCfg)
}
