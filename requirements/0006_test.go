package requirements

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test0006_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"doc file name and title IDs differ", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, docOf("docs/requirements/0002-adds-numbers.md", requirementDoc("0001.1")))
			sees(t, out, code, 1, "name the file NNNN-kebab-title.md and title it `# NNNN Title` with the same ID")
		}},
		{"two docs share one ID", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, with(goRepo(), map[string]string{"docs/requirements/0001-other.md": requirementDoc("0001.1")}))
			sees(t, out, code, 1, "requirement 0001 is used by")
		}},
		{"sections missing or empty", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, docOf("docs/requirements/0001-adds-numbers.md", requirementDoc("0001.1")+"\n## Out of scope\n\n- Subtraction.\n"))
			sees(t, out, code, 1, "sections must be exactly ## Why, ## Measurement, ## Requirements")
		}},
		{"measurement is not a numerator and denominator", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, docOf("docs/requirements/0001-adds-numbers.md", strings.Replace(requirementDoc("0001.1"), "\n-\n", "\nover\n", 1)))
			sees(t, out, code, 1, "Measurement needs three lines: numerator, `-`, denominator")
		}},
		{"acceptance ID repeated or foreign", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, reqs(acceptances("0001.1", "0001.1")))
			sees(t, out, code, 1, "0001.1 appears twice")
		}},
		{"doc holds more acceptances than the cap", func(t *testing.T) {
			t.Parallel()
			for _, row := range []struct {
				max, count int
				want       []string
			}{
				{0, 8, nil},
				{0, 9, []string{pairedDoc + ":49: 0001 holds 9 acceptances, at most 8; fold cases into table rows of one acceptance, or move a separate goal to its own doc"}},
				{10, 9, nil},
				{10, 11, []string{pairedDoc + ":57: 0001 holds 11 acceptances, at most 10;"}},
				{3, 4, []string{pairedDoc + ":29: 0001 holds 4 acceptances, at most 3;"}},
			} {
				t.Run(fmt.Sprintf("cap %d holds %d", row.max, row.count), func(t *testing.T) {
					t.Parallel()
					out, code := fitness(t, example(t, "happyRepo", cappedRepo(row.max, row.count)), nil, "requirements")
					sees(t, out, code, min(len(row.want), 1), row.want...)
				})
			}
		}},
		{"second given when then or wrong order", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, reqs(acceptances("0001.1")+"        - When I add them again\n"))
			sees(t, out, code, 1, "needs exactly one Given, one When, and one Then")
			out, code = judged(t, reqs("- 0001.1\n    - When I add them\n        - Given two numbers\n            - Then I get their sum\n"))
			sees(t, out, code, 1, "each acceptance is one Given, then one When, then one Then, each nested under the line above")
		}},
		{"steps chain up to seven beside a single step", func(t *testing.T) {
			t.Parallel()
			tests := map[string]string{"adder_test.go": testFile("adder", "Test0001_1"), "chain_test.go": caseFile("adder", "Test0001_2", stepCases("0001.2", 7)...)}
			out, code := fitness(t, example(t, "happyRepo", with(reqs(acceptances("0001.1")+chain("0001.2", 7)), tests)), nil, "requirements")
			sees(t, out, code, 0)
			out, code = judged(t, reqs(chain("0001.1", 8)))
			sees(t, out, code, 1, pairedDoc+":17: 0001.1 chains 8 steps, at most 7; move the rest to a new ID")
		}},
		{"misnumbered or single numbered steps fail", func(t *testing.T) {
			t.Parallel()
			for body, want := range map[string]string{
				strings.Replace(chain("0001.1", 3), "    3. When", "    4. When", 1): ":23: number this step `3.`; a chain numbers its steps `1.` to `7.` in order",
				strings.Replace(chain("0001.1", 2), "    2. When", "    - When", 1):  ":21: number this step `2.`",
				chain("0001.1", 1): ":19: a single step is the plain bullet `- When`; number steps only in a chain of two or more",
				strings.Replace(chain("0001.1", 2), "       - Then I get their sum\n", "", 1): ":19: step 1 needs exactly one Then nested under its When",
				strings.Replace(chain("0001.1", 2), "    2. When", "    2. Then", 1):          ":21: " + "each step is a numbered When nested under the Given",
			} {
				out, code := judged(t, reqs(body))
				sees(t, out, code, 1, pairedDoc+want)
			}
		}},
		{"two docs share a when", func(t *testing.T) {
			t.Parallel()
			other := "docs/requirements/0002-adds-again.md"
			adds := "- 0001.1\n    - Given Writer has two whole numbers\n        - When I add them\n            - Then I get their sum\n"
			again := "- 0002.1\n    - Given Writer has two whole numbers that carry\n        - When I add them\n            - Then I get their sum\n"
			out, code := judged(t, with(goRepo(), map[string]string{
				"adder_test.go": testFile("adder", "Test0001_1", "Test0002_1"),
				pairedDoc:       docWith(adds),
				other:           strings.Replace(docWith(again), "0001 Adds Numbers", "0002 Adds Again", 1),
			}))
			sees(t, out, code, 1, pairedDoc+" and "+other+` both cover "When I add them"`, "blob/main/fitness/adr-0004")
		}},
		{"stale source line fails then fixes", func(t *testing.T) {
			t.Parallel()
			doc := docWith("1. Adds\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n")
			sourced := withSource(doc, "`gone.go:99`")
			repo := example(t, "happyRepo", with(pairedRepo(), map[string]string{pairedDoc: sourced}))
			out, code := fitness(t, repo, nil, "requirements")
			sees(t, out, code, 1, pairedDoc+":15: delete this `Source:` line, or run `fitness requirements --fix`")
			out, code = fitness(t, repo, nil, "requirements", "--fix")
			sees(t, out, code, 0)
			if got, err := os.ReadFile(filepath.Join(repo, pairedDoc)); err != nil || string(got) != doc {
				t.Errorf("--fix must leave the doc without its Source line, got %v:\n%s", err, got)
			}
			git(t, repo, "add", "-A")
			out, code = fitness(t, repo, nil, "requirements")
			sees(t, out, code, 0)
			if strings.Contains(out, "Source") {
				t.Errorf("a doc without a Source line gets no error, got:\n%s", out)
			}
		}},
	})
}

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
			out, code := judged(t, with(goRepo(), map[string]string{"adder_test.go": ""}))
			sees(t, out, code, 1, "0001.1 has no test; name exactly one test Test0001_1")
			swiftOwnsNothing(t, xcTestFile())
			out, code = checked(t, with(scriptRepo(), map[string]string{"tests/adder.test.ts": ""}))
			sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
			out, code = checked(t, with(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest")}))
			sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
			out, code = judged(t, with(goRepo(), map[string]string{"more/more_test.go": testFile("more", "Test0001_1")}))
			sees(t, out, code, 1, "0001.1 is owned by 2 tests")
			out, code = judged(t, with(goRepo(), map[string]string{"adder_test.go": testFile("adder", "Test0001_1", "Test0001_2")}))
			sees(t, out, code, 1, "Test0001_2 names 0001.2, which no requirement defines")
		}},
		{"each chained step needs a named case", func(t *testing.T) {
			t.Parallel()
			chained := reqs(acceptances("0001.1") + chain("0001.2", 2))
			named := func(cases ...string) map[string]string {
				return with(chained, map[string]string{"chain_test.go": caseFile("adder", "Test0001_2", cases...)})
			}
			out, code := checked(t, named("0001.2.1 adds", "0001.2.2 doubles", "edge"))
			sees(t, out, code, 0)
			out, code = checked(t, named("0001.2.1 adds"))
			sees(t, out, code, 1, pairedDoc+":25: step 0001.2.2 has no test case; add a case to Test0001_2 whose name opens with 0001.2.2")
			out, code = checked(t, named("0001.2.1 adds", "0001.2.2 doubles", "0001.2.3 halves", "0001.1.1 sums"))
			sees(t, out, code, 1, `chain_test.go:8: case "0001.2.3 halves" names step 0001.2.3, which 0001.2 does not chain`,
				`chain_test.go:9: case "0001.1.1 sums" names step 0001.1.1, which 0001.2 does not chain`)
			out, code = checked(t, with(goRepo(), map[string]string{"adder_test.go": caseFile("adder", "Test0001_1", "0001.1.1 adds")}))
			sees(t, out, code, 1, `case "0001.1.1 adds" names step 0001.1.1, which 0001.1 does not chain`)
			for _, repo := range []map[string]string{scriptRepo(), swiftRepo()} {
				out, code = checked(t, with(repo, map[string]string{pairedDoc: docWith(chain("0001.1", 2))}))
				sees(t, out, code, 0)
			}
			out, _ = checked(t, with(reqs(chain("0001.1", 8)), map[string]string{"adder_test.go": caseFile("adder", "Test0001_1", stepCases("0001.1", 7)...)}))
			if strings.Contains(out, "does not chain") {
				t.Errorf("a malformed chain fails its good cases:\n%s", out)
			}
		}},
		{"test proving no requirement fails", func(t *testing.T) {
			t.Parallel()
			out, code := judged(t, with(goRepo(), map[string]string{"adder_test.go": testFile("adder", "Test0001_1", "TestAddsAgain")}))
			sees(t, out, code, 1, "TestAddsAgain proves no requirement")
			out, code = judged(t, with(reqs(chain("0001.1", 2)), map[string]string{"adder_test.go": testFile("adder", "Test0001_1", "Test0001_1_2")}))
			sees(t, out, code, 1, "Test0001_1_2 names step 0001.1.2; one test proves every step of its acceptance, so name it Test0001_1")
			out, code = judged(t, with(reqs(chain("0001.1", 2)), map[string]string{"adder_test.go": testFile("adder", "Test0001_1", "TestAdds0001_1_2")}))
			sees(t, out, code, 1, "TestAdds0001_1_2 proves no requirement")
			out, code = swiftJudged(t, map[string]string{swiftTest: swiftTestingFile("test0001_1", "testAddsAgain")})
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
		{"deleted ID brought back fails", func(t *testing.T) {
			t.Parallel()
			two := with(goRepo(), map[string]string{
				"adder_test.go":                          testFile("adder", "Test0001_1", "Test0001_2"),
				"docs/requirements/0001-adds-numbers.md": requirementDoc("0001.1", "0001.2"),
			})
			repo := example(t, "happyRepo", two)
			commit(t, repo, goRepo(), "docs(requirements): drop 0001.2")
			write(t, repo, two)
			out, code := fitness(t, repo, nil)
			sees(t, out, code, 1, "0001.2 was deleted before")
		}},
		{"deleted step brought back fails", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", goRepo())
			cased := func(steps int) map[string]string {
				return with(reqs(chain("0001.1", steps)), map[string]string{"adder_test.go": caseFile("adder", "Test0001_1", stepCases("0001.1", steps)...)})
			}
			commit(t, repo, cased(3), "docs(requirements): chain 0001.1 in three steps")
			out, code := fitness(t, repo, nil, "requirements")
			sees(t, out, code, 0)
			commit(t, repo, cased(2), "docs(requirements): drop step 0001.1.3")
			write(t, repo, cased(3))
			out, code = fitness(t, repo, nil, "requirements")
			sees(t, out, code, 1, pairedDoc+":23: 0001.1.3 was deleted before; IDs are never reused, so chain the step under a new acceptance ID")
		}},
		{"whole restructure reuses numbers, then counts removals afresh", func(t *testing.T) {
			t.Parallel()
			two := with(goRepo(), map[string]string{
				"adder_test.go":                          testFile("adder", "Test0001_1", "Test0001_2"),
				"docs/requirements/0001-adds-numbers.md": requirementDoc("0001.1", "0001.2"),
			})
			repo := example(t, "happyRepo", two)
			commit(t, repo, goRepo(), "docs(requirements): drop 0001.2")
			restructured := with(two, map[string]string{
				"docs/requirements/0001-adds-numbers.md": "",
				"docs/requirements/0001-sums-numbers.md": requirementDoc("0001.1", "0001.2"),
			})
			write(t, repo, restructured)
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, "requirements")
			sees(t, out, code, 0)
			commit(t, repo, restructured, "docs(requirements): restart the numbers")
			out, code = fitness(t, repo, nil, "requirements")
			sees(t, out, code, 0)
			dropped := map[string]string{"adder_test.go": goRepo()["adder_test.go"], "docs/requirements/0001-sums-numbers.md": requirementDoc("0001.1")}
			commit(t, repo, dropped, "docs(requirements): drop 0001.2 again")
			write(t, repo, restructured)
			out, code = fitness(t, repo, nil, "requirements")
			sees(t, out, code, 1, "0001.2 was deleted before")
		}},
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

