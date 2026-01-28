package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printHuman(checks []doctorCheck) {
	for _, c := range checks {
		switch c.Status {
		case "ok":
			fmt.Printf("✔ %s\n", c.Name)
		case "warn":
			fmt.Printf("⚠ %s: %s\n", c.Name, c.Message)
		case "error":
			fmt.Printf("❌ %s: %s\n", c.Name, c.Message)
		default:
			fmt.Printf("? %s\n", c.Name)
		}
	}
}
