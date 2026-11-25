package main

import (
	"slices"

	"github.com/may-journal/fitness-checks/internal/architecture"
	"github.com/may-journal/fitness-checks/internal/describe"
	"github.com/may-journal/fitness-checks/internal/packages"
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
	var errs, notes []string
	files := 0
	for _, dir := range packages.Roots() {
		found, noted, judged, err := judge(root, dir, slices.Contains(args, fixFlag))
		errs, notes, files = append(errs, found...), append(notes, noted...), files+judged
		if err != nil {
			return verdict(files, errs, notes), err
		}
	}
	return verdict(files, errs, notes), nil
}

// verdict passes with no errors and fails with any, carrying the notes.
func verdict(files int, errs, notes []string) checkkit.Result {
	result := checkkit.Pass(files)
	if len(errs) > 0 {
		result = checkkit.Fail(files, errs...)
	}
	result.Info = notes
	return result
}

// judge compares the journey view of the package root dir with its own
// docs, rewriting its marked region first when fix is set. It returns the
// findings, notes and number of files judged.
func judge(root, dir string, fix bool) (errs, notes []string, files int, err error) {
	if fix {
		if err := architecture.RewriteJourneyRegion(root, dir); err != nil {
			return nil, nil, 0, err
		}
	}
	model, err := architecture.LoadJourney(root, dir)
	errs, notes = model.JourneyErrors()
	errs = append(errs, model.JourneyRegionErrors()...)
	return errs, notes, model.JourneyFiles(), err
}
