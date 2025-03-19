package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SnowflakeServiceConfiguration_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given SnowflakeServiceConfiguration
		want  error
	}{
		{
			name: "fails when account is missing",
			given: SnowflakeServiceConfiguration{
				Username: "username",
				Password: "password",
			},
			want: errSnowflakeValidationNoAccount,
		},
		{
			name: "fails when username is missing",
			given: SnowflakeServiceConfiguration{
				Account:  "account",
				Password: "password",
			},
			want: errSnowflakeValidationNoUsername,
		},
		{
			name: "fails when password is missing",
			given: SnowflakeServiceConfiguration{
				Account:  "account",
				Username: "username",
			},
			want: errSnowflakeValidationNoPassword,
		},
		{
			name: "happy path",
			given: SnowflakeServiceConfiguration{
				Account:  "account",
				Username: "username",
				Password: "password",
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
