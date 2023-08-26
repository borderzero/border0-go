package service

import (
	"fmt"

	"github.com/borderzero/border0-go/lib/types/nilcheck"
)

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

// Validate validates the Configuration.
func (c *Configuration) Validate() error {
	switch c.ServiceType {

	case ServiceTypeDatabase:
		if nilcheck.AnyNotNil(c.HttpServiceConfiguration, c.SshServiceConfiguration, c.TlsServiceConfiguration) {
			return fmt.Errorf("service configuration for service type \"database\" can only have database service configuration defined")
		}
		if c.DatabaseServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"database\" must have database service configuration defined")
		}
		if err := c.DatabaseServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid database service configuration: %v", err)
		}
		return nil

	case ServiceTypeHttp:
		if nilcheck.AnyNotNil(c.DatabaseServiceConfiguration, c.SshServiceConfiguration, c.TlsServiceConfiguration) {
			return fmt.Errorf("service configuration for service type \"http\" can only have http service configuration defined")
		}
		if c.HttpServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"http\" must have http service configuration defined")
		}
		if err := c.HttpServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid http service configuration: %v", err)
		}
		return nil

	case ServiceTypeSsh:
		if nilcheck.AnyNotNil(c.HttpServiceConfiguration, c.DatabaseServiceConfiguration, c.TlsServiceConfiguration) {
			return fmt.Errorf("service configuration for service type \"ssh\" can only have ssh service configuration defined")
		}
		if c.SshServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"ssh\" must have ssh service configuration defined")
		}
		if err := c.SshServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid ssh service configuration: %v", err)
		}
		return nil

	case ServiceTypeTls:
		if nilcheck.AnyNotNil(c.HttpServiceConfiguration, c.DatabaseServiceConfiguration, c.SshServiceConfiguration) {
			return fmt.Errorf("service configuration for service type \"tls\" can only have tls service configuration defined")
		}
		if c.TlsServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"tls\" must have tls service configuration defined")
		}
		if err := c.TlsServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid tls service configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("service configuration has invalid service type \"%s\"", c.ServiceType)
	}
}
