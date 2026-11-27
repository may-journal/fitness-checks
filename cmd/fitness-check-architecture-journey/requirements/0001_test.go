// Package requirements proves the architecture-journey check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo with a small architecture.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// The shared helpers and architecture fixtures, named as the root
// requirement tests name them.
var (
	acChained         = acceptance.AcChained
	acFixture         = acceptance.AcFixture
	example           = acceptance.Example
	fitness           = acceptance.Fitness
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
	sees              = acceptance.Sees
	showsInfo         = acceptance.ShowsInfo
)

const (
	acJourney = acceptance.AcJourney
	acRoot    = acceptance.AcRoot
)

// TestMain keeps the host's git config out of every repo the tests make.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 labels naming every doc pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "07b8e753", func(_ req.A) {
			out, code := ajRun(t, ajFixture())
			sees(t, out, code, 0)
			passedWith(t, out, ajCheck, "2")
			files := ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0001, 0002)")
			out, code = ajRun(t, ajEdit(files, "-->|saves| done", "-->|saves 0002| done"))
			sees(t, out, code, 0)
			passedWith(t, out, ajCheck, "3")
			// A character code, an enclosing subgraph's ID and a renamed node all pass.
			for _, passing := range []map[string]string{
				ajEdit(ajFixture(), ajTagOpen, "open(2 Open #9829; the app<br/>Req 0001)"),
				ajEdit(ajFixture(), ajTagDone, "subgraph steps[Steps<br/>Req 0001]\n"+ajDone+"\nend"),
				ajEdit(ajFixture(), ajTagDone, "done[3 Record stored<br/>Req 0001]"),
				ajEdit(ajEdit(ajFixture(), ajTagDone, "stored[3 Record saved<br/>Req 0001]"), "|saves| done", "|saves| stored"),
			} {
				out, code := ajRun(t, passing)
				sees(t, out, code, 0)
			}
		})
	})
	t.Run("0001.1.2 a doc no diagram label names fails at the fence", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "7eaa0015", func(_ req.A) {
			out, code := ajRun(t, acFixture(acRoot))
			sees(t, out, code, 1, ajMissingDoc("0001"))
			out, code = ajRun(t, ajSecondDoc(ajFixture()))
			sees(t, out, code, 1, ajMissingDoc("0002"))
			out, code = ajRun(t, ajEdit(ajSecondDoc(ajFixture()), ajTagDone, "done[3 Record saved<br/>Req 0002]"))
			sees(t, out, code, 0)
			dropped := ajEdit(ajEdit(ajFixture(), ajTagDone, "done[3 Record saved<br/>Req 0002]"), ajTagOpen, "open(2 Open the app<br/>Req 0002)")
			out, code = ajRun(t, ajSecondDoc(dropped))
			sees(t, out, code, 1, ajMissingDoc("0001"))
			// A table row and prose outside the diagram cover nothing.
			files := acFixture(acRoot)
			files[ajView] += ajTable + "\nThe open step fulfils 0001.\n"
			out, code = ajRun(t, files)
			sees(t, out, code, 1, ajMissingDoc("0001"), "journey node open names no requirement", "journey node done names no requirement")
		})
	})
	t.Run("0001.1.3 a label naming an acceptance, unknown ID or nine IDs fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "d367abbe", func(a req.A) {
			for label, want := range map[string]string{
				"open(2 Open the app<br/>Req 0001.1)":     "architecture/05-user-journey.md:18: 0001.1 is an acceptance ID; name requirement doc 0001",
				"open(2 Open the app<br/>Req 0001, 0009)": "architecture/05-user-journey.md:18: no requirement doc 0009 in docs/requirements",
			} {
				out, code := ajRun(t, ajEdit(ajFixture(), ajTagOpen, label))
				sees(a.T(), out, code, 1, want)
			}
			out, code := ajRun(t, ajEdit(ajFixture(), ajTagOpen, "open(2 Open the app<br/>Req 0009)"))
			sees(a.T(), out, code, 1, "architecture/05-user-journey.md:18: no requirement doc 0009 in docs/requirements")
			a.True(!strings.Contains(out, "journey node open names no requirement"), "a mistyped ID also reported its step as naming none:\n%s", out)
			repo := example(t, "happyRepo", ajEdit(ajSecondDoc(ajFixture()), ajTagDone, "done[3 Record saved<br/>Req 0001, 0002]"))
			acceptance.Write(t, repo, map[string]string{ajDoc2: ""})
			out, code = fitness(t, repo, nil, ajCheck)
			sees(a.T(), out, code, 1, "architecture/05-user-journey.md:19: no requirement doc 0002 in docs/requirements")
			// A cited step reads as its doc once that step exists.
			files := ajEdit(ajEdit(acFixture(acRoot), ajOpen, "open(2 Open the app<br/>Req 0001.1.1)"), ajDone, "done[3 Record saved<br/>Req 0001.1.2]")
			files[ajDoc] = acChained
			out, code = ajRun(t, files)
			sees(a.T(), out, code, 0)
			out, code = ajRun(t, ajEdit(files, "Req 0001.1.2", "Req 0001.1.3"))
			sees(a.T(), out, code, 1, "architecture/05-user-journey.md:19: no step 0001.1.3 in docs/requirements")
		})
		req.Then(t, "0001.1.2.3.2", "db687e73", func(a req.A) {
			nine, eight := ajFixture(), ajFixture()
			for _, id := range []string{"0002", "0003", "0004", "0005", "0006", "0007", "0008", "0009"} {
				nine[strings.Replace(ajDoc2, "0002", id, 1)] = strings.ReplaceAll(nine[ajDoc], "0001", id)
				if id != "0009" {
					eight[strings.Replace(ajDoc2, "0002", id, 1)] = strings.ReplaceAll(eight[ajDoc], "0001", id)
				}
			}
			ids := "0001, 0002, 0003, 0004, 0005, 0006, 0007, 0008"
			out, code := ajRun(t, ajEdit(nine, ajTagOpen, "open(2 Open the app<br/>Req "+ids+", 0009)"))
			sees(a.T(), out, code, 1, "architecture/05-user-journey.md:18: label names 9 requirements; split the step so each names at most 8")
			out, code = ajRun(t, ajEdit(eight, ajTagOpen, "open(2 Open the app<br/>Req "+ids+")"))
			sees(a.T(), out, code, 0)
		})
	})
	t.Run("0001.1.4 a behavior node naming no ID fails at the node's line", func(t *testing.T) {
		t.Parallel()
		// An edge label naming 0001 still leaves the done node naming none.
		req.Then(t, "0001.1.2.4.1", "f02c5d4e", func(_ req.A) {
			for want, files := range map[string]map[string]string{
				"architecture/05-user-journey.md:18: journey node open names no requirement": ajEdit(ajFixture(), ajTagOpen, ajOpen),
				"journey node done names no requirement":                                     ajEdit(ajFixture(), ajTagDone, "subgraph steps[Steps]\n"+ajDone+"\nend"),
				"architecture/05-user-journey.md:19: journey node done names no requirement": ajEdit(ajEdit(ajFixture(), ajTagDone, ajDone), "-->|saves| done", "-->|saves 0001| done"),
			} {
				out, code := ajRun(t, files)
				sees(t, out, code, 1, want)
				if strings.Contains(out, "journey node user") || strings.Contains(out, "journey node steps") {
					t.Errorf("a persona or layout subgraph needs no ID:\n%s", out)
				}
			}
			out, code := ajRun(t, ajEdit(ajFixture(), ajTagDone, ajDone))
			sees(t, out, code, 1, "architecture/05-user-journey.md:19: journey node done names no requirement")
		})
	})
	t.Run("0001.1.5 steps reading as each other's docs fail and one alone adds a note", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "fe34b048", func(_ req.A) {
			crossed := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0002)"), ajTagDone, "done[3 Save the record<br/>Req 0001]")
			out, code := ajRun(t, crossed)
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: journey step open reads as requirement 0001 yet names 0002; name the doc whose behavior the step is, or reword the step",
				"architecture/05-user-journey.md:19: journey step done reads as requirement 0002 yet names 0001")
			straight := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagDone, "done[3 Save the record<br/>Req 0002]"), "-->|saves| done", "-->|saves 0001| done")
			out, code = ajNoNote(t, straight, "a step naming its own doc must add no note")
			sees(t, out, code, 0)
			grouped := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "subgraph open[2 Open the app<br/>Req 0002]\nend"), ajTagDone, "subgraph done[3 Save the record<br/>Req 0001]\nend")
			ajNoNote(t, grouped, "a layout subgraph is no step to judge")
			loose := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0002)"), "-->|saves| done", "-->|saves 0001| done")
			out, code = ajRun(t, loose)
			showsInfo(t, out, code, ajCheck, "architecture/05-user-journey.md:18: journey step open reads as requirement 0001 yet names 0002; check that the step names the doc whose behavior it is")
			// Sharing a title word with a doc it names, the step is never judged.
			both := ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0001, 0002)")
			out, code = ajNoNote(t, both, "a step naming the doc it reads as must add no note")
			sees(t, out, code, 0)
		})
	})
	t.Run("0001.1.6 no views or an unmoving scoped run passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "857dd667", func(a req.A) {
			out, code := ajRun(t, map[string]string{"src/main.go": "package main\n"})
			sees(a.T(), out, code, 0)
			passedWithNoFiles(a.T(), out, ajCheck)
		})
		req.Then(t, "0001.1.2.6.2", "fbf33601", func(a req.A) {
			for _, changed := range []string{"src/store/store_test.go", "CHANGELOG.md", "docs/guide.md"} {
				out, code := ajPullRequest(t, ajSecondDoc(ajFixture()), map[string]string{changed: "// Changed.\n"})
				sees(a.T(), out, code, 0)
				passedWithNoFiles(a.T(), out, ajCheck)
			}
		})
	})
	t.Run("0001.1.7 a missing doc stays reported when another file changes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.7.1", "b391d24a", func(_ req.A) {
			repo := example(t, "happyRepo", ajSecondDoc(ajFixture()))
			for _, changed := range []string{"src/store/store.go", ajDoc, ajView} {
				out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, ajCheck)
				sees(t, out, code, 1, ajMissingDoc("0002"))
			}
			for _, changed := range []string{"src/store/store.go", ajView} {
				files := ajSecondDoc(ajFixture())
				out, code := ajPullRequest(t, files, map[string]string{changed: files[changed] + "\n"})
				sees(t, out, code, 1, ajMissingDoc("0002"))
			}
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 regions holding their listed docs pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.1.1", "07b8e753", func(_ req.A) {
			out, code := ajRun(t, ajRegion(acFixture(acRoot), "0001", ajNode1))
			sees(t, out, code, 0)
			two := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", ajNode1))
			two = ajEdit(two, ajEnd, ajEnd+"    %% journey-nodes:start 0002\n    req0002[(3 Open<br/>Req 0002)]:::system\n"+ajEnd)
			out, code = ajRun(t, two)
			sees(t, out, code, 0)
			// A view with no markers stays hand-written.
			out, code = ajRun(t, ajFixture())
			sees(t, out, code, 0)
		})
	})
	t.Run("0001.2.2 a region differing from its docs fails naming the fix", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.2.1", "74540c91", func(_ req.A) {
			stale := ajSecondDoc(ajRegion(acFixture(acRoot), "0001 0002", ajNode1))
			out, code := ajRun(t, stale)
			sees(t, out, code, 1, "the generated journey nodes differ from the requirement docs; run `fitness architecture-journey --fix` to rewrite them")
			retitled := ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open The App<br/>Req 0001)]:::system")
			out, code = ajRun(t, retitled)
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: the generated journey nodes differ")
			swapped := ajEdit(ajSecondDoc(ajRegion(acFixture(acRoot), "0001", ajNode1)), "%% journey-nodes:start 0001", "%% journey-nodes:start 0002")
			out, code = ajRun(t, swapped)
			sees(t, out, code, 1, "the generated journey nodes differ")
		})
	})
	t.Run("0001.2.3 a doc listed in no region or two, or an unknown ID, fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.3.1", "5b016b35", func(a req.A) {
			unlisted := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", ajNode1))
			out, code := ajRun(t, unlisted)
			sees(a.T(), out, code, 1, "requirement doc 0002 is listed in no journey-nodes region")
			unknown := ajRegion(acFixture(acRoot), "0001 0009", ajNode1)
			out, code = ajRun(t, unknown)
			sees(a.T(), out, code, 1, "lists 0009, which is no requirement doc")
		})
		req.Then(t, "0001.2.1.3.2", "97c5032f", func(a req.A) {
			twice := ajEdit(ajRegion(acFixture(acRoot), "0001", ajNode1), ajEnd, ajEnd+"    %% journey-nodes:start 0001\n"+ajNode1+"\n"+ajEnd)
			out, code := ajRun(t, twice)
			sees(a.T(), out, code, 1, "requirement doc 0001 is listed in two journey-nodes regions")
		})
	})
	t.Run("0001.2.4 broken markers fail at their line", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.4.1", "df3bd461", func(_ req.A) {
			open := ajRegion(acFixture(acRoot), "0001", "")
			open[ajView] = strings.Replace(open[ajView], ajEnd, "", 1)
			out, code := ajRun(t, open)
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: journey view has `%% journey-nodes:start` with no end marker")
			nested := ajEdit(ajRegion(acFixture(acRoot), "0001", ajNode1), ajNode1, "    %% journey-nodes:start 0001\n"+ajNode1)
			out, code = ajRun(t, nested)
			sees(t, out, code, 1, "architecture/05-user-journey.md:19: journey view has `%% journey-nodes:start` before the last region's end marker")
			stray := ajEdit(ajFixture(), "user --> open", "%% journey-nodes:end\nuser --> open")
			out, code = ajRun(t, stray)
			sees(t, out, code, 1, "architecture/05-user-journey.md:20: journey view has `%% journey-nodes:end` with no start marker")
		})
	})
	t.Run("0001.2.5 fix rewrites every region and nothing else", func(t *testing.T) {
		t.Parallel()
		files := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", "    stale"))
		files = ajEdit(files, ajEnd, ajEnd+"    %% journey-nodes:start 0002\n    stale\n"+ajEnd)
		repo := example(t, "happyRepo", files)
		out, code := fitness(t, repo, nil, ajCheck, "--fix")
		want := strings.Replace(files[ajView], "    stale\n    %% journey-nodes:end\n    %% journey-nodes:start 0002\n    stale", "req0001[(2 Open<br/>Req 0001)]:::system\n    %% journey-nodes:end\n    %% journey-nodes:start 0002\n    req0002[(3 Open<br/>Req 0002)]:::system", 1)
		req.Then(t, "0001.2.1.5.1", "9d08f140", func(a req.A) {
			sees(a.T(), out, code, 0)
			body, err := os.ReadFile(filepath.Join(repo, ajView))
			a.NoError(err)
			a.Contains(string(body), "req0001[(2 Open<br/>Req 0001)]:::system")
			a.Contains(string(body), "req0002[(3 Open<br/>Req 0002)]:::system")
		})
		req.Then(t, "0001.2.1.5.2", "2340c439", func(a req.A) {
			ajViewIs(t, repo, want)
		})
	})
	t.Run("0001.2.6 fix writes nothing while a doc lacks its title", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.6.1", "5f85a67e", func(_ req.A) {
			files := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", "    stale"))
			files[ajDoc2] = strings.Replace(files[ajDoc2], "# 0002 Open\n", "", 1)
			repo := example(t, "happyRepo", files)
			out, code := fitness(t, repo, nil, ajCheck, "--fix")
			showsInfo(t, out, code, ajCheck, "docs/requirements/0002-save.md: needs `# NNNN Title` before the region can be written")
			ajViewIs(t, repo, files[ajView])
		})
	})
}

