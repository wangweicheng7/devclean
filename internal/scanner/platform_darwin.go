//go:build darwin

package scanner

func PlatformRules() []Rule {
	return []Rule{
		{
			Name:        "Xcode DerivedData",
			Description: "Xcode build cache",
			Paths: []string{
				"~/Library/Developer/Xcode/DerivedData",
			},
		},
		{
			Name:        "Node.js node_modules",
			Description: "Node dependency directories",
			Paths: []string{
				"**/node_modules",
			},
		},
		{
			Name:        "Flutter build cache",
			Description: "Flutter build artifacts",
			Paths: []string{
				"**/.dart_tool",
				"**/build",
			},
		},
	}
}
