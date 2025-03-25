package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateHttpServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *HttpServiceConfiguration
		expectedError error
	}{
		{
			name: "Happy case for http service type standard",
			configuration: &HttpServiceConfiguration{
				HttpServiceType: HttpServiceTypeStandard,
				StandardHttpServiceConfiguration: &StandardHttpServiceConfiguration{
					HostnameAndPort: HostnameAndPort{
						Hostname: "hello.com",
						Port:     443,
					},
					HostHeader: "whatever.com",
					Scheme:     "https",
				},
			},
			expectedError: nil,
		},
		{
			name: "Happy case for http service type connector file server",
			configuration: &HttpServiceConfiguration{
				HttpServiceType: HttpServiceTypeConnectorFileServer,
				FileServerHttpServiceConfiguration: &FileServerHttpServiceConfiguration{
					TopLevelDirectory: "/root",
				},
			},
			expectedError: nil,
		},
		{
			name:          "Should fail for http service type standard with missing config",
			configuration: &HttpServiceConfiguration{HttpServiceType: HttpServiceTypeStandard},
			expectedError: fmt.Errorf("http service configuration for http service type \"%s\" must have standard http service configuration defined", HttpServiceTypeStandard),
		},
		{
			name:          "Should fail for http service type connector-file-server with missing config",
			configuration: &HttpServiceConfiguration{HttpServiceType: HttpServiceTypeConnectorFileServer},
			expectedError: fmt.Errorf("http service configuration for http service type \"%s\" must have file server http service configuration defined", HttpServiceTypeConnectorFileServer),
		},
		{
			name: "Should fail for http service type standard with invalid config",
			configuration: &HttpServiceConfiguration{
				HttpServiceType:                  HttpServiceTypeStandard,
				StandardHttpServiceConfiguration: &StandardHttpServiceConfiguration{},
			},
			expectedError: errors.New("invalid standard http service configuration: host_header is a required field"),
		},
		{
			name: "Should fail for http service type connector-file-server with invalid config",
			configuration: &HttpServiceConfiguration{
				HttpServiceType:                    HttpServiceTypeConnectorFileServer,
				FileServerHttpServiceConfiguration: &FileServerHttpServiceConfiguration{},
			},
			expectedError: errors.New("invalid file server http service configuration: top_level_directory is a required field"),
		},
		{
			name: "Should fail for tls service type standard with extraneous config",
			configuration: &HttpServiceConfiguration{
				HttpServiceType: HttpServiceTypeStandard,
				StandardHttpServiceConfiguration: &StandardHttpServiceConfiguration{
					HostnameAndPort: HostnameAndPort{
						Hostname: "hello.com",
						Port:     443,
					},
					HostHeader: "whatever.com",
				},
				FileServerHttpServiceConfiguration: &FileServerHttpServiceConfiguration{}, // extra
			},
			expectedError: fmt.Errorf("http service type \"%s\" can only have standard http service configuration defined", TlsServiceTypeStandard),
		},
		{
			name: "Should fail for tls service type vpn with extraneous config",
			configuration: &HttpServiceConfiguration{
				HttpServiceType: HttpServiceTypeConnectorFileServer,
				FileServerHttpServiceConfiguration: &FileServerHttpServiceConfiguration{
					TopLevelDirectory: "/root",
				},
				StandardHttpServiceConfiguration: &StandardHttpServiceConfiguration{}, // extra
			},
			expectedError: fmt.Errorf("http service type \"%s\" can only have file server http service configuration defined", HttpServiceTypeConnectorFileServer),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}

}

func Test_ValidateStandardHttpServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *StandardHttpServiceConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when hostname and port are valid",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				HostHeader: "whatever.com",
				Scheme:     "http",
			},
			expectedError: nil,
		},
		{
			name: "Should fail when hostname-and-port is invalid",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Port: 443,
				},
				HostHeader: "whatever.com",
				Scheme:     "https",
			},
			expectedError: errors.New("hostname is a required field"),
		},
		{
			name: "Should fail when hostname header is invalid",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				Scheme: "https",
			},
			expectedError: errors.New("host_header is a required field"),
		},
		{
			name: "Should succeed when headers are valid",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				HostHeader: "whatever.com",
				Scheme:     "https",
				Headers: []Header{
					{"key1", "value1"},
					{"key2", "value2"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Should fail when header key is empty",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				Scheme:     "https",
				HostHeader: "whatever.com",
				Headers: []Header{
					{"", "value1"},
				},
			},
			expectedError: errors.New("headers key cannot be empty"),
		},
		{
			name: "Should fail when header key is invalid",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				Scheme:     "https",
				HostHeader: "whatever.com",
				Headers: []Header{
					{"!key1", "value1"},
				},
			},
			expectedError: errors.New("header key !key1 is an invalid header key"),
		},
		{
			name: "Should fail when header value is invalid",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				Scheme:     "https",
				HostHeader: "whatever.com",
				Headers: []Header{
					{"key1", "header \r\n value"},
				},
			},
			expectedError: errors.New("header value for key key1 with value header \r\n value contains invalid character: 13"),
		},
		{
			name: "Should fail when header key is too long",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				Scheme:     "https",
				HostHeader: "whatever.com",
				Headers: []Header{
					{string(make([]byte, 256)), "value1"},
				},
			},
			expectedError: errors.New("header key exceeds max length of 255"),
		},
		{
			name: "Should fail when header value is too long",
			configuration: &StandardHttpServiceConfiguration{
				HostnameAndPort: HostnameAndPort{
					Hostname: "hello.com",
					Port:     443,
				},
				Scheme:     "https",
				HostHeader: "whatever.com",
				Headers: []Header{
					{"key1", string(make([]byte, 4097))},
				},
			},
			expectedError: errors.New("header value for key key1 exceeds max length of 4096"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}

func Test_ValidateFileServerHttpServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *FileServerHttpServiceConfiguration
		expectedError error
	}{
		{
			name: "Should succeed when top level directory is valid",
			configuration: &FileServerHttpServiceConfiguration{
				TopLevelDirectory: "/root",
			},
			expectedError: nil,
		},
		{
			name:          "Should fail when top level directory is not present",
			configuration: &FileServerHttpServiceConfiguration{},
			expectedError: errors.New("top_level_directory is a required field"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedError, test.configuration.Validate())
		})
	}
}
