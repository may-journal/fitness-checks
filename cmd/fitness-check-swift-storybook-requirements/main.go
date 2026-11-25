package main

import (
	"github.com/may-journal/fitness-checks/internal/config"
	"github.com/may-journal/fitness-checks/internal/lang/swift/storybook"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Configured(docSource, config.Keys(defaults, configSource), checkkit.Check{
		Describe: checkkit.Describe{
			Name:     storybook.Name,
			Scope:    checkkit.ScopeFiles,
			Requires: []checkkit.Requirement{checkkit.RequiresRepo},
		},
		Run: run,
	})
}

// run judges the repo at root with its page limit per acceptance.
func run(root string, _ []string) (checkkit.Result, error) {
	return storybook.Run(root, maxPages(root))
}

// maxPages resolves the page limit: the .fitnessrc.json
// swiftStorybookRequirements.maxPagesPerAcceptance value when positive, else
// the default.
func maxPages(root string) int {
	var c fileConfig
	if config.Read(root, &c) == nil && c.SwiftStorybookRequirements.MaxPagesPerAcceptance > 0 {
		return c.SwiftStorybookRequirements.MaxPagesPerAcceptance
	}
	return defaults.SwiftStorybookRequirements.MaxPagesPerAcceptance
}
