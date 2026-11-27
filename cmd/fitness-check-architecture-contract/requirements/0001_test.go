// Package requirements proves the architecture-contract check's
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
	acChained         = acceptance.AcChained
	acEdit            = acceptance.AcEdit
	acFixture         = acceptance.AcFixture
	acView            = acceptance.AcView
	commit            = acceptance.Commit
	example           = acceptance.Example
	fitness           = acceptance.Fitness
	git               = acceptance.Git
	javaPom           = acceptance.JavaPom
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
	prEnv             = acceptance.PullRequestEnv
	sees              = acceptance.Sees
	with              = acceptance.With
)

const (
	acCheck     = acceptance.AcCheck
	acCode      = acceptance.AcCode
	acComponent = acceptance.AcComponent
	acContainer = acceptance.AcContainer
	acContext   = acceptance.AcContext
	acDoc       = acceptance.AcDoc
	acEvidence  = acceptance.AcEvidence
	acJourney   = acceptance.AcJourney
	acRoot      = acceptance.AcRoot
)

// TestMain keeps the host's git config out of every repo the tests make.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 no views pass for tests, tool configs or generated files unless product code needs them", func(t *testing.T) {
		t.Parallel()
		for _, files := range []map[string]string{nil, {
			"vite.config.ts":                        "export default {};\n",
			"tools/check.sh":                        "#!/bin/sh\nexit 0\n",
			"src/store/store_test.go":               "package store\n",
			"src/app.test.ts":                       "export {};\n",
			"src/__tests__/app.ts":                  "export {};\n",
			"pom.xml":                               javaPom(""),
			"src/test/java/AppTest.java":            "class AppTest {}\n",
			"src/integrationTest/java/AppFlow.java": "class AppFlow {}\n",
		}, {
			".gitattributes":  "src/gen/*.go linguist-generated\n",
			"src/gen/gen.go":  "package gen\n",
			"src/gen/gen2.go": "package gen\n",
		}} {
			out, code := acRun(t, files)
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, acCheck)
		}
		for _, file := range []string{"src/main.go", "Sources/App/App.swift", "src/app.ts", "src/app.js", "src/main/java/App.java"} {
			out, code := acRun(t, map[string]string{file: "// The app.\n"})
			sees(t, out, code, 1, "docs/architecture/ has no views, but product code such as "+file+" needs them; add the five views",
				"actors (01-system-context.md), containers (02-containers.md), components (03-components.md), code (04-code.md), and journey (05-user-journey.md)")
		}
	})
	t.Run("0001.1.2 a missing view or one with no diagram fails", func(t *testing.T) {
		t.Parallel()
		files := acFixture(acRoot)
		delete(files, acRoot+acJourney)
		out, code := acRun(t, files)
		sees(t, out, code, 1, "architecture/01-system-context.md:1: missing architecture view 05-user-journey.md")
		files = acFixture(acRoot)
		files[acRoot+acCode] = "# View\n\nNo diagram yet.\n"
		out, code = acRun(t, files)
		sees(t, out, code, 1, "architecture/04-code.md:1: view has no diagram outside its legends")
	})
	t.Run("0001.1.3 views under either root pass as one model", func(t *testing.T) {
		t.Parallel()
		for _, root := range []string{acRoot, "docs/architecture/"} {
			out, code := acRun(t, acFixture(root))
			sees(t, out, code, 0)
			passedWith(t, out, acCheck, "5")
		}
		split := acFixture(acRoot)
		docs := acFixture("docs/architecture/")
		for _, view := range []string{acCode, acJourney} {
			delete(split, acRoot+view)
			split["docs/architecture/"+view] = docs["docs/architecture/"+view]
		}
		out, code := acRun(t, split)
		sees(t, out, code, 0)
		passedWith(t, out, acCheck, "5")
		// Front matter, accessibility lines and Update-named IDs are styling.
		files := acFixture(acRoot)
		files[acRoot+acContext] = acView("---\ntitle: Context\n---\nflowchart LR\naccTitle: Context\nuser((1 User)):::persona\nUpdateApp[2 App]:::system\nuser -->|3| UpdateApp\nclassDef persona fill:#fff", "A person.", "The system.", "Uses it.")
		out, code = acRun(t, files)
		sees(t, out, code, 0)
		// A source no element claims is listed, not failed.
		out, code = acRun(t, with(acFixture(acRoot), map[string]string{"src/other/other.go": "package other\n"}))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.4 renumbered callouts keep diagram IDs passing", func(t *testing.T) {
		t.Parallel()
		for _, root := range []string{acRoot, "docs/architecture/"} {
			files := acFixture(root)
			files[root+acContext] = strings.NewReplacer("\"1 User\"", "\"21 User\"", "\"2 App\"", "\"22 App\"", "\"3\"", "\"23\"", "| 1 |", "| 21 |", "| 2 |", "| 22 |", "| 3 |", "| 23 |").Replace(files[root+acContext])
			out, code := acRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, acCheck, "5")
		}
	})
	t.Run("0001.1.5 only tests, changelog or docs changing in a pull request passes with no files", func(t *testing.T) {
		t.Parallel()
		broken := acEdit(acComponent, acEvidence, "#L900)")
		for _, changed := range []string{"src/store/store_test.go", "CHANGELOG.md", "docs/guide.md"} {
			out, code := acPullRequest(t, broken, map[string]string{changed: "// Changed.\n"})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, acCheck)
		}
	})
	t.Run("0001.1.6 a broken view stays reported when another file changes", func(t *testing.T) {
		t.Parallel()
		want := "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#L900"
		broken := acEdit(acComponent, acEvidence, "#L900)")
		repo := example(t, "happyRepo", broken)
		for _, changed := range []string{"README.md", "src/store/store.go", acRoot + acContext} {
			out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, acCheck)
			sees(t, out, code, 1, want)
		}
		for _, changed := range []string{"src/store/store.go", acRoot + acContext, acDoc} {
			files := acEdit(acComponent, acEvidence, "#L900)")
			out, code := acPullRequest(t, files, map[string]string{changed: files[changed] + "\n"})
			sees(t, out, code, 1, want)
		}
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 duplicate, invalid or undeclared IDs name the view line", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acContext, "System(app, \"2 App\")", "System(user, \"2 App\")", "architecture/01-system-context.md:6: duplicate ID user; first at architecture/01-system-context.md:5"},
			{acContext, "System(app,", "System(2app,", "architecture/01-system-context.md:6: invalid ID \"2app\""},
			{acContext, "Rel(user, app,", "Rel(user, ghost,", "architecture/01-system-context.md:7: unresolved ID ghost"},
			{acCode, "save --> Record", "save --> Missing", "architecture/04-code.md:9: unresolved ID Missing"},
		})
	})
	t.Run("0001.2.2 a container or component outside any boundary fails", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acComponent, "Container_Boundary(service, \"1 Service (src/cmd/service)\") {\n", "", "architecture/03-components.md:5: component entry needs an enclosing boundary"},
			{acContainer, "System_Boundary(app, \"2 App\") {\n", "", "architecture/02-containers.md:6: container service needs an enclosing boundary"},
		})
	})
	t.Run("0001.2.3 one person with two names or one name on two IDs fails", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acContainer, "Person(user, \"1 User\")", "Person(user, \"1 Customer\")", "architecture/02-containers.md:5: actor user is named \"Customer\" here and \"User\" at architecture/01-system-context.md:5"},
			{acContainer, "System_Boundary", "Person(member, \"5 User\")\nSystem_Boundary", "architecture/02-containers.md:6: actor name \"User\" is ambiguous: member here and user at architecture/01-system-context.md:5"},
		})
	})
	t.Run("0001.2.4 missing or overlapping source claims fail", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acComponent, "Component(entry, \"2 Entry\", \"Go\")", "Component(entry, \"2 Entry\", \"Go\", \"src/store/store.go\")", "architecture/03-components.md:7: overlapping source claims: store claims src/store and entry claims src/store/store.go at architecture/03-components.md:6"},
			{acComponent, "\"src/store\"", "\"src/stores\"", "architecture/03-components.md:7: unresolved source path src/stores"},
			{acComponent, "Keeps records.", "Keeps records in `src/nope*`.", "architecture/03-components.md:16: unresolved source path src/nope*"},
		})
		// A pattern matching a sibling, or text under no tracked folder, passes.
		for _, row := range []string{"Keeps records in `src/st*`.", "Keeps records, not `node_modules/.bin`."} {
			out, code := acRun(t, acEdit(acComponent, "Keeps records.", row))
			sees(t, out, code, 0)
		}
		// A bare root-level source is claimed only once tracked.
		files := acEdit(acComponent, "`src/cmd/service/main.go` starts it.", "`index.ts` starts it.")
		files[acRoot+acComponent] = strings.Replace(files[acRoot+acComponent], "Keeps records.", "Keeps `index.ts` records.", 1)
		out, code := acRun(t, files)
		sees(t, out, code, 0)
		out, code = acRun(t, with(files, map[string]string{"index.ts": "export const ready = true;\n"}))
		sees(t, out, code, 1, "overlapping source claims", "index.ts")
	})
	t.Run("0001.2.5 an evidence link to a missing file or line fails naming the row", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acComponent, acEvidence, "#L900)", "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#L900"},
			{acComponent, acEvidence, "#line3)", "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#line3"},
			{acComponent, acEvidence, "#L3-L1)", "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#L3-L1"},
			{acComponent, "service/main.go" + acEvidence, "service/gone.go" + acEvidence, "architecture/03-components.md:17: unresolved evidence src/cmd/service/gone.go#L3"},
		})
	})
	t.Run("0001.2.6 an unresolved requirement reference fails naming the reference line, leaving doc shape alone", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acJourney, "`0001.1`", "`0001.9`", "architecture/05-user-journey.md:28: unresolved requirement 0001.9 in docs/requirements"},
			{acJourney, "0001-open.md", "0002-open.md", "architecture/05-user-journey.md:29: unresolved requirement 0002 in docs/requirements"},
		})
		files := acFixture(acRoot)
		files[acDoc] = strings.ReplaceAll(files[acDoc], "opens\n-\n", "opens\n")
		out, code := acRun(t, files)
		sees(t, out, code, 0)
		if strings.Contains(out, acDoc) {
			t.Errorf("doc shape belongs to the requirements check, got:\n%s", out)
		}
		for step, want := range map[string]string{"`0001.1.2`": "", "`0001.1.3`": "architecture/05-user-journey.md:28: unresolved requirement 0001.1.3 in docs/requirements"} {
			files := acEdit(acJourney, "`0001.1`", step)
			files[acDoc] = acChained
			out, code := acRun(t, files)
			sees(t, out, code, min(len(want), 1), want)
		}
	})
	t.Run("0001.2.7 an unsupported Mermaid form or C4 macro fails", func(t *testing.T) {
		t.Parallel()
		acFails(t, []acBreak{
			{acCode, "classDiagram", "sequenceDiagram", "architecture/04-code.md:3: unsupported Mermaid form \"sequenceDiagram\""},
			{acContext, "Rel(user, app, \"3\")", "Rel(user, app, \"3\")\nShape(user, \"9\")", "architecture/01-system-context.md:8: unsupported C4 macro \"Shape\""},
		})
	})
}

// acBreak replaces old with new in view of the default fixture and expects
// the error want.
type acBreak struct{ view, old, new, want string }

// acFails runs each break and asserts it fails with its error.
func acFails(t *testing.T, breaks []acBreak) {
	t.Helper()
	for _, b := range breaks {
		out, code := acRun(t, acEdit(b.view, b.old, b.new))
		sees(t, out, code, 1, b.want)
	}
}

// acRun runs `fitness-install -- architecture-contract` on happyRepo with
// files written over it.
func acRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, acCheck)
}

// acPullRequest runs the check as a pull request into main that commits
// changed over happyRepo with files written over it.
func acPullRequest(t *testing.T, files, changed map[string]string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", files)
	git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	commit(t, repo, changed, "fix(app): change a file")
	return fitness(t, repo, prEnv, acCheck)
}
