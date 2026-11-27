package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func Test0007_1(t *testing.T) {
	t.Parallel()
	t.Run("0007.1.1 a test declared the way its language runs it owns its acceptance", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.1.1.1.1", "3c4f2fb1", func(_ req.A) {
			out, _ := swiftJudged(t, nil)
			acceptance.PassedWith(t, out, "requirements", "2")
			out, _ = swiftJudged(t, map[string]string{
				"Tests/AdderTests/AdderTests.swift": "import Testing\n\nstruct AdderSpec {\n  @Test(\"adds\") func test0001_1() {\n" + swiftThen(1, thenText) + "  }\n}\n",
			})
			acceptance.PassedWith(t, out, "requirements", "2")
			for _, owner := range []string{
				`test("0001.1 adds two numbers", () => {});`,
				`describe("adder", () => { it.each([1, 2])("0001.1 adds %i", (n) => { expect(n).toBeTruthy(); }); });`,
				"test.skip(`0001.1 adds two numbers`, () => {});",
				`test("0001.1 adds two numbers", () => { test.skip(process.env.CI === undefined, "0001.1 runs in CI"); });`,
			} {
				out, code := scriptChecked(t, append([]string{owner}, notTests...)...)
				acceptance.Sees(t, out, code, 0)
				acceptance.PassedWith(t, out, "requirements", "2")
			}
			jestFolder := map[string]string{"tests/adder.test.ts": "", "src/__tests__/adder.ts": vitestFile(`test("0001.1 adds two numbers", () => {});`)}
			out, code := checked(t, acceptance.With(scriptRepo(), jestFolder))
			acceptance.Sees(t, out, code, 0)
			acceptance.PassedWith(t, out, "requirements", "2")
			junitOwners(t)
		})
		req.Then(t, "0007.1.1.1.2", "d5fa66dc", func(_ req.A) {
			swiftOwnsNothing(t, xcTestFile()+"// func test0001_1() {}\n/* nested /* func test0001_1() {} */ */\nlet note = \"\\(\"func test0001_1() {}\")\"\nlet raw = #\"\"\"\n\" func test0001_1() {}\n\"\"\"#\n")
			swiftOwnsNothing(t, xcTestFile()+"\nstruct Helper {\n  func test0001_1() {}\n}\n")
			out, code := scriptChecked(t, notTests...)
			acceptance.Sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
			junitOwnsNothing(t)
		})
	})
	t.Run("0007.1.2 an acceptance with no such test fails naming the form to write", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.1.1.2.1", "e687c744", func(_ req.A) {
			swiftOwnsNothing(t, xcTestFile())
			out, code := checked(t, acceptance.With(scriptRepo(), map[string]string{"tests/adder.test.ts": ""}))
			acceptance.Sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
			out, code = checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest")}))
			acceptance.Sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
		})
	})
	t.Run("0007.1.3 a chain of steps in Node or Swift needs only its test", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.1.1.3.1", "8a958cdb", func(_ req.A) {
			for _, repo := range []map[string]string{
				acceptance.With(scriptRepo(), map[string]string{"tests/adder.test.ts": vitestFile(`test("0001.1 adds two numbers", () => {});`, nodeThen(2, "the total is even"))}),
				acceptance.With(swiftRepo(), map[string]string{swiftTest: strings.Replace(swiftTestingFile("test0001_1"), "  }\n}\n", swiftThen(2, "the total is even")+"  }\n}\n", 1)}),
			} {
				out, code := checked(t, acceptance.With(repo, map[string]string{acceptance.PairedDoc: acceptance.DocWith(acceptance.Chain("0001.1", 2))}))
				acceptance.Sees(t, out, code, 0)
			}
		})
	})
	t.Run("0007.1.4 a test named for no acceptance fails naming the form to rename it to", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.1.1.4.1", "4b1d9e22", func(_ req.A) {
			out, code := swiftJudged(t, map[string]string{swiftTest: swiftTestingFile("test0001_1", "testAddsAgain")})
			acceptance.Sees(t, out, code, 1, "testAddsAgain proves no requirement; rename it testNNNN_N")
			out, code = scriptChecked(t, `test("0001.1 adds two numbers", () => {});`, `it("adds again", () => {});`)
			acceptance.Sees(t, out, code, 1, `tests/adder.test.ts:4: it("adds again") proves no requirement; rename it test("NNNN.N …")`)
			out, code = checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest", "test0001_1", "addsAgain")}))
			acceptance.Sees(t, out, code, 1, "src/test/java/AdderTest.java:6: addsAgain proves no requirement; rename it testNNNN_N")
			unchangedTestFails(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "TestAddsAgain")}), "TestAddsAgain proves no requirement")
			for _, body := range []string{xcTestFile("testAddsAgain"), "import Testing\n@Test func addsAgain() {}\n"} {
				files := acceptance.With(swiftRepo(), map[string]string{"Tests/ExtraTests.swift": body})
				unchangedTestFails(t, files, "proves no requirement; rename it testNNNN_N")
			}
			files := acceptance.With(scriptRepo(), map[string]string{"tests/extra.spec.js": vitestFile(`it("adds again", () => {});`)})
			unchangedTestFails(t, files, `it("adds again") proves no requirement`)
			files = acceptance.With(junitRepo(), map[string]string{"src/test/java/ExtraTest.java": junitFile("ExtraTest", "addsAgain")})
			unchangedTestFails(t, files, "addsAgain proves no requirement")
		})
	})
}

