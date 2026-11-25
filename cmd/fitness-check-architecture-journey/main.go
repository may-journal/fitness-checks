package main

import (
	"slices"

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

// fixFlag rewrites the marked journey region from the requirement docs
// before the check judges.
const fixFlag = "--fix"

func run(root string, args []string) (checkkit.Result, error) {
	if slices.Contains(args, fixFlag) {
		if err := architecture.RewriteJourneyRegion(root); err != nil {
			return checkkit.Result{}, err
		}
	}
	model, err := architecture.Load(root)
	errs, notes := model.JourneyErrors()
	errs = append(errs, model.JourneyRegionErrors()...)
	result := checkkit.Pass(model.JourneyFiles())
	if len(errs) > 0 {
		result = checkkit.Fail(model.JourneyFiles(), errs...)
	}
	result.Info = notes
	return result, err
}
