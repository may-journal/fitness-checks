package main

import (
	"github.com/may-journal/fitness-checks/internal/architecture"
	"github.com/may-journal/fitness-checks/internal/describe"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: describe.RepoCheck("architecture-journey"),
		Run:      run,
	})
}

func run(root string, _ []string) (checkkit.Result, error) {
	model, err := architecture.Load(root)
	errs, notes := model.JourneyErrors()
	result := checkkit.Pass(model.JourneyFiles())
	if len(errs) > 0 {
		result = checkkit.Fail(model.JourneyFiles(), errs...)
	}
	result.Info = notes
	return result, err
}
