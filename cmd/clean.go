package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangweicheng7/devclean/internal/config"
	"github.com/wangweicheng7/devclean/internal/rules"
	"github.com/wangweicheng7/devclean/internal/scanner"
)

var (
	yesFlag   bool
	ruleNames string
)

type cleanItem struct {
	Rule  string   `json:"rule"`
	Size  int64    `json:"size"`
	Paths []string `json:"paths"`
}

type cleanOutput struct {
	Items     []cleanItem `json:"items"`
	TotalSize int64       `json:"total_size"`
	DryRun    bool        `json:"dry_run"`
}

var cleanJSON bool
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean development junk files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(onlyRules) > 0 && len(excludeRules) > 0 {
			return fmt.Errorf("--only and --exclude cannot be used together")
		}
		if cleanJSON && !dryRun {
			return fmt.Errorf("--json can only be used with --dry-run")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		allRules := scanner.PlatformRules()
		filtered := rules.Filter(allRules, onlyRules, excludeRules)
		s := scanner.New(filtered, cfg.Ignore)
		results, err := s.Scan()
		if err != nil {
			return err
		}
		if cleanJSON {
			var out cleanOutput
			out.DryRun = true

			for _, r := range results {
				out.Items = append(out.Items, cleanItem{
					Rule:  r.Rule.Name,
					Size:  r.Size,
					Paths: r.Paths,
				})
				out.TotalSize += r.Size
			}

			return printJSON(out)
		}

		enabled := map[string]bool{}
		if len(cfg.Rules) > 0 {
			for _, r := range cfg.Rules {
				enabled[strings.ToLower(r)] = true
			}
		}

		if ruleNames != "" {
			enabled = map[string]bool{}
			for _, r := range strings.Split(ruleNames, ",") {
				enabled[strings.ToLower(strings.TrimSpace(r))] = true
			}
		}
		filter := map[string]bool{}
		var targets []scanner.Result
		for _, r := range results {
			if len(filter) == 0 || filterMatch(filter, r.Rule.Name) {
				targets = append(targets, r)
			}
		}

		if len(targets) == 0 {
			fmt.Println("Nothing to clean.")
			return nil
		}

		var total int64
		for _, r := range targets {
			fmt.Printf("%-25s %6.2f MB\n",
				r.Rule.Name,
				float64(r.Size)/1024/1024,
			)
			total += r.Size
		}

		fmt.Println("--------------------------------")
		fmt.Printf("Total                     %6.2f MB\n",
			float64(total)/1024/1024,
		)

		if dryRun {
			fmt.Println("\n(dry-run) No files were deleted.")
			return nil
		}

		if !yesFlag {
			fmt.Print("\nProceed with cleaning? [y/N]: ")
			var input string
			fmt.Scanln(&input)
			if strings.ToLower(input) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		for _, r := range targets {
			_ = r.Clean(false)
		}

		fmt.Println("✅ Cleaning completed.")
		return nil
	},
}

func filterMatch(filter map[string]bool, ruleName string) bool {
	for k := range filter {
		if strings.Contains(strings.ToLower(ruleName), strings.ToLower(k)) {
			return true
		}
	}
	return false
}

func ruleEnabled(enabled map[string]bool, name string) bool {
	if len(enabled) == 0 {
		return true
	}
	for k := range enabled {
		if strings.Contains(strings.ToLower(name), k) {
			return true
		}
	}
	return false
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanJSON, "json", false, "Output clean result as JSON (dry-run only)")
	cleanCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "skip confirmation")
	cleanCmd.Flags().StringVar(&ruleNames, "rule", "", "only clean specified rules (comma-separated)")
	rootCmd.AddCommand(cleanCmd)
}
