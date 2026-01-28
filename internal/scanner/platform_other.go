//go:build !darwin && !linux && !windows

package scanner

func PlatformRules() []Rule {
	return []Rule{}
}
