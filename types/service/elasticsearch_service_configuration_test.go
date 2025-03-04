package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ElasticsearcherviceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testStandardConfig := &StandardElasticsearchServiceConfiguration{
		Protocol: "https",
		HostnameAndPort: HostnameAndPort{
			Hostname: "hostname",
			Port:     9902,
		},
		AuthenticationType: ElasticsearchAuthenticationTypeBasic,
		BasicAuthentication: &ElasticsearchServiceTypeBasicAuth{
			UsernameAndPassword: UsernameAndPassword{
				Username: "username",
				Password: "password",
			},
		},
	}

	tests := []struct {
		name  string
		given ElasticsearchServiceConfiguration
		want  error
	}{
		{
			name:  "elasticsearch service type is missing",
			given: ElasticsearchServiceConfiguration{
				// elasticsearch service type is missing
			},
			want: errors.New("elasticsearch service configuration has invalid elasticsearch service type \"\""),
		},
		{
			name: "invalid elasticsearch service type",
			given: ElasticsearchServiceConfiguration{
				ElasticsearchServiceType: "invalid",
			},
			want: errors.New("elasticsearch service configuration has invalid elasticsearch service type \"invalid\""),
		},
		{
			name: "standard type is picked, but standard config is missing",
			given: ElasticsearchServiceConfiguration{
				ElasticsearchServiceType: ElasticsearchServiceTypeStandard,
				// standard config is missing
			},
			want: errors.New("elasticsearch service configuration for service type \"standard\" must have standard elasticsearch service configuration defined"),
		},
		{
			name: "happy path - standard config",
			given: ElasticsearchServiceConfiguration{
				ElasticsearchServiceType:                  ElasticsearchServiceTypeStandard,
				StandardElasticsearchServiceConfiguration: testStandardConfig,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_StandardElasticsearchServiceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testProtocol := "https"
	testHostnameAndPort := HostnameAndPort{
		Hostname: "hostname",
		Port:     9902,
	}
	testBasicAuth := &ElasticsearchServiceTypeBasicAuth{
		UsernameAndPassword: UsernameAndPassword{
			Username: "username",
			Password: "password",
		},
	}

	tests := []struct {
		name  string
		given StandardElasticsearchServiceConfiguration
		want  error
	}{
		{
			name:  "elasticsearch protocol is missing",
			given: StandardElasticsearchServiceConfiguration{
				// elasticsearch protocol is missing
			},
			want: errors.New("protocol must be either \"http\" or \"https\""),
		},
		{
			name: "hostname is missing",
			given: StandardElasticsearchServiceConfiguration{
				Protocol: testProtocol,
				HostnameAndPort: HostnameAndPort{
					// hostname is missing
					Port: 9902,
				},
				AuthenticationType:  ElasticsearchAuthenticationTypeBasic,
				BasicAuthentication: testBasicAuth,
			},
			want: errors.New("hostname is a required field"),
		},
		{
			name: "invalid elasticsearch authentication type",
			given: StandardElasticsearchServiceConfiguration{
				Protocol:           testProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: "invalid",
			},
			want: errors.New("authentication_type must be \"basic\""),
		},
		{
			name: "basic auth is picked, but basic auth config is missing",
			given: StandardElasticsearchServiceConfiguration{
				Protocol:           testProtocol,
				HostnameAndPort:    testHostnameAndPort,
				AuthenticationType: ElasticsearchAuthenticationTypeBasic,
				// basic auth config is missing
			},
			want: errors.New("basic auth configuration must be provided when authentication_type is \"basic\""),
		},
		{
			name: "happy path - username and password auth",
			given: StandardElasticsearchServiceConfiguration{
				Protocol:            testProtocol,
				HostnameAndPort:     testHostnameAndPort,
				AuthenticationType:  ElasticsearchAuthenticationTypeBasic,
				BasicAuthentication: testBasicAuth,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}

func Test_ElasticsearchBasicdAuthConfiguration_Validate(t *testing.T) {
	t.Parallel()

	testUsername := "username"
	testPassword := "password"

	tests := []struct {
		name  string
		given ElasticsearchServiceTypeBasicAuth
		want  error
	}{
		{
			name: "username is missing",
			given: ElasticsearchServiceTypeBasicAuth{
				UsernameAndPassword: UsernameAndPassword{
					// username is missing
					Password: testPassword,
				},
			},
			want: errors.New("username is required"),
		},
		{
			name: "password is missing",
			given: ElasticsearchServiceTypeBasicAuth{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					// password is missing
				},
			},
			want: errors.New("password is required"),
		},
		{
			name: "happy path",
			given: ElasticsearchServiceTypeBasicAuth{
				UsernameAndPassword: UsernameAndPassword{
					Username: testUsername,
					Password: testPassword,
				},
			},
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.given.Validate()
			assert.Equal(t, test.want, err)
		})
	}
}
