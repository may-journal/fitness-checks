package requirements

import (
	"strconv"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// The Given files a user's repo adds to happyRepo. They live here, not as
// tracked example repos, so this repo's own checks never judge them.

// swiftRepo is a Swift package whose Swift Testing test test0001_1 owns
// 0001.1, beside a shared helper that imports apple-ui-testing-library.
func swiftRepo() map[string]string {
	return map[string]string{
		swiftScreens:                             "import UITestingLibrary\n\n/// Screens every test renders.\nenum Screens {}\n",
		"Package.swift":                          "// swift-tools-version:5.9\nimport PackageDescription\n\nlet package = Package(name: \"Adder\")\n",
		"Sources/Adder/Adder.swift":              "/// Adds numbers.\npublic enum Adder {\n  /// Returns the sum of a and b.\n  public static func add(_ a: Int, _ b: Int) -> Int { a + b }\n}\n",
		swiftTest:                                swiftTestingFile("test0001_1"),
		"docs/requirements/0001-adds-numbers.md": acceptance.DocWith(acceptance.Acceptances("0001.1")),
	}
}

// swiftTestingFile is a Swift Testing suite declaring one empty test per
// name; the test of 0001.1 tags its THEN.
func swiftTestingFile(names ...string) string {
	body := "import Testing\n\n@MainActor @Suite struct AdderSpec {\n"
	for _, n := range names {
		body += "  @Test func " + n + "() async throws {" + swiftBody(n) + "}\n"
	}
	return body + "}\n"
}

// xcTestFile is an XCTest case class declaring one empty test per name.
func xcTestFile(names ...string) string {
	body := "import XCTest\n\nfinal class AdderTests: XCTestCase {\n"
	for _, n := range names {
		body += "  func " + n + "() {" + swiftBody(n) + "}\n"
	}
	return body + "}\n"
}

// junitRepo is a Maven project whose JUnit test test0001_1 owns 0001.1.
func junitRepo() map[string]string {
	return map[string]string{
		"pom.xml":                                "<project/>\n",
		"src/main/java/Adder.java":               "/** Adds numbers. */\nclass Adder {\n  static int add(int a, int b) {\n    return a + b;\n  }\n}\n",
		"src/test/java/AdderTest.java":           junitFile("AdderTest", "test0001_1"),
		"docs/requirements/0001-adds-numbers.md": acceptance.DocWith(acceptance.Acceptances("0001.1")),
	}
}

// junitFile is a JUnit test class declaring one empty test per name.
func junitFile(class string, names ...string) string {
	body := "import org.junit.jupiter.api.Test;\n\nclass " + class + " {\n"
	for _, n := range names {
		body += "  @Test" + javaTagFor(n) + "\n  void " + n + "() {}\n"
	}
	return body + "}\n"
}

// scriptRepo is a TypeScript project whose vitest test titled "0001.1 …"
// owns 0001.1.
func scriptRepo() map[string]string {
	return map[string]string{
		"src/adder.ts":                           "/** Adds numbers. */\nexport const add = (a: number, b: number): number => a + b;\n",
		"tests/adder.test.ts":                    vitestFile(`test("0001.1 adds two numbers", () => {});`),
		"docs/requirements/0001-adds-numbers.md": acceptance.DocWith(acceptance.Acceptances("0001.1")),
	}
}

// vitestFile is a vitest test file holding the given statements.
func vitestFile(statements ...string) string {
	return "import { describe, expect, it, test } from \"vitest\";\n\n" + strings.Join(statements, "\n") + "\n" + nodeThen(1, thenText) + "\n"
}

// swiftScreens is the helper that imports apple-ui-testing-library.
const swiftScreens = "Tests/AdderTests/Screens.swift"

// libraryHint is the part of every test library failure that says what to
// use instead.
const libraryHint = "write it with Swift Testing and UITestingLibrary, hosted in the app, as https://github.com/may-journal/apple-ui-testing-library shows"

// swiftChecked runs only the requirements check on swiftRepo with files
// written over it.
func swiftChecked(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return checked(t, acceptance.With(swiftRepo(), files))
}

// librarySuite is a Swift Testing file whose test0001_1 runs body on
// apple-ui-testing-library, after the given extra imports.
func librarySuite(imports, body string) string {
	return "import Testing\nimport UITestingLibrary\n" + imports + "\n@MainActor @Test func test0001_1() async throws {\n" + body + swiftThen(1, thenText) + "}\n"
}

// systemSuite is an end-to-end XCTest class whose test0001_1 drives the
// app through UITestingLibrarySystem.
var systemSuite = "import UITestingLibrary\nimport UITestingLibrarySystem\nimport XCTest\n\n@MainActor\nfinal class SystemRequirements: XCTestCase {\n" +
	"  func test0001_1() throws {\n    let session = SystemSession(bundleIdentifier: \"org.example.app\", timeout: 10)\n" +
	"    session.launch()\n    let reset = [\"xcrun\", \"simctl\", \"erase\", \"all\"]\n    XCTAssertFalse(reset.isEmpty)\n    try session.taps(.button(\"Share\"))\n" + swiftThen(1, thenText) + "  }\n}\n"

// swiftJudged runs fitness on swiftRepo with files written over it.
func swiftJudged(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return suite(t, acceptance.With(swiftRepo(), files))
}

// swiftOwnsNothing is a test file whose test0001_1 is not a test.
func swiftOwnsNothing(t *testing.T, body string) {
	t.Helper()
	out, code := swiftJudged(t, map[string]string{"Tests/AdderTests/AdderTests.swift": body})
	acceptance.Sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
}

// scriptChecked runs the requirements check on scriptRepo with its test file
// replaced by statements.
func scriptChecked(t *testing.T, statements ...string) (string, int) {
	t.Helper()
	return checked(t, acceptance.With(scriptRepo(), map[string]string{"tests/adder.test.ts": vitestFile(statements...)}))
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

// junitOwners checks a JUnit test owns 0001.1 by its name and test
// annotation.
func junitOwners(t *testing.T) {
	t.Helper()
	for _, owner := range []string{
		"class AdderTest {\n  " + javaThen(1, thenText) + "\n  @ParameterizedTest\n  @ValueSource(ints = {1, 2})\n  void test0001_1(int n) {}\n}\n",
		"class AdderTest {\n  " + javaThen(1, thenText) + "\n  @TestFactory\n  Stream<DynamicTest> test0001_1() {\n    return Stream.empty();\n  }\n}\n",
		// TestNG runs a marked class's public methods, but not its setup.
		"@Test " + javaThen(1, thenText) + "\npublic class AdderTest {\n  public AdderTest() {}\n\n  @BeforeTest\n  public void setUp() {}\n\n  public void test0001_1() {}\n}\n",
	} {
		out, code := checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": owner}))
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "requirements", "2")
	}
}

// junitOwnsNothing checks a comment, string or helper method owns no
// acceptance.
func junitOwnsNothing(t *testing.T) {
	t.Helper()
	notJUnit := "class AdderTest {\n  // @Test void test0001_1() {}\n  String note = \"@Test void test0001_1() {}\";\n\n  void test0001_1() {}\n\n  @BeforeTest\n  void test0001_1() {}\n}\n"
	out, code := checked(t, acceptance.With(junitRepo(), map[string]string{"src/test/java/AdderTest.java": notJUnit}))
	acceptance.Sees(t, out, code, 1, "0001.1 has no test; name exactly one test test0001_1")
}

// passing is a check's verdict that judged one file.
const passing = `{"ok": true, "errors": [], "filesChecked": 1}`

// houseReadme is a folder check's README: a title, the rule, and the fix.
const houseReadme = "# House rule\n\nKeeps the house in order. Each room stays tidy.\n\n## Fix\n\nTidy the room the error names.\n"

// provenRepo is twoRepo whose requirement 0002 a check in fitness/req-0002
// proves, with no test of its own.
func provenRepo() map[string]string {
	return acceptance.With(twoRepo(), map[string]string{
		".fitnessrc.json":            `{"policy": "external", "checks": ["requirements", "req-0002"]}` + "\n",
		"reqs/0001_0002_test.go":     "",
		"reqs/0001_test.go":          acceptance.TestFile("reqs", "Test0001_1"),
		"fitness/req-0002/README.md": "# Doubles\n\nProves [requirement 0002](../../docs/requirements/0002-adds-again.md).\n",
		"fitness/req-0002/check":     acceptance.VerdictScript(passing),
	})
}

// unchangedTestFails checks ownership when staged and PR diffs leave tests alone.
func unchangedTestFails(t *testing.T, files map[string]string, want string) {
	t.Helper()
	for _, path := range []string{"README.md", "docs/requirements/0001-adds-numbers.md"} {
		for _, pr := range []bool{false, true} {
			repo := acceptance.Example(t, "happyRepo", files)
			acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
			acceptance.Write(t, repo, map[string]string{path: files[path] + "\n"})
			acceptance.Git(t, repo, "add", path)
			var env []string
			if pr {
				acceptance.Commit(t, repo, nil, "docs: change prose")
				env = []string{"GITHUB_ACTIONS=true", "GITHUB_BASE_REF=main"}
			}
			out, code := acceptance.Fitness(t, repo, env, "requirements")
			acceptance.Sees(t, out, code, 1, want)
		}
	}
}

// stepThenID is the ID of the n-th THEN of 0001.1, which has one GIVEN and one WHEN per step.
func stepThenID(n int) string { return "0001.1.1." + strconv.Itoa(n) + ".1" }

// swiftThen is the Swift activity tagging the n-th THEN of 0001.1, whose text is then.
func swiftThen(n int, then string) string {
	return "    XCTContext.runActivity(named: \"req " + stepThenID(n) + "#" + req.Hash(then) + "\") { _ in }\n"
}

// javaThen is the JUnit tag of the n-th THEN of 0001.1, whose text is then.
func javaThen(n int, then string) string {
	return "@Tag(\"req:" + stepThenID(n) + "#" + req.Hash(then) + "\")"
}

// nodeThen is the Node wrapper call tagging the n-th THEN of 0001.1, whose text is then.
func nodeThen(n int, then string) string {
	return "then(\"" + stepThenID(n) + "\", \"" + req.Hash(then) + "\", () => {});"
}

// swiftBody is the tag a Swift test body holds when the test proves 0001.1.
func swiftBody(name string) string {
	if name != "test0001_1" {
		return ""
	}
	return "\n" + swiftThen(1, thenText) + "  "
}

// javaTagFor is the tag a JUnit test holds on its @Test line when the test proves 0001.1.
func javaTagFor(name string) string {
	if name != "test0001_1" {
		return ""
	}
	return " " + javaThen(1, thenText)
}
