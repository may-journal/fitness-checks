package main

import (
	"github.com/may-journal/fitness-checks/internal/lang/swift/storybook"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{
			Name:     storybook.Name,
			Scope:    checkkit.ScopeFiles,
			Requires: []checkkit.Requirement{checkkit.RequiresRepo},
		},
		Run: storybook.Run,
	})
}
