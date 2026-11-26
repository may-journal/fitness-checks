package acceptance

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// The Given files a user's repo adds to an example repo. They live here, not
// as tracked example repos, so this repo's own checks never judge them.

// GoRepo is a Go module whose requirement 0001 has its one owning test.
func GoRepo() map[string]string {
	return map[string]string{
		"go.mod":                                 "module example.com/adder\n\ngo 1.22\n",
		"adder.go":                               "// Package adder adds numbers.\npackage adder\n\n// Add returns the sum of a and b.\nfunc Add(a, b int) int {\n\treturn a + b\n}\n",
		"adder_test.go":                          TestFile("adder", "Test0001_1"),
		"docs/requirements/0001-adds-numbers.md": RequirementDoc("0001.1"),
	}
}

// requirementTest matches a test name that owns an acceptance, such as
// Test0001_1. Group 1 is its ID.
var requirementTest = regexp.MustCompile(`^Test(\d{4}_\d+)$`)

// TestFile is a Go test file in pkg declaring one passing test per name. A
// test named for an acceptance holds the case of the one THEN that
// Acceptances renders for it.
func TestFile(pkg string, names ...string) string {
	body := "package " + pkg + "\n\nimport \"testing\"\n"
	for _, n := range names {
		m := requirementTest.FindStringSubmatch(n)
		if m == nil {
			body += "\nfunc " + n + "(t *testing.T) {}\n"
			continue
		}
		body += "\n" + CaseFunc(n, ThenCases(strings.ReplaceAll(m[1], "_", "."), 1)...)
	}
	return body
}

// RequirementDoc is requirement 0001 with one Given/When/Then per ID.
func RequirementDoc(ids ...string) string {
	return DocWith(Acceptances(ids...))
}

// Acceptances renders one numbered acceptance per ID at the place it numbers.
// Each level nests three spaces deeper, as Biome formats it.
func Acceptances(ids ...string) string {
	body := ""
	for _, id := range ids {
		body += named(id) + "      1. WHEN I add them\n         1. THEN I get their sum\n"
	}
	return body
}

// named is acceptance id's numbered name, which its place sets, and its one
// GIVEN.
func named(id string) string {
	_, n, _ := strings.Cut(id, ".")
	return n + ". Adds pair " + n + "\n   1. GIVEN two numbers\n"
}

// Chain renders id as one GIVEN and a numbered chain of steps, nested as
// Biome formats it. Step N's WHEN sits on line 2N+17 of DocWith.
func Chain(id string, steps int) string {
	body := named(id)
	for n, step := range chainSteps[:steps] {
		body += "      " + strconv.Itoa(n+1) + ". WHEN " + step[0] + "\n         1. THEN " + step[1] + "\n"
	}
	return body
}

// chainSteps are eight When and Then pairs, worded apart so no two steps
// read as duplicate lines.
var chainSteps = [][2]string{
	{"I add them", "I get their sum"},
	{"I double that sum", "the total is even"},
	{"I subtract the first number", "the second remains twice over"},
	{"I halve the result", "one whole number is left"},
	{"I square it", "a perfect square shows"},
	{"I negate the square", "its sign flips below zero"},
	{"I clear the display", "zero appears"},
	{"I undo once", "the negative square returns"},
}

// DocWith is requirement 0001 with the given Requirements, its Measurement
// ratio set apart by blank lines as Biome formats it.
func DocWith(reqs string) string {
	return "# 0001 Adds Numbers\n\n## Why\n\n" +
		"I can add two numbers without doing it by hand.\n\n## Measurement\n\nsums that are right\n\n-\n\nsums asked for\n\n" +
		"## Requirements\n\n" + reqs
}

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

// LoggedGit returns a PATH setting whose git logs each call's arguments to the
// returned file, one line each, before running the real git.
func LoggedGit(t *testing.T) (path, calls string) {
	t.Helper()
	real, err := exec.LookPath("git")
	MustDo(t, err)
	dir := t.TempDir()
	calls = filepath.Join(t.TempDir(), "git.log")
	Write(t, dir, map[string]string{"git": "#!/bin/sh\necho \"$*\" >> \"" + calls + "\"\nexec \"" + real + "\" \"$@\"\n"})
	return "PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH"), calls
}

