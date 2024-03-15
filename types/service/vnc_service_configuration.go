package service

import "fmt"

const (
	maxVncPasswordLength = 8
)

// VncServiceConfiguration represents service
// configuration for vnc services (fka sockets).
type VncServiceConfiguration struct {
	HostnameAndPort
	Password string `json:"password"`
}

// Validate validates the VncServiceConfiguration.
func (c *VncServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return fmt.Errorf("invalid hostname or port: %v", err)
	}
	if len(c.Password) > 8 {
		return fmt.Errorf("password can be at most %d characters long", maxVncPasswordLength)
	}
	return nil
}
