package service

import (
	"fmt"
	"strings"

	"github.com/borderzero/border0-go/lib/types/set"
)

const (
	// ExitNodeModeIPv4Only is the exit node mode for
	// routing only IPv4 traffic through the exit node.
	ExitNodeModeIPv4Only = "IPV4_ONLY"

	// ExitNodeModeIPv6Only is the exit node mode for
	// routing only IPv6 traffic through the exit node.
	ExitNodeModeIPv6Only = "IPV6_ONLY"

	// ExitNodeModeDualStack is the exit node mode for
	// routing all traffic (both IPv4 and IPv6) through
	// the exit node
	ExitNodeModeDualStack = "DUAL_STACK"
)

var (
	// ExitNodeModes is a set containing valid values of
	// exit node modes. Note that an empty value is also
	// valid (defaulting to dual stack).
	ExitNodeModes = set.New(
		ExitNodeModeIPv4Only,
		ExitNodeModeIPv6Only,
		ExitNodeModeDualStack,
	)

	// exitNodeModeErrFmt is a format string used to return a validation error for an invalid exit node mode.
	exitNodeModeErrFmt = fmt.Sprintf(
		"exit node mode \"%%s\" is not valid, must be one of [ %s ]",
		strings.Join(ExitNodeModes.Slice(), ", "),
	)
)

// ExitNodeServiceConfiguration represents service
// configuration for exit node services (fka sockets).
type ExitNodeServiceConfiguration struct {
	Mode string `json:"mode,omitempty"`
}

// Validate validates the ExitNodeServiceConfiguration.
func (c *ExitNodeServiceConfiguration) Validate() error {
	if c.Mode != "" {
		if !ExitNodeModes.Has(c.Mode) {
			return fmt.Errorf(exitNodeModeErrFmt, c.Mode)
		}
	}
	return nil
}