func Test0006_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"vet problem at module root fails", func(t *testing.T) {
			t.Parallel()
			out, code := goVet(t, example(t, "happyRepo", goVetModule(".", goVetBad)))
			sees(t, out, code, 1, "✖ "+goVetFinding)
		}},
		{"vet problem in subfolder module", func(t *testing.T) {
			t.Parallel()
			out, code := goVet(t, example(t, "happyRepo", goVetModule("api", goVetBad)))
			sees(t, out, code, 1, "✖ api: "+goVetFinding)
		}},
		{"clean module passes vet", func(t *testing.T) {
			t.Parallel()
			out, code := goVet(t, example(t, "happyRepo", goVetModule("api", "package m\n\nfunc F() int { return 1 }\n")))
			sees(t, out, code, 0, "All 1 checks passed")
			if !regexp.MustCompile(`│ passed\s+│ 1\s+│`).MatchString(row(out, "go-vet")) {
				t.Errorf("go-vet must pass counting 1 module, got row %q", row(out, "go-vet"))
			}
		}},
		{"vet with no go mod passes", func(t *testing.T) {
			t.Parallel()
			out, code := goVet(t, example(t, "happyRepo", map[string]string{"main.go": goVetBad}))
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "go-vet")
		}},
		{"vet skips testdata and vendor", func(t *testing.T) {
			t.Parallel()
			files := goVetModule("testdata/fixture", goVetBad)
			for name, body := range goVetModule("vendor/dep", goVetBad) {
				files[name] = body
			}
			out, code := goVet(t, example(t, "happyRepo", files))
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "go-vet")
		}},
		{"vet staged non Go change passes", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", goVetModule("api", goVetBad))
			write(t, repo, map[string]string{"README.md": readme + "\nMore text.\n"})
			git(t, repo, "add", "README.md")
			out, code := goVet(t, repo)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "go-vet")
		}},
		{"vet without toolchain shows info", func(t *testing.T) {
			t.Parallel()
			out, code := goVet(t, example(t, "happyRepo", goVetModule("api", goVetBad)), pathWithout(t, "go"))
			showsInfo(t, out, code, "go-vet", "Go not installed: see https://go.dev/dl")
		}},
		{"vet exit without file shows info", func(t *testing.T) {
			t.Parallel()
			out, code := goVet(t, example(t, "happyRepo", goVetModule(".", goVetBad)), crashingGo(t, "vet", "go: unexpected internal error"))
			showsInfo(t, out, code, "go-vet", "go vet ./... exited 2 without a finding: go: unexpected internal error")
		}},
		{"gofmt fails an unformatted file", func(t *testing.T) {
			t.Parallel()
			out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "bad.go": gofmtMessy})
			sees(t, out, code, 1, "bad.go: not gofmt-formatted (run: gofmt -w bad.go)")
		}},
		{"formatted files pass gofmt", func(t *testing.T) {
			t.Parallel()
			out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"gofmt skips testdata and vendor", func(t *testing.T) {
			t.Parallel()
			out, code := gofmtCheck(t, nil, map[string]string{
				"ok.go":                 gofmtTidy,
				"x/testdata/fixture.go": gofmtMessy,
				"vendor/lib/lib.go":     gofmtMessy,
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"gofmt with no Go files passes", func(t *testing.T) {
			t.Parallel()
			out, code := gofmtCheck(t, nil, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "gofmt")
		}},
		{"gofmt without toolchain shows info", func(t *testing.T) {
			t.Parallel()
			out, code := gofmtCheck(t, []string{pathWithout(t, "go", "gofmt")}, map[string]string{"ok.go": gofmtTidy})
			showsInfo(t, out, code, "gofmt", "Go not installed: see https://go.dev/dl")
		}},
		{"gofmt ignores untracked files", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"ok.go": gofmtTidy})
			write(t, repo, map[string]string{"scratch.go": gofmtMessy})
			out, code := fitness(t, repo, nil, "gofmt")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"gofmt crash shows info", func(t *testing.T) {
			t.Parallel()
			crashed := brokenTool(t, "gofmt", "echo 'fatal error: runtime: out of memory' >&2\nexit 2\n")
			out, code := gofmtCheck(t, []string{crashed}, map[string]string{"ok.go": gofmtTidy})
			showsInfo(t, out, code, "gofmt", "gofmt -l exited 2 without a finding: fatal error: runtime: out of memory")
			if strings.Contains(out, "not gofmt-formatted") {
				t.Errorf("a crash must not read as an unformatted file:\n%s", out)
			}
		}},
		{"gofmt syntax error fails", func(t *testing.T) {
			t.Parallel()
			out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "broken.go": "package m\n\nfunc {\n"})
			sees(t, out, code, 1, "broken.go:3:6: expected 'IDENT', found '{'")
			if strings.Contains(out, "not gofmt-formatted") {
				t.Errorf("a syntax error must fail as itself, not as an unformatted file:\n%s", out)
			}
		}},
	})
}

func Test0006_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"disable directives are listed", func(t *testing.T) {
			t.Parallel()
			legacy := "export const x = 1;\n" + noEslintDisableNext + "console.log(x);\n" +
				"const y = 2; // eslint-disable-line\n/* eslint-disable */\n"
			out, code := noEslintDisableCheck(t, map[string]string{"src/legacy.ts": legacy})
			sees(t, out, code, 1,
				"src/legacy.ts:2: eslint-disable-next-line",
				"src/legacy.ts:4: eslint-disable-line",
				"src/legacy.ts:5: eslint-disable")
		}},
		{"files without directives pass", func(t *testing.T) {
			t.Parallel()
			out, code := noEslintDisableCheck(t, map[string]string{
				"src/a.ts": "export const a = 1;\n",
				"src/b.js": "export const b = 2;\n",
			})
			sees(t, out, code, 0, "All 1 checks passed · 2 files scanned")
		}},
		{"directive text in strings fails", func(t *testing.T) {
			t.Parallel()
			files := map[string]string{}
			var want []string
			for _, ext := range noEslintDisableExts {
				files["src/f"+ext] = "export const hint = 'add eslint-disable-next-line to skip';\n"
				want = append(want, "src/f"+ext+":1: eslint-disable-next-line")
			}
			out, code := noEslintDisableCheck(t, files)
			sees(t, out, code, 1, append(want, "7 files scanned")...)
		}},
		{"directive text in Markdown passes", func(t *testing.T) {
			t.Parallel()
			out, code := noEslintDisableCheck(t, map[string]string{
				"docs/lint.md": "Never write " + noEslintDisableNext,
				"src/a.ts":     "export const a = 1;\n",
			})
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
		{"directive scan with no source passes", func(t *testing.T) {
			t.Parallel()
			out, code := noEslintDisableCheck(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "no-eslint-disable")
		}},
		{"node_modules and dist are scanned", func(t *testing.T) {
			t.Parallel()
			out, code := noEslintDisableCheck(t, map[string]string{
				"node_modules/pkg/index.js": "/* eslint-disable */\n",
				"dist/out.js":               "/* eslint-disable */\n",
				"src/clean.ts":              "export const a = 1;\n",
			})
			sees(t, out, code, 1,
				"node_modules/pkg/index.js:1: eslint-disable",
				"dist/out.js:1: eslint-disable",
				"3 files scanned")
		}},
		{"untracked directive file passes", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"src/a.ts": "export const a = 1;\n"})
			write(t, repo, map[string]string{"src/scratch.ts": noEslintDisableNext})
			out, code := noEslintDisableRun(t, repo)
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
		{"directive scan on staged file only", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"src/old.ts": noEslintDisableNext})
			write(t, repo, map[string]string{"src/new.ts": "export const a = 1;\n"})
			git(t, repo, "add", "src/new.ts")
			out, code := noEslintDisableRun(t, repo)
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
	})
}

func Test0006_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"unformatted Java files fail spotless", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"pom.xml": javaPom(spotlessMaven)})
			_, out, code := javaRun(t, repo, "spotless", map[string][]response{"mvn": {spotlessViolations}})
			sees(t, out, code, 1, javaSource+" is not formatted; run spotless:apply")
		}},
		{"multi module build passes spotless", func(t *testing.T) {
			t.Parallel()
			modules := map[string]string{
				"pom.xml":                         javaParent(spotlessMaven, "core", "web"),
				"core/pom.xml":                    javaPom(""),
				"core/src/main/java/Core.java":    "class Core {}\n",
				"core/src/test/resources/pom.xml": "<project/>\n",
				"web/pom.xml":                     javaPom(""),
				"web/src/main/java/web/Page.java": "class Page {}\n",
			}
			rp, out, code := javaRun(t, javaRepo(t, modules), "spotless", map[string][]response{"mvn": {spotlessClean}})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "spotless", "3")
			javaCalledOnce(t, rp, "mvn")
			// A nested build no pom lists as a module is a build root of its own.
			modules["tools/pom.xml"] = javaPom(spotlessMaven)
			rp, out, code = javaRun(t, javaRepo(t, modules), "spotless", map[string][]response{"mvn": {spotlessClean}})
			sees(t, out, code, 0, "All 1 checks passed")
			if calls := rp.calls("mvn"); len(calls) != 2 {
				t.Errorf("mvn calls = %q, want one per build root", calls)
			}
		}},
		{"spotless runs the build wrapper", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"pom.xml": javaPom(spotlessMaven)})
			javaWrapper(t, repo, "mvnw")
			rp, out, code := javaRun(t, repo, "spotless", map[string][]response{"mvn": {spotlessClean}, "mvnw": {spotlessClean}})
			sees(t, out, code, 0, "All 1 checks passed")
			javaCalledOnce(t, rp, "mvnw")
			if calls := rp.calls("mvn"); len(calls) != 0 {
				t.Errorf("the build tool on PATH ran %q, want the wrapper only", calls)
			}
		}},
		{"gradle runs without a daemon", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"build.gradle": "plugins { id 'java'; id 'com.diffplug.spotless' version '8.10.4' }\n"})
			rp, out, code := javaRun(t, repo, "spotless", map[string][]response{"gradle": {spotlessGradleClean}})
			sees(t, out, code, 0, "All 1 checks passed")
			if !rp.called("gradle", "--no-daemon", "spotlessCheck") {
				t.Errorf("gradle calls = %q, want spotlessCheck with --no-daemon", rp.calls("gradle"))
			}
		}},
		{"build without spotless gets snippet", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{
				"pom.xml":                         javaPom(""),
				"tools/pom.xml":                   javaPom(spotlessMaven),
				"src/test/resources/fake/pom.xml": javaPom(spotlessMaven),
			})
			_, out, code := javaRun(t, repo, "spotless", map[string][]response{"mvn": {spotlessClean}})
			sees(t, out, code, 1, "Spotless is not in the Maven build at .; add it:")
			javaSeesSnippet(t, out, "<artifactId>spotless-maven-plugin</artifactId>")
			if strings.Contains(read(out), "build at tools") {
				t.Errorf("the nested build declares Spotless, so only the root may fail:\n%s", out)
			}
		}},
		{"spotless with no Java build passes", func(t *testing.T) {
			t.Parallel()
			_, out, code := javaRun(t, example(t, "happyRepo", nil), "spotless", nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "spotless")
		}},
		{"gradle failure shows its causes", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"build.gradle": "plugins { id 'java'; id 'com.diffplug.spotless' version '8.10.4' }\n"})
			_, out, code := javaRun(t, repo, "spotless", map[string][]response{"gradle": {spotlessGradleNoJDK}})
			sees(t, out, code, 1, "Could not determine the dependencies of task ':spotlessJava'.")
			javaSeesSnippet(t, out, "> Cannot find a Java installation on your machine")
			javaSeesSnippet(t, out, "matching: {languageVersion=17,")
		}},
	})
}

