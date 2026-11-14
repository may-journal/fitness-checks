package main

import (
	"github.com/may-journal/fitness-checks/internal/lang/java/spotless"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "spotless", Scope: checkkit.ScopeFiles, Stage: checkkit.StageTests},
		Run:      spotless.Run,
	})
}
