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
	t.Run("0001.1.3.1.1 product code fails asking for all five", func(t *testing.T) {
		t.Parallel()
		for _, file := range []string{"src/main.go", "Sources/App/App.swift", "src/app.ts", "src/app.js", "src/main/java/App.java"} {
			out, code := acRun(t, map[string]string{file: "// The app.\n"})
			sees(t, out, code, 1, "docs/architecture/ has no views, but product code such as "+file+" needs them; add the five views",
				"actors (01-system-context.md), containers (02-containers.md), components (03-components.md), code (04-code.md), and journey (05-user-journey.md)")
		}
	})
	t.Run("0001.1.3.1.2 tests pass with no files", func(t *testing.T) {
		t.Parallel()
		out, code := acRun(t, map[string]string{
			"src/store/store_test.go":               "package store\n",
			"src/app.test.ts":                       "export {};\n",
			"src/__tests__/app.ts":                  "export {};\n",
			"src/test/java/AppTest.java":            "class AppTest {}\n",
			"src/integrationTest/java/AppFlow.java": "class AppFlow {}\n",
		})
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, acCheck)
	})
	t.Run("0001.1.3.1.3 tool configs pass with no files", func(t *testing.T) {
		t.Parallel()
		for _, files := range []map[string]string{nil, {
			"vite.config.ts": "export default {};\n",
			"tools/check.sh": "#!/bin/sh\nexit 0\n",
			"pom.xml":        javaPom(""),
		}} {
			out, code := acRun(t, files)
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, acCheck)
		}
	})
	t.Run("0001.1.3.1.4 generated files pass with no files", func(t *testing.T) {
		t.Parallel()
		out, code := acRun(t, map[string]string{
			".gitattributes":  "src/gen/*.go linguist-generated\n",
			"src/gen/gen.go":  "package gen\n",
			"src/gen/gen2.go": "package gen\n",
		})
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, acCheck)
	})
	t.Run("0001.1.3.2.1 it fails naming the view", func(t *testing.T) {
		t.Parallel()
		files := acFixture(acRoot)
		delete(files, acRoot+acJourney)
		out, code := acRun(t, files)
		sees(t, out, code, 1, "architecture/01-system-context.md:1: missing architecture view 05-user-journey.md")
	})
	t.Run("0001.1.3.2.2 a view with no diagram outside its legends fails the same way", func(t *testing.T) {
		t.Parallel()
		files := acFixture(acRoot)
		files[acRoot+acCode] = "# View\n\nNo diagram yet.\n"
		out, code := acRun(t, files)
		sees(t, out, code, 1, "architecture/04-code.md:1: view has no diagram outside its legends")
	})
	t.Run("0001.1.3.3.1 they pass as one model", func(t *testing.T) {
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
	t.Run("0001.1.3.4.1 diagram IDs keep their identity", func(t *testing.T) {
		t.Parallel()
		for _, root := range []string{acRoot, "docs/architecture/"} {
			files := acFixture(root)
			files[root+acContext] = strings.NewReplacer("\"1 User\"", "\"21 User\"", "\"2 App\"", "\"22 App\"", "\"3\"", "\"23\"", "| 1 |", "| 21 |", "| 2 |", "| 22 |", "| 3 |", "| 23 |").Replace(files[root+acContext])
			out, code := acRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, acCheck, "5")
		}
	})
	t.Run("0001.1.3.5.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		broken := acEdit(acComponent, acEvidence, "#L900)")
		for _, changed := range []string{"src/store/store_test.go", "CHANGELOG.md", "docs/guide.md"} {
			out, code := acPullRequest(t, broken, map[string]string{changed: "// Changed.\n"})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, acCheck)
		}
	})
	t.Run("0001.1.3.6.1 the broken view stays reported", func(t *testing.T) {
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
	t.Run("0001.2.2.1.1 it fails naming the view line", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acContext, "System(app, \"2 App\")", "System(user, \"2 App\")", "architecture/01-system-context.md:6: duplicate ID user; first at architecture/01-system-context.md:5"},
			{acContext, "System(app,", "System(2app,", "architecture/01-system-context.md:6: invalid ID \"2app\""},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
	})
	t.Run("0001.2.2.1.2 an edge naming an undeclared ID fails the same way", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acContext, "Rel(user, app,", "Rel(user, ghost,", "architecture/01-system-context.md:7: unresolved ID ghost"},
			{acCode, "save --> Record", "save --> Missing", "architecture/04-code.md:9: unresolved ID Missing"},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
	})
	t.Run("0001.2.2.2.1 it fails naming its line", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acComponent, "Container_Boundary(service, \"1 Service (src/cmd/service)\") {\n", "", "architecture/03-components.md:5: component entry needs an enclosing boundary"},
			{acContainer, "System_Boundary(app, \"2 App\") {\n", "", "architecture/02-containers.md:6: container service needs an enclosing boundary"},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
	})
	t.Run("0001.2.2.3.1 it fails naming both lines", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acContainer, "Person(user, \"1 User\")", "Person(user, \"1 Customer\")", "architecture/02-containers.md:5: actor user is named \"Customer\" here and \"User\" at architecture/01-system-context.md:5"},
			{acContainer, "System_Boundary", "Person(member, \"5 User\")\nSystem_Boundary", "architecture/02-containers.md:6: actor name \"User\" is ambiguous: member here and user at architecture/01-system-context.md:5"},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
	})
	t.Run("0001.2.2.4.1 it fails naming the claim lines", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acComponent, "\"src/store\"", "\"src/stores\"", "architecture/03-components.md:7: unresolved source path src/stores"},
			{acComponent, "Keeps records.", "Keeps records in `src/nope*`.", "architecture/03-components.md:16: unresolved source path src/nope*"},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
		// A pattern matching a sibling, or text under no tracked folder, passes.
		for _, row := range []string{"Keeps records in `src/st*`.", "Keeps records, not `node_modules/.bin`."} {
			out, code := acRun(t, acEdit(acComponent, "Keeps records.", row))
			sees(t, out, code, 0)
		}
	})
	t.Run("0001.2.2.4.2 two elements claiming one path or tracked root source fail the same way", func(t *testing.T) {
		t.Parallel()
		b := acBreak{acComponent, "Component(entry, \"2 Entry\", \"Go\")", "Component(entry, \"2 Entry\", \"Go\", \"src/store/store.go\")", "architecture/03-components.md:7: overlapping source claims: store claims src/store and entry claims src/store/store.go at architecture/03-components.md:6"}
		out, code := acRunBreak(t, b)
		sees(t, out, code, 1, b.want)
		// A bare root-level source is claimed only once tracked.
		files := acEdit(acComponent, "`src/cmd/service/main.go` starts it.", "`index.ts` starts it.")
		files[acRoot+acComponent] = strings.Replace(files[acRoot+acComponent], "Keeps records.", "Keeps `index.ts` records.", 1)
		out, code = acRun(t, files)
		sees(t, out, code, 0)
		out, code = acRun(t, with(files, map[string]string{"index.ts": "export const ready = true;\n"}))
		sees(t, out, code, 1, "overlapping source claims", "index.ts")
	})
	t.Run("0001.2.2.5.1 it fails naming the row", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acComponent, acEvidence, "#L900)", "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#L900"},
			{acComponent, acEvidence, "#line3)", "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#line3"},
			{acComponent, acEvidence, "#L3-L1)", "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#L3-L1"},
			{acComponent, "service/main.go" + acEvidence, "service/gone.go" + acEvidence, "architecture/03-components.md:17: unresolved evidence src/cmd/service/gone.go#L3"},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
	})
	t.Run("0001.2.2.6.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		b := acBreak{acJourney, "0001-open.md", "0002-open.md", "architecture/05-user-journey.md:29: unresolved requirement 0002 in docs/requirements"}
		out, code := acRunBreak(t, b)
		sees(t, out, code, 1, b.want)
		// Doc shape is left to the requirements check.
		files := acFixture(acRoot)
		files[acDoc] = strings.ReplaceAll(files[acDoc], "opens\n-\n", "opens\n")
		out, code = acRun(t, files)
		sees(t, out, code, 0)
		if strings.Contains(out, acDoc) {
			t.Errorf("doc shape belongs to the requirements check, got:\n%s", out)
		}
	})
	t.Run("0001.2.2.6.2 an acceptance or chained step that resolves to nothing fails the same way", func(t *testing.T) {
		t.Parallel()
		b := acBreak{acJourney, "`0001.1`", "`0001.9`", "architecture/05-user-journey.md:28: unresolved requirement 0001.9 in docs/requirements"}
		out, code := acRunBreak(t, b)
		sees(t, out, code, 1, b.want)
		for step, want := range map[string]string{"`0001.1.2`": "", "`0001.1.3`": "architecture/05-user-journey.md:28: unresolved requirement 0001.1.3 in docs/requirements"} {
			files := acEdit(acJourney, "`0001.1`", step)
			files[acDoc] = acChained
			out, code := acRun(t, files)
			sees(t, out, code, min(len(want), 1), want)
		}
	})
	t.Run("0001.2.2.7.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		for _, b := range []acBreak{
			{acCode, "classDiagram", "sequenceDiagram", "architecture/04-code.md:3: unsupported Mermaid form \"sequenceDiagram\""},
			{acContext, "Rel(user, app, \"3\")", "Rel(user, app, \"3\")\nShape(user, \"9\")", "architecture/01-system-context.md:8: unsupported C4 macro \"Shape\""},
		} {
			out, code := acRunBreak(t, b)
			sees(t, out, code, 1, b.want)
		}
	})
}

// acBreak replaces old with new in view of the default fixture and expects
// the error want.
type acBreak struct{ view, old, new, want string }

// acRunBreak runs the default fixture with the break applied.
func acRunBreak(t *testing.T, b acBreak) (string, int) {
	t.Helper()
	return acRun(t, acEdit(b.view, b.old, b.new))
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
