package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0007_1(t *testing.T) {
	t.Parallel()
	t.Run("0007.1.1.1.1 it owns its acceptance", func(t *testing.T) {
		t.Parallel()
		out, _ := swiftJudged(t, nil)
		acceptance.PassedWith(t, out, "requirements", "2")
		out, _ = swiftJudged(t, map[string]string{
			"Tests/AdderTests/AdderTests.swift": "import Testing\n\nstruct AdderSpec {\n  @Test(\"adds\") func test0001_1() {\n    XCTContext.runActivity(named: \"0001.1.1.1.1 I get their sum\") { _ in #expect(1 + 1 == 2) }\n  }\n}\n",
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
		junitOwns(t)
	})
	t.Run("0007.1.1.1.2 a name in a comment or a string or a helper owns nothing", func(t *testing.T) {
		t.Parallel()
		swiftOwnsNothing(t, xcTestFile()+"// func test0001_1() {}\n/* nested /* func test0001_1() {} */ */\nlet note = \"\\(\"func test0001_1() {}\")\"\nlet raw = #\"\"\"\n\" func test0001_1() {}\n\"\"\"#\n")
		swiftOwnsNothing(t, xcTestFile()+"\nstruct Helper {\n  func test0001_1() {}\n}\n")
		out, code := scriptChecked(t, notTests...)
		acceptance.Sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
		junitOwnsNothing(t)
	})
	t.Run("0007.1.1.2.1 it fails naming the test form to write", func(t *testing.T) {
		t.Parallel()
		swiftOwnsNothing(t, xcTestFile())
		out, code := checked(t, acceptance.With(scriptRepo(), map[string]string{"tests/adder.test.ts": ""}))
		acceptance.Sees(t, out, code, 1, `0001.1 has no test; name exactly one test test("0001.1 …")`)
		out, code = checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest")}))
		acceptance.Sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
	})
	t.Run("0007.1.1.3.1 it needs only its own test until the Node switch is on", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, acceptance.With(scriptRepo(), map[string]string{acceptance.PairedDoc: acceptance.DocWith(acceptance.Chain("0001.1", 2))}))
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0007.1.1.4.1 it fails naming the form to rename it to", func(t *testing.T) {
		t.Parallel()
		out, code := swiftJudged(t, map[string]string{swiftTest: swiftTestingFile("test0001_1", "testAddsAgain")})
		acceptance.Sees(t, out, code, 1, "testAddsAgain proves no requirement; rename it testNNNN_N")
		out, code = scriptChecked(t, `test("0001.1 adds two numbers", () => {});`, `it("adds again", () => {});`)
		acceptance.Sees(t, out, code, 1, `tests/adder.test.ts:4: it("adds again") proves no requirement; rename it test("NNNN.N …")`)
		out, code = checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": junitFile("AdderTest", "test0001_1", "addsAgain")}))
		acceptance.Sees(t, out, code, 1, "src/test/java/AdderTest.java:10: addsAgain proves no requirement; rename it testNNNN_N")
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
	t.Run("0007.1.1.5.1 it fails naming the THEN", func(t *testing.T) {
		t.Parallel()
		out, code := swiftJudged(t, map[string]string{swiftTest: "import Testing\n\n@MainActor @Suite struct AdderSpec {\n  @Test func test0001_1() async throws {}\n}\n"})
		acceptance.Sees(t, out, code, 1, "THEN 0001.1.1.1.1 has no test case; add a case to test0001_1 named \"0001.1.1.1.1 I get their sum\"")
		out, code = checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": "import org.junit.jupiter.api.Test;\n\nclass AdderTest {\n  @Test\n  void test0001_1() {}\n}\n"}))
		acceptance.Sees(t, out, code, 1, "THEN 0001.1.1.1.1 has no test case; add a case to test0001_1 named \"0001.1.1.1.1 I get their sum\"")
	})
}

func Test0007_2(t *testing.T) {
	t.Parallel()
	t.Run("0007.2.1.1.1 it fails naming the folder", func(t *testing.T) {
		t.Parallel()
		for _, repo := range []map[string]string{acceptance.GoRepo(), swiftRepo(), scriptRepo(), junitRepo()} {
			out, code := checked(t, acceptance.With(repo, map[string]string{"docs/requirements/0001-adds-numbers.md": ""}))
			acceptance.Sees(t, out, code, 1, "docs/requirements/ has no requirement docs")
		}
		out, code := checked(t, map[string]string{"src/app.js": "export const app = 1;\n"})
		acceptance.Sees(t, out, code, 1, "docs/requirements/ has no requirement docs")
	})
	t.Run("0007.2.1.2.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		for _, files := range []map[string]string{nil, {"vite.config.ts": "export default {};\n", ".eslintrc.cjs": "module.exports = {};\n"}} {
			out, code := checked(t, files)
			acceptance.Sees(t, out, code, 0)
			acceptance.PassedWithNoFiles(t, out, "requirements")
		}
	})
}

