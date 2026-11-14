package main

import (
	"errors"

	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "local-branch", Scope: checkkit.ScopeGitHub},
		Run:      withoutEngine,
	})
}

// withoutEngine reports that this program could not judge, since the engine
// supplies the PR and Plan judgments.
func withoutEngine(string, []string) (checkkit.Result, error) {
	return checkkit.Result{}, errors.New("local-branch runs only inside fitness, which supplies its PR and Plan judgments")
}
