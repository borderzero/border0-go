package service

// HostnameAndPort represents a host and port.
type HostnameAndPort struct {
	Hostname string `json:"hostname"`
	Port     uint16 `json:"port"`
}

// AwsCredentials represents aws credentials.
type AwsCredentials struct {
	AwsAccessKeyId     string `json:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
	AwsSessionToken    string `json:"aws_session_token,omitempty"`
	AwsProfile         string `json:"aws_profile,omitempty"`
	AwsRegion          string `json:"aws_region,omitempty"`
}

// UsernameAndPassword represents a username and password. Used for basic auth, for example, MySQL
// username and password in a database upstream configuration.
type UsernameAndPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TlsConfig represents a TLS configuration. Services can use this to configure TLS for securing
// connections to upstreams.
type TlsConfig struct {
	CaCertificate string `json:"ca_certificate"`
	Certificate   string `json:"certificate"`
	Key           string `json:"key"`
}
