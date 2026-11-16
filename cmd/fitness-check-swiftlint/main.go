package main

import (
	"github.com/may-journal/fitness-checks/internal/lang/swift/swiftlint"
	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/lang/swift"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{
			Name:  "swiftlint",
			Scope: checkkit.ScopeFiles,
			Files: match.Any(match.Exts(swift.Ext), match.Base(".swiftlint.yml")),
		},
		Run: swiftlint.Run,
	})
}
