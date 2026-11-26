package requirements

import (
	"fmt"
	"strings"
	"testing"
)

func Test0006_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"only the running test owns an ID", func(t *testing.T) {
			t.Parallel()
			out, _ := swiftJudged(t, nil)
			passedWith(t, out, "requirements", "2")
			out, _ = swiftJudged(t, map[string]string{
				"Tests/AdderTests/AdderTests.swift": "import Testing\n\nstruct AdderSpec {\n  @Test(\"adds\") func test0001_1() {}\n}\n",
			})
			passedWith(t, out, "requirements", "2")
			swiftOwnsNothing(t, xcTestFile()+"// func test0001_1() {}\n/* nested /* func test0001_1() {} */ */\nlet note = \"\\(\"func test0001_1() {}\")\"\nlet raw = #\"\"\"\n\" func test0001_1() {}\n\"\"\"#\n")
			swiftOwnsNothing(t, xcTestFile()+"\nstruct Helper {\n  func test0001_1() {}\n}\n")
			for _, owner := range []string{
				`test("0001.1 adds two numbers", () => {});`,
				`describe("adder", () => { it.each([1, 2])("0001.1 adds %i", (n) => { expect(n).toBeTruthy(); }); });`,
				"test.skip(`0001.1 adds two numbers`, () => {});",
				`test("0001.1 adds two numbers", () => { test.skip(process.env.CI === undefined, "0001.1 runs in CI"); });`,
			} {
				out, code := scriptChecked(t, append([]string{owner}, notTests...)...)
				sees(t, out, code, 0)
				passedWith(t, out, "requirements", "2")
			}
			jestFolder := map[string]string{"tests/adder.test.ts": "", "src/__tests__/adder.ts": vitestFile(`test("0001.1 adds two numbers", () => {});`)}
			out, code := checked(t, with(scriptRepo(), jestFolder))
			sees(t, out, code, 0)
			passedWith(t, out, "requirements", "2")
			out, code = scriptChecked(t, notTests...)
			sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
			junitOwnership(t)
		}},
		{"acceptance with no test or two tests fails naming the requirement", func(t *testing.T) {
			t.Parallel()
			swiftOwnsNothing(t, xcTestFile())
			out, code := checked(t, with(scriptRepo(), map[string]string{"tests/adder.test.ts": ""}))
			sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
			out, code = checked(t, with(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest")}))
			sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
		}},
		{"each chained step needs a named case", func(t *testing.T) {
			t.Parallel()
			for _, repo := range []map[string]string{scriptRepo(), swiftRepo()} {
				out, code := checked(t, with(repo, map[string]string{pairedDoc: docWith(chain("0001.1", 2))}))
				sees(t, out, code, 0)
			}
		}},
		{"test proving no requirement fails", func(t *testing.T) {
			t.Parallel()
			out, code := swiftJudged(t, map[string]string{swiftTest: swiftTestingFile("test0001_1", "testAddsAgain")})
			sees(t, out, code, 1, "testAddsAgain proves no requirement; rename it testNNNN_N")
			out, code = scriptChecked(t, `test("0001.1 adds two numbers", () => {});`, `it("adds again", () => {});`)
			sees(t, out, code, 1, `tests/adder.test.ts:4: it("adds again") proves no requirement; rename it test("NNNN.N …")`)
			out, code = checked(t, with(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest", "test0001_1", "addsAgain")}))
			sees(t, out, code, 1, "src/test/java/AdderTest.java:6: addsAgain proves no requirement; rename it testNNNN_N")
			unchangedTestFails(t, with(goRepo(), map[string]string{"adder_test.go": testFile("adder", "Test0001_1", "TestAddsAgain")}), "TestAddsAgain proves no requirement")
			for _, body := range []string{xcTestFile("testAddsAgain"), "import Testing\n@Test func addsAgain() {}\n"} {
				files := with(swiftRepo(), map[string]string{"Tests/ExtraTests.swift": body})
				unchangedTestFails(t, files, "proves no requirement; rename it testNNNN_N")
			}
			files := with(scriptRepo(), map[string]string{"tests/extra.spec.js": vitestFile(`it("adds again", () => {});`)})
			unchangedTestFails(t, files, `it("adds again") proves no requirement`)
			files = with(junitRepo(), map[string]string{"src/test/java/ExtraTest.java": junitFile("ExtraTest", "addsAgain")})
			unchangedTestFails(t, files, "addsAgain proves no requirement")
		}},
		{"code with no requirement docs fails", func(t *testing.T) {
			t.Parallel()
			for _, repo := range []map[string]string{goRepo(), swiftRepo(), scriptRepo(), junitRepo()} {
				out, code := checked(t, with(repo, map[string]string{"docs/requirements/0001-adds-numbers.md": ""}))
				sees(t, out, code, 1, "docs/requirements/ has no requirement docs")
			}
			out, code := checked(t, map[string]string{"src/app.js": "export const app = 1;\n"})
			sees(t, out, code, 1, "docs/requirements/ has no requirement docs")
		}},
		{"no covered code passes", func(t *testing.T) {
			t.Parallel()
			for _, files := range []map[string]string{nil, {"vite.config.ts": "export default {};\n", ".eslintrc.cjs": "module.exports = {};\n"}} {
				out, code := checked(t, files)
				sees(t, out, code, 0)
				passedWithNoFiles(t, out, "requirements")
			}
		}},
		{"Swift test on wrong framework fails", func(t *testing.T) {
			t.Parallel()
			out, code := swiftChecked(t, map[string]string{swiftTest: xcTestFile("test0001_1")})
			sees(t, out, code, 1, swiftTest+":1: imports XCTest; only an end-to-end test that imports UITestingLibrarySystem may use it, so "+libraryHint, swiftTest+":3: subclasses XCTestCase; "+libraryHint)
			out, code = swiftChecked(t, map[string]string{swiftTest: "import Testing\n@testable import XCTest\n\n@Test func test0001_1() {}\n"})
			sees(t, out, code, 1, swiftTest+":2: imports XCTest")
			// An end-to-end test drives the app through SystemSession, never XCUIApplication itself.
			out, code = swiftChecked(t, map[string]string{swiftTest: strings.Replace(systemSuite, "try session.taps", "XCUIApplication().launch()\n    try session.taps", 1)})
			sees(t, out, code, 1, "drives Xcode UI tests through XCUIApplication")
			out, code = swiftChecked(t, map[string]string{swiftTest: librarySuite("", "  let app = XCUIApplication()\n  app.launch()\n")})
			sees(t, out, code, 1, swiftTest+":5: drives Xcode UI tests through XCUIApplication, which launches the app in a separate process; "+libraryHint)
			boot := "  let shell = Process()\n  shell.arguments = [\"xcrun\", \"simctl\", \"boot\", \"booted\"]\n"
			out, code = swiftChecked(t, map[string]string{swiftTest: librarySuite("import Foundation\n", boot)})
			sees(t, out, code, 1, swiftTest+":7: boots a simulator through simctl outside the end-to-end mode; "+libraryHint)
			for _, framework := range []string{"ViewInspector", "SnapshotTesting", "KIF", "EarlGrey", "Hammer"} {
				out, code = swiftChecked(t, map[string]string{swiftTest: librarySuite("import "+framework+"\n", "")})
				sees(t, out, code, 1, swiftTest+":3: imports "+framework+", a UI test framework requirement tests may not use; "+libraryHint)
			}
			out, code = swiftChecked(t, map[string]string{swiftScreens: ""})
			sees(t, out, code, 1, swiftTest+": nothing in Tests/AdderTests/ imports UITestingLibrary; "+libraryHint)
			// A library import in a sibling folder does not reach the test.
			out, code = swiftChecked(t, map[string]string{swiftScreens: "", "Tests/OtherTests/Screens.swift": "import UITestingLibrary\n"})
			sees(t, out, code, 1, "nothing in Tests/AdderTests/ imports UITestingLibrary")
		}},
		{"Swift Testing on UITestingLibrary passes", func(t *testing.T) {
			t.Parallel()
			// Words in comments and strings are no import or call.
			quiet := "import Testing\n// import XCTest\n/* XCUIApplication() */\n@MainActor @Test(\"never shells out to simctl\") func test0001_1() async throws {\n  let note = \"XCUIApplication\"\n  _ = note\n}\n"
			for _, body := range []string{swiftTestingFile("test0001_1"), librarySuite("", ""), quiet} {
				out, code := swiftChecked(t, map[string]string{swiftTest: body})
				sees(t, out, code, 0)
				passedWith(t, out, "requirements", "2")
			}
			out, code := swiftChecked(t, map[string]string{swiftTest: systemSuite})
			sees(t, out, code, 0)
			passedWith(t, out, "requirements", "2")
		}},
		{"requirement proven by a check passes", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", provenRepo()), nil)
			sees(t, out, code, 0, "All 2 checks passed")
			passedWith(t, out, "req-0002", "1")
		}},
		{"requirement check without its doc fails", func(t *testing.T) {
			t.Parallel()
			orphan := with(provenRepo(), map[string]string{"fitness/req-0003/README.md": houseReadme, "fitness/req-0003/check": verdictScript(passing)})
			out, code := fitness(t, example(t, "happyRepo", orphan), nil, "requirements")
			sees(t, out, code, 1, "fitness/req-0003 proves requirement 0003, which no doc defines; write the doc or delete the check")
			out, code = stagedRun(t, provenRepo(), map[string]string{"fitness/req-0002/check": verdictScript(`{"ok": true, "errors": [], "filesChecked": 2}`)})
			sees(t, out, code, 1, "fitness/req-0002/check changed without its requirement doc docs/requirements/0002-adds-again.md; change both together")
		}},
		{"requirement check that judged no files fails", func(t *testing.T) {
			t.Parallel()
			empty := with(provenRepo(), map[string]string{"fitness/req-0002/check": verdictScript(`{"ok": true, "errors": [], "filesChecked": 0}`)})
			out, code := fitness(t, example(t, "happyRepo", empty), nil, "check", "req-0002")
			sees(t, out, code, 1, "req-0002 judged no files, so its scan is broken")
		}},
	})
}

