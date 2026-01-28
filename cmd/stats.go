package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/wangweicheng7/devclean/internal/config"
	"github.com/wangweicheng7/devclean/internal/rules"
	"github.com/wangweicheng7/devclean/internal/scanner"
)

type statsItem struct {
	Rule  string `json:"rule"`
	Paths int    `json:"paths"`
	Size  int64  `json:"size"`
}

type statsOutput struct {
	Rules     []statsItem `json:"rules"`
	TotalSize int64       `json:"total_size"`
}

var statsJSON bool
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show disk usage statistics by rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(onlyRules) > 0 && len(excludeRules) > 0 {
			return fmt.Errorf("--only and --exclude cannot be used together")
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
		if statsJSON {
			var out statsOutput

			for _, r := range results {
				out.Rules = append(out.Rules, statsItem{
					Rule:  r.Rule.Name,
					Paths: len(r.Paths),
					Size:  r.Size,
				})
				out.TotalSize += r.Size
			}

			return printJSON(out)
		}

		if len(results) == 0 {
			fmt.Println("✨ No junk found.")
			return nil
		}

		// sort by size desc
		sort.Slice(results, func(i, j int) bool {
			return results[i].Size > results[j].Size
		})

		var total int64

		fmt.Printf("%-16s %6s %12s\n", "Rule", "Paths", "Size")
		fmt.Println("----------------------------------------")

		for _, r := range results {
			fmt.Printf(
				"%-16s %6d %12s\n",
				r.Rule.Name,
				len(r.Paths),
				humanSize(r.Size),
			)
			total += r.Size
		}

		fmt.Println("----------------------------------------")
		fmt.Printf("%-16s %6s %12s\n", "Total", "", humanSize(total))

		return nil
	},
}

func init() {
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "Output stats as JSON")
	rootCmd.AddCommand(statsCmd)
}

func humanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
