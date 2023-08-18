package connector

// Metadata represents informational data about a connector.
type Metadata struct {
	AwsEc2IdentityMetadata *AwsEc2IdentityMetadata `json:"aws_ec2_identity_metadata,omitempty"`
}

// AwsEc2IdentityMetadata represents metadata for connectors running on AWS EC2 instances.
type AwsEc2IdentityMetadata struct {
	AwsAccountId        string `json:"aws_account_id,omitempty"`
	AwsRegion           string `json:"aws_region,omitempty"`
	AwsAvailabilityZone string `json:"aws_availability_zone,omitempty"`
	Ec2InstanceId       string `json:"ec2_instance_id,omitempty"`
	Ec2InstanceType     string `json:"ec2_instance_type,omitempty"`
	Ec2ImageId          string `json:"ec2_image_id,omitempty"`
	KernelId            string `json:"kernel_id,omitempty"`
	RamdiskId           string `json:"ramdisk_id,omitempty"`
	Architecture        string `json:"architecture,omitempty"`
	PrivateIpAddress    string `json:"private_ip_address,omitempty"`
}