func Test0006_7(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"failing Go test is named", func(t *testing.T) {
			t.Parallel()
			out, code := goTest(t, goTestModule("", goTestBroken))
			sees(t, out, code, 1, "--- FAIL: TestBroken", "m_test.go:5: boom")
		}},
		{"passing Go module passes", func(t *testing.T) {
			t.Parallel()
			out, code := goTest(t, goTestModule("", "func TestFine(t *testing.T) {}\n"))
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
		{"go test with no go mod passes", func(t *testing.T) {
			t.Parallel()
			out, code := goTest(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "go-test")
		}},
		{"go test skips testdata and vendor", func(t *testing.T) {
			t.Parallel()
			files := with(goTestModule("testdata/a/", goTestBroken), goTestModule("vendor/b/", goTestBroken))
			out, code := goTest(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "go-test")
		}},
		{"failing subfolder module shows path", func(t *testing.T) {
			t.Parallel()
			out, code := goTest(t, goTestModule("svc/", goTestBroken))
			sees(t, out, code, 1, "svc: --- FAIL: TestBroken", "svc: m_test.go:5: boom")
		}},
		{"Go test that does not compile fails", func(t *testing.T) {
			t.Parallel()
			out, code := goTest(t, goTestModule("", "func TestBroken(t *testing.T) { undefinedThing() }\n"))
			sees(t, out, code, 1, "m_test.go:5:33: undefined: undefinedThing")
		}},
		{"go test with only non Go staged", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", goTestModule("", goTestBroken))
			write(t, repo, map[string]string{"README.md": readme + "\nOne more line.\n"})
			git(t, repo, "add", "README.md")
			out, code := fitness(t, repo, nil, "go-test")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "go-test")
		}},
		{"go test missing toolchain shows info", func(t *testing.T) {
			t.Parallel()
			out, code := goTest(t, goTestModule("", "func TestFine(t *testing.T) {}\n"), pathWithout(t, "go"))
			showsInfo(t, out, code, "go-test", "Go not installed: see https://go.dev/dl")
			out, code = goTest(t, goTestModule("", goTestBroken), crashingGo(t, "test", "go: unexpected internal error"))
			showsInfo(t, out, code, "go-test", "go test ./... exited 2 without a finding: go: unexpected internal error")
		}},
		{"coverage exclude with no package json passes", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "vitest-coverage-exclude")
		}},
		{"empty coverage exclude passes", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{
				"package.json":     "{}\n",
				"vitest.config.js": coverageExcludeConfig("[]"),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"excluded test files and dist fail", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{
				"package.json":     "{}\n",
				"vitest.config.js": coverageExcludeConfig(`["**/*.test.ts", "dist"]`),
			})
			sees(t, out, code, 1, coverageExcludeRemove(`"**/*.test.ts"`), coverageExcludeRemove(`"dist"`))
		}},
		{"package json vitest exclude fails", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{
				"package.json": `{"vitest": {"test": {"coverage": {"exclude": ["src/foo.ts"]}}}}` + "\n",
			})
			sees(t, out, code, 1, coverageExcludeRemove(`"src/foo.ts"`))
		}},
		{"config file wins over package json", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{
				"package.json":     `{"vitest": {"coverage": {"exclude": ["src/foo.ts"]}}}` + "\n",
				"vitest.config.js": coverageExcludeConfig("[]"),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"only quoted exclude paths are named", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{
				"package.json": "{}\n",
				"vitest.config.mjs": `import { DTS_GLOB } from "./globs.js";` + "\n" +
					`export default { test: { coverage: { exclude: [DTS_GLOB, "src/foo.ts"] } } };` + "\n",
			})
			sees(t, out, code, 1, coverageExcludeRemove(`"src/foo.ts"`))
			if strings.Contains(out, "DTS_GLOB") {
				t.Errorf("an imported name must not be reported:\n%s", out)
			}
		}},
		{"shared config exclude fails", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{
				"package.json": "{}\n",
				"node_modules/@mayjournal/fitness-shared/config/vitest.config.js": coverageExcludeConfig(`["src/shared.ts"]`),
			})
			sees(t, out, code, 1, coverageExcludeRemove(`"src/shared.ts"`))
		}},
		{"built in shared config passes", func(t *testing.T) {
			t.Parallel()
			out, code := coverageExclude(t, map[string]string{"package.json": "{}\n"})
			sees(t, out, code, 0, "All 1 checks passed")
			if !coverageExcludeOneFile.MatchString(row(out, "vitest-coverage-exclude")) {
				t.Errorf("vitest-coverage-exclude must pass judging 1 file:\n%s", out)
			}
		}},
		{"coverage with no JS source passes", func(t *testing.T) {
			t.Parallel()
			out, code := vcfRun(t, example(t, "happyRepo", nil), nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "vitest-coverage-full")
			tooling := example(t, "happyRepo", map[string]string{"package.json": "{\"private\": true}\n", "vitest.config.mjs": vcfConfig("50"), "vite.config.ts": "export default {};\n"})
			out, code = vcfRun(t, tooling, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "vitest-coverage-full")
		}},
		{"threshold below full fails", func(t *testing.T) {
			t.Parallel()
			rp := standIn(t, map[string][]response{"vitest": {vcfFullRun}})
			out, code := vcfRun(t, vcfRepo(t, map[string]string{"vitest.config.mjs": vcfConfig("90")}), rp.env)
			sees(t, out, code, 1, "Vitest coverage thresholds must be 100 for branches, functions, lines, and statements.")
		}},
		{"shared config threshold below full", func(t *testing.T) {
			t.Parallel()
			shared := "node_modules/@mayjournal/fitness-shared/config/vitest.config.mjs"
			out, code := vcfRun(t, vcfRepo(t, map[string]string{shared: vcfConfig("90")}), nil)
			sees(t, out, code, 1, "Fitness-runner package must have Vitest coverage thresholds set to 100 for branches, functions, lines, and statements.")
		}},
		{"fully tested code passes", func(t *testing.T) {
			t.Parallel()
			rp := standIn(t, map[string][]response{"vitest": {vcfFullRun}})
			out, code := vcfRun(t, vcfRepo(t, nil), rp.env)
			sees(t, out, code, 0, "All 1 checks passed")
			vcfRanWithLocalConfig(t, rp)
		}},
		{"untested code shows coverage", func(t *testing.T) {
			t.Parallel()
			rp := standIn(t, map[string][]response{"vitest": {vcfShortRun}})
			repo := vcfRepo(t, map[string]string{"src/sub.ts": "export const sub = (a: number, b: number): number => a - b;\n"})
			out, code := vcfRun(t, repo, rp.env)
			sees(t, out, code, 1, "Coverage for statements (50%) does not meet global threshold (100%)")
			vcfRanWithLocalConfig(t, rp)
		}},
		{"missing vitest shows info", func(t *testing.T) {
			t.Parallel()
			out, code := vcfRun(t, vcfRepo(t, nil), nil)
			showsInfo(t, out, code, "vitest-coverage-full", "vitest not found in node_modules/.bin (walking up from root) or on PATH. Install it (npm install -D vitest)")
		}},
		{"vitest on path replaces broken binary", func(t *testing.T) {
			t.Parallel()
			rp := standIn(t, map[string][]response{"vitest": {vcfFullRun}})
			repo := vcfRepo(t, map[string]string{"node_modules/.bin/vitest": "#!/bin/sh\nexit 1\n"})
			mustDo(t, os.Chmod(filepath.Join(repo, "node_modules", ".bin", "vitest"), 0o644))
			out, code := vcfRun(t, repo, rp.env)
			sees(t, out, code, 0, "All 1 checks passed")
			vcfRanWithLocalConfig(t, rp)
		}},
		{"vitest exit without cause shows info", func(t *testing.T) {
			t.Parallel()
			rp := standIn(t, map[string][]response{"vitest": {{
				Match:  []string{"run --coverage"},
				Stderr: " MISSING DEPENDENCY  Cannot find dependency '@vitest/coverage-v8'\n",
				Exit:   1,
			}}})
			out, code := vcfRun(t, vcfRepo(t, nil), rp.env)
			showsInfo(t, out, code, "vitest-coverage-full", "vitest run --coverage exited 1 without a finding", "Cannot find dependency '@vitest/coverage-v8'")
			vcfRanWithLocalConfig(t, rp)
		}},
		{"build pinning no JDK fails", func(t *testing.T) {
			t.Parallel()
			out, code := javaVersionRun(t, map[string]string{"pom.xml": "<project/>\n"})
			sees(t, out, code, 1, "no JDK version is pinned; add a .java-version file naming the release",
				"JavaLanguageVersion.of(21)")
		}},
		{"version file at or below JDK passes", func(t *testing.T) {
			t.Parallel()
			out, code := javaVersionRun(t, map[string]string{"pom.xml": "<project/>\n", ".java-version": "17\n"})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"release above running JDK fails", func(t *testing.T) {
			t.Parallel()
			out, code := javaVersionRun(t, map[string]string{"pom.xml": javaVersionRelease("25")})
			sees(t, out, code, 1, "JDK 21.0.2 does not satisfy pom.xml (requires 25); install JDK 25 or point JAVA_HOME at one")
		}},
		{"maven release from a property", func(t *testing.T) {
			t.Parallel()
			out, code := javaVersionRun(t, map[string]string{
				"pom.xml": "<project><properties><java.version>25</java.version>" +
					"<maven.compiler.release>${java.version}</maven.compiler.release></properties></project>\n",
			})
			sees(t, out, code, 1, "does not satisfy pom.xml (requires 25)")
			// A parent from the local Maven repository maps the pom's property to
			// the compiler release, as framework parents do.
			home := t.TempDir()
			write(t, home, map[string]string{".m2/repository/" + javaVersionParentPom: javaVersionParent})
			child := "<project><parent><groupId>com.example</groupId><artifactId>parent</artifactId>" +
				"<version>1</version><relativePath/></parent><properties><java.version>25</java.version></properties></project>\n"
			javaVersionSees(t, javaRepo(t, map[string]string{"pom.xml": child}), home, "JDK 21.0.2 does not satisfy pom.xml (requires 25)")
			// The user's Maven settings may move the local repository.
			moved := t.TempDir()
			write(t, moved, map[string]string{".m2/settings.xml": "<settings><localRepository>${user.home}/maven</localRepository></settings>\n", "maven/" + javaVersionParentPom: javaVersionParent})
			javaVersionSees(t, javaRepo(t, map[string]string{"pom.xml": child}), moved, "does not satisfy pom.xml (requires 25)")
			// A pom outside the repo at the default relative path is no parent in
			// the repo, so the local repository's parent still counts.
			repo := javaRepo(t, map[string]string{"pom.xml": strings.Replace(child, "<relativePath/>", "", 1)})
			write(t, filepath.Dir(repo), map[string]string{"pom.xml": "<project><artifactId>parent</artifactId></project>\n"})
			javaVersionSees(t, repo, home, "does not satisfy pom.xml (requires 25)")
			// A module takes the property its parent in the repo sets.
			javaVersionSees(t, javaRepo(t, map[string]string{
				"pom.xml": "<project><artifactId>app</artifactId><modules><module>core</module></modules><properties><java.version>25</java.version></properties></project>\n",
				"core/pom.xml": "<project><parent><groupId>com.example</groupId><artifactId>app</artifactId><version>1</version></parent>" +
					"<properties><maven.compiler.release>${java.version}</maven.compiler.release></properties></project>\n",
			}), home, "does not satisfy core/pom.xml (requires 25)")
		}},
		{"pin naming no release fails", func(t *testing.T) {
			t.Parallel()
			out, code := javaVersionRun(t, map[string]string{"pom.xml": "<project/>\n", ".java-version": "system\n"})
			sees(t, out, code, 1, `.java-version pins JDK "system", which names no release`)
		}},
		{"no running JDK fails", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"pom.xml": javaVersionRelease("21")})
			out, code := fitness(t, repo, []string{pathWithout(t, "java"), "JAVA_HOME="}, "java-version")
			sees(t, out, code, 1, "Java not installed: install a JDK")
		}},
		{"java version with no build passes", func(t *testing.T) {
			t.Parallel()
			_, out, code := javaRun(t, example(t, "happyRepo", nil), "java-version", nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "java-version")
		}},
		{"gradle toolchain with no match fails", func(t *testing.T) {
			t.Parallel()
			work := func(release string, extra map[string]string) string { return javaVersionToolchain(t, release, extra) }
			// No installed JDK has the release and no resolver is configured. Gradle
			// never runs, so the check stays fast.
			javaVersionFails(t, work("30", nil), t.TempDir())
			// Gradle's own download folder counts, flat or as a macOS bundle one
			// folder down.
			javaVersionPasses(t, work("30", nil), javaVersionHome(t, map[string]string{".gradle/jdks/jdk-30/release": javaVersionJDK("30")}))
			javaVersionPasses(t, work("30", nil), javaVersionHome(t, map[string]string{".gradle/jdks/amazon-30/jdk-30/Contents/Home/release": javaVersionJDK("30")}))
			// Gradle always finds the JDK it runs on.
			javaVersionPasses(t, work("21", nil), t.TempDir())
			// An empty JAVA_HOME, home or environment name reads nothing from the
			// working folder, where a release file sits.
			cwd := work("30", map[string]string{"release": javaVersionJDK("30"), "gradle.properties": "org.gradle.java.installations.fromEnv=NO_SUCH_JDK_VARIABLE\n"})
			javaVersionFails(t, cwd, "")
			javaVersionPropertiesCases(t)
			javaVersionResolverCases(t)
			javaVersionVendorCases(t)
		}},
		{"failing Java tests are named", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"pom.xml": javaPom(javaTestMaven)})
			run := javaTestMavenRun(map[string]string{
				"target/surefire-reports/TEST-example.AppTest.xml": javaTestFailed,
				"target/site/jacoco/jacoco.xml":                    javaTestFull,
			})
			_, out, code := javaRun(t, repo, "java-test", map[string][]response{"mvn": {run}})
			sees(t, out, code, 1, "example.AppTest.subtracts failed: expected: <1> but was: <2>")
			if strings.Contains(read(out), "example.AppTest.adds") {
				t.Errorf("only the failing test may be named:\n%s", out)
			}
			// A report cut short names itself, never a missing test run.
			broken := javaTestMavenRun(map[string]string{"target/surefire-reports/TEST-example.AppTest.xml": "<testsuite><testcase"})
			_, out, code = javaRun(t, javaRepo(t, map[string]string{"pom.xml": javaPom(javaTestMaven)}), "java-test", map[string][]response{"mvn": {broken}})
			sees(t, out, code, 1, "target/surefire-reports/TEST-example.AppTest.xml is not a JUnit XML report; rerun the build")
		}},
		{"coverage below goal fails", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"pom.xml": javaPom(javaTestMaven)})
			run := javaTestMavenRun(map[string]string{
				"target/surefire-reports/TEST-example.AppTest.xml": javaTestPassed,
				"target/site/jacoco/jacoco.xml":                    javaTestCoverage("4", "1"),
			})
			_, out, code := javaRun(t, repo, "java-test", map[string][]response{"mvn": {run}})
			sees(t, out, code, 1, "target/site/jacoco/jacoco.xml: lines 4/5 covered; the goal is 100%")
			// A report without line counts can't show the goal is met.
			noLines := javaTestMavenRun(map[string]string{
				"target/surefire-reports/TEST-example.AppTest.xml": javaTestPassed,
				"target/site/jacoco/jacoco.xml":                    `<report name="app"><counter type="INSTRUCTION" missed="0" covered="9"/></report>`,
			})
			_, out, code = javaRun(t, javaRepo(t, map[string]string{"pom.xml": javaPom(javaTestMaven)}), "java-test", map[string][]response{"mvn": {noLines}})
			sees(t, out, code, 1, "target/site/jacoco/jacoco.xml counts no lines; compile with debug information")
		}},
		{"aggregate and uncovered module reports", func(t *testing.T) {
			t.Parallel()
			files := map[string]string{
				"pom.xml":                      javaParent(javaTestMaven, "core", "extra", "it", "kt", "coverage"),
				"core/pom.xml":                 javaTestModule("core"),
				"core/src/main/java/Core.java": "class Core {}\n",
				// Maven's aggregate covers the modules it holds a group for, so a
				// source another group lists leaves extra uncovered.
				"extra/pom.xml":                 javaTestModule("extra"),
				"extra/src/main/java/Core.java": "class Core {}\n",
				// A module of tests alone has no code of its own to cover.
				"it/pom.xml":                           javaTestModule("it"),
				"it/src/test/java/example/FlowIT.java": "class FlowIT {}\n",
				// A module of other languages has code the aggregate may miss.
				"kt/pom.xml":                         javaTestModule("kt"),
				"kt/src/main/kotlin/example/Util.kt": "package example\nobject Util\n",
				"coverage/pom.xml":                   javaTestModule("coverage"),
			}
			reports := map[string]string{
				"core/target/surefire-reports/TEST-example.CoreTest.xml": javaTestPassed,
				"core/target/site/jacoco/jacoco.xml":                     javaTestCoverage("2", "3"),
				"extra/target/site/jacoco/jacoco.xml":                    javaTestCoverage("3", "2"),
				"it/target/site/jacoco/jacoco.xml":                       `<report name="it"><counter type="INSTRUCTION" missed="0" covered="0"/></report>`,
				"kt/target/site/jacoco/jacoco.xml":                       javaTestCoverage("1", "4"),
				"coverage/target/site/jacoco-aggregate/jacoco.xml":       javaTestAggregate("4", "1", `<group name="core">`+javaTestListed("", "Core.java")+`</group>`),
			}
			_, out, code := javaRun(t, javaRepo(t, files), "java-test", map[string][]response{"mvn": {javaTestMavenRun(reports)}})
			sees(t, out, code, 1, "coverage/target/site/jacoco-aggregate/jacoco.xml: lines 4/5 covered", "extra/target/site/jacoco/jacoco.xml: lines 3/5 covered", "kt/target/site/jacoco/jacoco.xml: lines 1/5 covered")
			javaTestNotJudged(t, out, "core/target/site/jacoco/jacoco.xml", "it/target/site/jacoco/jacoco.xml")
			// An aggregate that does not parse covers what is unknown, so it is
			// named alone.
			reports["coverage/target/site/jacoco-aggregate/jacoco.xml"] = "<report"
			_, out, code = javaRun(t, javaRepo(t, files), "java-test", map[string][]response{"mvn": {javaTestMavenRun(reports)}})
			sees(t, out, code, 1, "coverage/target/site/jacoco-aggregate/jacoco.xml is not a JaCoCo XML report; rerun the build")
			javaTestNotJudged(t, out, "extra/target/site/jacoco/jacoco.xml")
			javaTestGradleAggregate(t)
		}},
		{"build that ran no test fails", func(t *testing.T) {
			t.Parallel()
			_, out, code := javaRun(t, javaRepo(t, map[string]string{"pom.xml": javaPom(javaTestMaven)}), "java-test", map[string][]response{"mvn": {javaTestMavenRun(nil)}})
			sees(t, out, code, 1, "the Maven build at . has product code such as "+javaSource+" and no tests; add tests under src/test/java")
			karate := javaRepo(t, map[string]string{
				"pom.xml":                        javaPom(javaTestMaven + "<dependency><groupId>io.karatelabs</groupId><artifactId>karate-junit6</artifactId></dependency>"),
				"src/test/resources/app.feature": "Feature: app\n",
			})
			_, out, code = javaRun(t, karate, "java-test", map[string][]response{"mvn": {javaTestMavenRun(nil)}})
			sees(t, out, code, 1, "the Maven build at . has test files but ran no test", "Karate: Karate runs feature files only from a test class with a @Karate.Test method")
		}},
		{"testcontainers without docker fails", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{
				"pom.xml":                         javaPom(javaTestMaven),
				"src/test/java/example/DbIT.java": "import org.testcontainers.containers.PostgreSQLContainer;\nclass DbIT {}\n",
			})
			down := response{Match: []string{"info"}, Stderr: "Cannot connect to the Docker daemon. Is the docker daemon running?\n", Exit: 1}
			rp, out, code := javaRun(t, repo, "java-test", map[string][]response{"mvn": {javaTestMavenRun(nil)}, "docker": {down}})
			sees(t, out, code, 1, "Testcontainers needs a running Docker; start Docker or set DOCKER_HOST, then rerun")
			if calls := rp.calls("mvn"); len(calls) != 0 {
				t.Errorf("the build ran %q, want no run without Docker", calls)
			}
			// A hung Docker fails fast rather than holding the check to its limit.
			hung := response{Match: []string{"info"}, Wait: 60}
			start := time.Now()
			_, out, code = javaRun(t, repo, "java-test", map[string][]response{"mvn": {javaTestMavenRun(nil)}, "docker": {hung}})
			sees(t, out, code, 1, "Testcontainers needs a running Docker")
			if took := time.Since(start); took > 30*time.Second {
				t.Errorf("a hung Docker held the check for %s, want under 30s", took)
			}
		}},
		{"build without jacoco reports gets snippet", func(t *testing.T) {
			t.Parallel()
			repo := javaRepo(t, map[string]string{"pom.xml": javaPom("")})
			_, out, code := javaRun(t, repo, "java-test", map[string][]response{"mvn": {javaTestMavenRun(nil)}})
			sees(t, out, code, 1, "JaCoCo is not in the Maven build at .; add it:")
			javaSeesSnippet(t, out, "<artifactId>jacoco-maven-plugin</artifactId>")
			noXML := response{Match: []string{"jacocoTestReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
				"build/test-results/test/TEST-example.AppTest.xml": javaTestPassed,
				"build/reports/jacoco/test/html/index.html":        "<html></html>\n",
			}}
			_, out, code = javaRun(t, javaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]response{"gradle": {noXML}})
			sees(t, out, code, 1, "no JaCoCo XML report after the Gradle build at . ran its tests; add:")
			javaSeesSnippet(t, out, "xml.required = true")
		}},
		{"java test with no build passes", func(t *testing.T) {
			t.Parallel()
			_, out, code := javaRun(t, example(t, "happyRepo", nil), "java-test", nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "java-test")
		}},
		{"passing tests use gradle build cache", func(t *testing.T) {
			t.Parallel()
			maven := javaTestMavenRun(map[string]string{
				"target/surefire-reports/TEST-example.AppTest.xml": javaTestPassed,
				"target/failsafe-reports/TEST-example.AppIT.xml":   javaTestPassed,
				"target/site/jacoco/jacoco.xml":                    javaTestFull,
			})
			rp, out, code := javaRun(t, javaRepo(t, map[string]string{"pom.xml": javaPom(javaTestMaven)}), "java-test", map[string][]response{"mvn": {maven}})
			sees(t, out, code, 0, "All 1 checks passed")
			javaCalledOnce(t, rp, "mvn")
			gradle := response{Match: []string{"test", "jacocoTestReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
				"build/test-results/test/TEST-example.AppTest.xml":          javaTestPassed,
				"build/test-results/integrationTest/TEST-example.AppIT.xml": javaTestPassed,
				"build/reports/jacoco/test/jacocoTestReport.xml":            javaTestFull,
			}}
			home := "GRADLE_USER_HOME=" + t.TempDir()
			rp, out, code = javaRun(t, javaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]response{"gradle": {gradle}}, home)
			sees(t, out, code, 0, "All 1 checks passed")
			javaCalledOnce(t, rp, "gradle")
			// Clean rules rerun the tasks on every Gradle version; --rerun needs 7.6.
			// The build cache then restores the reports of unchanged tests.
			if !rp.called("gradle", "cleanTest test cleanJacocoTestReport jacocoTestReport", "--build-cache") || rp.called("gradle", "--rerun") {
				t.Errorf("gradle calls = %q, want clean rules before each task, --build-cache and no --rerun", rp.calls("gradle"))
			}
			// Reports the cache restores are written during the run, so they count.
			rp, out, code = javaRun(t, javaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]response{"gradle": {javaTestCached}}, home)
			sees(t, out, code, 0, "All 1 checks passed")
			if !rp.called("gradle", "--build-cache") {
				t.Errorf("gradle calls = %q, want --build-cache", rp.calls("gradle"))
			}
			// Gradle reads any value but true as off, from the build or the user home.
			user := t.TempDir()
			mustDo(t, os.WriteFile(filepath.Join(user, "gradle.properties"), []byte("org.gradle.caching=False\n"), 0o644))
			javaTestUncached(t, "org.gradle.caching = no\n", gradle, home)
			javaTestUncached(t, "org.gradle.caching=true\n", gradle, "GRADLE_USER_HOME="+user)
		}},
	})
}

