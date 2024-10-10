package service

import (
	"fmt"
	"net/netip"

	"github.com/borderzero/border0-go/lib/types/set"
)

// SubnetRouterServiceConfiguration represents service
// configuration for router services (fka sockets).
type SubnetRouterServiceConfiguration struct {
	IPv4CIDRRanges       []string `json:"ipv4_cidr_ranges"`
	IPv6CIDRRanges       []string `json:"ipv6_cidr_ranges"`
	SourceIPPreservation bool     `json:"source_ip_preservation"`
}

// Validate validates the SubnetRouterServiceConfiguration.
func (c *SubnetRouterServiceConfiguration) Validate() error {
	v4set := set.New[string]()
	for _, cidr := range c.IPv4CIDRRanges {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			return fmt.Errorf("the cidr \"%s\" is not a valid IPv4 CIDR", cidr)
		}
		if prefix.Addr().Is4In6() {
			return fmt.Errorf("the cidr \"%s\" is an IPv4-mapped IPv6 CIDR and must be passed as an IPv6 CIDR", cidr)
		}
		if !prefix.Addr().Is4() {
			return fmt.Errorf("the cidr \"%s\" is not a valid IPv4 CIDR", cidr)
		}
		if v4set.Has(cidr) {
			return fmt.Errorf("duplicate IPv4 CIDR \"%s\"", cidr)
		}
		v4set.Add(cidr)
	}
	v6set := set.New[string]()
	for _, cidr := range c.IPv6CIDRRanges {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			return fmt.Errorf("the cidr \"%s\" is not a valid IPv6 CIDR", cidr)
		}
		if !prefix.Addr().Is6() {
			return fmt.Errorf("the cidr \"%s\" is not a valid IPv6 CIDR", cidr)
		}
		if v6set.Has(cidr) {
			return fmt.Errorf("duplicate IPv6 CIDR \"%s\"", cidr)
		}
		v6set.Add(cidr)
	}
	return nil
}
