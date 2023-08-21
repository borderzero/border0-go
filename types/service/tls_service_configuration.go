package service

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

func (c *TlsServiceConfiguration) Validate() error {
	// TODO
	return nil
}
