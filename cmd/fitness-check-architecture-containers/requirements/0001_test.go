// Package requirements proves the architecture-containers check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of every repo the tests make.
func TestMain(m *testing.M) { acceptance.Main(m) }

// The shared architecture fixtures, named as the root requirement tests name
// them, so the cases moved from there read unchanged.
const (
	accCheck     = acceptance.AccCheck
	accDevTools  = acceptance.AccDevTools
	accJavaMain  = acceptance.AccJavaMain
	accLibrary   = acceptance.AccLibrary
	accNoOwner   = acceptance.AccNoOwner
	acContainer  = acceptance.AcContainer
	accPackage   = acceptance.AccPackage
	accProject   = acceptance.AccProject
	accService   = acceptance.AccService
	acRoot       = acceptance.AcRoot
	accContainer = acRoot + acContainer
)

var (
	accPbxproj        = acceptance.AccPbxproj
	accEdit           = acceptance.AccEdit
	accJavaApp        = acceptance.AccJavaApp
	accKit            = acceptance.AccKit
	accNative         = acceptance.AccNative
	accRootCommand    = acceptance.AccRootCommand
	accRun            = acceptance.AccRun
	accXcode          = acceptance.AccXcode
	acEdit            = acceptance.AcEdit
	acFixture         = acceptance.AcFixture
	javaPom           = acceptance.JavaPom
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
	pbxproj           = acceptance.Pbxproj
	sees              = acceptance.Sees
	with              = acceptance.With
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it passes counting the container views", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, accXcode([2]string{"App", "application"}, [2]string{"Widget", "app-extension"}))
		sees(t, out, code, 0)
		passedWith(t, out, accCheck, "1")
	})
	t.Run("0001.1.2.1.2 a `*` or `.` path or its name works too", func(t *testing.T) {
		t.Parallel()
		checks := accEdit(with(acFixture(acRoot), map[string]string{
			"cmd/fitness-check-gofmt/main.go": "package main\n",
			"cmd/fitness-check-vet/main.go":   "package main\n",
		}), "}\nRel", "Container(checks, \"5 Checks\", \"Go\", \"cmd/fitness-check-*\")\n}\nRel")
		out, code := accRun(t, checks)
		sees(t, out, code, 0)
		// A Java build at the repo root beside the service, owned by the root
		// path alone.
		root := accEdit(with(acFixture(acRoot), map[string]string{
			"pom.xml":                  javaPom(""),
			"src/main/java/Start.java": "public class Start {\n  public static void main(String[] args) {}\n}\n",
		}), "}\nRel", "Container(start, \"5 Start\", \"Java\", \".\")\n}\nRel")
		out, code = accRun(t, root)
		sees(t, out, code, 0)
		out, code = accRun(t, accEdit(with(acFixture(acRoot), accRootCommand), "}\nRel", "Container(tool, \"5 Tool\", \"Go\")\n}\nRel"))
		sees(t, out, code, 0)
		out, code = accRun(t, accEdit(with(acFixture(acRoot), accJavaApp), "}\nRel", "Container(api, \"5 Api\", \"Java\")\n}\nRel"))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.2.1 it fails naming its entry point or target or deploy file", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			files map[string]string
			want  string
		}{
			{map[string]string{"src/cmd/api/main.go": "package main\n"}, "src/cmd/api/main.go:1: Go command api" + accNoOwner},
			{map[string]string{accPackage: accService}, "package.json:4: Node service portal" + accNoOwner},
			{map[string]string{accPackage: "{\n  \"name\": \"portal\",\n  \"bin\": \"cli.js\"\n}\n"}, "package.json:3: Node service portal" + accNoOwner},
			{map[string]string{accPackage: accDevTools, "Dockerfile": "FROM node\n"}, "Dockerfile:1: Node service portal" + accNoOwner},
			{accRootCommand, "main.go:1: Go command tool" + accNoOwner},
			{accJavaApp, "api/src/main/java/Api.java:4: Java application api" + accNoOwner},
			{with(accJavaApp, map[string]string{"pom.xml": javaPom(""), accJavaMain: "public class Api {\n  public static void main(\n      String[] args) {}\n}\n"}), accJavaMain + ":2: Java application api" + accNoOwner},
		}
		for _, c := range cases {
			out, code := accRun(t, with(acFixture(acRoot), c.files))
			sees(t, out, code, 1, c.want)
		}
		out, code := accRun(t, accEdit(accNative(), "Container(notify, \"5 Notify\", \"Swift\")\n", ""))
		sees(t, out, code, 1, "project.yml:17: extension target Notify"+accNoOwner)
	})
	t.Run("0001.1.2.2.2 two containers claiming it fail naming both claims", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, accEdit(acFixture(acRoot), "}\nRel", "Container(api, \"5 Api\", \"Go\", \"src/cmd\")\n}\nRel"))
		sees(t, out, code, 1, "src/cmd/service/main.go:1: Go command service is claimed by containers service at "+accContainer+":7 and api at "+accContainer+":8")
	})
	t.Run("0001.1.2.3.1 it fails naming the view line", func(t *testing.T) {
		t.Parallel()
		want := "architecture/02-containers.md:8: container widget matches no build target, entry point, or deploy file by name"
		out, code := accRun(t, accXcode([2]string{"App", "application"}))
		sees(t, out, code, 1, want)
		files := accNative()
		files[accProject] = strings.Replace(files[accProject], "  Widget:\n    templates: [Extension]\n    sources: # folders\n      - path: Widget\n        excludes:\n          - \"**/*.md\"\n", "", 1)
		out, code = accRun(t, files)
		sees(t, out, code, 1, want)
	})
	t.Run("0001.1.2.3.2 paths holding only source layers fail the same way", func(t *testing.T) {
		t.Parallel()
		files := accEdit(acFixture(acRoot), "}\nRel", "Container(store, \"5 Store\", \"Go\", \"src/store\")\n}\nRel")
		out, code := accRun(t, files)
		sees(t, out, code, 1, "architecture/02-containers.md:8: container store claims only source layers (src/store): no build target, entry point, or deploy file is under them; show them in the component view")
		files["src/store/store.go"] = "package main\n"
		out, code = accRun(t, files)
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.4.1 it fails naming the view line", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, acEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Node\", \"src/cmd/service\""))
		sees(t, out, code, 1, "architecture/02-containers.md:7: container service says \"Node\"")
		out, code = accRun(t, acEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Go 1.26\", \"src/cmd/service\""))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.4.2 it fails naming the product", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, acEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Node\", \"src/cmd/service\""))
		sees(t, out, code, 1, "but Go command service at src/cmd/service/main.go:1 is Go")
	})
	t.Run("0001.1.2.5.1 it fails naming the store", func(t *testing.T) {
		t.Parallel()
		files := accEdit(acFixture(acRoot), "}\nRel", "ContainerDb(db, \"5 Records\", \"SQLite\")\n}\nRel")
		out, code := accRun(t, files)
		sees(t, out, code, 1, "architecture/02-containers.md:8: store db names no schema, configuration, or source path")
		out, code = accRun(t, accEdit(files, "\"SQLite\")", "\"SQLite\", \"src/store\")"))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.5.2 an external system naming no platform service fails too", func(t *testing.T) {
		t.Parallel()
		files := with(accNative(), map[string]string{"App/Sync.swift": "import CloudKit\n"})
		out, code := accRun(t, files)
		sees(t, out, code, 1, "App/Sync.swift:1: CloudKit is used here but no external system in the containers view names it")
		out, code = accRun(t, accEdit(files, "}\nRel", "}\nSystem_Ext(icloud, \"7 Sync\", \"Apple\", \"CloudKit records\")\nRel"))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2.6.1 it fails naming that line", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, with(acFixture(acRoot), map[string]string{accPackage: "{\n  \"name\": \"portal\",\n"}))
		sees(t, out, code, 1, "package.json:1: cannot read package.json")
	})
	t.Run("0001.1.2.6.2 a target lacking a type or source folder fails too", func(t *testing.T) {
		t.Parallel()
		files := accNative()
		files[accProject] = strings.Replace(files[accProject], "    templates: [Extension]\n    sources: Notify\n", "    sources: Notify\n", 1)
		out, code := accRun(t, files)
		sees(t, out, code, 1, "project.yml:17: target Notify has no type")
		files = accNative()
		files[accProject] = strings.Replace(files[accProject], "sources: Notify", "sources: Missing", 1)
		out, code = accRun(t, files)
		sees(t, out, code, 1, "project.yml:17: target Notify names a source Missing that does not exist")
		out, code = accRun(t, accXcode([2]string{"App", "application"}, [2]string{"Widget", ""}))
		sees(t, out, code, 1, "App.xcodeproj/project.pbxproj:11: target Widget has no type")
	})
	library := func() map[string]string {
		return with(acFixture(acRoot), map[string]string{accPackage: accLibrary, "lib/index.ts": "export const kit = 1;\n"})
	}
	t.Run("0001.1.2.7.1 it owns the package ignoring trailing punctuation", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, library())
		sees(t, out, code, 1, "package.json:3: Node library @acme/kit"+accNoOwner)
		out, code = accRun(t, accKit(library(), "Kit", "Builds dist/esm/index.js."))
		sees(t, out, code, 0)
		out, code = accRun(t, accKit(library(), "Kit", "Wraps @acme/kit in dist/esm"))
		sees(t, out, code, 1, "package.json:3: Node library @acme/kit"+accNoOwner, "architecture/02-containers.md:8: container kit matches no build target, entry point, or deploy file by name")
	})
	t.Run("0001.1.2.7.2 a label or exact entry point word owns it too", func(t *testing.T) {
		t.Parallel()
		for _, kit := range [][2]string{{"Kit dist/esm", "Bundles"}, {"@acme/kit", "Bundles"}} {
			out, code := accRun(t, accKit(library(), kit[0], kit[1]))
			sees(t, out, code, 0)
		}
		for _, tools := range []string{
			"{\n  \"name\": \"tools\",\n  \"private\": true,\n  \"main\": \"lib/index.ts\"\n}\n",
			"{\n  \"name\": \"tools\",\n  \"main\": \"index.js\",\n  \"scripts\": {\n    \"test\": \"vitest\"\n  }\n}\n",
		} {
			out, code := accRun(t, with(library(), map[string]string{accPackage: tools}))
			sees(t, out, code, 0)
		}
		tracked := with(library(), map[string]string{accPackage: "{\n  \"name\": \"kit\",\n  \"exports\": \"./lib/index.ts\"\n}\n"})
		out, code := accRun(t, tracked)
		sees(t, out, code, 1, "package.json:3: Node library kit"+accNoOwner)
		out, code = accRun(t, accKit(tracked, "Kit", "Loads lib/index.ts"))
		sees(t, out, code, 0)
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.2.1.1 it passes with zero files", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, map[string]string{"src/main.go": "package main\n"})
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, accCheck)
	})
	t.Run("0001.2.2.2.1 it proves no product", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, accNative())
		sees(t, out, code, 0)
		passedWith(t, out, accCheck, "1")
		out, code = accRun(t, acFixture(acRoot))
		sees(t, out, code, 0)
		generated := accEdit(accNative(), "\"3 App\"", "\"3 Main\"")
		generated[accPbxproj] = pbxproj([2]string{"App", "application"})
		out, code = accRun(t, generated)
		sees(t, out, code, 0)
		out, code = accRun(t, with(acFixture(acRoot), map[string]string{accPackage: accDevTools, "src/gen/gen.go": "//go:build ignore\n\npackage main\n"}))
		sees(t, out, code, 0)
		out, code = accRun(t, with(acFixture(acRoot), with(accJavaApp, map[string]string{accJavaMain: "public class Api {\n  private static void main(String[] args) {}\n\n  static void main(String name) {}\n}\n"})))
		sees(t, out, code, 0)
	})
	t.Run("0001.2.2.2.2 a test folder or `scripts/` file proves none", func(t *testing.T) {
		t.Parallel()
		out, code := accRun(t, with(acFixture(acRoot), map[string]string{"test/fake/main.go": "package main\n", "scripts/gen/main.go": "package main\n"}))
		sees(t, out, code, 0)
	})
	t.Run("0001.2.2.3.1 it still reads that manifest", func(t *testing.T) {
		t.Parallel()
		files := accXcode([2]string{"App", "application"}, [2]string{"Widget", "app-extension"})
		files[".gitattributes"] = "*.pbxproj linguist-generated\n"
		out, code := accRun(t, files)
		sees(t, out, code, 0)
		files = accNative()
		files[".gitattributes"] = "project.yml linguist-generated\n"
		out, code = accRun(t, files)
		sees(t, out, code, 0)
		files = accEdit(accNative(), "Container(notify, \"5 Notify\", \"Swift\")\n", "")
		files[".gitattributes"] = "project.yml linguist-generated\n"
		out, code = accRun(t, files)
		sees(t, out, code, 1, "project.yml:17: extension target Notify"+accNoOwner)
	})
	t.Run("0001.2.2.4.1 it passes with zero files", func(t *testing.T) {
		t.Parallel()
		repo := accScoped(t, map[string]string{
			"src/store/store_test.go": "package store\n\nfunc helper() {}\n",
			"changelog/unreleased.md": "- Adds a helper.\n",
			"docs/guide.md":           "# Guide\n\nHow to run it.\n",
		})
		out, code := acceptance.Fitness(t, repo, acceptance.PullRequestEnv, accCheck)
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, accCheck)
	})
	t.Run("0001.2.2.5.1 it reads every view", func(t *testing.T) {
		t.Parallel()
		for _, change := range []map[string]string{
			{"src/cmd/service/main.go": "package main\n\nfunc main() { run() }\n\nfunc run() {}\n"},
			{accContainer: strings.Replace(acceptance.AcEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Node\", \"src/cmd/service\"")[accContainer], "Runs requests.", "Serves requests.", 1)},
		} {
			repo := accScoped(t, change)
			out, code := acceptance.Fitness(t, repo, acceptance.PullRequestEnv, accCheck)
			sees(t, out, code, 1, accContainer+":7: container service says \"Node\" but Go command service at src/cmd/service/main.go:1 is Go")
		}
	})
	t.Run("0001.2.2.5.2 it reads every product", func(t *testing.T) {
		t.Parallel()
		// A product left unowned on main stays reported when another file moves.
		unowned := with(acFixture(acRoot), map[string]string{"src/cmd/api/main.go": "package main\n"})
		for _, change := range []map[string]string{
			{"src/store/store.go": "package store\n\nfunc Save() {}\n"},
			{accContainer: strings.Replace(unowned[accContainer], "Runs requests.", "Serves requests.", 1)},
		} {
			out, code := acceptance.Fitness(t, accScoped(t, change, unowned), acceptance.PullRequestEnv, accCheck)
			sees(t, out, code, 1, "src/cmd/api/main.go:1: Go command api"+accNoOwner)
		}
	})
}

// accScoped makes base, or the fixture with the service's technology wrong,
// as main. It commits change on top, as a pull request would.
func accScoped(t *testing.T, change map[string]string, base ...map[string]string) string {
	t.Helper()
	files := acEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Node\", \"src/cmd/service\"")
	if len(base) > 0 {
		files = base[0]
	}
	repo := acceptance.Example(t, "happyRepo", files)
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	acceptance.Commit(t, repo, change, "fix(app): change it")
	return repo
}
