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

	// SshServiceTypeBuiltIn is the ssh service
	// type for the connector's built-in ssh service.
	SshServiceTypeBuiltIn = "built_in_ssh_service"
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
	// BuiltInSshServiceUsernameProviderPromptClient is the built-in ssh service
	// username provider option for prompting clients for the username.
	BuiltInSshServiceUsernameProviderPromptClient = "prompt_client"

	// BuiltInSshServiceUsernameProviderUseConnectorUser is the built-in ssh
	// service username provider option for using the connector's OS username.
	BuiltInSshServiceUsernameProviderUseConnectorUser = "use_connector_user"

	// BuiltInSshServiceUsernameProviderDefined is the built-in ssh
	// service username provider option for using an admin-defined username.
	BuiltInSshServiceUsernameProviderDefined = "defined"
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
	SsmTarget      string          `json:"ssm_target"`
	AwsCredentials *AwsCredentials `json:"aws_credentials,omitempty"`
}

// AwsEc2ICSshServiceConfiguration represents service configuration
// for aws ec2 instance connect ssh services (fka sockets).
type AwsEc2ICSshServiceConfiguration struct {
	HostnameAndPort                 // inherited
	AwsCredentials  *AwsCredentials `json:"aws_credentials,omitempty"`
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