func Test0006_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"test or doc without partner fails", func(t *testing.T) {
			t.Parallel()
			repo, env := hooked(t, pairedRepo())
			write(t, repo, map[string]string{pairedTest: edited(testFile("reqs", "Test0001_1"))})
			user(t, repo, env, "add", "-A")
			out, code := user(t, repo, env, "commit", "-m", "test(reqs): check sums twice")
			sees(t, out, code, 1, pairedTest+" changed without its requirement doc "+pairedDoc+"; change both together")
			out, code = prRun(t, pairedRepo(), map[string]string{pairedDoc: reworded(requirementDoc("0001.1"))}, "docs(requirements): reword 0001")
			sees(t, out, code, 1, pairedDoc+" changed without its requirement test "+pairedTest+"; change both together")
		}},
		{"test named for two requirements fails", func(t *testing.T) {
			t.Parallel()
			doc := "docs/requirements/0002-adds-again.md"
			test := map[string]string{twoTest: edited(twoRepo()[twoTest])}
			out, code := stagedRun(t, twoRepo(), test)
			sees(t, out, code, 1, twoTest+" changed without its requirement doc "+pairedDoc+" or "+doc)
			out, code = stagedRun(t, twoRepo(), with(test, map[string]string{doc: reworded(twoRepo()[doc])}))
			sees(t, out, code, 0)
		}},
		{"shared helper needs no partner", func(t *testing.T) {
			t.Parallel()
			out, code := stagedRun(t, pairedRepo(), map[string]string{"reqs/tools_test.go": edited(pairedRepo()["reqs/tools_test.go"])})
			sees(t, out, code, 0)
		}},
		{"merge or replay judges own changes", func(t *testing.T) {
			t.Parallel()
			keep := []string{"--no-edit"}
			for _, r := range [][][]string{
				{{"merge", "--no-edit", "side"}, keep, {"--amend", "-m", "chore(reqs): merge the reword"}},
				{{"cherry-pick", "side"}, keep},
				{{"revert", "--no-edit", "side"}, {"-m", "revert(reqs): undo the reword"}},
			} {
				out, code := replayed(t, r[0], r[1:]...)
				sees(t, out, code, 0)
			}
			out, code := fitness(t, mergedBase(t), prEnv, "requirements")
			sees(t, out, code, 0)
		}},
		{"amend judged since its parent", func(t *testing.T) {
			t.Parallel()
			repo, env := hooked(t, pairedRepo())
			write(t, repo, map[string]string{pairedTest: edited(testFile("reqs", "Test0001_1"))})
			userSteps(t, repo, env, []string{"add", "-A"}, []string{"commit", "-q", "--no-verify", "-m", "test(reqs): check sums twice"})
			write(t, repo, map[string]string{pairedDoc: reworded(requirementDoc("0001.1"))})
			user(t, repo, env, "add", "-A")
			out, code := user(t, repo, env, "commit", "--amend", "-m", "test(reqs): check sums twice")
			sees(t, out, code, 0, "test(reqs): check sums twice")
		}},
		{"fixing a failing test needs no doc", func(t *testing.T) {
			t.Parallel()
			failing := with(swiftRepo(), map[string]string{swiftTest: xcTestFile("test0001_1")})
			fixed := map[string]string{swiftTest: swiftTestingFile("test0001_1")}
			out, code := stagedRun(t, failing, fixed)
			sees(t, out, code, 0)
			out, code = prRun(t, failing, fixed, "test(adder): move to Swift Testing")
			sees(t, out, code, 0)
			repo, env := hooked(t, with(failing, map[string]string{".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n"}))
			write(t, repo, fixed)
			userSteps(t, repo, env, []string{"add", "-A"}, []string{"commit", "-q", "--no-verify", "-m", "test(adder): move to Swift Testing"})
			write(t, repo, map[string]string{swiftTest: edited(fixed[swiftTest])})
			user(t, repo, env, "add", "-A")
			out, code = user(t, repo, env, "commit", "--amend", "-m", "test(adder): move to Swift Testing")
			sees(t, out, code, 0, "test(adder): move to Swift Testing")
			out, code = stagedRun(t, swiftRepo(), map[string]string{swiftTest: edited(swiftTestingFile("test0001_1"))})
			sees(t, out, code, 1, swiftTest+" changed without its requirement doc "+pairedDoc+"; change both together")
		}},
		{"comment or format edits need no partner", func(t *testing.T) {
			t.Parallel()
			given := with(pairedRepo(), map[string]string{pairedDoc: unformatted.Replace(requirementDoc("0001.1"))})
			formatted := map[string]string{pairedDoc: requirementDoc("0001.1")}
			out, code := stagedRun(t, given, formatted)
			sees(t, out, code, 0)
			out, code = prRun(t, given, formatted, "style(docs): format requirement docs with Biome")
			sees(t, out, code, 0)
			out, code = stagedRun(t, given, with(formatted, map[string]string{pairedTest: edited(testFile("reqs", "Test0001_1"))}))
			sees(t, out, code, 0)
			sourced := with(pairedRepo(), map[string]string{pairedDoc: withSource(requirementDoc("0001.1"), "`adder.go:5`")})
			dropped := map[string]string{pairedDoc: requirementDoc("0001.1")}
			out, code = stagedRun(t, sourced, dropped)
			sees(t, out, code, 0)
			out, code = prRun(t, sourced, dropped, "docs(requirements): drop the Source line")
			sees(t, out, code, 0)
			out, code = stagedRun(t, sourced, map[string]string{pairedDoc: strings.Replace(requirementDoc("0001.1"), "sums asked for", "sums tried", 1)})
			sees(t, out, code, 1, pairedDoc+" changed without its requirement test")
			cited := map[string]string{pairedDoc: withSource(requirementDoc("0001.1"), "`adder.go:6`")}
			out, code = stagedRun(t, sourced, cited)
			sees(t, out, code, 1, pairedDoc+" changed without its requirement test")
			renamedPairs(t)
			commentedPairs(t)
		}},
		{"product code without tests and docs fails", func(t *testing.T) {
			t.Parallel()
			adder := map[string]string{"adder.go": rcAdder()}
			out, code := rcStaged(t, goRepo(), adder)
			sees(t, out, code, 1, "adder.go changed without a test and a requirement doc; change both with it")
			out, code = rcStaged(t, goRepo(), with(adder, map[string]string{"adder_test.go": edited(goRepo()["adder_test.go"])}))
			sees(t, out, code, 1, "adder.go changed without a requirement doc; change one with it")
			out, code = rcPull(t, goRepo(), adder)
			sees(t, out, code, 1, "adder.go changed without a test and a requirement doc")
			rcCodeWithComments(t)
		}},
		{"product code with test and doc passes", func(t *testing.T) {
			t.Parallel()
			out, code := rcStaged(t, goRepo(), map[string]string{
				"adder.go":      rcAdder(),
				"adder_test.go": edited(goRepo()["adder_test.go"]),
				pairedDoc:       reworded(goRepo()[pairedDoc]),
			})
			sees(t, out, code, 0, rcPassed)
		}},
		{"docs config or rename only passes", func(t *testing.T) {
			t.Parallel()
			adder := goRepo()["adder.go"]
			for _, files := range []map[string]string{
				{"README.md": readme + "\nMore words.\n"},
				{".fitnessrc.json": `{"proseBudget": {}}` + "\n"},
				{"CHANGELOG.md": "# Changelog\n\nNothing yet.\n"},
				{"adder.go": "", "sum.go": adder},
				{"adder.go": strings.Replace(adder, "\treturn a + b\n", "\n\treturn  a + b  \n", 1)},
			} {
				out, code := rcStaged(t, goRepo(), files)
				sees(t, out, code, 0, rcPassed)
			}
			rcCommentOnly(t)
		}},
		{"merge replaying product code passes", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", goRepo())
			git(t, repo, "checkout", "-q", "-b", "side")
			commit(t, repo, map[string]string{"adder.go": rcAdder()}, "feat(adder): double numbers")
			git(t, repo, "checkout", "-q", "-")
			commit(t, repo, map[string]string{"README.md": readme + "\nMore words.\n"}, "docs(readme): add words")
			git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "merge", "-q", "--no-commit", "--no-ff", "side")
			out, code := fitness(t, repo, nil, "requirements-changed")
			sees(t, out, code, 0, rcPassed)
		}},
		{"deleting unnamed product code passes", func(t *testing.T) {
			t.Parallel()
			given := with(goRepo(), map[string]string{"three.go": "package adder\n\n// Three is three.\nconst Three = 3\n"})
			out, code := rcStaged(t, given, map[string]string{"three.go": ""})
			sees(t, out, code, 0, rcPassed)
		}},
		{"deleting named product code fails", func(t *testing.T) {
			t.Parallel()
			named := strings.Replace(goRepo()[pairedDoc], "by hand.", "by hand, with `adder.go`.", 1)
			out, code := rcStaged(t, with(goRepo(), map[string]string{pairedDoc: named}), map[string]string{"adder.go": ""})
			sees(t, out, code, 1, "adder.go was deleted, leaving requirement doc "+pairedDoc+" and test adder_test.go unchanged; change or delete them with it")
		}},
		{"repo without requirement docs passes", func(t *testing.T) {
			t.Parallel()
			given := with(goRepo(), map[string]string{pairedDoc: ""})
			out, code := rcStaged(t, given, map[string]string{"adder.go": rcAdder()})
			sees(t, out, code, 0, rcPassed)
		}},
		{"requirements none trailer is refused", func(t *testing.T) {
			t.Parallel()
			repo, env := hooked(t, with(goRepo(), map[string]string{".fitnessrc.json": `{"policy": "external", "checks": ["requirements-changed"]}` + "\n"}))
			write(t, repo, map[string]string{"adder.go": rcAdder()})
			user(t, repo, env, "add", "-A")
			out, code := user(t, repo, env, "commit", "-m", "feat(adder): double numbers\n\nRequirements: none, a helper no user calls yet")
			sees(t, out, code, 1, "adder.go changed without a test and a requirement doc; change both with it")
		}},
	})
}

