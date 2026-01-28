package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type doctorCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // ok | warn | error
	Message string `json:"message,omitempty"`
}

type doctorOutput struct {
	Checks []doctorCheck `json:"checks"`
}

var doctorJSON bool
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check dev environment health",
	RunE: func(cmd *cobra.Command, args []string) error {
		checks := []doctorCheck{
			checkHomeWritable(),
			checkTmpSize(),
			checkGitConfig(),
		}

		if doctorJSON {
			return printJSON(doctorOutput{Checks: checks})
		}

		printHuman(checks)
		return nil
	},
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output result as JSON")
	rootCmd.AddCommand(doctorCmd)
}

func checkHomeWritable() doctorCheck {
	home, _ := os.UserHomeDir()
	test := filepath.Join(home, ".devclean_write_test")

	if err := os.WriteFile(test, []byte("ok"), 0644); err != nil {
		return doctorCheck{
			Name:    "home_writable",
			Status:  "error",
			Message: "home directory not writable",
		}
	}
	_ = os.Remove(test)

	return doctorCheck{
		Name:   "home_writable",
		Status: "ok",
	}
}

func checkGitConfig() doctorCheck {
	out, err := exec.Command("git", "config", "--global", "gc.auto").Output()
	if err != nil {
		return doctorCheck{
			Name:    "git_gc",
			Status:  "warn",
			Message: "git not configured or not found",
		}
	}

	return doctorCheck{
		Name:    "git_gc",
		Status:  "ok",
		Message: fmt.Sprintf("gc.auto = %s", strings.TrimSpace(string(out))),
	}
}