const (
	ajCheck   = "architecture-journey"
	ajView    = acRoot + acJourney
	ajDoc     = "docs/requirements/0001-open.md"
	ajDoc2    = "docs/requirements/0002-save.md"
	ajOpen    = "open(2 Open the app)"
	ajDone    = "done[3 Record saved]"
	ajTagOpen = "open(2 Open the app<br/>Req 0001)"
	ajTagDone = "done[3 Record saved<br/>Req 0001]"
	ajNode1   = "    req0001[(2 Open<br/>Req 0001)]:::system"
	ajEnd     = "    %% journey-nodes:end\n"
	ajMissing = "architecture/05-user-journey.md:15: requirement %s is missing from the journey diagram"
	ajTable   = "\n## Requirements\n\n| Requirement | Journey |\n| --- | --- |\n| [0001](../docs/requirements/0001-open.md) | `open`, `done` |\n"
)

// ajFixture is the contract fixture whose journey steps both name
// requirement 0001 in their labels. The diagram fence is line 15, the open
// step line 18, and the done step line 19.
func ajFixture() map[string]string {
	return ajEdit(ajEdit(acFixture(acRoot), ajOpen, ajTagOpen), ajDone, ajTagDone)
}

// ajSecondDoc adds requirement 0002, shaped like 0001.
func ajSecondDoc(files map[string]string) map[string]string {
	files[ajDoc2] = strings.ReplaceAll(files[ajDoc], "0001", "0002")
	return files
}

