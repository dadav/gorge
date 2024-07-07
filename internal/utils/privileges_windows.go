//go:build windows
// +build windows

package utils

import (
	"errors"
)

func IsRoot() bool {
	return false
}

func DropPrivileges(newUid, newGid string) error {
	return errors.New("cant drop privileges in windows")
}
