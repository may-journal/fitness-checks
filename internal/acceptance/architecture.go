package acceptance

import (
	"strconv"
	"strings"
	"testing"
)

// The architecture fixtures several packages' requirement tests share. They
// hold the five views, the runtimes built on them, and a requirement doc.

const (
	AcCheck     = "architecture-contract"
	AcRoot      = "architecture/"
	AcContext   = "01-system-context.md"
	AcContainer = "02-containers.md"
	AcComponent = "03-components.md"
	AcCode      = "04-code.md"
	AcJourney   = "05-user-journey.md"
	AcEvidence  = "#L3)"
)

// AcView writes a view as people already do, with a heading, a diagram and its
// callout table, and nothing added for the check. The diagram's Nth line is line
// N+3, and the first row sits on line count+8.
func AcView(diagram string, rows ...string) string {
	body := "# View\n\n```mermaid\n" + diagram + "\n```\n\n| # | Description | Why |\n| --- | --- | --- |\n"
	for i, row := range rows {
		body += "| " + strconv.Itoa(i+1) + " | " + row + " | Shared by later checks. |\n"
	}
	return body
}

// AcFixture is a small complete architecture under root. up leads from the
// root back to the repository, for links.
func AcFixture(root string) map[string]string {
	up := strings.Repeat("../", strings.Count(root, "/"))
	return map[string]string{
		root + AcContext: AcView("C4Context\nPerson(user, \"1 User\")\nSystem(app, \"2 App\")\nRel(user, app, \"3\")",
			"A person who uses the app.", "The system.", "Uses it.") +
			"\n```mermaid\nC4Context\ntitle Legend\nPerson(p, \"Person\")\n```\n",
		root + AcContainer: AcView("C4Container\nPerson(user, \"1 User\")\nSystem_Boundary(app, \"2 App\") {\nContainer(service, \"3 Service\", \"Go\", \"src/cmd/service\", $tags=\"v1\")\n}\nRel(user, service, \"4\")\nUpdateElementStyle(user, $fontColor=\"red\")",
			"Same person.", "The system.", "Runs requests.", "Calls it."),
		root + AcComponent: AcView("C4Component\nContainer_Boundary(service, \"1 Service (src/cmd/service)\") {\nComponent(entry, \"2 Entry\", \"Go\")\nComponent(store, \"3 Store\", \"Go\", \"src/store\")\n}\nRel(entry, store, \"4\")",
			"The service.", "`src/cmd/service/main.go` starts it.", "Keeps records.", "Saves through the store, see [source]("+up+"src/cmd/service/main.go"+AcEvidence+"."),
		root + AcCode: AcView("classDiagram\nclass Record[\"1 Record\"] {\n+ID\n}\nclass save[\"2 save()\"]\nsave --> Record : 3",
			"A saved record.", "Writes one record.", "Save returns the record.") +
			"\n```mermaid\nclassDiagram\nclass Example\nParent <|-- Child : inherits\n```\n",
		root + AcJourney: "## Legend\n\n```mermaid\nflowchart TB\nperson((1 Person)):::persona\nclassDef persona fill:#fff\n```\n\n| # | Description | Why |\n| --- | --- | --- |\n| 1 | Persona | A legend row. |\n\n" +
			AcView("flowchart TB\nuser((1 User)):::persona\nopen(2 Open the app):::action\ndone[3 Record saved]:::screen\nuser --> open\nopen -->|saves| done\nclassDef persona fill:#fff",
				"Same person.", "Opens it, `0001.1`.", "Shows the record, see [0001]("+up+"docs/requirements/0001-open.md)."),
		"src/cmd/service/main.go": "package main\n\nfunc main() {}\n",
		"src/store/store.go":      "package store\n",
		"src/store/store_test.go": "package store\n",
		"tools/check.sh":          "#!/bin/sh\nexit 0\n",
		AcDoc:                     AcDocWith("User opens the app"),
	}
}

// AcEdit replaces old with new in one view of the default fixture.
func AcEdit(view, old, new string) map[string]string {
	files := AcFixture(AcRoot)
	files[AcRoot+view] = strings.Replace(files[AcRoot+view], old, new, 1)
	return files
}

const (
	AccCheck    = "architecture-containers"
	AccNoOwner  = " has no container in the containers view"
	AccProject  = "project.yml"
	AccPackage  = "package.json"
	AccPbxproj  = "App.xcodeproj/project.pbxproj"
	AccDevTools = "{\n  \"name\": \"portal\",\n  \"private\": true,\n  \"bin\": null,\n  \"devDependencies\": {\n    \"prettier\": \"3.0.0\"\n  }\n}\n"
	AccService  = "{\n  \"name\": \"portal\",\n  \"scripts\": {\n    \"start\": \"node server.js\"\n  }\n}\n"
)

