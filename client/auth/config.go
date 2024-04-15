package auth

import (
	"fmt"

	"github.com/borderzero/border0-go/lib/osutil"
)

// Config represents authentication configuration.
type Config struct {
	tokenStorageFilePath string // file to write token to
	tokenWritingEnabled  bool   // whether the border0 client should attempt to write tokens to the token storage file path
	browserEnabled       bool   // whether the border0 client should attempt opening the default browser for completing device authorization flow.

	legacyAuth bool   // whether to use programmatic authentication
	email      string // DEPRECATED: programmatic authentication email
	password   string // DEPRECATED: programmatic authentication password
}

// GetTokenStorageFilePath is the getter for the token storage file path.
func (c *Config) GetTokenStorageFilePath() string { return c.tokenStorageFilePath }

// ShouldWriteTokensToDisk is the getter for the token writing disabled boolean.
func (c *Config) ShouldWriteTokensToDisk() bool { return c.tokenWritingEnabled }

// ShouldTryOpeningBrowser is the getter for the disabled browser boolean.
func (c *Config) ShouldTryOpeningBrowser() bool { return c.browserEnabled }

// ShouldUseLegacyAuthentication is the getter for the legacy (programmatic) authentication boolean.
func (c *Config) ShouldUseLegacyAuthentication() bool { return c.legacyAuth }

// GetEmail is the getter for the email.
func (c *Config) GetEmail() string { return c.email }

// GetEmail is the getter for the password.
func (c *Config) GetPassword() string { return c.password }

// Option represents an authentication option
type Option func(*Config)

// GetConfig returns a configuration object populated with the given options.
func GetConfig(opts ...Option) (*Config, error) {
	c := &Config{
		tokenStorageFilePath: "", // will only try populating if not set after applying opts
		tokenWritingEnabled:  true,
		browserEnabled:       true,
		legacyAuth:           false,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.tokenStorageFilePath == "" {
		hd, err := osutil.GetUserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("no token storage filepath provided and failed to get user home directory: %v", err)
		}
		c.tokenStorageFilePath = fmt.Sprintf("%s/.border0/token", hd)
	}

	return c, nil
}

// WithTokenStorageFilePath is the authentication option to set the token storage file path.
func WithTokenStorageFilePath(filePath string) Option {
	return func(c *Config) { c.tokenStorageFilePath = filePath }
}

// WithTokenWritingDisabled is the authentication option to toggle whether tokens
// acquired through the device authorization flow should be written to disk.
func WithTokenWriting(enabled bool) Option {
	return func(c *Config) { c.tokenWritingEnabled = enabled }
}

// WithOpenBrowser is the authentication option to toggle whether the
// browser should be opened for the device authorization flow URL.
func WithOpenBrowser(enabled bool) Option {
	return func(c *Config) { c.browserEnabled = enabled }
}

// WithLegacyCredentials is the authentication option to set legacy credentials.
func WithLegacyCredentials(email, password string) Option {
	return func(c *Config) { c.legacyAuth = true; c.email = email; c.password = password }
}
