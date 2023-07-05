package types

const (
	UpstreamConnectionTypeSSH              = "ssh"
	UpstreamConnectionTypeAwsEC2Connection = "aws_ec2_connect"
	UpstreamConnectionTypeAwsSSM           = "aws_ssm"
	UpstreamConnectionTypeDatabase         = "database"

	UpstreamAuthenticationTypeUsernamePassword = "username_password"
	UpstreamAuthenticationTypeBorder0Cert      = "border0_cert"
	UpstreamAuthenticationTypeSSHPrivateKey    = "ssh_private_key"
)

type ConnectorServiceUpstreamConfig struct {
	BaseUpstreamDetails
	UpstreamConnectionType string `json:"upstream_connection_type"`

	SSHConfiguration      *SSHConfiguration      `json:"ssh_configuration,omitempty"`
	DatabaseConfiguration *DatabaseConfiguration `json:"database_configuration,omitempty"`
}

type BaseUpstreamDetails struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type BasicCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AWSConfiguration struct {
	Region         string          `json:"region"`
	InstanceID     string          `json:"instance_id"`
	AwsProfile     string          `json:"aws_profile"`
	AwsCredentials *AwsCredentials `json:"aws_credentials,omitempty"`
}

type SSHConfiguration struct {
	UpstreamAuthenticationType string                `json:"upstream_authentication_type,omitempty"`
	SSHPrivateKeyDetails       *SSHPrivateKeyDetails `json:"ssh_private_key_details,omitempty"`
	BasicCredentials           *BasicCredentials     `json:"basic_credentials,omitempty"`
}

type SSHPrivateKeyDetails struct {
	Key string `json:"key"`
}

type DatabaseConfiguration struct {
	UpstreamAuthenticationType string `json:"upstream_authentication_type"`
}
