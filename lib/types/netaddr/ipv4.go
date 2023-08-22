package netaddr

import "regexp"

// IsCIDRv4 returns true if a given string is an IPv4 CIDR.
func IsCIDRv4(target string) bool {
	return regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{1,2}$`).MatchString(target)
}

// IsIPv4 returns true if a given string is an IPv4 adddress.
func IsIPv4(target string) bool {
	return regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`).MatchString(target)
}