func Test0006_7(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{})
}

// judged runs the full suite on happyRepo plus files and returns the output.
func judged(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil)
}

// unformatted is requirementDoc before a formatter rewrote its quotes, wrapping,
// blank lines and nesting.
var unformatted = strings.NewReplacer(
	"two numbers without", "two numbers\nwithout",
	`['../../.fitnessrc.json']`, `["../../.fitnessrc.json"]`,
	"\n\n-\n\n", "\n-\n",
	"\n  - Given", "\n    - Given",
	"\n    - When", "\n        - When",
	"\n      - Then", "\n            - Then",
)

// commentedPairs proves a test or doc changed only in comments needs no
// partner change. A code change beside the comments still needs one.
func commentedPairs(t *testing.T) {
	t.Helper()
	test := testFile("reqs", "Test0001_1")
	rewritten := strings.Replace(withComment(test), "func Test0001_1", "// Test0001_1 proves the sum.\nfunc Test0001_1", 1)
	noted := map[string]string{pairedDoc: requirementDoc("0001.1") + "\n<!-- Sums stay small. -->\n"}
	for _, files := range []map[string]string{{pairedTest: rewritten}, noted} {
		out, code := stagedRun(t, pairedRepo(), files)
		sees(t, out, code, 0)
		out, code = prRun(t, pairedRepo(), files, "docs(reqs): reword comments")
		sees(t, out, code, 0)
	}
	out, code := stagedRun(t, pairedRepo(), map[string]string{pairedTest: edited(rewritten)})
	sees(t, out, code, 1, pairedTest+" changed without its requirement doc "+pairedDoc+"; change both together")
	out, code = prRun(t, pairedRepo(), map[string]string{pairedTest: edited(rewritten)}, "test(reqs): check sums twice")
	sees(t, out, code, 1, pairedTest+" changed without its requirement doc "+pairedDoc+"; change both together")
	out, code = stagedRun(t, pairedRepo(), with(noted, map[string]string{pairedTest: edited(rewritten)}))
	sees(t, out, code, 1, pairedTest+" changed without its requirement doc "+pairedDoc+"; change both together")
	out, code = stagedRun(t, pairedRepo(), map[string]string{pairedTest: rewritten, pairedDoc: reworded(requirementDoc("0001.1"))})
	sees(t, out, code, 1, pairedDoc+" changed without its requirement test "+pairedTest+"; change both together")
	out, code = stagedRun(t, pairedRepo(), map[string]string{"reqs/sums_test.go": rewritten, pairedTest: ""})
	sees(t, out, code, 0)
	shown := with(pairedRepo(), map[string]string{pairedDoc: requirementDoc("0001.1") + "\n```html\n<!-- a -->\n```\n"})
	out, code = stagedRun(t, shown, map[string]string{pairedDoc: requirementDoc("0001.1") + "\n```html\n<!-- b -->\n```\n"})
	sees(t, out, code, 1, pairedDoc+" changed without its requirement test")
}

