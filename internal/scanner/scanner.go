package scanner

import (
	"os"
	"path/filepath"
)

type Result struct {
	Rule  Rule
	Size  int64
	Paths []string
}

type Scanner struct {
	Rules []Rule
}

func New(rules []Rule) *Scanner {
	return &Scanner{Rules: rules}
}

func (s *Scanner) Scan() ([]Result, error) {
	var results []Result

	for _, rule := range s.Rules {
		var total int64
		var hitPaths []string

		for _, p := range rule.Paths {
			expanded := expandPath(p)
			matches, _ := filepath.Glob(expanded)

			for _, m := range matches {
				info, err := os.Lstat(m)
				if err != nil {
					continue
				}

				// 不跟随 symlink
				if info.Mode()&os.ModeSymlink != 0 {
					continue
				}

				size := dirSize(m)
				if size > 0 {
					total += size
					hitPaths = append(hitPaths, m)
				}
			}
		}

		if total > 0 {
			results = append(results, Result{
				Rule:  rule,
				Size:  total,
				Paths: hitPaths,
			})
		}
	}

	return results, nil
}

func (r Result) Clean(dryRun bool) error {
	for _, p := range r.Paths {
		// 再次安全校验
		if isUnsafePath(p) {
			continue
		}

		if dryRun {
			continue
		}

		_ = os.RemoveAll(p)
	}
	return nil
}

func expandPath(path string) string {
	if path[:1] == "~" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}

func dirSize(path string) int64 {
	var size int64

	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		size += info.Size()
		return nil
	})

	return size
}

func isUnsafePath(path string) bool {
	base := filepath.Base(path)
	if base == ".git" {
		return true
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return true
	}

	home, _ := os.UserHomeDir()
	if abs == "/" || abs == home {
		return true
	}

	return false
}
