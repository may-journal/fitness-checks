package main

import (
	"github.com/may-journal/fitness-checks/internal/lang/java/checkstyle"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "checkstyle", Scope: checkkit.ScopeFiles, Stage: checkkit.StageTests},
		Run:      checkstyle.Run,
	})
}
