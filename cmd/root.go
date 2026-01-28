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
	dryRun  bool
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "devclean",
	Short: "Developer environment cleaner for macOS",
	Long: `devclean is a macOS-oriented CLI tool to clean
development caches, build artifacts and other safe-to-remove files.

Default mode is dry-run (no files will be deleted).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nothing to do. Try 'devclean scan' or 'devclean clean'.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", true, "Preview what will be cleaned (default true)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
