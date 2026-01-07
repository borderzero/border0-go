package service

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	// maxRdpUsernameLength is the maximum length for RDP usernames.
	// Set to 256 to accommodate UPN format (user@domain.com).
	maxRdpUsernameLength = 256

	// maxRdpPasswordLength is the maximum length for RDP passwords.
	// Windows supports passwords up to 256 characters.
	maxRdpPasswordLength = 256

	// maxRdpDomainLength is the maximum length for RDP domain names.
	// Set to 253 to accommodate FQDN format.
	maxRdpDomainLength = 253
)

// RdpServiceConfiguration represents service
// configuration for rdp services (fka sockets).
type RdpServiceConfiguration struct {
	HostnameAndPort
	UsernameAndPassword
	Domain string `json:"domain,omitempty"`
}

// Validate validates the RdpServiceConfiguration.
func (c *RdpServiceConfiguration) Validate() error {
	if err := c.HostnameAndPort.Validate(); err != nil {
		return fmt.Errorf("invalid hostname or port: %v", err)
	}
	if err := c.validateCredentials(); err != nil {
		return err
	}
	return nil
}

// validateCredentials validates the username, password, and domain fields.
func (c *RdpServiceConfiguration) validateCredentials() error {
	// Check mutual dependency: if one of username/password is provided, both must be provided
	if c.Username != "" && c.Password == "" {
		return fmt.Errorf("password is required when username is provided")
	}
	if c.Password != "" && c.Username == "" {
		return fmt.Errorf("username is required when password is provided")
	}

	// Check domain dependency: if domain is provided, username must be provided
	if c.Domain != "" && c.Username == "" {
		return fmt.Errorf("username is required when domain is provided")
	}

	// Validate username if provided
	if c.Username != "" {
		if len(c.Username) > maxRdpUsernameLength {
			return fmt.Errorf("username exceeds maximum length of %d characters", maxRdpUsernameLength)
		}
		if containsControlCharacters(c.Username) {
			return fmt.Errorf("username contains invalid control characters")
		}
	}

	// Validate password if provided
	if c.Password != "" {
		if len(c.Password) > maxRdpPasswordLength {
			return fmt.Errorf("password exceeds maximum length of %d characters", maxRdpPasswordLength)
		}
		if containsControlCharacters(c.Password) {
			return fmt.Errorf("password contains invalid control characters")
		}
	}

	// Validate domain if provided
	if c.Domain != "" {
		if len(c.Domain) > maxRdpDomainLength {
			return fmt.Errorf("domain exceeds maximum length of %d characters", maxRdpDomainLength)
		}
		if containsControlCharacters(c.Domain) {
			return fmt.Errorf("domain contains invalid control characters")
		}
	}

	return nil
}

// containsControlCharacters checks if a string contains any control characters.
func containsControlCharacters(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) && !strings.ContainsRune("\t\n\r", r) {
			return true
		}
	}
	return false
}