// Occurrences counts how often want appears in the file at path.
func Occurrences(path, want string) int {
	raw, _ := os.ReadFile(path)
	return strings.Count(string(raw), want)
}

// CaseFile is a Go test file in pkg declaring test name. It runs one case
// per case name, each ending in one assertion call.
func CaseFile(pkg, name string, cases ...string) string {
	return "package " + pkg + "\n\nimport \"testing\"\n\n" + CaseFunc(name, cases...)
}

// CaseFiles is a Go test file in pkg declaring each of funcs, as CaseFunc
// writes them.
func CaseFiles(pkg string, funcs ...string) string {
	return "package " + pkg + "\n\nimport \"testing\"\n\n" + strings.Join(funcs, "\n")
}

// CaseFunc is the Go test function name, which runs one case per case name,
// each ending in one assertion call.
func CaseFunc(name string, cases ...string) string {
	body := "func " + name + "(t *testing.T) {\n"
	for _, c := range cases {
		body += "\tt.Run(" + strconv.Quote(c) + ", func(t *testing.T) {\n\t\tif t.Name() == \"\" {\n\t\t\tt.Fatal(\"unnamed\")\n\t\t}\n\t})\n"
	}
	return body + "}\n"
}

// ThenCases names one case for each of the first whens WHENs of acceptance
// id, as Acceptances and Chain render it. A name is the ID plus the words.
func ThenCases(id string, whens int) []string {
	var out []string
	for n := 1; n <= whens; n++ {
		out = append(out, id+".1."+strconv.Itoa(n)+".1 "+chainSteps[n-1][1])
	}
	return out
}

// PairedDoc is the path of requirement 0001 in GoRepo.
const PairedDoc = "docs/requirements/0001-adds-numbers.md"

// PairedRepo is GoRepo with requirement 0001's test in reqs/0001_test.go,
// named for it, beside a shared helper that proves no requirement. Only the
// requirements check runs, so commits meet no other check.
func PairedRepo() map[string]string {
	return With(GoRepo(), map[string]string{
		".fitnessrc.json":    `{"policy": "external", "checks": ["requirements"]}` + "\n",
		"adder_test.go":      "",
		"reqs/0001_test.go":  TestFile("reqs", "Test0001_1"),
		"reqs/tools_test.go": "package reqs\n\n// two is a helper the tests share.\nfunc two() int { return 2 }\n",
	})
}

// DocOf is GoRepo with requirement 0001 replaced by body at path.
func DocOf(path, body string) map[string]string {
	return With(GoRepo(), map[string]string{PairedDoc: "", path: body})
}

// Reqs is GoRepo with requirement 0001's Requirements replaced.
func Reqs(body string) map[string]string {
	return DocOf(PairedDoc, DocWith(body))
}

// WithSource is doc with the Source line older docs end their Measurement
// with, citing cite. It sits on line 19.
func WithSource(doc, cite string) string {
	return strings.Replace(doc, "## Requirements", "Source: "+cite+"\n\n## Requirements", 1)
}

// Monorepo is GoRepo with packages pkgs/a and pkgs/b, each with its own
// requirement 0001 and test, named in .fitnessrc.json when packages is set.
func Monorepo(packages bool) map[string]string {
	rc := `{"policy": "external", "checks": ["requirements"]}`
	if packages {
		rc = `{"policy": "external", "checks": ["requirements"], "packages": ["pkgs/*"]}`
	}
	repo := With(GoRepo(), map[string]string{".fitnessrc.json": rc + "\n"})
	for _, name := range []string{"a", "b"} {
		repo["pkgs/"+name+"/"+name+".go"] = "// Package " + name + " adds numbers.\npackage " + name + "\n"
		repo[PackageDoc(name)] = RequirementDoc("0001.1")
		repo["pkgs/"+name+"/requirements/0001_test.go"] = TestFile("requirements", "Test0001_1")
	}
	return repo
}

// PackageDoc is the path of requirement 0001 in package pkgs/<name>.
func PackageDoc(name string) string {
	return "pkgs/" + name + "/docs/requirements/0001-adds-numbers.md"
}