// AccRootCommand is a Go command at the repository root, which go build
// names after the module path's last element.
var AccRootCommand = map[string]string{"go.mod": "module example.com/tool\n\ngo 1.26\n", "main.go": "package main\n"}

// AccNative is a native app with Node tools whose XcodeGen project puts the
// Notify target on line 17 of project.yml.
func AccNative() map[string]string {
	files := AcFixture(AcRoot)
	files["src/cmd/service/main.go"] = "package service\n\n// Service is a library.\n"
	files[AcRoot+AcContainer] = AcView("C4Container\nPerson(user, \"1 User\")\nSystem_Boundary(app, \"2 App\") {\nContainer(main, \"3 App\", \"Swift 6, SwiftUI\", \"App/\")\nContainer(widget, \"4 Widget\", \"Swift\")\nContainer(notify, \"5 Notify\", \"Swift\")\n}\nRel(user, main, \"6\")",
		"Same person.", "The app.", "The app target.", "The widget target.", "The notification target.", "Opens it.")
	files[AccProject] = "name: App\ntargetTemplates:\n  Extension:\n    type: app-extension\n    sources: [\"${target_name}\"]\ntargets:\n  App:\n    type: application # the app\n    platform: iOS\n    sources: [App]\n  Widget:\n    templates: [Extension]\n    sources: # folders\n      - path: Widget\n        excludes:\n          - \"**/*.md\"\n  Notify:\n    templates: [Extension]\n    sources: Notify\n  AppTests:\n    type: bundle.unit-test\n    platform: iOS\n    sources: [Tests]\n"
	files["App/App.swift"] = "import SwiftUI\n"
	files["Widget/Widget.swift"] = "import WidgetKit\n"
	files["Notify/Notify.swift"] = "import UserNotifications\n"
	files["Tests/AppTests.swift"] = "import XCTest\n"
	files[AccPackage] = AccDevTools
	return files
}

// AccXcode is the native app with an Xcode project instead of XcodeGen,
// declaring the app and the widget.
func AccXcode(targets ...[2]string) map[string]string {
	files := AccNative()
	delete(files, AccProject)
	files[AcRoot+AcContainer] = strings.Replace(files[AcRoot+AcContainer], "Container(notify, \"5 Notify\", \"Swift\")\n", "", 1)
	files[AccPbxproj] = Pbxproj(targets...)
	return files
}

// Pbxproj renders native targets as Xcode writes them; an empty type leaves
// the productType line out.
func Pbxproj(targets ...[2]string) string {
	body := "// !$*UTF8*$!\n{\n\tobjects = {\n"
	for _, t := range targets {
		body += "\t\tA /* " + t[0] + " */ = {\n\t\t\tisa = PBXNativeTarget;\n\t\t\tname = " + t[0] + ";\n"
		if t[1] != "" {
			body += "\t\t\tproductType = \"com.apple.product-type." + t[1] + "\";\n"
		}
		body += "\t\t};\n"
	}
	return body + "\t};\n}\n"
}

// AccRun runs `fitness-install -- architecture-containers` on happyRepo with
// files written over it.
func AccRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return Fitness(t, Example(t, "happyRepo", files), env, AccCheck)
}

// AccEdit replaces old with new in the containers view of files.
func AccEdit(files map[string]string, old, new string) map[string]string {
	files[AcRoot+AcContainer] = strings.Replace(files[AcRoot+AcContainer], old, new, 1)
	return files
}

// AccJavaMain is the product source of AccJavaApp.
const AccJavaMain = "api/src/main/java/Api.java"

// AccJavaApp is a Maven project in api/ whose product source declares a
// main method below a commented-out one.
var AccJavaApp = map[string]string{
	"api/pom.xml":                  JavaPom(""),
	AccJavaMain:                    "// static void main(String[] args) proves nothing here.\npublic class Api {\n\n  public static void main(String[] args) {}\n}\n",
	"api/src/test/java/ApiIT.java": "class ApiIT {\n  public static void main(String[] args) {}\n}\n",
}

// AccLibrary is an npm library shipping only built files, with main on line
// 3 and a build script making them.
const AccLibrary = "{\n  \"name\": \"@acme/kit\",\n  \"main\": \"dist/cjs/index.js\",\n  \"module\": \"dist/esm/index.js\",\n  \"exports\": {\n    \".\": {\n      \"types\": \"./dist/types/index.d.ts\",\n      \"import\": \"./dist/esm/index.js\"\n    }\n  },\n  \"files\": [\"dist\"],\n  \"scripts\": {\n    \"build\": \"rollup -c\"\n  }\n}\n"

