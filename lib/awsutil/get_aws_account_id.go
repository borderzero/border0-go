package awsutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// AwsAccountIdFromConfig returns the aws account id given an aws config.
// It makes a call to AWS Session Token Service's (STS) "GetCallerIdentity"
// API endpoint -- which does not require any IAM permissions to call.
func AwsAccountIdFromConfig(
	ctx context.Context,
	cfg aws.Config,
	timeout time.Duration,
) (string, error) {
	// new AWS STS service client
	stsService := sts.NewFromConfig(cfg)

	// new context object with a timeout
	getCallerIdentityContext, getCallerIdentityContextCancel := context.WithTimeout(ctx, timeout)
	defer getCallerIdentityContextCancel()

	// perform STS GetCallerIdentity API call to get caller identity details
	getCallerIdentityOutput, err := stsService.GetCallerIdentity(getCallerIdentityContext, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("failed to AWS account ID via the AWS STS API: %w", err)
	}

	// ensure response contains account id
	awsAccountId := aws.ToString(getCallerIdentityOutput.Account)
	if awsAccountId == "" {
		return "", fmt.Errorf("the AWS STS API returned an empty AWS account ID")
	}

	// success!
	return awsAccountId, nil
}
