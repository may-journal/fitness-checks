package main

import (
	"github.com/may-journal/fitness-checks/internal/lang/swift/localized"
	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/conf"
	"github.com/may-journal/fitness-runner/go/pkg/lang/swift"
)

func main() {
	program.Configured(docSource, localized.Keys, checkkit.Check{
		Describe: checkkit.Describe{
			Name:     localized.Name,
			Scope:    checkkit.ScopeFiles,
			Files:    match.Any(match.Exts(swift.Ext, ".xcstrings"), match.Base(conf.FileName)),
			Needs:    match.Exts(".xcstrings"),
			Requires: []checkkit.Requirement{checkkit.RequiresRepo},
		},
		Run: localized.Run,
	})
}