// AccKit adds a kit container with label and description to the library.
func AccKit(files map[string]string, label, description string) map[string]string {
	return AccEdit(With(files, nil), "}\nRel", "Container(kit, \"5 "+label+"\", \"TypeScript\", \""+description+"\")\n}\nRel")
}

// AcDoc is the fixture's requirement doc.
const AcDoc = "docs/requirements/0001-open.md"

// AcDocWith is requirement 0001 with one acceptance whose Given is given.
// The Given sits on line 16.
func AcDocWith(given string) string {
	return "# 0001 Open\n\n## Why\n\nI open the app.\n\n## Measurement\n\nopens\n-\ntries\n\n## Requirements\n\n- 0001.1\n    - Given " + given + "\n        - When it starts\n            - Then it opens\n"
}

// AcChained is the fixture doc whose 0001.1 chains two steps, so `0001.1.2`
// cites its second.
var AcChained = strings.Replace(AcDocWith("User opens the app"), "        - When it starts\n            - Then it opens\n",
	"        1. When it starts\n            - Then it opens\n        2. When it opens\n            - Then it shows the record\n", 1)

const (
	AoCheck    = "architecture-components"
	AoView     = AcRoot + AcComponent
	AoMain     = "src/cmd/service/main.go"
	AoStore    = "src/store/store.go"
	AoHome     = "App/Screens/HomeScreen.swift"
	AoEdge     = "Rel(entry, store, \"4\")"
	AoSaves    = "Saves through the store, see [source](../src/cmd/service/main.go#L3)."
	AoEntry    = "Component(entry, \"2 Entry\", \"Go\")\n"
	AoModel    = "Component(model, \"3 Model\", \"Swift\", \"App/Model\")\n"
	AoNoOwner  = ":1: no component in the components view names this source"
	AoNotDrawn = ", but " + AoView + " draws no "
	AoImports  = "package main\n\nimport \"example.com/app/src/store\"\n\nfunc main() { store.Save() }\n"
	aoScreen   = "import SwiftUI\n\nstruct HomeScreen {\n    let store = Store()\n}\n"
	AoDrawn    = "Rel(ui, model, \"8\")"
	AoFrames   = ".claude/skills/build-log/scripts/Frames.swift"
	AoSeed     = ".codex/scripts/Seed.swift"
)

// AoTooling adds component 9, owning the agent tooling folder, to the
// native app's components view, before its drawn edge.
func AoTooling(component, row string) map[string]string {
	files := AoEdit(AoNative(), AoDrawn, component+"\n"+AoDrawn)
	files[AoView] += "| 9 | " + row + " | Agent tooling meets the product bar. |\n"
	return files
}

// AoComponents is the components view starting on line 4, with inside added to
// the service boundary and after and rows appended.
func AoComponents(inside, after string, rows ...string) string {
	diagram := "C4Component\nContainer_Boundary(service, \"1 Service (src/cmd/service)\") {\n" + AoEntry + "Component(store, \"3 Store\", \"Go\", \"src/store\")\n" + inside + "}\n" + AoEdge + after
	return AcView(diagram, append([]string{"The service.", "`src/cmd/service` starts it.", "Keeps records.", AoSaves}, rows...)...)
}

// AoFixture is the Go architecture whose entry imports the store, so the
// drawn edge and its evidence on line 3 both hold.
func AoFixture() map[string]string {
	return With(AcFixture(AcRoot), map[string]string{
		"go.mod": "module example.com/app\n\ngo 1.26\n",
		AoMain:   AoImports,
		AoStore:  "package store\n\n// Save writes one record.\nfunc Save() {}\n",
		AoView:   AoComponents("", ""),
	})
}

// AoNative is the native app with screens reading a model, a widget, and
// a notification extension, each under its container's boundary. The
// diagram starts on line 4; the model sits on line 7.
func AoNative() map[string]string {
	files := AccNative()
	for _, gone := range []string{"App/App.swift", AoMain, AoStore, "src/store/store_test.go"} {
		files[gone] = ""
	}
	files[AoHome] = aoScreen
	files["App/Model/Store.swift"] = "import Foundation\n\nfinal class Store {\n    var records: [String] = []\n}\n"
	files[AoView] = AcView("C4Component\nContainer_Boundary(main, \"1 App (App/)\") {\nComponent(ui, \"2 Screens\", \"Swift\", \"App/Screens\")\n"+AoModel+"}\nContainer_Boundary(widget, \"4 Widget\") {\nComponent(glance, \"5 Glance\", \"Swift\", \"Widget/\")\n}\nContainer_Boundary(notify, \"6 Notify\") {\nComponent(alert, \"7 Alert\", \"Swift\", \"Notify/\")\n}\nRel(ui, model, \"8\")",
		"The app.", "Screens.", "The model.", "The widget.", "Draws the glance.", "The notification extension.", "Posts alerts.", "Screens read the model.")
	return files
}

