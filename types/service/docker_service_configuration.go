package service

import (
	"fmt"
)

// DockerServiceConfiguration represents service
// configuration for docker services (fka sockets).
type DockerServiceConfiguration struct {
	DockerHost string `json:"docker_host"`
}

// Validate validates the DockerServiceConfiguration.
func (c *DockerServiceConfiguration) Validate() error {
	if c.DockerHost == "" {
		return fmt.Errorf("docker_host is a required field")
	}
	return nil
}
