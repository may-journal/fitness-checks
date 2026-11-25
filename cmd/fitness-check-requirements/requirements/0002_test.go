package requirements

// cspell:ignore noprefix

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// packageDoc is the path of requirement 0001 in package pkgs/<name>.
func packageDoc(name string) string {
	return "pkgs/" + name + "/" + rootDoc
}

// owning is package pkgs/<name>'s code with its requirement 0001 holding ids,
// each proven by its own test. The folder name may hold any character.
func owning(name string, ids ...string) map[string]string {
	var tests []string
	for _, id := range ids {
		tests = append(tests, "Test"+strings.ReplaceAll(id, ".", "_"))
	}
	return map[string]string{
		"pkgs/" + name + "/adder.go":                  "// Package adder adds numbers.\npackage adder\n",
		packageDoc(name):                              acceptance.RequirementDoc(ids...),
		"pkgs/" + name + "/requirements/0001_test.go": acceptance.TestFile("requirements", tests...),
	}
}

func Test0002_1(t *testing.T) {
	t.Parallel()
	rc := map[string]string{".fitnessrc.json": `{"policy": "external", "checks": ["requirements"], "packages": ["pkgs/*"]}` + "\n"}
	// Git quotes a package path that holds a quote, and its header holds spaces.
	const a = `a "b`
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), acceptance.With(rc, owning(a, "0001.1", "0001.2"))))
	acceptance.Commit(t, repo, owning(a, "0001.1"), "docs(requirements): drop 0001.2 from pkgs/a")
	// A user's prefix setting must not hide which package a diff belongs to.
	acceptance.Git(t, repo, "config", "diff.noprefix", "true")
	// Neither pkgs/b nor the root held 0001.2, so each may take that number.
	acceptance.Write(t, repo, acceptance.With(owning(a, "0001.1", "0001.2"), owning("b", "0001.1", "0001.2")))
	acceptance.Write(t, repo, map[string]string{rootDoc: acceptance.RequirementDoc("0001.1", "0001.2"), "adder_test.go": acceptance.TestFile("adder", "Test0001_1", "Test0001_2")})
	path, calls := acceptance.LoggedGit(t)
	out, code := acceptance.Fitness(t, repo, []string{path}, "requirements")
	acceptance.Sees(t, out, code, 1, packageDoc(a)+":21: 0001.2 was deleted before")
	for _, doc := range []string{packageDoc("b"), rootDoc} {
		if strings.Contains(out, doc+":21: 0001.2 was deleted before") {
			t.Errorf("%s never held 0001.2, so pkgs/a's deletion must not count against it:\n%s", doc, out)
		}
	}
	if n := acceptance.Occurrences(calls, "-U100000"); n != 1 {
		t.Errorf("requirements read the ID history of 3 package roots with %d git log calls, want 1", n)
	}
}