func Test0007_2(t *testing.T) {
	t.Parallel()
	t.Run("0007.2.1 code with no requirement doc fails naming the folder", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.2.1.1.1", "9076e521", func(_ req.A) {
			for _, repo := range []map[string]string{acceptance.GoRepo(), swiftRepo(), scriptRepo(), junitRepo()} {
				out, code := checked(t, acceptance.With(repo, map[string]string{"docs/requirements/0001-adds-numbers.md": ""}))
				acceptance.Sees(t, out, code, 1, "docs/requirements/ has no requirement docs")
			}
			out, code := checked(t, map[string]string{"src/app.js": "export const app = 1;\n"})
			acceptance.Sees(t, out, code, 1, "docs/requirements/ has no requirement docs")
		})
	})
	t.Run("0007.2.2 a repo with no covered code passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.2.1.2.1", "857dd667", func(_ req.A) {
			for _, files := range []map[string]string{nil, {"vite.config.ts": "export default {};\n", ".eslintrc.cjs": "module.exports = {};\n"}} {
				out, code := checked(t, files)
				acceptance.Sees(t, out, code, 0)
				acceptance.PassedWithNoFiles(t, out, "requirements")
			}
		})
	})
}

func Test0007_3(t *testing.T) {
	t.Parallel()
	t.Run("0007.3.1 a Swift test on another framework fails naming the line then the library", func(t *testing.T) {
		t.Parallel()
		boot := "  let shell = Process()\n  shell.arguments = [\"xcrun\", \"simctl\", \"boot\", \"booted\"]\n"
		type failing struct {
			files map[string]string
			lines []string // the line each failure names
			hint  bool     // whether the failure names the library to use
		}
		cases := []failing{
			{map[string]string{swiftTest: xcTestFile("test0001_1")}, []string{swiftTest + ":1: imports XCTest; only an end-to-end test that imports UITestingLibrarySystem may use it, so ", swiftTest + ":3: subclasses XCTestCase; "}, true},
			{map[string]string{swiftTest: "import Testing\n@testable import XCTest\n\n@Test func test0001_1() {}\n"}, []string{swiftTest + ":2: imports XCTest"}, false},
			// An end-to-end test drives the app through SystemSession, never XCUIApplication itself.
			{map[string]string{swiftTest: strings.Replace(systemSuite, "try session.taps", "XCUIApplication().launch()\n    try session.taps", 1)}, []string{"drives Xcode UI tests through XCUIApplication"}, false},
			{map[string]string{swiftTest: librarySuite("", "  let app = XCUIApplication()\n  app.launch()\n")}, []string{swiftTest + ":5: drives Xcode UI tests through XCUIApplication, which launches the app in a separate process; "}, true},
			{map[string]string{swiftTest: librarySuite("import Foundation\n", boot)}, []string{swiftTest + ":7: boots a simulator through simctl outside the end-to-end mode; "}, true},
			{map[string]string{swiftScreens: ""}, []string{swiftTest + ": nothing in Tests/AdderTests/ imports UITestingLibrary; "}, true},
			// A library import in a sibling folder does not reach the test.
			{map[string]string{swiftScreens: "", "Tests/OtherTests/Screens.swift": "import UITestingLibrary\n"}, []string{"nothing in Tests/AdderTests/ imports UITestingLibrary"}, false},
		}
		for _, framework := range []string{"ViewInspector", "SnapshotTesting", "KIF", "EarlGrey", "Hammer"} {
			cases = append(cases, failing{map[string]string{swiftTest: librarySuite("import "+framework+"\n", "")}, []string{swiftTest + ":3: imports " + framework + ", a UI test framework requirement tests may not use; "}, true})
		}
		outs := make([]string, len(cases))
		codes := make([]int, len(cases))
		for i, c := range cases {
			outs[i], codes[i] = swiftChecked(t, c.files)
		}
		req.Then(t, "0007.3.1.1.1", "aabd0f94", func(_ req.A) {
			for i, c := range cases {
				acceptance.Sees(t, outs[i], codes[i], 1, c.lines...)
			}
		})
		req.Then(t, "0007.3.1.1.2", "3cb9ec2f", func(_ req.A) {
			for i, c := range cases {
				if c.hint {
					for _, line := range c.lines {
						acceptance.Sees(t, outs[i], codes[i], 1, line+libraryHint)
					}
				}
			}
		})
	})
	t.Run("0007.3.2 a Swift Testing test on UITestingLibrary passes", func(t *testing.T) {
		t.Parallel()
		// Words in comments and strings are no import or call.
		req.Then(t, "0007.3.1.2.1", "07b8e753", func(_ req.A) {
			quiet := "import Testing\n// import XCTest\n/* XCUIApplication() */\n@MainActor @Test(\"never shells out to simctl\") func test0001_1() async throws {\n  let note = \"XCUIApplication\"\n  _ = note\n" + swiftThen(1, thenText) + "}\n"
			for _, body := range []string{swiftTestingFile("test0001_1"), librarySuite("", ""), quiet} {
				out, code := swiftChecked(t, map[string]string{swiftTest: body})
				acceptance.Sees(t, out, code, 0)
				acceptance.PassedWith(t, out, "requirements", "2")
			}
			out, code := swiftChecked(t, map[string]string{swiftTest: systemSuite})
			acceptance.Sees(t, out, code, 0)
			acceptance.PassedWith(t, out, "requirements", "2")
		})
	})
}

