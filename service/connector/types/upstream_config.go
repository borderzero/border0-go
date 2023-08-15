package types

import "github.com/borderzero/border0-go/service"

const (
	// UpstreamConnectionTypeSSH represents a SSH type of upstream connection.
	UpstreamConnectionTypeSSH = "ssh"
	// UpstreamConnectionTypeAwsEC2Connection represents a AWS EC2 Connect type of upstream connection.
	UpstreamConnectionTypeAwsEC2Connection = "aws_ec2_connect"
	// UpstreamConnectionTypeAwsSSM represents a AWS SSM type of upstream connection.
	UpstreamConnectionTypeAwsSSM = "aws_ssm"
	// UpstreamConnectionTypeBuiltInSshServer represents the Border0 built-in SSH server.
	UpstreamConnectionTypeBuiltInSshServer = "built_in_ssh_server"
	// UpstreamConnectionTypeDatabase represents a database type of upstream connection.
	UpstreamConnectionTypeDatabase = "database"
)

// Upstream authentication type constants.
const (
	// UpstreamAuthenticationTypeUsernamePassword represents a username password type of upstream authentication.
	UpstreamAuthenticationTypeUsernamePassword = "username_password"
	// UpstreamAuthenticationTypeBorder0Cert represents a Border0 certificate type of upstream authentication.
	UpstreamAuthenticationTypeBorder0Cert = "border0_cert"
	// UpstreamAuthenticationTypeSSHPrivateKey represents a SSH private key type of upstream authentication.
	UpstreamAuthenticationTypeSSHPrivateKey = "ssh_private_key"
)

// Built-in ssh server username provider constants.
const (
	// UsernameProviderPromptClient specifies that the username will be prompted-for to clients.
	UsernameProviderPromptClient = "prompt_client"
	// UsernameProviderUseConnectorUser specifies that the username will be derived from the connector's user.
	UsernameProviderUseConnectorUser = "use_connector_user"
	// UsernameProviderDefined specifies that the username will be defined in the socket upstream data by admins.
	UsernameProviderDefined = "defined"
)

// ConnectorServiceUpstreamConfig represents a configuration of a connector service upstream.
type ConnectorServiceUpstreamConfig struct {
	// BaseUpstreamDetails contains basic details of the upstream connection.
	BaseUpstreamDetails
	// UpstreamConnectionType specifies the type of the upstream connection.
	UpstreamConnectionType string `json:"upstream_connection_type"`
	// SSHConfiguration is optional and represents a configuration for a SSH connection.
	SSHConfiguration *SSHConfiguration `json:"ssh_configuration,omitempty"`
	// DatabaseConfiguration is optional and represents a configuration for a database connection.
	DatabaseConfiguration *service.DatabaseServiceConfiguration `json:"database_service_configuration,omitempty"`
}

// BaseUpstreamDetails represents basic details of an upstream connection.
type BaseUpstreamDetails struct {
	// Hostname of the upstream connection.
	Hostname string `json:"hostname"`
	// Port of the upstream connection.
	Port int `json:"port"`
}

// BasicCredentials represents a basic username-password pair.
type BasicCredentials struct {
	// Username for the credentials.
	Username string `json:"username,omitempty"`
	// Password for the credentials.
	Password string `json:"password,omitempty"`
}

// AWSConfiguration represents an AWS configuration for upstream connection.
type AWSConfiguration struct {
	// Region of the AWS configuration.
	Region string `json:"region"`
	// InstanceID of the AWS configuration.
	InstanceID string `json:"instance_id"`
	// AwsCredentials are optional and represent AWS credentials for the configuration.
	AwsCredentials *AwsCredentials `json:"aws_credentials,omitempty"`
}

// AwsCredentials represents AWS SSM details
type AwsSSMDetails struct {
	AWSConfiguration
}

// AwsCredentials represents AWS EC2 Connect details
type AwsEC2ConnectDetails struct {
	AWSConfiguration
}

// SSHConfiguration represents a configuration for a SSH connection.
type SSHConfiguration struct {
	// UpstreamAuthenticationType specifies the type of authentication for the SSH connection.
	UpstreamAuthenticationType string `json:"upstream_authentication_type,omitempty"`
	// AwsSSMDetails are optional and represent AWS SSM details for SSH connection.
	AwsSSMDetails *AwsSSMDetails `json:"aws_ssm_details,omitempty"`
	// AwsEC2ConnectDetails are optional and represent AWS EC2 Connect details for SSH connection.
	AwsEC2ConnectDetails *AwsEC2ConnectDetails `json:"aws_ec2_connect_details,omitempty"`
	// SSHPrivateKeyDetails are optional and represent a private key details for SSH connection.
	SSHPrivateKeyDetails *SSHPrivateKeyDetails `json:"ssh_private_key_details,omitempty"`
	// Border0 cert is optional and represents a Border0 certificate for SSH connection.
	Border0CertificateDetails *Border0CertificateDetails `json:"border0_certificate_details,omitempty"`
	// BuiltInSshServerDetails is optional and represents details for a built-in ssh server ssh connection.
	BuiltInSshServerDetails *BuiltInSshServerDetails `json:"built_in_ssh_server_details,omitempty"`
	// BasicCredentials are optional and represent a username-password pair for SSH connection.
	BasicCredentials *BasicCredentials `json:"basic_credentials,omitempty"`
}

// Border0CertificateDetails represents details of a Border0 certificate.
type Border0CertificateDetails struct {
	Username string `json:"username,omitempty"`
}

// BuiltInSshServerDetails represents details to use for a built-in ssh server socket.
type BuiltInSshServerDetails struct {
	UsernameProvider string `json:"username_provider,omitempty"`
	Username         string `json:"username,omitempty"`
}

// SSHPrivateKeyDetails represents details of a SSH private key.
type SSHPrivateKeyDetails struct {
	// Key is the SSH private key.
	Key string `json:"key"`
	// Username is the SSH username
	Username string `json:"username,omitempty"`
}
