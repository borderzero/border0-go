package common

import (
	"errors"
	"fmt"

	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/borderzero/border0-go/lib/types/regex"
	"github.com/borderzero/border0-go/lib/types/set"
)

const (
	externalVarPattern                   = `^\$\{(from:).+\}$`
	awsAccessKeyIdPattern                = `^AKIA[0-9A-Z]{16}$`
	awsSecretAccessKeyPattern            = `^[A-Za-z0-9]{40}$`
	awsProfilePattern                    = `^[a-zA-Z0-9-_]+$`
	awsSessionTokenConservativeMaxLength = 2048
)

// AwsCredentials represents aws credentials.
type AwsCredentials struct {
	AwsAccessKeyId     *string `json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey *string `json:"aws_secret_access_key,omitempty"`
	AwsSessionToken    *string `json:"aws_session_token,omitempty"`
	AwsProfile         *string `json:"aws_profile,omitempty"`
}

// Validate validates the AwsCredentials.
func (c *AwsCredentials) Validate() error {
	awsAccessKeyId := pointer.ValueOrZero(c.AwsAccessKeyId)
	awsSecretAccessKey := pointer.ValueOrZero(c.AwsSecretAccessKey)
	awsSessionToken := pointer.ValueOrZero(c.AwsSessionToken)
	awsProfile := pointer.ValueOrZero(c.AwsProfile)

	hasAwsAccessKeyId := awsAccessKeyId != ""
	hasAwsSecretAccessKey := awsSecretAccessKey != ""
	hasAwsSessionToken := awsSessionToken != ""
	hasAwsProfile := awsProfile != ""

	if hasAwsAccessKeyId && !hasAwsSecretAccessKey {
		return errors.New("aws_secret_access_key is required when aws_access_key_id is provided")
	}
	if !hasAwsAccessKeyId && hasAwsSecretAccessKey {
		return errors.New("aws_access_key_id is required when aws_secret_access_key is provided")
	}
	if hasAwsAccessKeyId && hasAwsSecretAccessKey {
		if !regex.MatchAny(awsAccessKeyId, []string{externalVarPattern, awsAccessKeyIdPattern}...) {
			return errors.New("invalid aws_access_key_id")
		}
		if !regex.MatchAny(awsSecretAccessKey, []string{externalVarPattern, awsSecretAccessKey}...) {
			return errors.New("invalid aws_secret_access_key")
		}
	}

	if hasAwsSessionToken {
		if !hasAwsAccessKeyId || !hasAwsSecretAccessKey {
			return errors.New("both aws_access_key_id and aws_secret_access_key are required when aws_session_token is provided")
		}
		// note: aws session token has no documented pattern we can validate against... they are typically
		// several hundred characters long and there isn't a defined maximum length documented anywhere.
		// However, we can choose a conservative size to ensure we don't allow the field to be unbounded.
		if len(awsSessionToken) > awsSessionTokenConservativeMaxLength {
			return fmt.Errorf("aws_session_token too long (%d characters)", len(awsSessionToken))
		}
	}

	if hasAwsProfile {
		if !regex.MatchAny(awsProfile, []string{externalVarPattern, awsProfilePattern}...) {
			return errors.New("invalid aws_profile")
		}
	}

	return nil
}

// GetValidAwsRegions returns a set of valid aws regions.
func GetValidAwsRegions() set.Set[string] {
	return set.New(
		"af-south-1",     // Africa (Cape Town).
		"ap-east-1",      // Asia Pacific (Hong Kong).
		"ap-northeast-1", // Asia Pacific (Tokyo).
		"ap-northeast-2", // Asia Pacific (Seoul).
		"ap-northeast-3", // Asia Pacific (Osaka).
		"ap-south-1",     // Asia Pacific (Mumbai).
		"ap-south-2",     // Asia Pacific (Hyderabad).
		"ap-southeast-1", // Asia Pacific (Singapore).
		"ap-southeast-2", // Asia Pacific (Sydney).
		"ap-southeast-3", // Asia Pacific (Jakarta).
		"ap-southeast-4", // Asia Pacific (Melbourne).
		"ca-central-1",   // Canada (Central).
		"eu-central-1",   // Europe (Frankfurt).
		"eu-central-2",   // Europe (Zurich).
		"eu-north-1",     // Europe (Stockholm).
		"eu-south-1",     // Europe (Milan).
		"eu-south-2",     // Europe (Spain).
		"eu-west-1",      // Europe (Ireland).
		"eu-west-2",      // Europe (London).
		"eu-west-3",      // Europe (Paris).
		"il-central-1",   // Israel (Tel Aviv).
		"me-central-1",   // Middle East (UAE).
		"me-south-1",     // Middle East (Bahrain).
		"sa-east-1",      // South America (Sao Paulo).
		"us-east-1",      // US East (N. Virginia).
		"us-east-2",      // US East (Ohio).
		"us-west-1",      // US West (N. California).
		"us-west-2",      // US West (Oregon).
	)
}

// ValidateAwsRegions validates that a list of strings is
// valid aws regions in the default AWS partition.
func ValidateAwsRegions(regions ...string) error {
	validRegions := GetValidAwsRegions()

	for _, region := range regions {
		if !validRegions.Has(region) {
			return fmt.Errorf("region \"%s\" is not a valid aws region", region)
		}
	}

	return nil
}
