package service

import (
	"fmt"

	"github.com/borderzero/border0-go/lib/types/netaddr"
	"github.com/borderzero/border0-go/lib/types/null"
)

const (
	// TlsServiceTypeStandard is the tls
	// service type for standard tls services.
	TlsServiceTypeStandard = "standard"

	// TlsServiceTypeVpn is the tls service
	// type for the connector's built-in vpn.
	TlsServiceTypeVpn = "vpn"

	// TlsServiceTypeHttpProxy is the tls service type
	// for the connector's built-in http (forward) proxy.
	TlsServiceTypeHttpProxy = "http_proxy"
)

// TlsServiceConfiguration represents service
// configuration for tls services (fka sockets).
type TlsServiceConfiguration struct {
	TlsServiceType string `json:"tls_service_type,omitempty"`

	// mutually exclusive fields below
	StandardTlsServiceConfiguration  *StandardTlsServiceConfiguration  `json:"standard_tls_service_configuration,omitempty"`
	VpnTlsServiceConfiguration       *VpnTlsServiceConfiguration       `json:"vpn_tls_service_configuration,omitempty"`
	HttpProxyTlsServiceConfiguration *HttpProxyTlsServiceConfiguration `json:"http_proxy_tls_service_configuration,omitempty"`
}

// StandardTlsServiceConfiguration represents service
// configuration for standard tls services (fka sockets).
type StandardTlsServiceConfiguration struct {
	HostnameAndPort
}

// VpnTlsServiceConfiguration represents service
// configuration for vpn services services over the tls socket.
type VpnTlsServiceConfiguration struct {
	VpnSubnet string   `json:"vpn_subnet"`
	Routes    []string `json:"routes,omitempty"`
}

// HttpProxyTlsServiceConfiguration represents service
// configuration for http proxy services over the tls socket.
type HttpProxyTlsServiceConfiguration struct {
	HostAllowlist []string `json:"host_allowlist,omitempty"`
}

// Validate validates the TlsServiceConfiguration.
func (c *TlsServiceConfiguration) Validate() error {
	switch c.TlsServiceType {

	case TlsServiceTypeStandard:
		if !null.All(c.VpnTlsServiceConfiguration, c.HttpProxyTlsServiceConfiguration) {
			return fmt.Errorf(
				"tls service type \"%s\" can only have standard tls service configuration defined",
				TlsServiceTypeStandard)
		}
		if c.StandardTlsServiceConfiguration == nil {
			return fmt.Errorf(
				"tls service configuration for tls service type \"%s\" must have standard tls service configuration defined",
				TlsServiceTypeStandard,
			)
		}
		if err := c.StandardTlsServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid standard tls service configuration: %v", err)
		}
		return nil

	case TlsServiceTypeVpn:
		if !null.All(c.StandardTlsServiceConfiguration, c.HttpProxyTlsServiceConfiguration) {
			return fmt.Errorf(
				"tls service type \"%s\" can only have vpn tls service configuration defined",
				TlsServiceTypeVpn)
		}
		if c.VpnTlsServiceConfiguration == nil {
			return fmt.Errorf(
				"tls service configuration for tls service type \"%s\" must have vpn tls service configuration defined",
				TlsServiceTypeVpn,
			)
		}
		if err := c.VpnTlsServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid vpn tls service configuration: %v", err)
		}
		return nil

	case TlsServiceTypeHttpProxy:
		if !null.All(c.StandardTlsServiceConfiguration, c.VpnTlsServiceConfiguration) {
			return fmt.Errorf(
				"tls service type \"%s\" can only have http proxy tls service configuration defined",
				TlsServiceTypeHttpProxy)
		}
		if c.HttpProxyTlsServiceConfiguration == nil {
			return fmt.Errorf(
				"tls service configuration for tls service type \"%s\" must have http proxy tls service configuration defined",
				TlsServiceTypeHttpProxy,
			)
		}
		if err := c.HttpProxyTlsServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid http proxy tls service configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("tls service configuration has invalid tls service type \"%s\"", c.TlsServiceType)
	}
}

// Validate validates the StandardTlsServiceConfiguration.
func (c *StandardTlsServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates the VpnTlsServiceConfiguration.
func (c *VpnTlsServiceConfiguration) Validate() error {
	if c.VpnSubnet == "" {
		return fmt.Errorf("vpn_subnet is a required field")
	}
	if !netaddr.IsCIDRv4(c.VpnSubnet) {
		return fmt.Errorf("vpn_subnet \"%s\" is not valid IPv4 CIDR notation", c.VpnSubnet)
	}
	for i, route := range c.Routes {
		if !netaddr.IsCIDRv4(route) {
			return fmt.Errorf("routes[%d] (\"%s\") is not valid IPv4 CIDR notation", i, route)
		}
	}
	return nil
}

// Validate validates the HttpProxyTlsServiceConfiguration.
func (c *HttpProxyTlsServiceConfiguration) Validate() error {
	// nothing to validate
	return nil
}
