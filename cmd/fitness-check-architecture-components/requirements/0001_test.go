// Package requirements proves the architecture-components check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain keeps the host's git config out of every repo the tests make.
func TestMain(m *testing.M) { acceptance.Main(m) }

// The shared architecture fixtures, named as the root requirement tests name
// them, so the cases moved from there read unchanged.
const (
	acCheck    = acceptance.AcCheck
	aoCheck    = acceptance.AoCheck
	aoDrawn    = acceptance.AoDrawn
	aoEdge     = acceptance.AoEdge
	aoEntry    = acceptance.AoEntry
	aoFrames   = acceptance.AoFrames
	aoHome     = acceptance.AoHome
	aoImports  = acceptance.AoImports
	aoMain     = acceptance.AoMain
	aoModel    = acceptance.AoModel
	aoNoOwner  = acceptance.AoNoOwner
	aoNotDrawn = acceptance.AoNotDrawn
	aoSaves    = acceptance.AoSaves
	aoScripts  = acceptance.AoScripts
	aoSeed     = acceptance.AoSeed
	aoStore    = acceptance.AoStore
	aoView     = acceptance.AoView
)

var (
	aoComponents      = acceptance.AoComponents
	aoEdit            = acceptance.AoEdit
	aoFixture         = acceptance.AoFixture
	aoNative          = acceptance.AoNative
	aoResolved        = acceptance.AoResolved
	aoRun             = acceptance.AoRun
	aoTooling         = acceptance.AoTooling
	aoWorker          = acceptance.AoWorker
	example           = acceptance.Example
	fitness           = acceptance.Fitness
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
	sees              = acceptance.Sees
	with              = acceptance.With
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a product source with no component owner fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "70819bcf", func(a req.A) {
			out, code := aoRun(t, aoUnowned())
			sees(a.T(), out, code, 1, "src/api/api.go"+aoNoOwner, "src/api/api.config.ts"+aoNoOwner, "index.mjs"+aoNoOwner)
		})
		req.Then(t, "0001.1.2.1.2", "1cddce84", func(a req.A) {
			tooling := aoNative()
			tooling[aoFrames] = "import Foundation\n\nlet frameRate = 30\n"
			out, code := aoRun(t, tooling)
			sees(a.T(), out, code, 1, aoFrames+aoNoOwner)
			tooling = aoTooling("Component(frames, \"9 Build log frames\", \"Swift\", \".claude/skills/build-log\")", "Captures build log frames.")
			tooling[aoFrames] = "import Foundation\n\nlet frameRate = 30\n"
			out, code = aoRun(t, tooling)
			sees(a.T(), out, code, 0)
		})
	})
	t.Run("0001.1.2 tools need no owner but started files do", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "8d80efba", func(a req.A) {
			out, _ := aoRun(t, aoUnowned())
			aoNeedNoOwner(a, out, "rollup.config.mjs", ".eslintrc")
		})
		req.Then(t, "0001.1.2.2.2", "909d1b23", func(a req.A) {
			product := []string{"tools/build.mjs", "scripts/server.mjs", "lib/index.js", "pub/a.mjs"}
			files := with(aoFixture(), map[string]string{
				"package.json":         aoScripts,
				"scripts/analyze.mjs":  "export const size = 1;\n",
				"scripts/lib/stats.ts": "export const mean = 0;\n",
				"tools/gen.ts":         "export const gen = 1;\n",
				"seed.mjs":             "export const rows = [];\n",
			})
			var want []string
			for _, file := range product {
				files[file] = "export const run = 1;\n"
				want = append(want, file+aoNoOwner)
			}
			out, code := aoRun(t, files)
			sees(a.T(), out, code, 1, want...)
			aoNeedNoOwner(a, out, "scripts/analyze", "scripts/lib", "tools/gen.ts", "seed.mjs")
			files = aoEdit(files, "\"src/store\")", "\"src/store\", \"and "+strings.Join(product, " ")+"\")")
			out, code = aoRun(t, files)
			sees(a.T(), out, code, 0)
			passedWith(a.T(), out, aoCheck, "1")
		})
	})
	t.Run("0001.1.3 a code span or label claims a root source by bare name", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "7d9df1a0", func(_ req.A) {
			files := aoEdit(aoUnowned(), "\"src/store\")", "\"src/store\", \"and src/api\")")
			files = aoEdit(files, aoEntry, "Component(entry, \"2 Entry\", \"Go\", \"Loads index.mjs, per package.json\")\n")
			files = aoEdit(files, "| `src/cmd/service` starts it. |", "| `src/cmd/service` starts it, per `package.json`. |")
			files = aoEdit(files, "| Keeps records. |", "| Keeps records, `index.mjs` and `missing.mjs` export them, per `package.json`. |")
			out, code := aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
			out, code = fitness(t, example(t, "happyRepo", files), nil, acCheck)
			sees(t, out, code, 0)
			files["src/api/api.go"], files["src/api/api.config.ts"] = "", ""
			out, code = fitness(t, example(t, "happyRepo", files), nil, acCheck)
			sees(t, out, code, 1, aoView+":7: unresolved source path src/api")
		})
	})
	t.Run("0001.1.4 one source claimed by two components fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.4.1", "1c83ba4c", func(_ req.A) {
			second := "\n```mermaid\nC4Component\nContainer_Boundary(jobs, \"5 Jobs\") {\nComponent(cache, \"6 Cache\", \"Go\", \"src/store\")\n}\n```\n\n| # | Description | Why |\n| --- | --- | --- |\n| 5 | Jobs. | Later. |\n| 6 | Cache. | Later. |\n"
			files := aoFixture()
			files[aoView] += second
			out, code := aoRun(t, files)
			sees(t, out, code, 1, aoStore+":1: components store at "+aoView+":7 and cache at "+aoView+":22 both name this source")
			files = aoEdit(files, "\"Go\", \"src/store\")\n}\n```", "\"Go\", \"src/cache\")\n}\n```")
			files["src/cache/cache.go"] = "package cache\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
		})
	})
	t.Run("0001.1.5 a reference missing from the view fails naming the view", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "7ee31d5c", func(a req.A) {
			out, code := aoRun(t, aoFixture())
			sees(a.T(), out, code, 0)
			passedWith(a.T(), out, aoCheck, "1")
			want := aoMain + ":3: entry imports src/store, owned by store" + aoNotDrawn + "entry -> store edge"
			out, code = aoRun(t, aoEdit(aoFixture(), aoEdge, ""))
			sees(a.T(), out, code, 1, want)
			out, code = aoRun(t, aoNative())
			sees(a.T(), out, code, 0)
			out, code = aoRun(t, aoEdit(aoNative(), "Rel(ui, model, \"8\")", ""))
			sees(a.T(), out, code, 1, aoHome+":4: ui names Store, owned by model"+aoNotDrawn+"ui -> model edge")
			tooling := aoTooling("Component(seed, \"9 Seed\", \"Swift\", \".codex/scripts\")", "Seeds sample records.")
			tooling[aoSeed] = "import Foundation\n\nlet sample = Store()\n"
			out, code = aoRun(t, tooling)
			sees(a.T(), out, code, 1, aoSeed+":3: seed names Store, owned by model"+aoNotDrawn+"seed -> model edge")
			tooling = aoEdit(tooling, aoDrawn, aoDrawn+"\nRel(seed, model, \"10\")")
			tooling[aoView] += "| 10 | Seed fills the store. | Agent tooling meets the product bar. |\n"
			out, code = aoRun(t, tooling)
			sees(a.T(), out, code, 0)
			// An edge from the enclosing boundary covers entry, and an edge no
			// reference uses is allowed.
			boundary := aoEdit(aoFixture(), aoEdge, "Rel(service, store, \"4\")")
			boundary = aoEdit(boundary, aoSaves, "Saves through the store.")
			out, code = aoRun(t, boundary)
			sees(a.T(), out, code, 0)
			unused := aoEdit(aoFixture(), aoEdge, aoEdge+"\nRel(store, entry, \"5\")")
			unused[aoView] += "| 5 | Store calls back. | Allowed unused. |\n"
			out, code = aoRun(t, unused)
			sees(a.T(), out, code, 0)
		})
		req.Then(t, "0001.1.2.5.2", "e200e9d2", func(a req.A) {
			want := aoMain + ":3: entry imports src/store, owned by store" + aoNotDrawn + "entry -> store edge"
			out, code := aoRun(t, aoEdit(aoFixture(), aoEdge, "Rel(store, entry, \"4\")"))
			sees(a.T(), out, code, 1, want+"; "+aoView+":9 draws store -> entry")
		})
	})
	t.Run("0001.1.6 only a build condition keeps a hidden reference", func(t *testing.T) {
		t.Parallel()
		files := with(aoFixture(), map[string]string{
			"src/api/api.go":           "package api\n\nfunc Ping() {}\n",
			"src/cmd/service/linux.go": "//go:build linux\n\npackage main\n\nimport \"example.com/app/src/api\"\n\nfunc ping() { api.Ping() }\n",
		})
		api := "Component(api, \"5 API\", \"Go\", \"src/api\")\n"
		files[aoView] = aoComponents(api, "", "Answers pings.")
		linuxOut, linuxCode := aoRun(t, files)
		files[aoView] = aoComponents(api, "\nRel(entry, api, \"6\")", "Answers pings.", "Pings the API.")
		drawnOut, drawnCode := aoRun(t, files)
		native := aoEdit(aoNative(), "Rel(ui, model, \"8\")", "")
		native[aoHome] = "import SwiftUI\n\nstruct HomeScreen {\n#if DEBUG // tools\n    let store = Store()\n#endif\n}\n"
		debugOut, debugCode := aoRun(t, native)
		native[aoHome] = "import SwiftUI\n\nstruct HomeScreen {\n#if DEBUG\n    let debug = true\n#endif // DEBUG\n    let store = Store()\n}\n"
		afterOut, afterCode := aoRun(t, native)
		native[aoView] = aoNative()[aoView]
		fixedOut, fixedCode := aoRun(t, native)
		hidden := aoEdit(aoNative(), "Rel(ui, model, \"8\")", "")
		hidden[aoHome] = "import SwiftUI\n\nstruct HomeScreen {\n    // let store = Store()\n    let name = \"Store()\"\n}\n"
		hidden["App/Screens/Generated.swift"] = "let cached = Store()\n"
		hidden[".gitattributes"] = "App/Screens/Generated.swift linguist-generated\n"
		generatedOut, generatedCode := aoRun(t, hidden)
		delete(hidden, ".gitattributes")
		shownOut, shownCode := aoRun(t, hidden)
		req.Then(t, "0001.1.2.6.1", "80f69079", func(a req.A) {
			sees(a.T(), generatedOut, generatedCode, 0)
			sees(a.T(), shownOut, shownCode, 1, "App/Screens/Generated.swift:1: ui names Store, owned by model"+aoNotDrawn+"ui -> model edge")
			a.True(!strings.Contains(shownOut, aoHome), "a comment or string made a reference:\n%s", shownOut)
		})
		req.Then(t, "0001.1.2.6.2", "8b1f4e7a", func(a req.A) {
			sees(a.T(), drawnOut, drawnCode, 0)
			sees(a.T(), afterOut, afterCode, 1, aoHome+":7: ui names Store, owned by model"+aoNotDrawn+"ui -> model edge")
			sees(a.T(), fixedOut, fixedCode, 0)
		})
		req.Then(t, "0001.1.2.6.3", "df75b1bc", func(a req.A) {
			sees(a.T(), linuxOut, linuxCode, 1, "src/cmd/service/linux.go:5: entry imports src/api under go:build linux, owned by api"+aoNotDrawn+"entry -> api edge")
			sees(a.T(), debugOut, debugCode, 1, aoHome+":5: ui names Store under #if DEBUG, owned by model"+aoNotDrawn+"ui -> model edge")
		})
	})
	t.Run("0001.1.7 an unparseable product source fails naming the source it cannot parse", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.7.1", "5c79e016", func(_ req.A) {
			out, code := aoRun(t, with(aoFixture(), map[string]string{aoStore: "package store\n\nfunc Save( {}\n"}))
			sees(t, out, code, 1, aoStore+":1: cannot parse "+aoStore+":3:12: expected ')', found '{'")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 no architecture views pass with zero files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.2.1.1", "91ecdad2", func(_ req.A) {
			out, code := aoRun(t, map[string]string{"src/main.go": "package main\n"})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, aoCheck)
		})
	})
	t.Run("0001.2.2 swift names reach the right declaration", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.2.2.1", "49fb0c32", func(_ req.A) {
			out, code := aoRun(t, aoResolved("        if save { store.save(all: []) }\n"))
			sees(t, out, code, 1, aoHome+":7: ui names Store.save(all:), owned by sync"+aoNotDrawn+"ui -> sync edge",
				"App/Model/Saving.swift:8: model names Store.save(all:), owned by sync"+aoNotDrawn+"model -> sync edge")
			out, code = aoRun(t, aoResolved("        if save { load(into: store) }\n"))
			sees(t, out, code, 1, aoHome+":7: ui names load(into:), owned by sync"+aoNotDrawn+"ui -> sync edge")
			files := aoResolved("        func load(into store: Store) {}\n        if save { load(into: store) }\n")
			files["App/Model/Saving.swift"] = "extension Store {\n    func save(_ record: String) {}\n}\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
		})
	})
	t.Run("0001.2.3 edge evidence must come from the first component", func(t *testing.T) {
		t.Parallel()
		row := aoView + ":17: "
		cases := []struct{ old, new, want string }{
			{"cmd/service/main.go#L3)", "store/store.go#L1)", "evidence src/store/store.go#L1 is owned by store, not entry"},
			{"#L3)", "#L1)", "evidence src/cmd/service/main.go#L1 names nothing store owns"},
			{"../src/cmd/service/main.go#L3)", "../tools/check.sh)", "unsupported evidence tools/check.sh; Go and Swift sources only"},
		}
		req.Then(t, "0001.2.2.3.1", "504812e9", func(a req.A) {
			for _, c := range []int{0, 2} {
				out, code := aoRun(t, aoEdit(aoFixture(), cases[c].old, cases[c].new))
				sees(a.T(), out, code, 1, row+cases[c].want)
			}
		})
		req.Then(t, "0001.2.2.3.2", "cc1bdc0c", func(a req.A) {
			out, code := aoRun(t, aoEdit(aoFixture(), cases[1].old, cases[1].new))
			sees(a.T(), out, code, 1, row+cases[1].want)
			files := aoFixture()
			files[aoView] = aoComponents("", "\nSystem_Ext(mailer, \"5 Mailer\", \"SMTP relay\")\nRel(entry, mailer, \"6\")", "Relays mail.", "Sends through the relay, see [source](../src/cmd/service/main.go#L7).")
			files[aoMain] = aoImports + "\nfunc send(c rune) bool { return c == '\"' } // smtp relay\n"
			out, code = aoRun(t, files)
			sees(a.T(), out, code, 1, aoView+":21: evidence src/cmd/service/main.go#L7 does not name Mailer outside comments")
			files[aoMain] = aoImports + "var body = `raw\n` // smtp relay\n"
			out, code = aoRun(t, files)
			sees(a.T(), out, code, 1, aoView+":21: evidence src/cmd/service/main.go#L7 does not name Mailer outside comments")
			files[aoMain] = aoImports + "\nfunc send() string { return \"https://smtp\" }\n"
			out, code = aoRun(t, files)
			sees(a.T(), out, code, 0)
			renamed := aoEdit(aoFixture(), "#L3)", "#L5)")
			renamed[aoStore] = "package records\n\nfunc Save() {}\n"
			renamed["src/store/tools.go"] = "//go:build ignore\n\npackage main\n"
			renamed[aoMain] = strings.Replace(aoImports, "store.Save()", "records.Save()", 1)
			out, code = aoRun(t, renamed)
			sees(a.T(), out, code, 0)
			renamed[aoMain] = strings.Replace(aoImports, "store.Save()", "_ = 1", 1)
			out, code = aoRun(t, renamed)
			sees(a.T(), out, code, 1, row+"evidence src/cmd/service/main.go#L5 names nothing store owns")
		})
	})
	t.Run("0001.2.4 a component reached by another runtime fails", func(t *testing.T) {
		t.Parallel()
		shared := func() map[string]string {
			return aoEdit(aoWorker(), "Component(store, \"3 Store\", \"Go\", \"src/store\")\n}", "}\nSystem_Boundary(shared, \"Shared\") {\nComponent(store, \"3 Store\", \"Go\", \"src/store\")\n}")
		}
		files := aoEdit(aoNative(), "Rel(ui, model, \"8\")", "Rel(ui, model, \"8\")\nRel(glance, model, \"9\")")
		files["Widget/Widget.swift"] = "import WidgetKit\n\nlet shared = Store()\n"
		req.Then(t, "0001.2.2.4.1", "4bbe20e3", func(a req.A) {
			out, code := aoRun(t, aoWorker())
			sees(a.T(), out, code, 1, "src/cmd/worker/main.go:3: Go command worker reaches store here, but store is under container service at "+aoView+":5, not worker; a component two runtimes share belongs outside every container boundary")
			moved := aoEdit(aoEdit(shared(), aoEntry, ""), "Component(job,", aoEntry+"Component(job,")
			out, code = aoRun(t, moved)
			sees(a.T(), out, code, 1, aoMain+":1: Go command service is built from entry, which is under container worker at "+aoView+":11, not service")
			out, code = aoRun(t, files)
			sees(a.T(), out, code, 1, "Widget/Widget.swift:3: extension target Widget reaches model here, but model is under container main at "+aoView+":5, not widget")
		})
		req.Then(t, "0001.2.2.4.2", "5e48c625", func(a req.A) {
			out, code := aoRun(t, shared())
			sees(a.T(), out, code, 0)
			passedWith(a.T(), out, aoCheck, "1")
			out, code = aoRun(t, aoEdit(files, aoModel+"}", "}\nSystem_Boundary(shared, \"Shared\") {\n"+aoModel+"}"))
			sees(a.T(), out, code, 0)
		})
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 only tests or changelog or docs moving passes with zero files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.3.1.1.1", "91ecdad2", func(_ req.A) {
			repo := aoScoped(t, map[string]string{
				"src/store/store_test.go": "package store\n\nfunc helper() {}\n",
				"changelog/unreleased.md": "- Adds a helper.\n",
				"docs/guide.md":           "# Guide\n\nHow to run it.\n",
			})
			out, code := fitness(t, repo, acceptance.PullRequestEnv, aoCheck)
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, aoCheck)
		})
	})
	t.Run("0001.3.2 a view or a product source moving reads the whole model and every product source", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.3.1.2.1", "34ab9dea", func(a req.A) {
			for _, change := range []map[string]string{
				{aoStore: "package store\n\n// Save writes one record.\nfunc Save() {}\n\n// Load reads one record.\nfunc Load() {}\n"},
				{aoView: strings.Replace(aoEdit(aoFixture(), aoEdge, "")[aoView], "Keeps records.", "Keeps every record.", 1)},
			} {
				out, code := fitness(t, aoScoped(t, change), acceptance.PullRequestEnv, aoCheck)
				sees(a.T(), out, code, 1, aoMain+":3: entry imports src/store, owned by store"+aoNotDrawn+"entry -> store edge")
			}
		})
		req.Then(t, "0001.3.1.2.2", "7b5a492f", func(a req.A) {
			// A source left unowned on main stays reported when another file moves.
			unowned := with(aoFixture(), map[string]string{"src/api/api.go": "package api\n"})
			for _, change := range []map[string]string{
				{aoStore: "package store\n\n// Save writes one record.\nfunc Save() {}\n\n// Load reads one record.\nfunc Load() {}\n"},
				{aoView: strings.Replace(unowned[aoView], "Keeps records.", "Keeps every record.", 1)},
			} {
				out, code := fitness(t, aoScoped(t, change, unowned), acceptance.PullRequestEnv, aoCheck)
				sees(a.T(), out, code, 1, "src/api/api.go"+aoNoOwner)
			}
		})
	})
}