func Test0007_3(t *testing.T) {
	t.Parallel()
	findings := []struct {
		files map[string]string
		line  string
	}{
		{map[string]string{swiftTest: xcTestFile("test0001_1")}, swiftTest + ":1: imports XCTest; only an end-to-end test that imports UITestingLibrarySystem may use it, so "},
		{map[string]string{swiftTest: xcTestFile("test0001_1")}, swiftTest + ":3: subclasses XCTestCase; "},
		{map[string]string{swiftTest: "import Testing\n@testable import XCTest\n\n@Test func test0001_1() {}\n"}, swiftTest + ":2: imports XCTest"},
		// An end-to-end test drives the app through SystemSession, never XCUIApplication itself.
		{map[string]string{swiftTest: strings.Replace(systemSuite, "try session.taps", "XCUIApplication().launch()\n    try session.taps", 1)}, "drives Xcode UI tests through XCUIApplication"},
		{map[string]string{swiftTest: librarySuite("", "  let app = XCUIApplication()\n  app.launch()\n")}, swiftTest + ":5: drives Xcode UI tests through XCUIApplication, which launches the app in a separate process; "},
		{map[string]string{swiftTest: librarySuite("import Foundation\n", "  let shell = Process()\n  shell.arguments = [\"xcrun\", \"simctl\", \"boot\", \"booted\"]\n")}, swiftTest + ":7: boots a simulator through simctl outside the end-to-end mode; "},
		{map[string]string{swiftScreens: ""}, swiftTest + ": nothing in Tests/AdderTests/ imports UITestingLibrary; "},
		// A library import in a sibling folder does not reach the test.
		{map[string]string{swiftScreens: "", "Tests/OtherTests/Screens.swift": "import UITestingLibrary\n"}, "nothing in Tests/AdderTests/ imports UITestingLibrary"},
	}
	for _, framework := range []string{"ViewInspector", "SnapshotTesting", "KIF", "EarlGrey", "Hammer"} {
		findings = append(findings, struct {
			files map[string]string
			line  string
		}{map[string]string{swiftTest: librarySuite("import "+framework+"\n", "")}, swiftTest + ":3: imports " + framework + ", a UI test framework requirement tests may not use; "})
	}
	t.Run("0007.3.1.1.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		for _, f := range findings {
			out, code := swiftChecked(t, f.files)
			acceptance.Sees(t, out, code, 1, f.line)
		}
	})
	t.Run("0007.3.1.1.2 it fails naming the library to use", func(t *testing.T) {
		t.Parallel()
		for _, f := range findings {
			out, code := swiftChecked(t, f.files)
			acceptance.Sees(t, out, code, 1, libraryHint)
		}
	})
	t.Run("0007.3.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		// Words in comments and strings are no import or call.
		quiet := "import Testing\n// import XCTest\n/* XCUIApplication() */\n@MainActor @Test(\"never shells out to simctl\") func test0001_1() async throws {\n  let note = \"XCUIApplication\"\n  _ = note\n  XCTContext.runActivity(named: \"0001.1.1.1.1 I get their sum\") { _ in #expect(note.isEmpty == false) }\n}\n"
		for _, body := range []string{swiftTestingFile("test0001_1"), librarySuite("", "  XCTContext.runActivity(named: \"0001.1.1.1.1 I get their sum\") { _ in #expect(1 + 1 == 2) }\n"), quiet} {
			out, code := swiftChecked(t, map[string]string{swiftTest: body})
			acceptance.Sees(t, out, code, 0)
			acceptance.PassedWith(t, out, "requirements", "2")
		}
		out, code := swiftChecked(t, map[string]string{swiftTest: systemSuite})
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "requirements", "2")
	})
}

func Test0007_4(t *testing.T) {
	t.Parallel()
	t.Run("0007.4.1.1.1 the acceptances need no test", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", provenRepo()), nil)
		acceptance.Sees(t, out, code, 0, "All 2 checks passed")
		acceptance.PassedWith(t, out, "req-0002", "1")
	})
	t.Run("0007.4.1.2.1 it fails naming the check", func(t *testing.T) {
		t.Parallel()
		orphan := acceptance.With(provenRepo(), map[string]string{"fitness/req-0003/README.md": houseReadme, "fitness/req-0003/check": acceptance.VerdictScript(passing)})
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", orphan), nil, "requirements")
		acceptance.Sees(t, out, code, 1, "fitness/req-0003 proves requirement 0003, which no doc defines; write the doc or delete the check")
		out, code = stagedRun(t, provenRepo(), map[string]string{"fitness/req-0002/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 2}`)})
		acceptance.Sees(t, out, code, 1, "fitness/req-0002/check changed without its requirement doc docs/requirements/0002-adds-again.md; change both together")
	})
	t.Run("0007.4.1.3.1 it fails since its scan is broken", func(t *testing.T) {
		t.Parallel()
		empty := acceptance.With(provenRepo(), map[string]string{"fitness/req-0002/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 0}`)})
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", empty), nil, "check", "req-0002")
		acceptance.Sees(t, out, code, 1, "req-0002 judged no files, so its scan is broken")
	})
}
