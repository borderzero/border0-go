package osutil

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
)

// GetUserHomeDir returns the home directory of the current user
// If the process is run as sudo, it will return the home directory of the sudo user
// ie. the orginal user (not root), based on the SUDO_USER env var
func GetUserHomeDir() (string, error) {
	var homedir string

	// check we're using SUDO
	// We can do this by check the SUDO_USER env var
	sudoUsername := os.Getenv("SUDO_USER")
	if sudoUsername != "" {
		if runtime.GOOS == "darwin" {
			// This is because of:
			// https://github.com/golang/go/issues/24383
			// os/user: LookupUser() doesn't find users on macOS when compiled with CGO_ENABLED=0
			// So we'll just hard code for MACOS
			homedir = "/Users/" + sudoUsername
		} else {
			currentUser, err := user.Lookup(sudoUsername)
			if err != nil {
				return "", fmt.Errorf("couldn't get user details: %w", err)
			}
			homedir = currentUser.HomeDir
		}
	} else {
		var err error
		homedir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("couldn't get user home dir: %w", err)
		}
	}
	return homedir, nil
}