// renamedPairs proves a renamed doc keeping its number and wording needs no test
// change, unlike one with new wording.
func renamedPairs(t *testing.T) {
	t.Helper()
	const renamed = "docs/requirements/0001-sums-numbers.md"
	moved := map[string]string{pairedDoc: "", renamed: requirementDoc("0001.1")}
	out, code := stagedRun(t, pairedRepo(), moved)
	sees(t, out, code, 0)
	out, code = prRun(t, pairedRepo(), moved, "docs(requirements): rename the sum requirement")
	sees(t, out, code, 0)
	out, code = stagedRun(t, pairedRepo(), map[string]string{pairedDoc: "", renamed: reworded(requirementDoc("0001.1"))})
	sees(t, out, code, 1, renamed+" changed without its requirement test")
}

// biomePackage is a JavaScript project that installs Biome at version %s.
const biomePackage = "{\"name\": \"app\", \"private\": true, \"devDependencies\": {\"@biomejs/biome\": \"%s\"}}\n"

// biomeFlags are the flags fitness passes Biome before the paths, with %s
// the cap that lifts its default of 20 diagnostics.
const biomeFlags = "--reporter=github --max-diagnostics=%s --no-errors-on-unmatched --files-ignore-unknown=true"

// biomeMajors pairs declared version, installed version and cap: Biome 1.8
// takes at most 65535, and `latest` hides the version.
var biomeMajors = [][3]string{{"2.4.15", "2.4.15", "none"}, {"^1.9.4", "1.9.4", "none"}, {"latest", "1.8.3", "65535"}}

