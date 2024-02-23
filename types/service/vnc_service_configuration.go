package service

import "fmt"

// VncServiceConfiguration represents service
// configuration for vnc services (fka sockets).
type VncServiceConfiguration struct {
	HostnameAndPort
}

// Validate validates the VncServiceConfiguration.
func (c *VncServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return fmt.Errorf("invalid hostname or port: %v", err)
	}
	return nil
}