func Test0006_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"one page per acceptance passes", func(t *testing.T) {
			t.Parallel()
			view := storybookView("#Preview(\"0001.1 day\", traits: .sizeThatFitsLayout) { BadgeView() }\n")
			pages := "let pages = [\n    BookPage(title: \"0001.2\") { BadgeView() },\n]\n"
			out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view, "App/Pages.swift": pages})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "swift-storybook-requirements", "2")
		}},
		{"two pages sharing an acceptance fail", func(t *testing.T) {
			t.Parallel()
			view := storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n\n#Preview(\"0001.1 night\") { BadgeView() }\n")
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
			sees(t, out, code, 1, "App/BadgeView.swift:9 - page \"0001.1 night\" repeats acceptance 0001.1 from line 7; draw each variant as a labeled row on one page",
				"App/Shown.swift:4 - page \"0001.1 shown\" repeats acceptance 0001.1 from App/BadgeView.swift:7")
		}},
		{"preview and BookPage sharing an acceptance fail", func(t *testing.T) {
			t.Parallel()
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.2 empty\") { BadgeView() }\n")})
			sees(t, out, code, 1, "App/Shown.swift:5 - page \"0001.2 shown\" repeats acceptance 0001.2 from App/BadgeView.swift:7; draw each variant as a labeled row on one page")
		}},
		{"configured pages per acceptance", func(t *testing.T) {
			t.Parallel()
			two := storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n\n#Preview(\"0001.1 night\") { BadgeView() }\n")
			three := two + "\n#Preview(\"0001.1 dusk\") { BadgeView() }\n"
			raised := `{"swiftStorybookRequirements": {"maxPagesPerAcceptance": 2}}` + "\n"
			out, code := storybookCheck(t, map[string]string{".fitnessrc.json": raised, storybookShown: storybookPages, "App/BadgeView.swift": two})
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = storybookCheck(t, map[string]string{".fitnessrc.json": raised, storybookShown: storybookPages, "App/BadgeView.swift": three})
			sees(t, out, code, 1, "App/BadgeView.swift:11 - page \"0001.1 dusk\" repeats acceptance 0001.1 from line 7")
			zero := `{"swiftStorybookRequirements": {"maxPagesPerAcceptance": 0}}` + "\n"
			out, code = storybookCheck(t, map[string]string{".fitnessrc.json": zero, storybookShown: storybookPages, "App/BadgeView.swift": two})
			sees(t, out, code, 1, "App/BadgeView.swift:9 - page \"0001.1 night\" repeats acceptance 0001.1 from line 7")
		}},
		{"unnamed preview or page fails", func(t *testing.T) {
			t.Parallel()
			pages := "let pages = [BookPage(\"Pages.swift\", 1, content: Card(title: \"0001.1 day\"))]\n"
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview {\n    BadgeView()\n}\n"), "App/Pages.swift": pages})
			sees(t, out, code, 1, "App/BadgeView.swift:7 - page has no name; name it by the acceptance it shows, such as #Preview(\"0001.1 label\")",
				"App/Pages.swift:1 - page has no name")
		}},
		{"name without acceptance ID fails", func(t *testing.T) {
			t.Parallel()
			pages := "let pages = [BookPage(title: \"BadgeView: empty\") { BadgeView() }]\n"
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"day\") { BadgeView() }\n"), "App/Pages.swift": pages})
			sees(t, out, code, 1, "App/BadgeView.swift:7 - page \"day\" does not start with an acceptance ID; name it such as \"0001.1 day\"",
				"App/Pages.swift:1 - page \"BadgeView: empty\" does not start with an acceptance ID")
		}},
		{"chained acceptance keeps one page", func(t *testing.T) {
			t.Parallel()
			chained := strings.Replace(storybookDoc, "        - When the photo shows\n            - Then it has no badge\n",
				"        1. When the photo shows\n            - Then it has no badge\n        2. When the tester taps it\n            - Then the badge stays hidden\n", 1)
			out, code := storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained})
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained, "App/BadgeView.swift": storybookView("#Preview(\"0001.2.2 tapped\") { BadgeView() }\n")})
			sees(t, out, code, 1, `App/BadgeView.swift:7 - page "0001.2.2 tapped" names step 0001.2.2; one page shows every step of an acceptance, so name it "0001.2 tapped"`)
			out, code = storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained, "App/BadgeView.swift": storybookView("#Preview(\"0001.9.2 tapped\") { BadgeView() }\n")})
			sees(t, out, code, 1, `App/BadgeView.swift:7 - page "0001.9.2 tapped" does not start with an acceptance ID`)
		}},
		{"page named by undefined ID fails", func(t *testing.T) {
			t.Parallel()
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.9 day\") { BadgeView() }\n")})
			sees(t, out, code, 1, "App/BadgeView.swift:7 - page \"0001.9 day\" names 0001.9, which no doc in docs/requirements defines")
		}},
		{"duplicate page name fails", func(t *testing.T) {
			t.Parallel()
			pages := "\nlet pages = [BookPage(title: \"0001.1 day\") { BadgeView() }]\n"
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n"), "App/Pages.swift": pages})
			sees(t, out, code, 1, "App/Pages.swift:2 - page \"0001.1 day\" repeats the name of the page at App/BadgeView.swift:7")
		}},
		{"variable or interpolated name fails", func(t *testing.T) {
			t.Parallel()
			pages := "let name = \"0001.1 day\"\nlet pages = [BookPage(title: name, file: #file) { BadgeView() }]\n"
			view := storybookView("#Preview(\"0001.1 \\(mode)\") { BadgeView() }\n")
			out, code := storybookCheck(t, map[string]string{"App/Pages.swift": pages, "App/BadgeView.swift": view})
			sees(t, out, code, 1, "App/Pages.swift:2 - page name name is not a plain literal; write the acceptance name, such as \"0001.1 label\"",
				"App/BadgeView.swift:7 - page name \"0001.1 \\(mode)\" is not a plain literal")
		}},
		{"no pages to judge passes", func(t *testing.T) {
			t.Parallel()
			view := storybookView("// #Preview { BadgeView() }\n/* #Preview { BadgeView() } */\nlet note = \"#Preview { BadgeView() }\"\n")
			out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "swift-storybook-requirements", "2")
			out, code = storybookCheck(t, map[string]string{storybookShown: ""})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "swift-storybook-requirements")
			out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview { BadgeView() }\n"), "docs/requirements/0001-badge.md": ""})
			sees(t, out, code, 0, "All 1 checks passed")
			if strings.Contains(out, "page has no name") {
				t.Errorf("want no finding without requirement docs:\n%s", out)
			}
		}},
		{"repo without Storybook passes", func(t *testing.T) {
			t.Parallel()
			view := storybookView("import StorybookHistoryKit\n// import Storybook\nlet note = \"\"\"\nimport StorybookKit\n\"\"\"\n\n#Preview {\n    BadgeView()\n}\n")
			out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view})
			showsInfo(t, out, code, "swift-storybook-requirements", "no Storybook in this repo: no Swift file imports a Storybook module or declares a BookPage, so acceptances need no pages")
			if strings.Contains(out, "has no Storybook page") || strings.Contains(out, "page has no name") {
				t.Errorf("want no finding in a repo without Storybook:\n%s", out)
			}
		}},
		{"Storybook import forms are judged", func(t *testing.T) {
			t.Parallel()
			for _, imports := range []string{
				"@testable import StorybookKit\n",
				"import SwiftUIStorybook\n",
				"#if DEBUG\n    @_spi(Storybook, Internal) import StorybookSandbox\n#endif\n",
				"import SwiftUI; import Storybook\n",
				"let pages = [BookPage(title: \"0001.1 day\") { Text(verbatim: \"\") }]\n",
			} {
				out, code := storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: imports})
				sees(t, out, code, 1, "docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
			}
		}},
		{"acceptance with no page fails", func(t *testing.T) {
			t.Parallel()
			catalog := "import Storybook\n\nenum AppStorybook {}\n"
			out, code := storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: catalog, "App/BadgeView.swift": storybookView("")})
			sees(t, out, code, 1, "docs/requirements/0001-badge.md:15 - acceptance 0001.1 has no Storybook page; add one named \"0001.1 <label>\"",
				"docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
			out, code = storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: catalog, "App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n")})
			sees(t, out, code, 1, "docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
			if strings.Contains(out, "acceptance 0001.1") {
				t.Errorf("want no finding for the acceptance a page names:\n%s", out)
			}
		}},
		{"literal in SwiftUI text fails", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{
				"App/ShareView.swift": localizedView,
				"App/StateView.swift": localizedStateView,
				"App/CountView.swift": localizedCountView,
				"App/HelpView.swift":  "import SwiftUI\n\nstruct HelpView: View {\n    var body: some View {\n        Text(\"\"\"\n            Share a \"photo\"\n            with a friend\n            \"\"\")\n    }\n}\n",
			})
			sees(t, out, code, 1,
				`App/StateView.swift:8 - "Shared" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference`,
				`App/StateView.swift:8 - "Private" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference`,
				`App/StateView.swift:9 - "Nobody" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference`,
				`App/StateView.swift:11 - "Photos" passed to Label skips the String Catalog; use a Words reference for catalog entry "noun.photo"`,
				`App/CountView.swift:7 - "photo" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference`,
				`App/CountView.swift:7 - "photos" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference`,
				"App/HelpView.swift:5 - \"\"\"\n            Share a \"photo\"\n            with a friend\n            \"\"\" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference",
				`App/ShareView.swift:7 - "Share" passed to Button skips the String Catalog; use a Words reference for catalog entry "verb.share"`,
				`App/ShareView.swift:8 - "Library" passed to navigationTitle skips the String Catalog; add it to the String Catalog and use a Words reference`,
				`App/ShareView.swift:9 - "\(count) shared" passed to Text skips the String Catalog; add it to the String Catalog and use a Words reference`)
		}},
		{"references and values pass", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{"App/PhotoView.swift": `import SwiftUI

struct PhotoView: View {
    let caption: String
    let count: Int
    let titles: [String: String]

    var body: some View {
        VStack {
            Label(Words.Noun.photo, systemImage: "photo")
            Text(caption)
            Text("·")
            Text("\(count)")
            Text("\(titles["first"] ?? Words.Noun.photo)")
            Button(Words.Verb.share) {}
                .accessibilityLabel(Words.Verb.share)
        }
    }
}
`})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"test typing catalog label fails", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{"Tests/Requirements/ShareTests.swift": `import Testing

struct ShareTests {
    func tapsShare(_ user: User) async throws {
        try await user.taps(.button("Share"))
    }
}
`})
			sees(t, out, code, 1,
				`Tests/Requirements/ShareTests.swift:5 - "Share" types catalog entry "verb.share"; read it through a Words reference`)
		}},
		{"test captions and prices pass", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{"Tests/Requirements/Support/Guest.swift": `struct Guest {
    let user: User

    func sharesCaptionedPhoto() async throws {
        try await user.writes("Swim at the lake")
        try await user.sees(.text("$9.99"))
        try await user.sees(.image("Sparkle"))
        try await user.doesNotSee(.button("Archive"))
        try await user.taps(.button(Words.Verb.share))
        let label = "Share"
        try await user.writes(label)
        try await user.adds(Photo(title: "Photos", caption: "Share"))
    }
}
`})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"allow pattern with reason passes", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{
				".fitnessrc.json": `{"localizedText": {"reference": "Words", "allow": [
  {"pattern": "Beta [0-9]+", "reason": "build tags stay in English for testers"}
]}}` + "\n",
				"App/BuildTag.swift": "import SwiftUI\n\nstruct BuildTag: View {\n    var body: some View { Text(\"Beta 3\") }\n}\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"allow pattern without reason fails", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{
				".fitnessrc.json": `{"localizedText": {"reference": "Words", "allow": [
  {"pattern": "Beta [0-9]+"},
  {"pattern": "(Beta", "reason": "build tags stay in English for testers"}
]}}` + "\n",
			})
			sees(t, out, code, 1,
				`localizedText.allow pattern "Beta [0-9]+" needs a reason`,
				`localizedText.allow pattern "(Beta" is not a regular expression`)
		}},
		{"app without string catalog passes", func(t *testing.T) {
			t.Parallel()
			out, code := localizedText(t, map[string]string{
				"App/Localizable.xcstrings": "",
				"App/ShareView.swift":       localizedView,
			})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "localized-text")
		}},
	})
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

