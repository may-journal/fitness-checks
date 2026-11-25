package acceptance

import (
	"os"
	"os/exec"
	"path/filepath"
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

// TestFile is a Go test file in pkg declaring one passing test per name.
func TestFile(pkg string, names ...string) string {
	body := "package " + pkg + "\n\nimport \"testing\"\n"
	for _, n := range names {
		body += "\nfunc " + n + "(t *testing.T) {}\n"
	}
	return body
}

// RequirementDoc is requirement 0001 with one Given/When/Then per ID.
func RequirementDoc(ids ...string) string {
	return DocWith(Acceptances(ids...))
}

// Acceptances renders one well-formed chain per ID, nested two spaces deeper
// at each level as Biome formats it.
func Acceptances(ids ...string) string {
	body := ""
	for _, id := range ids {
		body += "- " + id + "\n  - Given two numbers\n    - When I add them\n      - Then I get their sum\n"
	}
	return body
}

// Chain renders id as one Given and a numbered chain of steps, nested as
// Biome formats it. Step N's When sits on line 2N+17 of DocWith.
func Chain(id string, steps int) string {
	body := "- " + id + "\n  - Given two numbers\n"
	for n, step := range chainSteps[:steps] {
		body += "    " + strconv.Itoa(n+1) + ". When " + step[0] + "\n       - Then " + step[1] + "\n"
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
