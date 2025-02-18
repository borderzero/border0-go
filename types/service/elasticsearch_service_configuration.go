package service

import (
	"errors"
	"fmt"
)

const (
	// ElasticsearchServiceTypeStandard is the elasticsearch
	// service type for standard elasticsearch services.
	ElasticsearchServiceTypeStandard = "standard"

	// ElasticsearchAuthenticationTypeBasic is the basic authentication type for elasticsearch.
	ElasticsearchAuthenticationTypeBasic = "basic"
)

// ElasticsearchServiceConfiguration represents service
// configuration for elasticsearch services (fka sockets).
type ElasticsearchServiceConfiguration struct {
	ElasticsearchServiceType string `json:"elasticsearch_service_type"`

	// mutually exclusive fields below
	StandardElasticSeachServiceConfiguration *StandardElasticSeachServiceConfiguration `json:"standard_elasticsearch_service_configuration,omitempty"`
}

// StandardElasticSeachServiceConfiguration represents service
// configuration for standard elasticsearch services (fka sockets).
type StandardElasticSeachServiceConfiguration struct {
	HostnameAndPort        // inherited
	Protocol        string `json:"protocol"`
	HostHeader      string `json:"host_header,omitempty"`

	AuthenticationType  string                             `json:"authentication_type"`
	BasicAuthentication *ElasticsearchServiceTypeBasicAuth `json:"basic_auth_configuration,omitempty"`
}

// ElasticsearchServiceTypeBasicAuth represents basic auth configuration that based on username and password.
type ElasticsearchServiceTypeBasicAuth struct {
	UsernameAndPassword
}

// Validate validates the ElasticSearchServiceConfiguration.
func (c *ElasticsearchServiceConfiguration) Validate() error {
	switch c.ElasticsearchServiceType {

	case ElasticsearchServiceTypeStandard:
		if c.StandardElasticSeachServiceConfiguration == nil {
			return fmt.Errorf(
				"elasticsearch service configuration for service type \"%s\" must have standard elasticsearch service configuration defined",
				ElasticsearchServiceTypeStandard,
			)
		}
		if err := c.StandardElasticSeachServiceConfiguration.Validate(); err != nil {
			return fmt.Errorf("invalid standard elasticsearch service configuration: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("http service configuration has invalid http service type \"%s\"", c.ElasticsearchServiceType)
	}
}

// Validate validates the ElasticsearchServiceTypeStandard.
func (c *StandardElasticSeachServiceConfiguration) Validate() error {
	switch c.Protocol {
	case "http", "https":
	default:
		return fmt.Errorf("protocol must be either \"http\" or \"https\"")
	}
	if err := c.HostnameAndPort.Validate(); err != nil {
		return err
	}

	switch c.AuthenticationType {
	case ElasticsearchAuthenticationTypeBasic:
		if c.BasicAuthentication == nil {
			return fmt.Errorf("basic auth configuration must be provided when authentication_type is \"%s\"", ElasticsearchAuthenticationTypeBasic)
		}
		if err := c.BasicAuthentication.Validate(); err != nil {
			return fmt.Errorf("invalid basic auth configuration: %v", err)
		}
	default:
		return fmt.Errorf("authentication_type must be \"%s\"", ElasticsearchAuthenticationTypeBasic)
	}
	return nil
}

func (c ElasticsearchServiceTypeBasicAuth) Validate() error {
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
