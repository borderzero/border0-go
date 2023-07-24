package types

// ConnectorMetadata represents informational data about a connector.
type ConnectorMetadata struct {
	AwsEc2IdentityMetadata *AwsEc2IdentityMetadata `json:"aws_ec2_identity_metadata,omitempty"`
}

// AwsEc2IdentityMetadata represents metadata for connectors running on AWS EC2 instances.
type AwsEc2IdentityMetadata struct {
	AwsAccountId        string `json:"aws_account_id"`
	AwsRegion           string `json:"aws_region"`
	AwsAvailabilityZone string `json:"aws_availability_zone"`
	Ec2InstanceId       string `json:"ec2_instance_id"`
	Ec2InstanceType     string `json:"ec2_instance_type"`
	Ec2ImageId          string `json:"ec2_image_id"`
	KernelId            string `json:"kernel_id"`
	RamdiskId           string `json:"radisk_id"`
	Architecture        string `json:"architecture"`
	PrivateIpAddress    string `json:"private_ip_address"`
}