// withSource is doc with the Source line older docs end their Measurement
// with, citing cite. It sits on line 19.
func withSource(doc, cite string) string {
	return strings.Replace(doc, "## Requirements", "Source: "+cite+"\n\n## Requirements", 1)
}

// docOf is goRepo with requirement 0001 replaced by body at path.
func docOf(path, body string) map[string]string {
	return with(goRepo(), map[string]string{"docs/requirements/0001-adds-numbers.md": "", path: body})
}

// cappedRepo is pairedRepo whose 0001 holds count tested acceptances, under
// a positive limit as requirements.maxAcceptances or else the default.
func cappedRepo(limit, count int) map[string]string {
	var ids, tests []string
	for n := 1; n <= count; n++ {
		ids = append(ids, fmt.Sprintf("0001.%d", n))
		tests = append(tests, fmt.Sprintf("Test0001_%d", n))
	}
	rc := `{"policy": "external", "checks": ["requirements"]}`
	if limit > 0 {
		rc = fmt.Sprintf(`{"policy": "external", "checks": ["requirements"], "requirements": {"maxAcceptances": %d}}`, limit)
	}
	return with(pairedRepo(), map[string]string{
		".fitnessrc.json":   rc + "\n",
		"reqs/0001_test.go": testFile("reqs", tests...),
		pairedDoc:           requirementDoc(ids...),
	})
}

