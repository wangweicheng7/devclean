package rules

import (
	"strings"

	"github.com/wangweicheng7/devclean/internal/scanner"
)

func Filter(
	all []scanner.Rule,
	only []string,
	exclude []string,
) []scanner.Rule {
	if len(only) == 0 && len(exclude) == 0 {
		return all
	}

	nameMap := make(map[string]scanner.Rule)
	for _, r := range all {
		nameMap[strings.ToLower(r.Name)] = r
	}

	var result []scanner.Rule

	// --only 优先生效
	if len(only) > 0 {
		for _, name := range only {
			if r, ok := nameMap[strings.ToLower(name)]; ok {
				result = append(result, r)
			}
		}
		return result
	}

	// --exclude
	excludeSet := make(map[string]struct{})
	for _, n := range exclude {
		excludeSet[strings.ToLower(n)] = struct{}{}
	}

	for _, r := range all {
		if _, banned := excludeSet[strings.ToLower(r.Name)]; !banned {
			result = append(result, r)
		}
	}

	return result
}
