package service

import (
	"fmt"

	"github.com/borderzero/border0-go/lib/types/set"
)

// HostnameAndPort represents a host and port.
type HostnameAndPort struct {
	Hostname string `json:"hostname"`
	Port     uint16 `json:"port"`
}

// UsernameAndPassword represents a username and password. Used for basic auth, for example, MySQL
// username and password in a database upstream configuration.
type UsernameAndPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TlsConfig represents a TLS configuration. Services can use this to configure TLS for securing
// connections to upstreams.
type TlsConfig struct {
	CaCertificate string `json:"ca_certificate"`
	Certificate   string `json:"certificate"`
	Key           string `json:"key"`
}

// validateUsernameWithProvider validates a username_provider, username pair.
func validateUsernameWithProvider(
	usernameProvider string,
	username string,
	allowedEmptyUserProviders set.Set[string],
) error {
	if usernameProvider == UsernameProviderDefined || usernameProvider == "" {
		if username == "" {
			return fmt.Errorf("username must be provided when username_provider is \"%s\"", UsernameProviderDefined)
		}
		return nil
	}
	if !allowedEmptyUserProviders.Has(usernameProvider) {
		return fmt.Errorf("username_provider %s is not valid", usernameProvider)
	}
	if username != "" {
		return fmt.Errorf("username must be empty when username_provider is %s", usernameProvider)
	}
	return nil
}