// AoRun runs `fitness-install -- architecture-components` on happyRepo with
// files written over it.
func AoRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return Fitness(t, Example(t, "happyRepo", files), env, AoCheck)
}

// AoEdit replaces old with new in the components view of files.
func AoEdit(files map[string]string, old, new string) map[string]string {
	files[AoView] = strings.Replace(files[AoView], old, new, 1)
	return files
}

// AoWorker adds a worker command whose job reads the store, its boundary
// opening on line 10, or 11 once the store moves.
func AoWorker() map[string]string {
	files := AoFixture()
	files[AcRoot+AcContainer] = strings.Replace(files[AcRoot+AcContainer], "}\nRel", "Container(worker, \"5 Worker\", \"Go\", \"src/cmd/worker\")\n}\nRel", 1)
	files[AoView] = AoComponents("", "\nContainer_Boundary(worker, \"5 Worker (src/cmd/worker)\") {\nComponent(job, \"6 Job\", \"Go\", \"src/cmd/worker\")\n}\nRel(job, store, \"7\")", "The worker.", "Runs jobs.", "Reads records.")
	files["src/cmd/worker/main.go"] = AoImports
	return files
}

// AoResolved is the native app with a sync component beside the model, each
// declaring a load overload and save extension.
func AoResolved(line7 string) map[string]string {
	files := AoEdit(AoNative(), AoModel+"}", AoModel+"Component(sync, \"9 Sync\", \"Swift\", \"App/Sync\")\n}")
	files = AoEdit(files, "Rel(ui, model, \"8\")", "Rel(ui, model, \"8\")\nRel(sync, model, \"10\")")
	files[AoView] += "| 9 | Syncs records. | Shared by later checks. |\n| 10 | Sync fills the store. | Shared by later checks. |\n"
	files["App/Model/Load.swift"] = "func load(from path: String) -> Store {\n    Store()\n}\n"
	files["App/Model/Saving.swift"] = "#if os(iOS)\nextension Store {\n#else\nextension Store: Sendable {\n#endif\n    static func < (lhs: Store, rhs: Store) -> Bool { true }\n" +
		"    func save(_ record: String) {\n        if let kept = [record].first(where: { !$0.isEmpty }) { save(all: [kept]) }\n    }\n}\n"
	files["App/Sync/Load.swift"] = "func load(into store: Store) {\n    store.records = []\n}\n"
	files["App/Sync/Saving.swift"] = "extension Store {\n    func save(all records: [String]) {}\n}\n"
	files[AoHome] = "import SwiftUI\n\nstruct HomeScreen {\n    func refresh(save: Bool) {\n        let store = load(from: \"records\")\n        store.save(\"first\")\n" + line7 + "    }\n}\n"
	return files
}

// AoScripts is a package whose scripts run tools past runner flags, beside a
// start script, a published entry, folder and exports pattern.
const AoScripts = "{\n  \"main\": \"tools/build.mjs\",\n  \"files\": [\"lib\"],\n  \"exports\": {\n    \"./x/*\": \"./pub/*.mjs\"\n  },\n  \"scripts\": {\n" +
	"    \"analyze\": \"node scripts/analyze.mjs\",\n    \"gen\": \"npx tsx tools/gen.ts && echo done\",\n    \"seed\": \"MODE=dev node -r dotenv/config seed.mjs\",\n" +
	"    \"build\": \"node tools/build.mjs\",\n    \"start\": \"node --require dotenv/config scripts/server.mjs\",\n    \"dev\": \"node --watch scripts/server.mjs\",\n" +
	"    \"check\": \"node lib/index.js --check\",\n    \"pub\": \"node pub/a.mjs\"\n  }\n}\n"

// ArchitectureChecks are the six checks that read the architecture model.
var ArchitectureChecks = []string{"architecture-actors", "architecture-code", "architecture-components", "architecture-containers", "architecture-contract", "architecture-journey"}

// With returns base with overrides applied.
func With(base map[string]string, overrides map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return out
}

// JavaPom is a Maven build that declares plugins.
func JavaPom(plugins string) string {
	return "<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId>" +
		"<artifactId>app</artifactId><version>1</version><build><plugins>" + plugins + "</plugins></build></project>\n"
}
