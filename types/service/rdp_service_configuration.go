package service

import "fmt"

// RdpServiceConfiguration represents service
// configuration for rdp services (fka sockets).
type RdpServiceConfiguration struct {
	HostnameAndPort
}

// Validate validates the RdpServiceConfiguration.
func (c *RdpServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return fmt.Errorf("invalid hostname or port: %v", err)
	}
	return nil
}