// biomeTooOld is a Biome before 1.8, whose CLI rejects `--reporter=github`,
// declared under a spec that hides its version.
const biomeTooOld = "1.7.3"

// biomeRejects runs check on files with Biome 1.7.3 installed under
// `latest`, and expects it to fail naming that version without running Biome.
func biomeRejects(t *testing.T, files map[string]string, check string) {
	t.Helper()
	files["package.json"] = fmt.Sprintf(biomePackage, "latest")
	rp := standIn(t, map[string][]response{"biome": {biomeVersion(biomeTooOld)}})
	out, code := fitness(t, example(t, "happyRepo", files), rp.env, check)
	sees(t, out, code, 1, "Biome "+biomeTooOld+" found; Biome 1.8 or later is required")
	if rp.called("biome", "--reporter=github") {
		t.Errorf("Biome %s cannot take --reporter=github, yet it ran: %q", biomeTooOld, rp.calls("biome"))
	}
}

// biomeVersion is real `biome --version` output for version v.
func biomeVersion(v string) response {
	return response{Match: []string{"--version"}, Stdout: "Version: " + v + "\n"}
}

// biomeMany is more findings than Biome's default cap of 20 diagnostics.
const biomeMany = 25

// swiftTest is the path of swiftRepo's requirement test file, and
// swiftScreens of the helper that imports apple-ui-testing-library.
const (
	swiftTest    = "Tests/AdderTests/AdderTests.swift"
	swiftScreens = "Tests/AdderTests/Screens.swift"
)