// ajRun runs `fitness-install -- architecture-journey` on happyRepo with
// files written over it.
func ajRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, ajCheck)
}

// ajNoNote runs ajRun and asserts no step reads as another doc, failing
// with why otherwise.
func ajNoNote(t *testing.T, files map[string]string, why string) (string, int) {
	t.Helper()
	out, code := ajRun(t, files)
	if strings.Contains(out, "reads as requirement") {
		t.Errorf("%s:\n%s", why, out)
	}
	return out, code
}

// ajEdit replaces old with new in the journey view of a fixture.
func ajEdit(files map[string]string, old, new string) map[string]string {
	files[ajView] = strings.Replace(files[ajView], old, new, 1)
	return files
}

// ajRegion replaces the fixture's two steps with a marked region listing
// ids and holding nodes, wired by the same edges.
func ajRegion(files map[string]string, ids, nodes string) map[string]string {
	files = ajEdit(files, ajOpen+":::action\ndone", "%% journey-nodes:start "+ids+"\n"+nodes+"\n"+ajEnd+"done")
	files = ajEdit(files, ajDone+":::screen\n", "")
	files = ajEdit(files, "user --> open\nopen -->|saves| done", "user --> req0001\nreq0001 --> req0001")
	return files
}

// ajPullRequest runs the check as a pull request into main that commits
// changed over happyRepo with files written over it.
func ajPullRequest(t *testing.T, files, changed map[string]string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", files)
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	acceptance.Commit(t, repo, changed, "fix(app): change a file")
	return fitness(t, repo, acceptance.PullRequestEnv, ajCheck)
}

// ajMissingDoc is the error for requirement id missing from the journey.
func ajMissingDoc(id string) string {
	return strings.Replace(ajMissing, "%s", id, 1)
}

// ajViewIs asserts the journey view in repo reads want.
func ajViewIs(t *testing.T, repo, want string) {
	t.Helper()
	req.Asserted(t)
	body, err := os.ReadFile(filepath.Join(repo, ajView))
	if err != nil || string(body) != want {
		t.Errorf("journey view after --fix: %v\n%s\nwant:\n%s", err, body, want)
	}
}
