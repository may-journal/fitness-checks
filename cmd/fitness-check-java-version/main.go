package main

import (
	"github.com/may-journal/fitness-checks/internal/lang/java/version"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "java-version", Scope: checkkit.ScopeFiles},
		Run:      version.Run,
	})
}
