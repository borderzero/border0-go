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

	// ServiceTypeVnc is the service type for vnc services (fka sockets).
	ServiceTypeVnc = "vnc"

	// ServiceTypeVpn is the service type for vpn services (fka sockets).
	ServiceTypeVpn = "vpn"

	// ServiceTypeRdp is the service type for rdp services (fka sockets).
	ServiceTypeRdp = "rdp"

	// ServiceTypeDocker is the service type for docker services (fka sockets).
	ServiceTypeDocker = "docker"
)

// Configuration represents upstream service configuration.
type Configuration struct {
	ServiceType string `json:"service_type"`

	DatabaseServiceConfiguration *DatabaseServiceConfiguration `json:"database_service_configuration,omitempty"`
	HttpServiceConfiguration     *HttpServiceConfiguration     `json:"http_service_configuration,omitempty"`
	SshServiceConfiguration      *SshServiceConfiguration      `json:"ssh_service_configuration,omitempty"`
	TlsServiceConfiguration      *TlsServiceConfiguration      `json:"tls_service_configuration,omitempty"`
	VncServiceConfiguration      *VncServiceConfiguration      `json:"vnc_service_configuration,omitempty"`
	VpnServiceConfiguration      *VpnServiceConfiguration      `json:"vpn_service_configuration,omitempty"`
	RdpServiceConfiguration      *RdpServiceConfiguration      `json:"rdp_service_configuration,omitempty"`
	DockerServiceConfiguration   *DockerServiceConfiguration   `json:"docker_service_configuration,omitempty"`
}

type validatable interface {
	Validate() error
}

// Validate validates the Configuration.
func (c *Configuration) Validate() error {
	all := map[string]validatable{
		ServiceTypeDatabase: c.DatabaseServiceConfiguration,
		ServiceTypeHttp:     c.HttpServiceConfiguration,
		ServiceTypeSsh:      c.SshServiceConfiguration,
		ServiceTypeTls:      c.TlsServiceConfiguration,
		ServiceTypeVnc:      c.VncServiceConfiguration,
		ServiceTypeVpn:      c.VpnServiceConfiguration,
		ServiceTypeRdp:      c.RdpServiceConfiguration,
		ServiceTypeDocker:   c.DockerServiceConfiguration,
	}

	if currentConfig, ok := all[c.ServiceType]; ok {
		otherConfigs := []any{}
		for serviceType, serviceTypeConfig := range all {
			if serviceType != c.ServiceType {
				otherConfigs = append(otherConfigs, serviceTypeConfig)
			}
		}
		if nilcheck.AnyNotNil(otherConfigs...) {
			return fmt.Errorf("service configuration for service type \"%s\" can only have %s service configuration defined", c.ServiceType, c.ServiceType)
		}
		if nilcheck.Nil(currentConfig) {
			return fmt.Errorf("service configuration for service type \"%s\" must have %s service configuration defined", c.ServiceType, c.ServiceType)
		}
		if err := currentConfig.Validate(); err != nil {
			return fmt.Errorf("invalid %s service configuration: %v", c.ServiceType, err)
		}
		return nil
	}

	return fmt.Errorf("service configuration has invalid service type \"%s\"", c.ServiceType)
}

// ConnectorServiceConfiguration includes both the connector socket and upstream service configuration
type ConnectorServiceConfiguration struct {
	ConnectorAuthenticationEnabled bool          `json:"connector_authentication_enabled"`
	EndToEndEncryptionEnabled      bool          `json:"end_to_end_encryption_enabled"`
	RecordingEnabled               bool          `json:"recording_enabled"`
	Upstream                       Configuration `json:"upstream_configuration"`
}

// Validate validates the ConnectorServiceConfiguration.
func (c *ConnectorServiceConfiguration) Validate() error {
	if err := c.Upstream.Validate(); err != nil {
		return fmt.Errorf("invalid upstream configuration: %w", err)
	}
	return nil
}
