package service

// ExitNodeServiceConfiguration represents service
// configuration for exit node services (fka sockets).
type ExitNodeServiceConfiguration struct {
	// NOTE(@adrianosela): uncomment if needed later along with
	// the comments in exit_node_service_configuraiton_test.go
	//
	// DisableIPv4 bool `json:"disable_ipv4,omitempty"`
	// DisableIPv6 bool `json:"disable_ipv6,omitempty"`
}

// Validate validates the ExitNodeServiceConfiguration.
func (c *ExitNodeServiceConfiguration) Validate() error { return nil }
