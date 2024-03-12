package service

import (
	"fmt"

	"github.com/borderzero/border0-go/lib/types/netaddr"
)

// VpnServiceConfiguration represents service
// configuration for vpn services (fka sockets).
type VpnServiceConfiguration struct {
	DHCPPoolSubnet   string   `json:"dhcp_pool_subnet"`
	AdvertisedRoutes []string `json:"advertised_routes,omitempty"`
}

// Validate validates the VpnServiceConfiguration.
func (c *VpnServiceConfiguration) Validate() error {
	if c.DHCPPoolSubnet == "" {
		return fmt.Errorf("dhcp_pool_subnet is a required field")
	}
	if !netaddr.IsCIDRv4(c.DHCPPoolSubnet) {
		return fmt.Errorf("dhcp_pool_subnet \"%s\" is not valid IPv4 CIDR notation", c.DHCPPoolSubnet)
	}
	for i, route := range c.AdvertisedRoutes {
		if !netaddr.IsCIDRv4(route) {
			return fmt.Errorf("advertised_routes[%d] (\"%s\") is not valid IPv4 CIDR notation", i, route)
		}
	}
	return nil
}