// libraryHint is the part of every test library failure that says what to
// use instead.
const libraryHint = "write it with Swift Testing and UITestingLibrary, hosted in the app, as https://github.com/may-journal/apple-ui-testing-library shows"

// swiftChecked runs only the requirements check on swiftRepo with files
// written over it.
func swiftChecked(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return checked(t, with(swiftRepo(), files))
}

// librarySuite is a Swift Testing file whose test0001_1 runs body on
// apple-ui-testing-library, after the given extra imports.
func librarySuite(imports, body string) string {
	return "import Testing\nimport UITestingLibrary\n" + imports + "\n@MainActor @Test func test0001_1() async throws {\n" + body + "}\n"
}

// systemSuite is an end-to-end XCTest class whose test0001_1 drives the
// app through UITestingLibrarySystem.
const systemSuite = "import UITestingLibrary\nimport UITestingLibrarySystem\nimport XCTest\n\n@MainActor\nfinal class SystemRequirements: XCTestCase {\n" +
	"  func test0001_1() throws {\n    let session = SystemSession(bundleIdentifier: \"org.example.app\", timeout: 10)\n" +
	"    session.launch()\n    let reset = [\"xcrun\", \"simctl\", \"erase\", \"all\"]\n    XCTAssertFalse(reset.isEmpty)\n    try session.taps(.button(\"Share\"))\n  }\n}\n"