// reqs is goRepo with requirement 0001's Requirements replaced.
func reqs(body string) map[string]string {
	return docOf("docs/requirements/0001-adds-numbers.md", docWith(body))
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

// goTest runs `fitness-install -- go-test` on happyRepo with files written
// over it and env added.
func goTest(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, "go-test")
}

// goTestModule is a Go module in dir whose test file holds body.
func goTestModule(dir, body string) map[string]string {
	return map[string]string{
		dir + "go.mod":    "module example.com/m\n\ngo 1.22\n",
		dir + "m_test.go": "package m\n\nimport \"testing\"\n\n" + body,
	}
}

// goTestBroken is a test that fails with the message "boom" on line 5.
const goTestBroken = "func TestBroken(t *testing.T) { t.Error(\"boom\") }\n"

// goVetMod is a minimal go.mod for a Given's module.
const goVetMod = "module m\n\ngo 1.21\n"

// goVetBad is Go source the compiler accepts but go vet flags.
const goVetBad = "package m\n\nimport \"fmt\"\n\nfunc F() { fmt.Printf(\"%d\\n\", \"x\") }\n"

// goVetFinding is what go vet reports for goVetBad.
const goVetFinding = `m.go:5:24: fmt.Printf format %d has arg "x" of wrong type string`

// goVetModule returns the files of a module in dir holding src.
func goVetModule(dir, src string) map[string]string {
	return map[string]string{dir + "/go.mod": goVetMod, dir + "/m.go": src}
}

// goVet runs `fitness-install -- go-vet` in repo with env added.
func goVet(t *testing.T, repo string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, repo, env, "go-vet")
}

// gofmtCheck runs `fitness-install -- gofmt` on happyRepo with files written
// over it and env added.
func gofmtCheck(t *testing.T, env []string, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, "gofmt")
}

// gofmtTidy is a Go file gofmt leaves as it is; gofmtMessy is one it rewrites.
const (
	gofmtTidy  = "package m\n\nfunc F() int { return 1 }\n"
	gofmtMessy = "package m\nfunc   G( ) int {return 2}\n"
)

// noEslintDisableRun runs `fitness-install -- no-eslint-disable <args>` in repo.
func noEslintDisableRun(t *testing.T, repo string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, repo, nil, append([]string{"no-eslint-disable"}, args...)...)
}

// noEslintDisableCheck runs the check on happyRepo with files written over it.
func noEslintDisableCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return noEslintDisableRun(t, example(t, "happyRepo", files))
}

// noEslintDisableNext is a next-line directive comment.
const noEslintDisableNext = "// eslint-disable-next-line no-console\n"

// noEslintDisableExts are the extensions the check scans.
var noEslintDisableExts = []string{".cts", ".mts", ".ts", ".cjs", ".js", ".mjs", ".tsx"}

// coverageExclude runs `fitness-install -- vitest-coverage-exclude` on
// happyRepo with files written over it.
func coverageExclude(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "vitest-coverage-exclude")
}

// coverageExcludeRemove is the failure line naming one quoted entry.
func coverageExcludeRemove(quoted string) string {
	return "Vitest coverage exclude must be empty, so coverage judges every file; remove: " + quoted
}

// coverageExcludeConfig is a CommonJS Vitest config with exclude as given.
func coverageExcludeConfig(exclude string) string {
	return "module.exports = { test: { coverage: { exclude: " + exclude + " } } };\n"
}

// coverageExcludeOneFile matches a results row that passed judging 1 file.
var coverageExcludeOneFile = regexp.MustCompile(`│ passed\s+│ 1\s+│`)

// vcfConfig is a local Vitest config holding the given coverage thresholds
// over src/**/*.ts, the way a consumer repo writes one.
func vcfConfig(branches string) string {
	return "export default {\n  test: {\n    coverage: {\n      include: ['src/**/*.ts'],\n      provider: 'v8',\n" +
		"      thresholds: { branches: " + branches + ", functions: 100, lines: 100, statements: 100 },\n" +
		"    },\n    include: ['src/**/*.test.ts'],\n  },\n};\n"
}

// vcfProject is a small JavaScript project whose one function is fully tested.
var vcfProject = map[string]string{
	"package.json":      `{"type": "module"}` + "\n",
	"vitest.config.mjs": vcfConfig("100"),
	"src/add.ts":        "export const add = (a: number, b: number): number => a + b;\n",
	"src/add.test.ts": "import { expect, test } from 'vitest';\nimport { add } from './add';\n\n" +
		"test('adds', () => {\n  expect(add(1, 2)).toBe(3);\n});\n",
}

// vcfRepo is happyRepo with vcfProject and then files written over it.
func vcfRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	all := map[string]string{}
	for _, set := range []map[string]string{vcfProject, files} {
		for name, body := range set {
			all[name] = body
		}
	}
	return example(t, "happyRepo", all)
}

func vcfRun(t *testing.T, repo string, env []string) (string, int) {
	t.Helper()
	return fitness(t, repo, env, "vitest-coverage-full")
}

// vcfTable is the coverage table Vitest's v8 provider prints after a run.
func vcfTable(rows string) string {
	rule := "----------|---------|----------|---------|---------|-------------------\n"
	return " % Coverage report from v8\n" + rule +
		"File      | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s \n" + rule + rows + rule
}

// vcfFullRun is what `vitest run --coverage` prints when every line ran.
var vcfFullRun = response{
	Match: []string{"run --coverage"},
	Stdout: "\n RUN  v3.2.4 /repo\n\n \u2713 src/add.test.ts (1 test) 2ms\n\n" +
		" Test Files  1 passed (1)\n      Tests  1 passed (1)\n   Start at  10:15:02\n   Duration  412ms\n\n" +
		vcfTable("All files |     100 |      100 |     100 |     100 |                   \n"+
			" add.ts   |     100 |      100 |     100 |     100 |                   \n"),
}

// vcfShortRun is what `vitest run --coverage` prints when sub.ts never ran:
// the table on stdout, then one threshold error per metric on stderr.
var vcfShortRun = response{
	Match: []string{"run --coverage"},
	Stdout: "\n RUN  v3.2.4 /repo\n\n \u2713 src/add.test.ts (1 test) 2ms\n\n" +
		" Test Files  1 passed (1)\n      Tests  1 passed (1)\n   Start at  10:15:02\n   Duration  418ms\n\n" +
		vcfTable("All files |      50 |      100 |      50 |      50 |                   \n"+
			" add.ts   |     100 |      100 |     100 |     100 |                   \n"+
			" sub.ts   |       0 |      100 |       0 |       0 | 1                 \n"),
	Stderr: "ERROR: Coverage for lines (50%) does not meet global threshold (100%)\n" +
		"ERROR: Coverage for functions (50%) does not meet global threshold (100%)\n" +
		"ERROR: Coverage for statements (50%) does not meet global threshold (100%)\n",
	Exit: 1,
}

// vcfRanWithLocalConfig fails t unless fitness ran vitest once with
// `run --coverage` alone, leaving Vitest to load the repo's own config.
func vcfRanWithLocalConfig(t *testing.T, rp replays) {
	t.Helper()
	if calls := rp.calls("vitest"); len(calls) != 1 || calls[0] != "run --coverage" {
		t.Errorf("vitest calls = %q, want one \"run --coverage\"", calls)
	}
}

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

