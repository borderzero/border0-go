package service

const (
	// SshServiceTypeStandard is the ssh
	// service type for standard ssh services.
	SshServiceTypeStandard = "standard"

	// SshServiceTypeAwsSsm is the ssh service
	// type for aws session manager ssh services.
	SshServiceTypeAwsSsm = "aws_ssm"

	// SshServiceTypeAwsEc2InstanceConnect is the ssh service
	// type for aws ec2 instance connect ssh services.
	SshServiceTypeAwsEc2InstanceConnect = "aws_ec2_instance_connect"

	// SshServiceTypeConnectorBuiltIn is the ssh service
	// type for the connector's built-in ssh service.
	SshServiceTypeConnectorBuiltIn = "connector_built_in_ssh_service"
)

const (
	// SsmTargetTypeEc2 is the ssm target type for ec2 targets.
	SsmTargetTypeEc2 = "ec2"

	// SsmTargetTypeEcs is the ssm target type for ecs targets.
	SsmTargetTypeEcs = "ecs"
)

const (
	// StandardSshServiceAuthenticationTypeUsernameAndPassword is the standard ssh
	// service authentication type for authenticating with a username and password.
	StandardSshServiceAuthenticationTypeUsernameAndPassword = "username_and_password"

	// StandardSshServiceAuthenticationTypePrivateKey is the standard ssh
	// service authentication type for authenticating with a private key.
	StandardSshServiceAuthenticationTypePrivateKey = "private_key"

	// StandardSshServiceAuthenticationTypeBorder0Certificate is the standard ssh
	// service authentication type for authenticating with a border0-signed certificate.
	StandardSshServiceAuthenticationTypeBorder0Certificate = "border0_certificate"
)

const (
	// UsernameProviderDefined is the username provider
	// option for using an admin-defined (static) username.
	UsernameProviderDefined = "defined"

	// UsernameProviderPromptClient is username provider option
	// for prompting connecting clients for the username.
	UsernameProviderPromptClient = "prompt_client"

	// UsernameProviderUseConnectorUser is username provider
	// option for using the connector's OS username.
	//
	// NOTE: This option can only be used as the username
	// provider for connector built-in ssh services.
	UsernameProviderUseConnectorUser = "use_connector_user"
)

// SshServiceConfiguration represents service
// configuration for shell services (fka sockets).
type SshServiceConfiguration struct {
	SshServiceType string `json:"ssh_service_type"`

	// mutually exclusive fields below
	StandardSshServiceConfiguration *StandardSshServiceConfiguration `json:"standard_ssh_service_configuration,omitempty"`
	AwsSsmSshServiceConfiguration   *AwsSsmSshServiceConfiguration   `json:"aws_ssm_ssh_service_configuration,omitempty"`
	AwsEc2ICSshServiceConfiguration *AwsEc2ICSshServiceConfiguration `json:"aws_ec2ic_ssh_service_configuration,omitempty"`
	BuiltInSshServiceConfiguration  *BuiltInSshServiceConfiguration  `json:"built_in_ssh_service_configuration,omitmepty"`
}

// StandardSshServiceConfiguration represents service
// configuration for standard ssh services (fka sockets).
type StandardSshServiceConfiguration struct {
	SshAuthenticationType string `json:"ssh_authentication_type"`

	HostnameAndPort // inherited

	// mutually exclusive fields below
	UsernameAndPasswordAuthConfiguration *UsernameAndPasswordAuthConfiguration `json:"username_and_password_auth_configuration,omitempty"`
	PrivateKeyAuthConfiguration          *PrivateKeyAuthConfiguration          `json:"private_key_auth_configuration,omitempty"`
	Border0CertificateAuthConfiguration  *Border0CertificateAuthConfiguration  `json:"border0_certificate_auth_configuration,omitempty"`
}

// AwsSsmSshServiceConfiguration represents service
// configuration for aws ssm ssh services (fka sockets).
type AwsSsmSshServiceConfiguration struct {
	SsmTargetType string `json:"ssm_target_type"`

	// mutually exclusive fields below
	AwsSsmEc2TargetConfiguration *AwsSsmEc2TargetConfiguration `json:"aws_ssm_ec2_target_configuration,omitempty"`
	AwsSsmEcsTargetConfiguration *AwsSsmEcsTargetConfiguration `json:"aws_ssm_ecs_target_configuration,omitempty"`
}

// AwsSsmEc2TargetConfiguration represents service configuration
// for aws ssm ssh services (fka sockets) that have EC2 instances
// as their ssm target.
type AwsSsmEc2TargetConfiguration struct {
	Ec2InstanceId     string          `json:"ec2_instance_id"`
	Ec2InstanceRegion string          `json:"ec2_instance_region"`
	AwsCredentials    *AwsCredentials `json:"aws_credentials,omitempty"`
}

// AwsSsmEcsTargetConfiguration represents service configuration
// for aws ssm ssh services (fka sockets) that have ECS services
// as their ssm target.
type AwsSsmEcsTargetConfiguration struct {
	EcsClusterRegion string          `json:"ecs_cluster_region"`
	EcsClusterName   string          `json:"ecs_cluster_name"`
	EcsServiceName   string          `json:"ecs_service_name"`
	AwsCredentials   *AwsCredentials `json:"aws_credentials,omitempty"`
}

// AwsEc2ICSshServiceConfiguration represents service configuration
// for aws ec2 instance connect ssh services (fka sockets).
type AwsEc2ICSshServiceConfiguration struct {
	HostnameAndPort                   // inherited
	Ec2InstanceId     string          `json:"ec2_instance_id"`
	Ec2InstanceRegion string          `json:"ec2_instance_region"`
	AwsCredentials    *AwsCredentials `json:"aws_credentials,omitempty"`
}

// BuiltInSshServiceConfiguration represents the service configuration
// for the connector built-in ssh services (fka sockets).
type BuiltInSshServiceConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
}

// UsernameAndPasswordAuthConfiguration represents authentication configuration
// for standard ssh services that require a username and password for authentication.
type UsernameAndPasswordAuthConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password"`
}

// PrivateKeyAuthConfiguration represents authentication configuration
// for standard ssh services that require a private key for authentication.
type PrivateKeyAuthConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
	PrivateKey       string `json:"private_key"`
}

// UsernameAndPasswordAuthConfiguration represents authentication configuration
// for standard ssh services that require a border0-signed certificate for authentication.
type Border0CertificateAuthConfiguration struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
}
