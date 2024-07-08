//go:build !windows
// +build !windows

package utils

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func IsRoot() bool {
	return os.Geteuid() == 0
}

func DropPrivileges(newUid, newGid string) error {
	if newUid == "" {
		return errors.New("user option is empty, cant drop privileges")
	}

	if newGid == "" {
		return errors.New("group option is unset, cant drop privileges")
	}

	gid, err := strconv.Atoi(newGid)
	if err != nil {
		g, err := user.LookupGroup(newGid)
		if err != nil {
			return err
		}
		gid, err = strconv.Atoi(g.Gid)
		if err != nil {
			return err
		}
	}

	if err = syscall.Setgid(gid); err != nil {
		return err
	}

	uid, err := strconv.Atoi(newUid)
	if err != nil {
		u, err := user.Lookup(newUid)
		if err != nil {
			return err
		}
		uid, err = strconv.Atoi(u.Uid)
		if err != nil {
			return err
		}
	}
	if err = syscall.Setuid(uid); err != nil {
		return err
	}

	return nil
}
