package main

import (
	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-checks/internal/mdfilename"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/mdx"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{
			Name:     "markdown-filename-camel-case",
			Scope:    checkkit.ScopeFiles,
			Files:    match.Exts(mdx.Exts...),
			Requires: []checkkit.Requirement{checkkit.RequiresRepo},
		},
		Run: run,
	})
}

func run(root string, _ []string) (checkkit.Result, error) {
	return mdfilename.Run(root, mdfilename.Camel), nil
}
