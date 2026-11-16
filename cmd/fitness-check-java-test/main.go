package main

import (
	"github.com/may-journal/fitness-checks/internal/describe"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/lang/java/tests"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: describe.JavaBuild("java-test", checkkit.StageTests),
		Run:      tests.Run,
	})
}
