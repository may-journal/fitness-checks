// Package requirements proves docs/requirements as a user meets them: each test
// runs the installer at FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// The helpers every requirement test folder shares live in acceptance, so a
// package's own requirement tests use the same ones.
var (
	example           = acceptance.Example
	commit            = acceptance.Commit
	write             = acceptance.Write
	git               = acceptance.Git
	mustDo            = acceptance.MustDo
	fitness           = acceptance.Fitness
	execute           = acceptance.Execute
	started           = acceptance.Started
	userEnv           = acceptance.UserEnv
	installer         = acceptance.Installer
	read              = acceptance.Read
	sees              = acceptance.Sees
	row               = acceptance.Row
	passedWithNoFiles = acceptance.PassedWithNoFiles
	showsInfo         = acceptance.ShowsInfo
	prEnv             = acceptance.PullRequestEnv
	ansi              = acceptance.ANSI
	names             = acceptance.Names
	resultRow         = acceptance.ResultRow
	goRepo            = acceptance.GoRepo
	testFile          = acceptance.TestFile
	with              = acceptance.With
	loggedGit         = acceptance.LoggedGit
	occurrences       = acceptance.Occurrences
	verdictScript     = acceptance.VerdictScript
	bareEntryText     = acceptance.BareEntryText
)

// upstreamChangelog is a CHANGELOG.md as an upstream project keeps it.
const upstreamChangelog = acceptance.UpstreamChangelog

var (
	remote = acceptance.Remote
	user   = acceptance.User
)
