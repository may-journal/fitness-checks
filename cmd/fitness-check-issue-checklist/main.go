package main

import (
	"fmt"
	"path"

	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/bodycheck"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/checklist"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "issue-checklist", Scope: checkkit.ScopeBlob, Files: match.Any(match.Under(issueBodies), match.Exact(checkkit.BlobPath)), ContextInlineArg: "--body"},
		Run:      run,
	})
}

// issueBodies is the folder of every issue body's virtual path, whatever its
// kind, since a closed issue of any kind keeps its checklist.
var issueBodies = path.Dir(checkkit.IssueBodyPath("")) + "/"

func run(_ string, args []string) (checkkit.Result, error) {
	return bodycheck.RunIf(args, bodycheck.Closing, func(body string) []string {
		var errs []string
		for _, item := range checklist.Unchecked(body) {
			errs = append(errs, fmt.Sprintf("unchecked item: %q", item))
		}
		return errs
	})
}
