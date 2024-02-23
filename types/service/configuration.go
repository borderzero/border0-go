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

	// ServiceTypeTcp is the service type for tcp services (fka sockets).
	ServiceTypeTcp = "tcp"

	// ServiceTypeTls is the service type for tls services (fka sockets).
	ServiceTypeTls = "tls"

	// ServiceTypeVnc is the service type for vnc services (fka sockets).
	ServiceTypeVnc = "vnc"

	// ServiceTypeRdp is the service type for rdp services (fka sockets).
	ServiceTypeRdp = "rdp"
)

// Configuration represents upstream service configuration.
type Configuration struct {
	ServiceType string `json:"service_type"`

	DatabaseServiceConfiguration *DatabaseServiceConfiguration `json:"database_service_configuration,omitempty"`
	HttpServiceConfiguration     *HttpServiceConfiguration     `json:"http_service_configuration,omitempty"`
	SshServiceConfiguration      *SshServiceConfiguration      `json:"ssh_service_configuration,omitempty"`
	TcpServiceConfiguration      *TcpServiceConfiguration      `json:"tcp_service_configuration,omitempty"`
	TlsServiceConfiguration      *TlsServiceConfiguration      `json:"tls_service_configuration,omitempty"`
	VncServiceConfiguration      *VncServiceConfiguration      `json:"vnc_service_configuration,omitempty"`
	RdpServiceConfiguration      *RdpServiceConfiguration      `json:"rdp_service_configuration,omitempty"`
}

// Validate validates the Configuration.
func (c *Configuration) Validate(allowExperimentalFeatures bool) error {
	switch c.ServiceType {

	case ServiceTypeDatabase:
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeDatabase)...) {
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
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeHttp)...) {
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
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeSsh)...) {
			return fmt.Errorf("service configuration for service type \"ssh\" can only have ssh service configuration defined")
		}
		if c.SshServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"ssh\" must have ssh service configuration defined")
		}
		if err := c.SshServiceConfiguration.Validate(allowExperimentalFeatures); err != nil {
			return fmt.Errorf("invalid ssh service configuration: %v", err)
		}
		return nil

	case ServiceTypeTcp:
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeTcp)...) {
			return fmt.Errorf("service configuration for service type \"tcp\" can only have tcp service configuration defined")
		}
		if c.TcpServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"tcp\" must have tcp service configuration defined")
		}
		if err := c.TcpServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid tcp service configuration: %v", err)
		}
		return nil

	case ServiceTypeVnc:
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeVnc)...) {
			return fmt.Errorf("service configuration for service type \"%s\" can only have %s service configuration defined", ServiceTypeVnc, ServiceTypeVnc)
		}
		if c.VncServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"%s\" must have %s service configuration defined", ServiceTypeVnc, ServiceTypeVnc)
		}
		if err := c.VncServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid %s service configuration: %v", ServiceTypeVnc, err)
		}
		return nil

	case ServiceTypeRdp:
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeRdp)...) {
			return fmt.Errorf("service configuration for service type \"%s\" can only have %s service configuration defined", ServiceTypeRdp, ServiceTypeRdp)
		}
		if c.RdpServiceConfiguration == nil {
			return fmt.Errorf("service configuration for service type \"%s\" must have %s service configuration defined", ServiceTypeRdp, ServiceTypeRdp)
		}
		if err := c.RdpServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid %s service configuration: %v", ServiceTypeRdp, err)
		}
		return nil

	// deprecated: now referred to as tcp sockets
	case ServiceTypeTls:
		if nilcheck.AnyNotNil(allConfigsExcept(c, ServiceTypeTls)...) {
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

// ConnectorServiceConfiguration includes both the connector socket and upstream service configuration
type ConnectorServiceConfiguration struct {
	ConnectorAuthenticationEnabled bool          `json:"connector_authentication_enabled"`
	EndToEndEncryptionEnabled      bool          `json:"end_to_end_encryption_enabled"`
	RecordingEnabled               bool          `json:"recording_enabled"`
	Upstream                       Configuration `json:"upstream_configuration"`
}

// Validate validates the ConnectorServiceConfiguration.
func (c *ConnectorServiceConfiguration) Validate(allowExperimentalFeatures bool) error {
	if err := c.Upstream.Validate(allowExperimentalFeatures); err != nil {
		return fmt.Errorf("invalid upstream configuration: %w", err)
	}
	return nil
}

func allConfigsExcept(c *Configuration, svcType string) []any {
	all := []any{}

	if svcType != ServiceTypeDatabase {
		all = append(all, c.DatabaseServiceConfiguration)
	}
	if svcType != ServiceTypeHttp {
		all = append(all, c.HttpServiceConfiguration)
	}
	if svcType != ServiceTypeSsh {
		all = append(all, c.SshServiceConfiguration)
	}
	if svcType != ServiceTypeTcp {
		all = append(all, c.TcpServiceConfiguration)
	}
	if svcType != ServiceTypeTls {
		all = append(all, c.TlsServiceConfiguration)
	}
	if svcType != ServiceTypeVnc {
		all = append(all, c.VncServiceConfiguration)
	}
	if svcType != ServiceTypeRdp {
		all = append(all, c.RdpServiceConfiguration)
	}

	return all
}
