package service

import (
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"github.com/borderzero/border0-go/lib/types/set"
)

// SubnetRouterServiceConfiguration represents service
// configuration for subnet router services (fka sockets).
type SubnetRouterServiceConfiguration struct {
	IPv4CIDRRanges []string `json:"ipv4_cidr_ranges"`
	IPv6CIDRRanges []string `json:"ipv6_cidr_ranges"`
	DNSPatterns    []string `json:"dns_patterns"`
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
	dnsset := set.New[string]()
	for _, pattern := range c.DNSPatterns {
		if err := validateDNSPattern(pattern); err != nil {
			return fmt.Errorf("invalid DNS pattern \"%s\": %w", pattern, err)
		}
		normalized := strings.ToLower(pattern)
		if dnsset.Has(normalized) {
			return fmt.Errorf("duplicate DNS pattern \"%s\"", pattern)
		}
		dnsset.Add(normalized)
	}
	return nil
}

var (
	// dnsLabelRegex matches a valid DNS label (1-63 chars, alphanumeric or hyphen, cannot start/end with hyphen)
	dnsLabelRegex = `[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?`
	// dnsPatternRegex matches a valid DNS name or wildcard pattern
	// - Optional wildcard prefix: *.
	// - At least two labels separated by dots
	// - Total domain part (excluding wildcard): up to 253 chars
	dnsPatternRegex = regexp.MustCompile(`^(\*\.)?` + dnsLabelRegex + `(\.` + dnsLabelRegex + `)+$`)
)

// validateDNSPattern validates that a string is either a valid DNS name or a valid DNS wildcard pattern.
// Valid patterns include:
//   - Fully qualified domain names (e.g., "example.com", "subdomain.example.com")
//   - Wildcard patterns (e.g., "*.example.com", "*.subdomain.example.com")
func validateDNSPattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("DNS pattern cannot be empty")
	}

	// Check if pattern matches DNS name or wildcard format
	if !dnsPatternRegex.MatchString(pattern) {
		return fmt.Errorf("must be a valid DNS name or wildcard pattern (e.g., example.com or *.example.com)")
	}

	// Check total domain length (RFC 1035: max 253 characters)
	domainPart := strings.TrimPrefix(pattern, "*.")
	if len(domainPart) > 253 {
		return fmt.Errorf("domain name exceeds maximum length of 253 characters")
	}

	return nil
}
