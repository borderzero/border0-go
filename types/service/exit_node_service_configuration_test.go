package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateExitNodeServiceConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configuration *ExitNodeServiceConfiguration
		expectError   bool
		expectedError error
	}{
		{
			name:          "Happy case for exit node with empty configuration",
			configuration: &ExitNodeServiceConfiguration{},
			expectError:   false,
		},
		{
			name:          "Happy case for exit node with mode dual-stack",
			configuration: &ExitNodeServiceConfiguration{Mode: ExitNodeModeDualStack},
			expectError:   false,
		},

		{
			name:          "Happy case for exit node with mode ipv4 only",
			configuration: &ExitNodeServiceConfiguration{Mode: ExitNodeModeIPv4Only},
			expectError:   false,
		},
		{
			name:          "Fail validation when exit node mode is invalid",
			configuration: &ExitNodeServiceConfiguration{Mode: "invalid"},
			expectError:   true,
			expectedError: fmt.Errorf(exitNodeModeErrFmt, "invalid"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.configuration.Validate()
			if test.expectError {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