// aoUnowned is the Go fixture with an API folder, a root module and root
// tool configs that no component claims.
func aoUnowned() map[string]string {
	return with(aoFixture(), map[string]string{
		"src/api/api.go":        "package api\n",
		"src/api/api.config.ts": "export const port = 8080;\n",
		"index.mjs":             "export const ready = true;\n",
		"rollup.config.mjs":     "export default { input: \"index.mjs\" };\n",
		".eslintrc.cjs":         "module.exports = {};\n",
		"package.json":          "{}\n",
	})
}

// aoNeedNoOwner asserts out names none of the tools, since a tool needs no
// owner.
func aoNeedNoOwner(a req.A, out string, tools ...string) {
	for _, tool := range tools {
		a.True(!strings.Contains(out, tool), "%s needs no owner:\n%s", tool, out)
	}
}

// aoScoped makes base, or the Go fixture with its drawn edge removed, as
// main. It commits change on top, as a pull request would.
func aoScoped(t *testing.T, change map[string]string, base ...map[string]string) string {
	t.Helper()
	files := aoEdit(aoFixture(), aoEdge, "")
	if len(base) > 0 {
		files = base[0]
	}
	repo := example(t, "happyRepo", files)
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	acceptance.Commit(t, repo, change, "fix(app): change it")
	return repo
}
