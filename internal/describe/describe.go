// Package describe builds the describe metadata shared by groups of checks.
// It sits apart from package match, so a check that only needs patterns
// links none of the language modules these builders read.
package describe

import (
	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/lang/golang"
	"github.com/may-journal/fitness-runner/go/pkg/lang/java"
)

// RepoCheck describes a files check named name that judges the whole
// tracked tree once a run reaches it.
func RepoCheck(name string) checkkit.Describe {
	return checkkit.Describe{Name: name, Scope: checkkit.ScopeFiles, Requires: []checkkit.Requirement{checkkit.RequiresRepo}}
}

// GoModule describes a test-stage check that runs the Go toolchain over a
// module. Go tests read testdata, embedded files and docs, so any changed
// file reaches it, in a repo with a go.mod.
func GoModule(name string) checkkit.Describe {
	d := RepoCheck(name)
	d.Stage = checkkit.StageTests
	d.Needs = match.Base(golang.Manifest)
	return d
}

// JavaBuild describes a check that needs a Maven or Gradle build, in the
// stage given, which is empty for the fast stage.
func JavaBuild(name, stage string) checkkit.Describe {
	d := RepoCheck(name)
	d.Stage = stage
	d.Needs = match.Any(match.Base(java.Maven.Files...), match.Base(java.Gradle.Files...))
	return d
}
