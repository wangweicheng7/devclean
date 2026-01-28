package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/wangweicheng7/devclean/internal/config"
	"github.com/wangweicheng7/devclean/internal/scanner"
)

type scanItem struct {
	Rule  string   `json:"rule"`
	Size  int64    `json:"size"`
	Paths []string `json:"paths"`
}

type scanOutput struct {
	Results   []scanItem `json:"results"`
	TotalSize int64      `json:"total_size"`
}

var scanJSON bool
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan development junk files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Platform: %s\n\n", runtime.GOOS)

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := scanner.New(scanner.PlatformRules(), cfg.Ignore)
		results, err := s.Scan()
		if err != nil {
			return err
		}
		if scanJSON {
			var out scanOutput

			for _, r := range results {
				out.Results = append(out.Results, scanItem{
					Rule:  r.Rule.Name,
					Size:  r.Size,
					Paths: r.Paths,
				})
				out.TotalSize += r.Size
			}

			return printJSON(out)
		}

		var total int64
		for _, r := range results {
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

		return nil
	},
}

func init() {
	scanCmd.Flags().BoolVar(&scanJSON, "json", false, "Output scan result as JSON")
	rootCmd.AddCommand(scanCmd)
}
