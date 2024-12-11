package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateExitNodeServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *ExitNodeServiceConfiguration
		expectError   bool
	}{
		{
			name:          "Happy case for exit node with no network protocol disabled",
			configuration: &ExitNodeServiceConfiguration{},
			expectError:   false,
		},

		// NOTE(@adrianosela): uncomment if needed later
		//
		// {
		// 	name:          "Happy case for exit node with ipv4 disabled",
		// 	configuration: &ExitNodeServiceConfiguration{DisableIPv4: true},
		// 	expectError:   false,
		// },
		// {
		// 	name:          "Happy case for exit node with ipv6 disabled",
		// 	configuration: &ExitNodeServiceConfiguration{DisableIPv6: true},
		// 	expectError:   false,
		// },
		// {
		// 	name:          "Happy case for exit node with ipv4 and ipv6 disabled",
		// 	configuration: &ExitNodeServiceConfiguration{DisableIPv4: true, DisableIPv6: true},
		// 	expectError:   false,
		// },
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.configuration.Validate()
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
