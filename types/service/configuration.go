package service

const (
	// ServiceTypeDatabase is the service type for database services (fka sockets).
	ServiceTypeDatabase = "database"

	// ServiceTypeHttp is the service type for http services (fka sockets).
	ServiceTypeHttp = "http"

	// ServiceTypeSsh is the service type for ssh services (fka sockets).
	ServiceTypeSsh = "ssh"

	// ServiceTypeTls is the service type for tls services (fka sockets).
	ServiceTypeTls = "tls"
)

// Configuration represents service configuration.
type Configuration struct {
	ServiceType string `json:"service_type"`

	DatabaseServiceConfiguration *DatabaseServiceConfiguration `json:"database_service_configuration,omitempty"`
	HttpServiceConfiguration     *HttpServiceConfiguration     `json:"http_service_configuration,omitempty"`
	SshServiceConfiguration      *SshServiceConfiguration      `json:"ssh_service_configuration,omitempty"`
	TlsServiceConfiguration      *TlsServiceConfiguration      `json:"tls_service_configuration,omitempty"`
}
