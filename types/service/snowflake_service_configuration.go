package service

import (
	"errors"
)

var (
	errSnowflakeValidationNoAccount  = errors.New("account is required")
	errSnowflakeValidationNoUsername = errors.New("username is required")
	errSnowflakeValidationNoPassword = errors.New("password is required")
)

// SnowflakeServiceConfiguration represents service
// configuration for snowflake services (fka sockets).
type SnowflakeServiceConfiguration struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ensures that the `SnowflakeServiceConfiguration` has the required fields.
func (config SnowflakeServiceConfiguration) Validate() error {
	if config.Account == "" {
		return errSnowflakeValidationNoAccount
	}
	if config.Username == "" {
		return errSnowflakeValidationNoUsername
	}
	if config.Password == "" {
		return errSnowflakeValidationNoPassword
	}
	return nil
}
