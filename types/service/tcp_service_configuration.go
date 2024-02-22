package service

import "fmt"

// TcpServiceConfiguration represents service
// configuration for tcp services (fka sockets).
type TcpServiceConfiguration struct {
	HostnameAndPort
}

// Validate validates the TcpServiceConfiguration.
func (c *TcpServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return fmt.Errorf("invalid hostname or port: %v", err)
	}
	return nil
}
