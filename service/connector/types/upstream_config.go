package types

const (
	// UpstreamConnectionTypeSSH represents a SSH type of upstream connection.
	UpstreamConnectionTypeSSH = "ssh"
	// UpstreamConnectionTypeAwsEC2Connection represents a AWS EC2 Connect type of upstream connection.
	UpstreamConnectionTypeAwsEC2Connection = "aws_ec2_connect"
	// UpstreamConnectionTypeAwsSSM represents a AWS SSM type of upstream connection.
	UpstreamConnectionTypeAwsSSM = "aws_ssm"
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

// ConnectorServiceUpstreamConfig represents a configuration of a connector service upstream.
type ConnectorServiceUpstreamConfig struct {
	// BaseUpstreamDetails contains basic details of the upstream connection.
	BaseUpstreamDetails
	// UpstreamConnectionType specifies the type of the upstream connection.
	UpstreamConnectionType string `json:"upstream_connection_type"`
	// SSHConfiguration is optional and represents a configuration for a SSH connection.
	SSHConfiguration *SSHConfiguration `json:"ssh_configuration,omitempty"`
	// DatabaseConfiguration is optional and represents a configuration for a database connection.
	DatabaseConfiguration *DatabaseConfiguration `json:"database_configuration,omitempty"`
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
	// AwsProfile of the AWS configuration.
	AwsProfile string `json:"aws_profile"`
	// AwsCredentials are optional and represent AWS credentials for the configuration.
	AwsCredentials *AwsCredentials `json:"aws_credentials,omitempty"`
}

// SSHConfiguration represents a configuration for a SSH connection.
type SSHConfiguration struct {
	// UpstreamAuthenticationType specifies the type of authentication for the SSH connection.
	UpstreamAuthenticationType string `json:"upstream_authentication_type,omitempty"`
	// SSHPrivateKeyDetails are optional and represent a private key details for SSH connection.
	SSHPrivateKeyDetails *SSHPrivateKeyDetails `json:"ssh_private_key_details,omitempty"`
	// BasicCredentials are optional and represent a username-password pair for SSH connection.
	BasicCredentials *BasicCredentials `json:"basic_credentials,omitempty"`
}

// SSHPrivateKeyDetails represents details of a SSH private key.
type SSHPrivateKeyDetails struct {
	// Key is the SSH private key.
	Key string `json:"key"`
}

// DatabaseConfiguration represents a configuration for a database connection.
type DatabaseConfiguration struct {
	// UpstreamAuthenticationType specifies the type of authentication for the database connection.
	UpstreamAuthenticationType string `json:"upstream_authentication_type"`
}
