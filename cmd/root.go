/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	dryRun  bool
)

// rootCmd 是 devclean 的根命令
var rootCmd = &cobra.Command{
	Use:   "devclean",
	Short: "Developer environment cleaner for macOS",
	Long: `devclean is a macOS-oriented developer environment cleanup tool.

It helps you safely clean:
  - build caches
  - dependency directories
  - temporary files
  - logs

WITHOUT touching git history or source code.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 是 CLI 入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"enable verbose output",
	)

	rootCmd.PersistentFlags().BoolVar(
		&dryRun,
		"dry-run",
		false,
		"show what would be cleaned without deleting",
	)
}
