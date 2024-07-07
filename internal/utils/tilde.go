package utils

import (
	"os/user"
	"strings"
)

// ExpandTilde replaces ~ with the homedir of the current user
func ExpandTilde(path string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return strings.Replace(path, "~", u.HomeDir, 1), nil
}
