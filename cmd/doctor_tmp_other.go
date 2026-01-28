//go:build !darwin
// +build !darwin

package cmd

func checkTmpSize() doctorCheck {
	return doctorCheck{
		Name:    "tmp_size",
		Status:  "ok",
		Message: "check skipped on this platform",
	}
}
