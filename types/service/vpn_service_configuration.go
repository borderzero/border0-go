package service

import (
	"fmt"

	"github.com/borderzero/border0-go/lib/types/netaddr"
)

// VpnServiceConfiguration represents service
// configuration for vpn services (fka sockets).
type VpnServiceConfiguration struct {
	VpnSubnet string   `json:"vpn_subnet"`
	Routes    []string `json:"routes,omitempty"`
}

// Validate validates the VpnServiceConfiguration.
func (c *VpnServiceConfiguration) Validate() error {
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
