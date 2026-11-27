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
	t.Run("0001.1.1 an actor's name or ID opening a GIVEN passes", func(t *testing.T) {
		t.Parallel()
		for _, given := range []string{"User" + aaOpens, "user" + aaOpens, "User", "User has an app"} {
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 0)
			passedWith(t, out, aaCheck, "1")
		}
		// Host is declared but named by no Given.
		out, code := aaRun(t, aaWith("User"+aaOpens, aaRel, aaHost))
		sees(t, out, code, 0)
		passedWith(t, out, aaCheck, "1")
	})
	t.Run("0001.1.2 no declared actor opening a GIVEN fails naming the line and the actors", func(t *testing.T) {
		t.Parallel()
		for _, given := range []string{"no network", "Users open the app", "Admin" + aaOpens, "USER" + aaOpens, "a user" + aaOpens} {
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
		}
	})
	t.Run("0001.1.3 actors joined by and pass until one names none", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaWith("User and Host open the app", aaRel, aaHost))
		sees(t, out, code, 0)
		given := "User and Admin open the app"
		out, code = aaRun(t, aaWith(given, aaRel, aaHost))
		sees(t, out, code, 1, aaNoActor(given, " names no actor at \"Admin open the app\"", "User, Host"))
		out, code = aaRun(t, aaWith("User and", aaRel, aaHost))
		sees(t, out, code, 1, aaNoActor("User and", " names no actor at \"\"", "User, Host"))
	})
	t.Run("0001.1.4 an ID spelling another actor's name is ambiguous", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaWith("User"+aaOpens, aaRel, aaRel+"\nPerson(User, \"4 Host\")"))
		sees(t, out, code, 1, aaGiven+"actor \"User\" is ambiguous: user at "+aaContext+":5 and User at "+aaContext+":8")
	})
	t.Run("0001.1.5 only a context view's C4 Person or Person_Ext are actors", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, aaWith("Guest"+aaOpens, aaRel, aaRel+"\nPerson_Ext(guest, \"4 Guest\")"))
		sees(t, out, code, 0)
		for _, given := range []string{"App" + aaOpens, "app" + aaOpens, "Person" + aaOpens, "p" + aaOpens} {
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
		}
		// Admin is a person in the containers view only.
		files := aaFixture("Admin" + aaOpens)
		files[acRoot+acContainer] = strings.Replace(files[acRoot+acContainer], "Person(user, \"1 User\")", "Person(user, \"1 User\")\nPerson(admin, \"9 Admin\")", 1)
		out, code = aaRun(t, files)
		sees(t, out, code, 1, aaNoActor("Admin"+aaOpens, " starts with no actor", "User"))
		given := "User" + aaOpens
		out, code = aaRun(t, aaWith(given, "C4Context\nPerson(user, \"1 User\")\nSystem(app, \"2 App\")\nRel(user, app, \"3\")", "flowchart LR\nuser((1 User)):::person\napp[2 App]:::system\nuser -->|3| app"))
		sees(t, out, code, 1, aaGiven+"\""+given+"\" starts with no actor; no actor is declared at "+aaContext+":3")
	})
	t.Run("0001.1.6 a repo with no architecture views passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := aaRun(t, map[string]string{"src/main.go": "package main\n", acDoc: acDocWith("no network")})
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, aaCheck)
	})
	t.Run("0001.1.7 doc shape and missing views are left to other checks", func(t *testing.T) {
		t.Parallel()
		files := aaFixture("User" + aaOpens)
		// A bullet acceptance is read for no GIVEN, so even its unknown actor is left to requirements.
		bullets := "- 0001.1\n  - Given Admin" + aaOpens + "\n    - When it starts\n      - Then it opens\n"
		head, _, found := strings.Cut(files[acDoc], "1. Starts")
		if !found {
			t.Fatalf("the fixture doc no longer opens its acceptance with 1. Starts:\n%s", files[acDoc])
		}
		for _, doc := range []string{
			strings.Replace(files[acDoc], "      1. WHEN it starts\n", "", 1),
			head + bullets,
		} {
			files[acDoc] = doc
			out, code := aaRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aaCheck, "1")
			if strings.Contains(out, acDoc) {
				t.Errorf("doc shape belongs to the requirements check, got:\n%s", out)
			}
		}
		files = aaFixture("User" + aaOpens)
		delete(files, acRoot+acContainer)
		out, code := aaRun(t, files)
		sees(t, out, code, 0)
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a pull request changing only tests or changelogs passes with no files", func(t *testing.T) {
		t.Parallel()
		// The base already names an undeclared Host, which the change leaves alone.
		repo := example(t, "happyRepo", aaFixture("Host"+aaOpens))
		git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		commit(t, repo, map[string]string{"src/store/store_test.go": "package store\n\n// Checked.\n", "changelog/fix.md": "Fixed.\n"}, "test(store): touch the store test")
		out, code := fitness(t, repo, prEnv, aaCheck)
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, aaCheck)
	})
	t.Run("0001.2.2 any other run fails each GIVEN naming a removed actor", func(t *testing.T) {
		t.Parallel()
		given := "Host" + aaOpens
		out, code := aaRun(t, aaWith(given, aaRel, aaHost))
		sees(t, out, code, 0)
		repo := example(t, "happyRepo", aaFixture(given))
		for _, changed := range []string{"README.md", "src/store/store.go", acRoot + acContext} {
			out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, aaCheck)
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