// spotlessMaven declares Spotless in a Maven build.
const spotlessMaven = "<plugin><groupId>com.diffplug.spotless</groupId><artifactId>spotless-maven-plugin</artifactId></plugin>"

// spotlessClean is the build's answer when Spotless finds nothing to rewrite.
var spotlessClean = response{Match: []string{"spotless:check"}, Stdout: "[INFO] BUILD SUCCESS\n"}

// spotlessGradleClean is Gradle's answer for a clean spotlessCheck.
var spotlessGradleClean = response{Match: []string{"spotlessCheck"}, Stdout: "> Task :spotlessCheck\n\nBUILD SUCCESSFUL in 2s\n"}

// spotlessViolations is Maven's answer when Spotless would rewrite the
// Java source: the file under the violations header, then its diff.
var spotlessViolations = response{
	Match: []string{"spotless:check"},
	Stdout: "[INFO] Scanning for projects...\n" +
		"[ERROR] Failed to execute goal com.diffplug.spotless:spotless-maven-plugin:3.10.4:check (default-cli) on project app: The following files had format violations:\n" +
		"[ERROR]     " + javaSource + "\n" +
		"[ERROR]         @@ -1 +1 @@\n" +
		"[ERROR]         -class App {  }\n" +
		"[ERROR]         +class App {}\n" +
		"[ERROR] Run 'mvn spotless:apply' to fix these violations.\n",
	Exit: 1,
}

// spotlessGradleNoJDK is Gradle's answer, saved from a real run, when the
// toolchain needs a JDK the machine lacks.
var spotlessGradleNoJDK = response{
	Match: []string{"spotlessCheck"},
	Stdout: "FAILURE: Build failed with an exception.\n\n* What went wrong:\n" +
		"Could not determine the dependencies of task ':spotlessJava'.\n" +
		"> Could not resolve all dependencies for configuration ':compileClasspath'.\n" +
		"   > Failed to calculate the value of task ':compileJava' property 'javaCompiler'.\n" +
		"      > Cannot find a Java installation on your machine (Mac OS X 27.0.1 aarch64) matching: {languageVersion=17, vendor=any vendor, implementation=vendor-specific, nativeImageCapable=false}. Toolchain download repositories have not been configured.\n" +
		"\n* Try:\n> Run with --stacktrace option to get the stack trace.\n\nBUILD FAILED in 1s\n",
	Exit: 1,
}

// cspell:ignore jdks foojay adoptium hotspot openj

// javaVersionRun runs `fitness-install -- java-version` on a Java repo with
// files written over it and a JDK 21 running.
func javaVersionRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	_, out, code := javaRun(t, javaRepo(t, files), "java-version", nil)
	return out, code
}

// javaVersionRelease is a Maven build whose compiler release is release.
func javaVersionRelease(release string) string {
	return "<project><properties><maven.compiler.release>" + release + "</maven.compiler.release></properties></project>\n"
}

// javaVersionParentPom is where a local Maven repository keeps the parent
// pom com.example:parent:1.
const javaVersionParentPom = "com/example/parent/1/parent-1.pom"

// javaVersionParent maps java.version to the compiler release, as
// framework parents do.
const javaVersionParent = "<project><properties><java.version>17</java.version>" +
	"<maven.compiler.release>${java.version}</maven.compiler.release></properties></project>\n"

// javaVersionSees runs java-version on repo with home as the user's home
// and asserts it fails with want.
func javaVersionSees(t *testing.T, repo, home, want string) {
	t.Helper()
	rp := standIn(t, map[string][]response{"java": {javaJDK}})
	out, code := fitness(t, repo, append(rp.env, "JAVA_HOME=", "HOME="+home), "java-version")
	sees(t, out, code, 1, want)
}

// javaVersionNoGradle is a Gradle stand-in that fails if it ever runs, since
// the check must find toolchain JDKs without starting Gradle.
var javaVersionNoGradle = response{Stderr: "gradle must not run\n", Exit: 1}

// javaVersionMissing30 is the failure for a toolchain release no installed
// JDK has. No machine holds a JDK 30, so fixture folders decide the result.
const javaVersionMissing30 = "build.gradle asks for a JDK 30 toolchain, which is not installed; install JDK 30"

// javaVersionToolchain is a Gradle build whose toolchain asks for release.
func javaVersionToolchain(t *testing.T, release string, extra map[string]string) string {
	t.Helper()
	return javaVersionBuild(t, "languageVersion = JavaLanguageVersion.of("+release+")", extra)
}

// javaVersionBuild is a Gradle build whose toolchain block holds body.
func javaVersionBuild(t *testing.T, body string, extra map[string]string) string {
	t.Helper()
	files := map[string]string{"build.gradle": "java { toolchain { " + body + " } }\n"}
	for name, text := range extra {
		files[name] = text
	}
	return javaRepo(t, files)
}

// javaVersionJDK is the release file of an installed JDK of release.
func javaVersionJDK(release string) string {
	return javaVersionVendorJDK(release, "Eclipse Adoptium", "Hotspot")
}

// javaVersionVendorJDK is the release file of a JDK of release from vendor
// running the JVM variant.
func javaVersionVendorJDK(release, vendor, variant string) string {
	return "JAVA_VERSION=\"" + release + ".0.1\"\nIMPLEMENTOR=\"" + vendor + "\"\nJVM_VARIANT=\"" + variant + "\"\n"
}

// javaVersionGradleRun runs java-version on repo with home as the user's
// home and a Gradle stand-in, failing the test if Gradle ran.
func javaVersionGradleRun(t *testing.T, repo, home string) (string, int) {
	t.Helper()
	rp := standIn(t, map[string][]response{"java": {javaJDK}, "gradle": {javaVersionNoGradle}})
	out, code := fitness(t, repo, append(rp.env, "JAVA_HOME=", "GRADLE_USER_HOME=", "HOME="+home), "java-version")
	if calls := rp.calls("gradle"); len(calls) != 0 {
		t.Errorf("gradle calls = %q, want none", calls)
	}
	return out, code
}

// javaVersionPasses asserts java-version passes repo with home.
func javaVersionPasses(t *testing.T, repo, home string) {
	t.Helper()
	out, code := javaVersionGradleRun(t, repo, home)
	sees(t, out, code, 0, "All 1 checks passed")
}

// javaVersionFails asserts java-version fails repo with home, naming a
// missing JDK 30.
func javaVersionFails(t *testing.T, repo, home string) {
	t.Helper()
	out, code := javaVersionGradleRun(t, repo, home)
	sees(t, out, code, 1, javaVersionMissing30)
}

// javaVersionHome is a user home holding the files.
func javaVersionHome(t *testing.T, files map[string]string) string {
	t.Helper()
	home := t.TempDir()
	write(t, home, files)
	return home
}

// javaVersionFolder is a folder holding a JDK 30 release file.
func javaVersionFolder(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, map[string]string{"release": javaVersionJDK("30")})
	return dir
}

// javaVersionResolver is a settings script that registers a download
// repository.
const javaVersionResolver = "toolchainManagement { jvm { javaRepositories { repository('foojay') { resolverClass = FoojayResolver } } } }\n"

// javaVersionPropertiesCases covers the gradle.properties forms, build root
// and precedence.
func javaVersionPropertiesCases(t *testing.T) {
	jdks := javaVersionFolder(t)
	paths := func(text string) map[string]string { return map[string]string{"gradle.properties": text} }
	// gradle.properties lists JDK folders in every properties form: an equals
	// sign, a colon or a space, and a list continued across lines.
	for _, text := range []string{
		"org.gradle.java.installations.paths=" + jdks + "\n",
		"org.gradle.java.installations.paths : " + jdks + "\n",
		"org.gradle.java.installations.paths " + jdks + "\n",
		"# the list\norg.gradle.java.installations.paths=/no/such,\\\n    " + jdks + "\n",
	} {
		javaVersionPasses(t, javaVersionToolchain(t, "30", paths(text)), t.TempDir())
	}
	// The list splits on commas only, so a folder name may hold a colon.
	colon := filepath.Join(t.TempDir(), "a:b")
	write(t, colon, map[string]string{"release": javaVersionJDK("30")})
	javaVersionPasses(t, javaVersionToolchain(t, "30", paths("org.gradle.java.installations.paths="+colon+"\n")), t.TempDir())
	// A build in a folder reads its own gradle.properties, not the repo's.
	nested := javaRepo(t, map[string]string{
		"app/build.gradle":  "java { toolchain { languageVersion = JavaLanguageVersion.of(30) } }\n",
		"gradle.properties": "org.gradle.java.installations.paths=" + jdks + "\n",
	})
	out, code := javaVersionGradleRun(t, nested, t.TempDir())
	sees(t, out, code, 1, "app/build.gradle asks for a JDK 30 toolchain")
	// The Gradle user home's gradle.properties overrides the build's.
	repo := javaVersionToolchain(t, "30", paths("org.gradle.java.installations.paths="+jdks+"\n"))
	javaVersionFails(t, repo, javaVersionHome(t, map[string]string{".gradle/gradle.properties": "org.gradle.java.installations.paths=\n"}))
}

// javaVersionResolverCases covers what counts as a toolchain resolver.
func javaVersionResolverCases(t *testing.T) {
	settings := func(text string) map[string]string { return map[string]string{"settings.gradle": text} }
	build := func(extra map[string]string) string { return javaVersionToolchain(t, "30", extra) }
	// A registered download repository, or the Foojay plugin, lets Gradle
	// download the JDK.
	javaVersionPasses(t, build(settings(javaVersionResolver)), t.TempDir())
	javaVersionPasses(t, build(settings("plugins { id 'org.gradle.toolchains.foojay-resolver-convention' version '1.0.0' }\n")), t.TempDir())
	// A resolver in an init script of the Gradle user home passes too.
	javaVersionPasses(t, build(nil), javaVersionHome(t, map[string]string{".gradle/init.d/resolver.gradle": javaVersionResolver}))
	// An empty block registers none, and a commented-out one is no block.
	javaVersionFails(t, build(settings("toolchainManagement { jvm { } }\n")), t.TempDir())
	javaVersionFails(t, build(settings("// id 'org.gradle.toolchains.foojay-resolver-convention'\n/* "+javaVersionResolver+" */\n")), t.TempDir())
	// Turning auto-download off wins over a resolver.
	off := settings(javaVersionResolver)
	off["gradle.properties"] = "org.gradle.java.installations.auto-download=false\n"
	javaVersionFails(t, build(off), t.TempDir())
}

// javaVersionVendorCases covers vendor and implementation constraints.
func javaVersionVendorCases(t *testing.T) {
	home := func(vendor, variant string) string {
		return javaVersionHome(t, map[string]string{".gradle/jdks/jdk-30/release": javaVersionVendorJDK("30", vendor, variant)})
	}
	adoptium := func() string {
		return javaVersionBuild(t, "languageVersion = JavaLanguageVersion.of(30); vendor = JvmVendorSpec.ADOPTIUM", nil)
	}
	// A toolchain that names a vendor needs a JDK from it.
	javaVersionPasses(t, adoptium(), home("Eclipse Adoptium", "Hotspot"))
	javaVersionFails(t, adoptium(), home("Amazon.com Inc.", "Hotspot"))
	// A build that asks for the J9 implementation fails too.
	j9 := func() string {
		return javaVersionBuild(t, "languageVersion = JavaLanguageVersion.of(30); implementation = JvmImplementation.J9", nil)
	}
	javaVersionPasses(t, j9(), home("IBM Corporation", "Openj9"))
	javaVersionFails(t, j9(), home("IBM Corporation", "Hotspot"))
}

// cspell:ignore jacoco testcontainers karatelabs opentest sessioninfo Postgre

