package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Result struct {
	Rule  Rule
	Size  int64
	Paths []string
}

type Scanner struct {
	Rules  []Rule
	Ignore []string
}

func New(rules []Rule, ignore []string) *Scanner {
	return &Scanner{
		Rules:  rules,
		Ignore: ignore,
	}
}

func (s *Scanner) Scan() ([]Result, error) {
	if len(s.Rules) == 0 {
		return nil, nil
	}

	workerCount := runtime.NumCPU()
	if workerCount > len(s.Rules) {
		workerCount = len(s.Rules)
	}

	ruleCh := make(chan Rule)
	resultCh := make(chan Result)

	var wg sync.WaitGroup

	// workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rule := range ruleCh {
				var (
					totalSize int64
					paths     []string
				)

				for _, pattern := range rule.Paths {
					matches, err := filepath.Glob(expandPath(pattern))
					if err != nil {
						continue
					}

					for _, m := range matches {
						if isIgnored(m, s.Ignore) {
							continue
						}

						info, err := os.Lstat(m)
						if err != nil {
							continue
						}

						totalSize += info.Size()
						paths = append(paths, m)
					}
				}

				if len(paths) > 0 {
					resultCh <- Result{
						Rule:  rule,
						Size:  totalSize,
						Paths: paths,
					}
				}
			}
		}()
	}

	// feed rules
	go func() {
		for _, r := range s.Rules {
			ruleCh <- r
		}
		close(ruleCh)
	}()

	// close result channel
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// collect
	var results []Result
	for r := range resultCh {
		results = append(results, r)
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

func isIgnored(path string, ignore []string) bool {
	for _, p := range ignore {
		exp := expandPath(p)
		if ok, _ := filepath.Match(exp, path); ok {
			return true
		}
	}
	return false
}