// swiftJudged runs fitness on swiftRepo with files written over it.
func swiftJudged(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return judged(t, with(swiftRepo(), files))
}

// swiftOwnsNothing is a test file whose test0001_1 is not a test.
func swiftOwnsNothing(t *testing.T, body string) {
	t.Helper()
	out, code := swiftJudged(t, map[string]string{"Tests/AdderTests/AdderTests.swift": body})
	sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
}

// checked runs only the requirements check on happyRepo plus files, so a
// TypeScript repo meets no JavaScript tool.
func checked(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "requirements")
}

// scriptChecked runs the requirements check on scriptRepo with its test file
// replaced by statements.
func scriptChecked(t *testing.T, statements ...string) (string, int) {
	t.Helper()
	return checked(t, with(scriptRepo(), map[string]string{"tests/adder.test.ts": vitestFile(statements...)}))
}

// notTests mention 0001.1 in ways vitest runs no test for.
var notTests = []string{
	`// test("0001.1 in a comment", () => {});`,
	`const note = "test('0001.1 in a string', () => {})";`,
	`describe("0001.1 a describe title", () => {});`,
	`class Matcher { test(value: string) { return /0001\.1/.test(value); } }`,
	`class Typed { test(value: string): boolean { return value === "0001.1"; } }`,
	`interface Check { test(value: string): boolean; }`,
}

// junitOwnership checks a JUnit test owns 0001.1 by its name and test
// annotation, and a comment, string or helper method does not.
func junitOwnership(t *testing.T) {
	t.Helper()
	for _, owner := range []string{
		"class AdderTest {\n  @ParameterizedTest\n  @ValueSource(ints = {1, 2})\n  void test0001_1(int n) {}\n}\n",
		"class AdderTest {\n  @TestFactory\n  Stream<DynamicTest> test0001_1() {\n    return Stream.empty();\n  }\n}\n",
		// TestNG runs a marked class's public methods, but not its setup.
		"@Test\npublic class AdderTest {\n  public AdderTest() {}\n\n  @BeforeTest\n  public void setUp() {}\n\n  public void test0001_1() {}\n}\n",
	} {
		out, code := checked(t, with(junitRepo(), map[string]string{"src/test/java/AdderTest.java": owner}))
		sees(t, out, code, 0)
		passedWith(t, out, "requirements", "2")
	}
	notJUnit := "class AdderTest {\n  // @Test void test0001_1() {}\n  String note = \"@Test void test0001_1() {}\";\n\n  void test0001_1() {}\n\n  @BeforeTest\n  void test0001_1() {}\n}\n"
	out, code := checked(t, with(junitRepo(), map[string]string{"src/test/java/AdderTest.java": notJUnit}))
	sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
}

// cspell:ignore mvnw diffplug

// rcStaged writes files over a fresh copy of given, stages them, and runs
// requirements-changed on the staged commit.
func rcStaged(t *testing.T, given, files map[string]string) (string, int) {
	t.Helper()
	return stagedCheck(t, given, files, "requirements-changed")
}