// javaTestMaven declares JaCoCo in a Maven build.
const javaTestMaven = "<plugin><groupId>org.jacoco</groupId><artifactId>jacoco-maven-plugin</artifactId></plugin>"

// javaTestGradle is a Gradle build that applies JaCoCo.
const javaTestGradle = "plugins { id 'java'; id 'jacoco' }\n"

// javaTestPassed is a JUnit XML report of one passing test.
const javaTestPassed = `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="AppTest" tests="1" failures="0" errors="0" skipped="0"><testcase name="adds" classname="example.AppTest" time="0.01"/></testsuite>
`

// javaTestFailed is a JUnit XML report of one passing and one failing test.
const javaTestFailed = `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="AppTest" tests="2" failures="1" errors="0" skipped="0">
<testcase name="adds" classname="example.AppTest" time="0.01"/>
<testcase name="subtracts" classname="example.AppTest" time="0.01"><failure message="expected: &lt;1&gt; but was: &lt;2&gt;" type="org.opentest4j.AssertionFailedError">org.opentest4j.AssertionFailedError: expected: &lt;1&gt; but was: &lt;2&gt;
	at example.AppTest.subtracts(AppTest.java:12)</failure></testcase>
</testsuite>
`

// javaTestCoverage is a JaCoCo XML report whose own totals hold covered
// and missed lines and cover every branch.
func javaTestCoverage(covered, missed string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><!DOCTYPE report PUBLIC "-//JACOCO//DTD Report 1.1//EN" "report.dtd">` +
		`<report name="app"><sessioninfo id="s" start="1" dump="2"/><package name="example"><counter type="LINE" missed="0" covered="1"/></package>` +
		`<counter type="LINE" missed="` + missed + `" covered="` + covered + `"/><counter type="BRANCH" missed="0" covered="2"/></report>
`
}

// javaTestFull is a JaCoCo report with every line and branch covered.
var javaTestFull = javaTestCoverage("5", "0")

// javaTestMavenRun is Maven's answer to verify: it writes files and passes,
// as it does past failing tests.
func javaTestMavenRun(files map[string]string) response {
	return response{Match: []string{"verify"}, Files: files, Stdout: "[INFO] BUILD SUCCESS\n"}
}

// javaTestGradleAggregate proves Gradle's aggregate, which lists packages
// without module groups, covers a project when it lists one of its
// sources by declared package.
func javaTestGradleAggregate(t *testing.T) {
	t.Helper()
	files := map[string]string{
		"settings.gradle":              "include 'core', 'extra', 'locked'\n",
		"build.gradle":                 "plugins { id 'java'; id 'jacoco-report-aggregation' }\n",
		"core/build.gradle":            javaTestGradle,
		"core/src/main/java/Core.java": "class Core {}\n",
		// A class the aggregate's filters leave out keeps core covered.
		"core/src/main/java/example/Port.java": "package example;\ninterface Port {}\n",
		// The same file name in another package is another source.
		"extra/build.gradle":                    javaTestGradle,
		"extra/src/main/java/example/Core.java": "/* a; b */\npackage example;\nclass Core {}\n",
		// A source that cannot be read is not listed.
		"locked/build.gradle":                    javaTestGradle,
		"locked/src/main/java/example/Core.java": "package example;\nclass Core {}\n",
	}
	repo := javaRepo(t, files)
	mustDo(t, os.Chmod(filepath.Join(repo, "locked/src/main/java/example/Core.java"), 0))
	run := response{Match: []string{"testCodeCoverageReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
		"core/build/test-results/test/TEST-example.CoreTest.xml":                 javaTestPassed,
		"core/build/reports/jacoco/test/jacocoTestReport.xml":                    javaTestCoverage("2", "3"),
		"extra/build/reports/jacoco/test/jacocoTestReport.xml":                   javaTestCoverage("3", "2"),
		"locked/build/reports/jacoco/test/jacocoTestReport.xml":                  javaTestCoverage("1", "4"),
		"build/reports/jacoco/testCodeCoverageReport/testCodeCoverageReport.xml": javaTestAggregate("5", "0", javaTestListed("", "Core.java")),
	}}
	_, out, code := javaRun(t, repo, "java-test", map[string][]response{"gradle": {run}})
	sees(t, out, code, 1, "extra/build/reports/jacoco/test/jacocoTestReport.xml: lines 3/5 covered", "locked/build/reports/jacoco/test/jacocoTestReport.xml: lines 1/5 covered")
	javaTestNotJudged(t, out, "core/build/reports/jacoco/test/jacocoTestReport.xml")
}

// javaTestModule is a Maven module's build file with its artifact id.
func javaTestModule(id string) string {
	return "<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId><artifactId>" + id + "</artifactId><version>1</version></project>\n"
}

// javaTestAggregate is an aggregate JaCoCo report with the given line
// totals that lists what listed holds: groups, or packages.
func javaTestAggregate(covered, missed, listed string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><report name="coverage"><sessioninfo id="s" start="1" dump="2"/>` + listed +
		`<counter type="LINE" missed="` + missed + `" covered="` + covered + `"/></report>
`
}

// javaTestListed is a report's package listing one source file.
func javaTestListed(pkg, source string) string {
	return `<package name="` + pkg + `"><sourcefile name="` + source + `"/></package>`
}

// javaTestNotJudged fails when the output names any of the reports.
func javaTestNotJudged(t *testing.T, out string, reports ...string) {
	t.Helper()
	for _, r := range reports {
		if strings.Contains(read(out), r) {
			t.Errorf("%s may not be judged:\n%s", r, out)
		}
	}
}

// javaTestCached is a Gradle run that restores the tests and coverage of
// an unchanged build from its build cache, rewriting their reports.
var javaTestCached = response{Match: []string{"test", "jacocoTestReport"}, Stdout: "> Task :test FROM-CACHE\n> Task :jacocoTestReport FROM-CACHE\nBUILD SUCCESSFUL in 1s\n", Files: map[string]string{
	"build/test-results/test/TEST-example.AppTest.xml": javaTestPassed,
	"build/reports/jacoco/test/jacocoTestReport.xml":   javaTestFull,
}}

// javaTestUncached asserts a Gradle build with props in its
// gradle.properties, run with env, reruns its tests without the build
// cache.
func javaTestUncached(t *testing.T, props string, gradle response, env string) {
	t.Helper()
	repo := javaRepo(t, map[string]string{"build.gradle": javaTestGradle, "gradle.properties": props})
	rp, out, code := javaRun(t, repo, "java-test", map[string][]response{"gradle": {gradle}}, env)
	sees(t, out, code, 0, "All 1 checks passed")
	if !rp.called("gradle", "cleanTest test") || rp.called("gradle", "--build-cache") {
		t.Errorf("gradle calls = %q, want the tests run without --build-cache", rp.calls("gradle"))
	}
}

// localizedCatalog is the fixture app's String Catalog: a live verb, a plural
// noun, a stale word, and Xcode's empty key.
const localizedCatalog = `{
  "sourceLanguage" : "en",
  "strings" : {
    "" : { "shouldTranslate" : false },
    "noun.photo" : {
      "localizations" : { "en" : { "variations" : { "plural" : {
        "one" : { "stringUnit" : { "state" : "translated", "value" : "Photo" } },
        "other" : { "stringUnit" : { "state" : "translated", "value" : "Photos" } }
      } } } }
    },
    "verb.archive" : {
      "extractionState" : "stale",
      "localizations" : { "en" : { "stringUnit" : { "state" : "translated", "value" : "Archive" } } }
    },
    "verb.share" : {
      "extractionState" : "extracted_with_value",
      "localizations" : { "en" : { "stringUnit" : { "state" : "new", "value" : "Share" } } }
    }
  },
  "version" : "1.0"
}
`

// localizedWords is the fixture app's reference type, which reads each
// label from the catalog by key.
const localizedWords = `import Foundation

enum Words {
    enum Verb {
        static var share: String {
            String(localized: "verb.share", defaultValue: "Share", comment: "Verb: share a photo")
        }
    }

    enum Noun {
        static var photo: String {
            String(localized: "noun.photo", defaultValue: "Photo", comment: "Noun: one photo")
        }
    }
}
`

// localizedApp is the fixture app with a String Catalog, its Words
// reference, and a config naming that reference, with files written over
// it. An empty body deletes a fixture file.
func localizedApp(t *testing.T, files map[string]string) string {
	t.Helper()
	app := map[string]string{
		".fitnessrc.json":           `{"localizedText": {"reference": "Words"}}` + "\n",
		"App/Localizable.xcstrings": localizedCatalog,
		"App/Words.swift":           localizedWords,
	}
	maps.Copy(app, files)
	return example(t, "happyRepo", app)
}

// localizedText runs `fitness-install -- localized-text` on the fixture
// app with files written over it.
func localizedText(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, localizedApp(t, files), nil, "localized-text")
}

// localizedView is a view that passes literals to a button and a title.
const localizedView = `import SwiftUI

struct ShareView: View {
    let count: Int

    var body: some View {
        Button("Share") {}
            .navigationTitle("Library")
        Text("\(count) shared")
    }
}
`

// localizedStateView passes literals as ternary and ?? branches.
const localizedStateView = `import SwiftUI

struct StateView: View {
    let isShared: Bool
    let name: String?

    var body: some View {
        Text(isShared ? "Shared" : "Private")
        Text(name ?? "Nobody")
        Text(isShared ? Words.Verb.share : Words.Noun.photo)
        Label("Photos", systemImage: isShared ? "person" : "lock")
    }
}
`

// localizedCountView passes literals as branches inside an interpolation.
const localizedCountView = `import SwiftUI

struct CountView: View {
    let count: Int

    var body: some View {
        Text("\(count == 1 ? "photo" : "photos")")
        Text("\(count)\(Words.Noun.photo)")
    }
}
`

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

// storybookDoc is the fixture app's requirement doc, defining 0001.1 and
// 0001.2.
const storybookDoc = `# 0001 Badge

## Why

Testers see which photos are new.

## Measurement

badges seen
-
photos opened

## Requirements

- 0001.1
    - Given Tester opens a new photo
        - When the photo shows
            - Then it has a badge
- 0001.2
    - Given Tester opens an old photo
        - When the photo shows
            - Then it has no badge
`

// storybookShown is the fixture app's page file that imports Storybook and
// shows both acceptances. A test's own pages are all that it judges.
const storybookShown = "App/Shown.swift"

// storybookPages is a storybookShown body that shows only 0001.2, so a test
// draws the pages of 0001.1 itself.
const storybookPages = "import Storybook\n\nlet shown = [BookPage(title: \"0001.2 shown\") { BadgeView() }]\n"

// storybookCatalog is a file that imports Storybook and declares no page, as
// an app with no pages yet has.
const storybookCatalog = "App/AppStorybook.swift"

// storybookCheck runs swift-storybook-requirements on happyRepo with the
// fixture requirement doc, storybookShown, and files written over them. An
// empty body deletes a fixture file.
func storybookCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	app := map[string]string{
		"docs/requirements/0001-badge.md": storybookDoc,
		storybookShown:                    "import Storybook\n\nlet shown = [\n    BookPage(title: \"0001.1 shown\") { BadgeView() },\n    BookPage(title: \"0001.2 shown\") { BadgeView() },\n]\n",
	}
	maps.Copy(app, files)
	return fitness(t, example(t, "happyRepo", app), nil, "swift-storybook-requirements")
}

// storybookView wraps previews in a SwiftUI view file whose first preview
// starts on line 7.
func storybookView(previews string) string {
	return "import SwiftUI\n\nstruct BadgeView: View {\n    var body: some View { Text(verbatim: \"\") }\n}\n\n" + previews
}
