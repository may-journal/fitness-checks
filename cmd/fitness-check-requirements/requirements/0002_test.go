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

// historyRun runs requirements on a monorepo whose package pkgs/a dropped
// 0001.2 and whose other packages never held it, with git logging its calls.
func historyRun(t *testing.T) (out, calls string) {
	t.Helper()
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
	out, _ = acceptance.Fitness(t, repo, []string{path}, "requirements")
	return out, calls
}

func Test0002_1(t *testing.T) {
	t.Parallel()
	t.Run("0002.1.2.1.1 one git log serves them all", func(t *testing.T) {
		t.Parallel()
		_, calls := historyRun(t)
		if n := acceptance.Occurrences(calls, "-U100000"); n != 1 {
			t.Errorf("requirements read the ID history of 3 package roots with %d git log calls, want 1", n)
		}
	})
	t.Run("0002.1.2.1.2 each package keeps its own deleted IDs", func(t *testing.T) {
		t.Parallel()
		out, _ := historyRun(t)
		acceptance.Sees(t, out, 1, 1, packageDoc(`a "b`)+":21: 0001.2 was deleted before")
		for _, doc := range []string{packageDoc("b"), rootDoc} {
			if strings.Contains(out, doc+":21: 0001.2 was deleted before") {
				t.Errorf("%s never held 0001.2, so pkgs/a's deletion must not count against it:\n%s", doc, out)
			}
		}
		if n := strings.Count(acceptance.Read(out), "0001.2 was deleted before"); n != 1 {
			t.Errorf("only pkgs/a held 0001.2, so one doc must fail, got %d:\n%s", n, out)
		}
	})
}

func Test0002_2(t *testing.T) {
	t.Parallel()
	judged := func(t *testing.T, files map[string]string) (string, int) {
		t.Helper()
		return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(acceptance.Monorepo(true), files)), nil, "requirements")
	}
	t.Run("0002.2.1.1.1 each package with code owns its docs", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, nil)
		acceptance.Sees(t, out, code, 0)
		out, code = judged(t, map[string]string{packageDoc("b"): "", "pkgs/b/requirements/0001_test.go": ""})
		acceptance.Sees(t, out, code, 1, "pkgs/b/docs/requirements/ has no requirement docs")
	})
	t.Run("0002.2.1.1.2 each package with code owns its tests", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, map[string]string{packageDoc("a"): acceptance.RequirementDoc("0001.1", "0001.2")})
		acceptance.Sees(t, out, code, 1, packageDoc("a")+":21: 0001.2 has no test")
	})
	t.Run("0002.2.1.1.3 a package without code needs none", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, map[string]string{"pkgs/c/README.md": "# Shared assets\n"})
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0002.2.1.2.1 the number repeats only across packages", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, map[string]string{"pkgs/a/docs/requirements/0001-other.md": acceptance.RequirementDoc("0001.1")})
		acceptance.Sees(t, out, code, 1, "requirement 0001 is used by "+packageDoc("a")+", pkgs/a/docs/requirements/0001-other.md")
		if strings.Contains(out, "docs/requirements/0001-adds-numbers.md, pkgs/") {
			t.Errorf("the repo root's 0001 must not clash with a package's:\n%s", out)
		}
	})
	t.Run("0002.2.1.3.1 the repo is one package", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.Monorepo(false)), nil, "requirements")
		acceptance.Sees(t, out, code, 1, "0001.1 is owned by 3 tests")
	})
}

func Test0002_3(t *testing.T) {
	t.Parallel()
	named := "1. Adds numbers\n   1. GIVEN Developer has two numbers\n      1. WHEN the `semantic-commit` check reads them\n         1. THEN it gets their sum\n"
	given := func(rootDoc, packageDoc string) map[string]string {
		return acceptance.With(acceptance.With(acceptance.Monorepo(true), owning("fitness-check-semantic-commit", "0001.1")), map[string]string{
			acceptance.PairedDoc: acceptance.DocWith(rootDoc),
			"pkgs/fitness-check-semantic-commit/docs/requirements/0001-adds-numbers.md": acceptance.DocWith(packageDoc),
			"pkgs/fitness-check-semantic-commit/requirements/0001_test.go":              acceptance.CaseFile("requirements", "Test0001_1", "0001.1.1.1.1 it gets their sum"),
		})
	}
	t.Run("0002.3.1.1.1 it fails naming the acceptance with the check", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, given(named, named))
		acceptance.Sees(t, out, code, 1, "acceptance 0001.1 names the semantic-commit check; move it to that check's own requirement doc")
	})
	t.Run("0002.3.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		plain := strings.Replace(named, "the `semantic-commit` check", "he", 1)
		removed := plain + "2. Removed in 2026.11.26.0400: moved to the semantic-commit check's own requirement\n"
		out, code := checked(t, given(removed, named))
		acceptance.Sees(t, out, code, 0)
	})
}