func Test0007_4(t *testing.T) {
	t.Parallel()
	t.Run("0007.4.1 a requirement a passing check proves needs no test", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.4.1.1.1", "850ffc60", func(_ req.A) {
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", provenRepo()), nil)
			acceptance.Sees(t, out, code, 0, "All 2 checks passed")
			acceptance.PassedWith(t, out, "req-0002", "1")
		})
	})
	t.Run("0007.4.2 a check with no doc or a change without its doc fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.4.1.2.1", "e011985a", func(_ req.A) {
			orphan := acceptance.With(provenRepo(), map[string]string{"fitness/req-0003/README.md": houseReadme, "fitness/req-0003/check": acceptance.VerdictScript(passing)})
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", orphan), nil, "requirements")
			acceptance.Sees(t, out, code, 1, "fitness/req-0003 proves requirement 0003, which no doc defines; write the doc or delete the check")
			out, code = stagedRun(t, provenRepo(), map[string]string{"fitness/req-0002/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 2}`)})
			acceptance.Sees(t, out, code, 1, "fitness/req-0002/check changed without its requirement doc docs/requirements/0002-adds-again.md; change both together")
		})
	})
	t.Run("0007.4.3 a check that judged no files fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0007.4.1.3.1", "4f26e4fc", func(_ req.A) {
			empty := acceptance.With(provenRepo(), map[string]string{"fitness/req-0002/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 0}`)})
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", empty), nil, "check", "req-0002")
			acceptance.Sees(t, out, code, 1, "req-0002 judged no files, so its scan is broken")
		})
	})
}