// rcPull commits files on a pull request branch off main in a fresh copy of
// given, and runs requirements-changed on it.
func rcPull(t *testing.T, given, files map[string]string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", given)
	git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	commit(t, repo, files, "feat(adder): double numbers")
	return fitness(t, repo, prEnv, "requirements-changed")
}

// rcAdder is goRepo's adder.go with one more function, a product change.
func rcAdder() string {
	return goRepo()["adder.go"] + "\n// Double returns twice a.\nfunc Double(a int) int {\n\treturn a + a\n}\n"
}

// rcPassed is what a passing single-check run prints.
const rcPassed = "All 1 checks passed"

// rcCommented is goRepo's adder.go with every comment rewritten and one
// added after code. Its body is indented anew, as a comment edit may leave it.
func rcCommented() string {
	return "// Package adder sums numbers.\npackage adder\n\n// Add returns a plus b,\n// in either order.\nfunc Add(a, b int) int {\n    return a + b // the sum\n}\n"
}

// rcScript is a shell script, product code, with comments to rewrite.
const rcScript = "#!/bin/sh\n# greet says hello.\necho \"# not a comment\" # says it\necho done\n"

// rcGiven is goRepo with a shell script, a file behind a build tag, and a
// cgo preamble, which the toolchain reads as code.
func rcGiven() map[string]string {
	return with(goRepo(), map[string]string{
		"greet.sh": rcScript,
		"tag.go":   "//go:build darwin\n\npackage adder\n",
		"cgo.go":   "package adder\n\n// int three(void) { return 3; }\nimport \"C\"\n",
	})
}

// rcCommentOnly proves that edits to comments alone, in Go and shell,
// pass staged and on a pull request.
func rcCommentOnly(t *testing.T) {
	t.Helper()
	files := map[string]string{
		"adder.go": rcCommented(),
		"greet.sh": "#!/bin/sh\n# greet prints a greeting.\n\necho \"# not a comment\" # prints it\n  echo done\n",
	}
	out, code := rcStaged(t, rcGiven(), files)
	sees(t, out, code, 0, rcPassed)
	out, code = rcPull(t, rcGiven(), files)
	sees(t, out, code, 0, rcPassed)
}

// rcCodeWithComments proves that a change to code, a directive, or a cgo
// preamble beside comment edits still fails, staged and in a pull.
// Joined lines and a changed string count as code, and so does a test or
// doc changed only in its comments.
func rcCodeWithComments(t *testing.T) {
	t.Helper()
	greeting := strings.Replace(rcScript, "says hello", "prints a greeting", 1)
	for _, c := range []struct{ file, body string }{
		{"adder.go", strings.Replace(rcCommented(), "a + b", "b + a", 1)},
		{"greet.sh", strings.Replace(greeting, "# not a comment\"", "# still not a comment\"", 1)},
		{"greet.sh", strings.Replace(greeting, "# says it\necho", "echo", 1)},
		{"tag.go", "//go:build linux\n\n// Package adder builds on Linux.\npackage adder\n"},
		{"cgo.go", "package adder\n\n// int three(void) { return 4; }\nimport \"C\"\n"},
	} {
		out, code := rcStaged(t, rcGiven(), map[string]string{c.file: c.body})
		sees(t, out, code, 1, c.file+" changed without a test and a requirement doc; change both with it")
		out, code = rcPull(t, rcGiven(), map[string]string{c.file: c.body})
		sees(t, out, code, 1, c.file+" changed without a test and a requirement doc")
	}
	partners := map[string]string{
		"adder.go":      rcAdder(),
		"adder_test.go": withComment(goRepo()["adder_test.go"]),
		pairedDoc:       goRepo()[pairedDoc] + "\n<!-- Sums stay small. -->\n",
	}
	out, code := rcStaged(t, rcGiven(), partners)
	sees(t, out, code, 1, "adder.go changed without a test and a requirement doc; change both with it")
}

// provenRepo is twoRepo whose requirement 0002 a check in fitness/req-0002
// proves, with no test of its own.
func provenRepo() map[string]string {
	return with(twoRepo(), map[string]string{
		".fitnessrc.json":            `{"policy": "external", "checks": ["requirements", "req-0002"]}` + "\n",
		"reqs/0001_0002_test.go":     "",
		"reqs/0001_test.go":          testFile("reqs", "Test0001_1"),
		"fitness/req-0002/README.md": "# Doubles\n\nProves [requirement 0002](../../docs/requirements/0002-adds-again.md).\n",
		"fitness/req-0002/check":     verdictScript(passing),
	})
}
