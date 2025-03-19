package service

import (
	"fmt"
	"net/netip"

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

	// ServiceTypeKubernetes is the service type for kubernetes services (fka sockets).
	ServiceTypeKubernetes = "kubernetes"

	// DEPRECATED_ServiceTypeSubnetRoutes is the service type for subnet routes services (fka sockets).
	// This has been deprecated in favour of "subnet_router".
	//
	// Remove after february 2025
	DEPRECATED_ServiceTypeSubnetRoutes = "subnet_routes"

	// ServiceTypeSubnetRouter is the service type for subnet router services (fka sockets).
	ServiceTypeSubnetRouter = "subnet_router"

	// ServiceTypeExitNode is the service type for exit node services (fka sockets).
	ServiceTypeExitNode = "exit_node"

	// ServiceTypeElasticsearch is the service type for Elasticsearch services (fka sockets).
	ServiceTypeElasticsearch = "elasticsearch"

	// ServiceTypeSnowflake is the service type for Snowflake services (fka sockets).
	ServiceTypeSnowflake = "snowflake"
)

// Configuration represents upstream service configuration.
type Configuration struct {
	ServiceType string `json:"service_type"`

	DatabaseServiceConfiguration      *DatabaseServiceConfiguration      `json:"database_service_configuration,omitempty"`
	HttpServiceConfiguration          *HttpServiceConfiguration          `json:"http_service_configuration,omitempty"`
	SshServiceConfiguration           *SshServiceConfiguration           `json:"ssh_service_configuration,omitempty"`
	TlsServiceConfiguration           *TlsServiceConfiguration           `json:"tls_service_configuration,omitempty"`
	VncServiceConfiguration           *VncServiceConfiguration           `json:"vnc_service_configuration,omitempty"`
	VpnServiceConfiguration           *VpnServiceConfiguration           `json:"vpn_service_configuration,omitempty"`
	RdpServiceConfiguration           *RdpServiceConfiguration           `json:"rdp_service_configuration,omitempty"`
	KubernetesServiceConfiguration    *KubernetesServiceConfiguration    `json:"kubernetes_service_configuration,omitempty"`
	SubnetRouterServiceConfiguration  *SubnetRouterServiceConfiguration  `json:"subnet_router_service_configuration,omitempty"`
	ExitNodeServiceConfiguration      *ExitNodeServiceConfiguration      `json:"exit_node_service_configuration,omitempty"`
	ElasticsearchServiceConfiguration *ElasticsearchServiceConfiguration `json:"elasticsearch_service_configuration,omitempty"`
	SnowflakeServiceConfiguration     *SnowflakeServiceConfiguration     `json:"snowflake_service_configuration,omitempty"`

	// remove after february 2025
	DEPRECATED_SubnetRoutesServiceConfiguration *SubnetRouterServiceConfiguration `json:"subnet_routes_service_configuration,omitempty"`
}

type validatable interface {
	Validate() error
}

// Validate validates the Configuration.
func (c *Configuration) Validate() error {
	all := map[string]validatable{
		ServiceTypeDatabase:      c.DatabaseServiceConfiguration,
		ServiceTypeHttp:          c.HttpServiceConfiguration,
		ServiceTypeSsh:           c.SshServiceConfiguration,
		ServiceTypeTls:           c.TlsServiceConfiguration,
		ServiceTypeVnc:           c.VncServiceConfiguration,
		ServiceTypeVpn:           c.VpnServiceConfiguration,
		ServiceTypeRdp:           c.RdpServiceConfiguration,
		ServiceTypeKubernetes:    c.KubernetesServiceConfiguration,
		ServiceTypeSubnetRouter:  c.SubnetRouterServiceConfiguration,
		ServiceTypeExitNode:      c.ExitNodeServiceConfiguration,
		ServiceTypeElasticsearch: c.ElasticsearchServiceConfiguration,
		ServiceTypeSnowflake:     c.SnowflakeServiceConfiguration,

		// remove after february 2025
		DEPRECATED_ServiceTypeSubnetRoutes: c.DEPRECATED_SubnetRoutesServiceConfiguration,
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
	PrivateNetworkIPv4             *string       `json:"private_network_ipv4"`
	PrivateNetworkIPv6             *string       `json:"private_network_ipv6"`
}

// Validate validates the ConnectorServiceConfiguration.
func (c *ConnectorServiceConfiguration) Validate() error {
	if err := c.Upstream.Validate(); err != nil {
		return fmt.Errorf("invalid upstream configuration: %w", err)
	}

	if c.PrivateNetworkIPv4 != nil {
		if err := validateIP(c.PrivateNetworkIPv4, "ipv4"); err != nil {
			return err
		}
	}

	if c.PrivateNetworkIPv6 != nil {
		if err := validateIP(c.PrivateNetworkIPv6, "ipv6"); err != nil {
			return err
		}
	}

	return nil
}

func validateIP(ipStr *string, version string) error {
	if ipStr == nil {
		return nil
	}
	addr, err := netip.ParseAddr(*ipStr)
	if err != nil {
		return fmt.Errorf("invalid IP address: %s", *ipStr)
	}
	switch version {
	case "ipv4":
		if !addr.Is4() {
			return fmt.Errorf("expected an IPv4 address: %s", *ipStr)
		}
	case "ipv6":
		if !addr.Is6() {
			return fmt.Errorf("expected an IPv6 address: %s", *ipStr)
		}
	default:
		return fmt.Errorf("unknown IP version: %s", version)
	}
	return nil
}
