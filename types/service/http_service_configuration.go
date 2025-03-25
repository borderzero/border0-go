package service

import (
	"fmt"
	"regexp"

	"github.com/borderzero/border0-go/lib/types/nilcheck"
)

const (
	// HttpServiceTypeStandard is the http
	// service type for standard http services.
	HttpServiceTypeStandard = "standard"

	// HttpServiceTypeConnectorFileServer is the http service
	// type for the connector's built-in file webserver.
	HttpServiceTypeConnectorFileServer = "connector_file_server"
)

const (
	MaxHeaderKeyLength   = 255
	MaxHeaderValueLength = 4096
)

var validHeaderKey = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

// HttpServiceConfiguration represents service
// configuration for http services (fka sockets).
type HttpServiceConfiguration struct {
	HttpServiceType string `json:"http_service_type"`

	// mutually exclusive fields below
	StandardHttpServiceConfiguration   *StandardHttpServiceConfiguration   `json:"standard_http_service_configuration,omitempty"`
	FileServerHttpServiceConfiguration *FileServerHttpServiceConfiguration `json:"fileserver_http_service_configuration,omitempty"`
}

// Headers represents a list of key-value headers.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// StandardHttpServiceConfiguration represents service
// configuration for standard http services (fka sockets).
type StandardHttpServiceConfiguration struct {
	HostnameAndPort          // inherited
	Scheme          string   `json:"scheme"`
	HostHeader      string   `json:"host_header"`
	Headers         []Header `json:"headers,omitempty"`
}

// FileServerHttpServiceConfiguration represents service
// configuration for the connector built-in file webserver.
type FileServerHttpServiceConfiguration struct {
	TopLevelDirectory string `json:"top_level_directory,omitempty"`
}

// Validate validates the HttpServiceConfiguration.
func (c *HttpServiceConfiguration) Validate() error {
	switch c.HttpServiceType {

	case HttpServiceTypeStandard:
		if nilcheck.AnyNotNil(c.FileServerHttpServiceConfiguration) {
			return fmt.Errorf(
				"http service type \"%s\" can only have standard http service configuration defined",
				HttpServiceTypeStandard,
			)
		}
		if c.StandardHttpServiceConfiguration == nil {
			return fmt.Errorf(
				"http service configuration for http service type \"%s\" must have standard http service configuration defined",
				HttpServiceTypeStandard,
			)
		}
		if err := c.StandardHttpServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid standard http service configuration: %v", err)
		}
		return nil

	case HttpServiceTypeConnectorFileServer:
		if nilcheck.AnyNotNil(c.StandardHttpServiceConfiguration) {
			return fmt.Errorf(
				"http service type \"%s\" can only have file server http service configuration defined",
				HttpServiceTypeConnectorFileServer,
			)
		}
		if c.FileServerHttpServiceConfiguration == nil {
			return fmt.Errorf(
				"http service configuration for http service type \"%s\" must have file server http service configuration defined",
				HttpServiceTypeConnectorFileServer,
			)
		}
		if err := c.FileServerHttpServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid file server http service configuration: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("http service configuration has invalid http service type \"%s\"", c.HttpServiceType)
	}
}

// Validate validates the StandardHttpServiceConfiguration.
func (c *StandardHttpServiceConfiguration) Validate() error {
	if c.HostHeader == "" {
		return fmt.Errorf("host_header is a required field")
	}
	if err := c.HostnameAndPort.Validate(); err != nil {
		return err
	}

	// validate headers
	for _, header := range c.Headers {
		if len(header.Key) == 0 {
			return fmt.Errorf("headers key cannot be empty")
		}

		if len(header.Key) > MaxHeaderKeyLength {
			return fmt.Errorf("header key exceeds max length of %d", MaxHeaderKeyLength)
		}

		if !validHeaderKey.MatchString(header.Key) {
			return fmt.Errorf("header key %s is an invalid header key", header.Key)
		}

		if len(header.Value) > MaxHeaderValueLength {
			return fmt.Errorf("header value for key %s exceeds max length of %d", header.Key, MaxHeaderValueLength)
		}

		for _, ch := range header.Value {
			if ch < 32 || ch == 127 { // ASCII control characters not allowed
				return fmt.Errorf("header value for key %s with value %s contains invalid character: %v", header.Key, header.Value, ch)
			}
		}
	}

	return nil
}

// Validate validates the FileServerHttpServiceConfiguration.
func (c *FileServerHttpServiceConfiguration) Validate() error {
	if c.TopLevelDirectory == "" {
		return fmt.Errorf("top_level_directory is a required field")
	}
	return nil
}
