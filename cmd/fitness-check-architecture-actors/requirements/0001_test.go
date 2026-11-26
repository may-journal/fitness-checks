// Package requirements proves the architecture-actors check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo with a small architecture.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// The shared helpers and architecture fixtures, named as the root
// requirement tests name them.
var (
	acDocWith         = acceptance.AcDocWith
	acFixture         = acceptance.AcFixture
	commit            = acceptance.Commit
	example           = acceptance.Example
	git               = acceptance.Git
	prEnv             = acceptance.PullRequestEnv
	fitness           = acceptance.Fitness
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
	sees              = acceptance.Sees
)

const (
	acContext   = acceptance.AcContext
	acContainer = acceptance.AcContainer
	acDoc       = acceptance.AcDoc
	acRoot      = acceptance.AcRoot
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it passes counting the docs read", func(t *testing.T) {
		t.Parallel()
		for _, given := range []string{"User" + aaOpens, "user" + aaOpens, "User", "User has an app"} {
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 0)
			passedWith(t, out, aaCheck, "1")
		}
	})
	t.Run("0001.1.2.1.2 it passes though some actors start no Given", func(t *testing.T) {
		t.Parallel()
		// Host is declared but named by no Given.
		out, code := aaRun(t, aaWith("User"+aaOpens, aaRel, aaHost))
		sees(t, out, code, 0)
		passedWith(t, out, aaCheck, "1")
	})
	t.Run("0001.1.2.2.1 it fails naming the Given line", func(t *testing.T) {
		t.Parallel()
		given := "no network"
		out, code := aaRun(t, aaFixture(given))
		sees(t, out, code, 1, aaGiven+"\""+given+"\" starts with no actor")
	})
	t.Run("0001.1.2.2.2 it names the context line", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaFixture("no network"))
		sees(t, out, code, 1, "actors declared at "+aaContext+":3")
	})
	t.Run("0001.1.2.2.3 it names the declared actors", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaFixture("no network"))
		sees(t, out, code, 1, aaContext+":3 are User")
	})
	t.Run("0001.1.2.2.4 a plural or another case or an undeclared name fails the same way", func(t *testing.T) {
		t.Parallel()
		for _, given := range []string{"no network", "Users open the app", "Admin" + aaOpens, "USER" + aaOpens, "a user" + aaOpens} {
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
		}
	})
	t.Run("0001.1.2.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaWith("User and Host open the app", aaRel, aaHost))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.3.2 it fails at a word after `and` that names no actor", func(t *testing.T) {
		t.Parallel()
		given := "User and Admin open the app"
		out, code := aaRun(t, aaWith(given, aaRel, aaHost))
		sees(t, out, code, 1, aaNoActor(given, " names no actor at \"Admin open the app\"", "User, Host"))
		out, code = aaRun(t, aaWith("User and", aaRel, aaHost))
		sees(t, out, code, 1, aaNoActor("User and", " names no actor at \"\"", "User, Host"))
	})
	ambiguous := func(t *testing.T) (string, int) {
		return aaRun(t, aaWith("User"+aaOpens, aaRel, aaRel+"\nPerson(User, \"4 Host\")"))
	}
	t.Run("0001.1.2.4.1 it fails as ambiguous naming the Given line", func(t *testing.T) {
		t.Parallel()
		out, code := ambiguous(t)
		sees(t, out, code, 1, aaGiven+"actor \"User\" is ambiguous")
	})
	t.Run("0001.1.2.4.2 it names both context lines", func(t *testing.T) {
		t.Parallel()
		out, code := ambiguous(t)
		sees(t, out, code, 1, "user at "+aaContext+":5 and User at "+aaContext+":8")
	})
	t.Run("0001.1.2.5.1 a system or a legend fails", func(t *testing.T) {
		t.Parallel()
		for _, given := range []string{"App" + aaOpens, "app" + aaOpens, "Person" + aaOpens, "p" + aaOpens} {
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
		}
	})
	t.Run("0001.1.2.5.2 a flowchart `person` node fails", func(t *testing.T) {
		t.Parallel()
		given := "User" + aaOpens
		out, code := aaRun(t, aaWith(given, "C4Context\nPerson(user, \"1 User\")\nSystem(app, \"2 App\")\nRel(user, app, \"3\")", "flowchart LR\nuser((1 User)):::person\napp[2 App]:::system\nuser -->|3| app"))
		sees(t, out, code, 1, aaGiven+"\""+given+"\" starts with no actor; no actor is declared at "+aaContext+":3")
	})
	t.Run("0001.1.2.5.3 another view's person fails", func(t *testing.T) {
		t.Parallel()
		// Admin is a person in the containers view only.
		files := aaFixture("Admin" + aaOpens)
		files[acRoot+acContainer] = strings.Replace(files[acRoot+acContainer], "Person(user, \"1 User\")", "Person(user, \"1 User\")\nPerson(admin, \"9 Admin\")", 1)
		out, code := aaRun(t, files)
		sees(t, out, code, 1, aaNoActor("Admin"+aaOpens, " starts with no actor", "User"))
	})
	t.Run("0001.1.2.5.4 only the context's C4 `Person` or `Person_Ext` are actors", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaWith("Guest"+aaOpens, aaRel, aaRel+"\nPerson_Ext(guest, \"4 Guest\")"))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.6.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, map[string]string{"src/main.go": "package main\n", acDoc: acDocWith("no network")})
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, aaCheck)
	})
	// A bullet acceptance is read for no GIVEN, so even its unknown actor is left to requirements.
	shaped := func(t *testing.T, bullets bool) map[string]string {
		files := aaFixture("User" + aaOpens)
		if !bullets {
			files[acDoc] = strings.Replace(files[acDoc], "      1. WHEN it starts\n", "", 1)
			return files
		}
		head, _, found := strings.Cut(files[acDoc], "1. Starts")
		if !found {
			t.Fatalf("the fixture doc no longer opens its acceptance with 1. Starts:\n%s", files[acDoc])
		}
		files[acDoc] = head + "- 0001.1\n  - Given Admin" + aaOpens + "\n    - When it starts\n      - Then it opens\n"
		return files
	}
	t.Run("0001.1.2.7.1 it reports none of them", func(t *testing.T) {
		t.Parallel()
		for _, bullets := range []bool{false, true} {
			out, code := aaRun(t, shaped(t, bullets))
			sees(t, out, code, 0)
			if strings.Contains(out, acDoc) {
				t.Errorf("doc shape belongs to the requirements check, got:\n%s", out)
			}
		}
	})
	t.Run("0001.1.2.7.2 it leaves them to `requirements`", func(t *testing.T) {
		t.Parallel()
		for _, bullets := range []bool{false, true} {
			out, code := aaRun(t, shaped(t, bullets))
			sees(t, out, code, 0)
			passedWith(t, out, aaCheck, "1")
		}
	})
	t.Run("0001.1.2.7.3 it leaves them to `architecture-contract`", func(t *testing.T) {
		t.Parallel()
		files := aaFixture("User" + aaOpens)
		delete(files, acRoot+acContainer)
		out, code := aaRun(t, files)
		sees(t, out, code, 0)
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	pullRequest := func(t *testing.T) (string, int) {
		// The base already names an undeclared Host, which the change leaves alone.
		repo := example(t, "happyRepo", aaFixture("Host"+aaOpens))
		git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		commit(t, repo, map[string]string{"src/store/store_test.go": "package store\n\n// Checked.\n", "changelog/fix.md": "Fixed.\n"}, "test(store): touch the store test")
		return fitness(t, repo, prEnv, aaCheck)
	}
	t.Run("0001.2.1.1.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := pullRequest(t)
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, aaCheck)
	})
	t.Run("0001.2.1.1.2 such changes cannot move the views", func(t *testing.T) {
		t.Parallel()
		out, code := pullRequest(t)
		sees(t, out, code, 0)
		if strings.Contains(out, acDoc) {
			t.Errorf("a test and changelog change cannot move the views, yet the base's Given was reported:\n%s", out)
		}
	})
	anyRun := func(t *testing.T, changed string) (string, int) {
		return fitness(t, example(t, "happyRepo", aaFixture("Host"+aaOpens)), []string{"FITNESS_CHANGED_FILES=" + changed}, aaCheck)
	}
	t.Run("0001.2.1.2.1 it reads every doc", func(t *testing.T) {
		t.Parallel()
		for _, changed := range []string{"README.md", "src/store/store.go", acRoot + acContext} {
			out, code := anyRun(t, changed)
			sees(t, out, code, 1, aaGiven)
		}
	})
	t.Run("0001.2.1.2.2 it fails each Given naming that actor", func(t *testing.T) {
		t.Parallel()
		given := "Host" + aaOpens
		out, code := aaRun(t, aaWith(given, aaRel, aaHost))
		sees(t, out, code, 0)
		for _, changed := range []string{"README.md", "src/store/store.go", acRoot + acContext} {
			out, code := anyRun(t, changed)
			sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
		}
	})
}

const (
	aaCheck   = "architecture-actors"
	aaGiven   = acDoc + ":16: Given "
	aaContext = "architecture/01-system-context.md"
	aaRel     = "Rel(user, app, \"3\")"
	aaHost    = aaRel + "\nPerson(host, \"4 Host\")"
	aaOpens   = " opens the app"
)

// aaFixture is the default architecture with its one Given replaced.
func aaFixture(given string) map[string]string {
	files := acFixture(acRoot)
	files[acDoc] = acDocWith(given)
	return files
}

// aaWith is aaFixture with one edit to the context diagram.
func aaWith(given, oldContext, newContext string) map[string]string {
	files := aaFixture(given)
	files[acRoot+acContext] = strings.Replace(files[acRoot+acContext], oldContext, newContext, 1)
	return files
}

// aaRun runs `fitness-install -- architecture-actors` on happyRepo with files
// written over it.
func aaRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, aaCheck)
}

// aaNoActor is the failure for a Given whose text names no actor among
// the declared names.
func aaNoActor(given, where, actors string) string {
	return aaGiven + "\"" + given + "\"" + where + "; actors declared at " + aaContext + ":3 are " + actors
}
