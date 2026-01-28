package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check dev environment health",
	Run: func(cmd *cobra.Command, args []string) {
		checkHomeWritable()
		checkTmpSize()
		checkGitConfig()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func checkHomeWritable() {
	home, _ := os.UserHomeDir()
	test := filepath.Join(home, ".devclean_write_test")

	if err := os.WriteFile(test, []byte("ok"), 0644); err != nil {
		fmt.Println("❌ Home directory not writable")
		return
	}
	_ = os.Remove(test)
	fmt.Println("✔ Home directory writable")
}

func checkTmpSize() {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/tmp", &stat); err != nil {
		return
	}
	size := int64(stat.Blocks) * int64(stat.Bsize)
	if size > 10*1024*1024*1024 {
		fmt.Printf("⚠ /tmp is large (%s)\n", humanSize(size))
	}
}

func checkGitConfig() {
	out, err := exec.Command("git", "config", "--global", "gc.auto").Output()
	if err != nil {
		return
	}
	fmt.Printf("ℹ git gc.auto = %s", out)
}
