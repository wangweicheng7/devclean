//go:build darwin
// +build darwin

package cmd

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func checkTmpSize() doctorCheck {
	var stat unix.Statfs_t
	if err := unix.Statfs("/tmp", &stat); err != nil {
		return doctorCheck{
			Name:    "tmp_size",
			Status:  "warn",
			Message: "cannot stat /tmp",
		}
	}

	size := int64(stat.Blocks) * int64(stat.Bsize)
	if size > 10*1024*1024*1024 {
		return doctorCheck{
			Name:    "tmp_size",
			Status:  "warn",
			Message: fmt.Sprintf("/tmp is large (%s)", humanSize(size)),
		}
	}

	return doctorCheck{
		Name:   "tmp_size",
		Status: "ok",
	}
}
