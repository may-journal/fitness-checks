package requirements

// cspell:ignore noprefix

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
	const pkgA = `a "b`
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), acceptance.With(rc, owning(pkgA, "0001.1", "0001.2"))))
	acceptance.Commit(t, repo, owning(pkgA, "0001.1"), "docs(requirements): drop 0001.2 from pkgs/a")
	// A user's prefix setting must not hide which package a diff belongs to.
	acceptance.Git(t, repo, "config", "diff.noprefix", "true")
	// Neither pkgs/b nor the root held 0001.2, so each may take that number.
	acceptance.Write(t, repo, acceptance.With(owning(pkgA, "0001.1", "0001.2"), owning("b", "0001.1", "0001.2")))
	acceptance.Write(t, repo, map[string]string{rootDoc: acceptance.RequirementDoc("0001.1", "0001.2"), "adder_test.go": acceptance.TestFile("adder", "Test0001_1", "Test0001_2")})
	path, calls := acceptance.LoggedGit(t)
	out, code := acceptance.Fitness(t, repo, []string{path}, "requirements")
	req.Then(t, "0002.1.2.1.1", "d36524df", func(a req.A) {
		a.Equal(1, acceptance.Occurrences(calls, "-U100000"))
	})
	req.Then(t, "0002.1.2.1.2", "146b5c0c", func(a req.A) {
		acceptance.Sees(t, out, code, 1, packageDoc(pkgA)+":21: 0001.2 was deleted before")
		for _, doc := range []string{packageDoc("b"), rootDoc} {
			a.True(!strings.Contains(out, doc+":21: 0001.2 was deleted before"), "%s never held 0001.2, so pkgs/a's deletion must not count against it:\n%s", doc, out)
		}
		a.Equal(1, strings.Count(acceptance.Read(out), "0001.2 was deleted before"))
	})
}

func Test0002_2(t *testing.T) {
	t.Parallel()
	t.Run("0002.2.1 each package with code owns its docs plus tests", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.Monorepo(true)), nil, "requirements")
		req.Then(t, "0002.2.1.1.1", "2a1efc8a", func(_ req.A) {
			acceptance.Sees(t, out, code, 0)
			out, code = acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(acceptance.Monorepo(true), map[string]string{
				packageDoc("b"): "", "pkgs/b/requirements/0001_test.go": "",
			})), nil, "requirements")
			acceptance.Sees(t, out, code, 1, "pkgs/b/docs/requirements/ has no requirement docs")
		})
		req.Then(t, "0002.2.1.1.2", "ad307dee", func(_ req.A) {
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(acceptance.Monorepo(true), map[string]string{
				packageDoc("a"): acceptance.RequirementDoc("0001.1", "0001.2"),
			})), nil, "requirements")
			acceptance.Sees(t, out, code, 1, packageDoc("a")+":21: 0001.2 has no test")
		})
		req.Then(t, "0002.2.1.1.3", "9bfa20ab", func(_ req.A) {
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(acceptance.Monorepo(true), map[string]string{"pkgs/c/README.md": "# Shared assets\n"})), nil, "requirements")
			acceptance.Sees(t, out, code, 0)
		})
	})
	t.Run("0002.2.2 a number repeats only across packages", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0002.2.1.2.1", "a8e18db1", func(_ req.A) {
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(acceptance.Monorepo(true), map[string]string{
				"pkgs/a/docs/requirements/0001-other.md": acceptance.RequirementDoc("0001.1"),
			})), nil, "requirements")
			acceptance.Sees(t, out, code, 1, "requirement 0001 is used by "+packageDoc("a")+", pkgs/a/docs/requirements/0001-other.md")
			if strings.Contains(out, "docs/requirements/0001-adds-numbers.md, pkgs/") {
				t.Errorf("the repo root's 0001 must not clash with a package's:\n%s", out)
			}
		})
	})
	t.Run("0002.2.3 a repo naming no packages is one package", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0002.2.1.3.1", "5ada5e4e", func(_ req.A) {
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.Monorepo(false)), nil, "requirements")
			acceptance.Sees(t, out, code, 1, "0001.1 is owned by 3 tests")
		})
	})
}

func Test0002_3(t *testing.T) {
	t.Parallel()
	named := "1. Adds numbers\n   1. GIVEN Developer has two numbers\n      1. WHEN the `semantic-commit` check reads them\n         1. THEN it gets their sum\n"
	given := func(rootDoc, packageDoc string) map[string]string {
		return acceptance.With(acceptance.With(acceptance.Monorepo(true), owning("fitness-check-semantic-commit", "0001.1")), map[string]string{
			acceptance.PairedDoc: acceptance.DocWith(rootDoc),
			"pkgs/fitness-check-semantic-commit/docs/requirements/0001-adds-numbers.md": acceptance.DocWith(packageDoc),
		})
	}
	t.Run("0002.3.1 a root acceptance naming a single check fails naming both", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0002.3.1.1.1", "9261fd7f", func(_ req.A) {
			out, code := checked(t, given(named, acceptance.Acceptances("0001.1")))
			acceptance.Sees(t, out, code, 1, "acceptance 0001.1 names the semantic-commit check; move it to that check's own requirement doc")
		})
	})
	t.Run("0002.3.2 a package doc or a Removed line naming a check passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0002.3.1.2.1", "07b8e753", func(_ req.A) {
			plain := strings.Replace(named, "the `semantic-commit` check", "he", 1)
			removed := plain + "2. Removed in 2026.11.26.0400: moved to the semantic-commit check's own requirement\n"
			gets := func(pkg string) string {
				return strings.Replace(testHeader, "adder", pkg, 1) + "\nfunc Test0001_1(t *testing.T) {\n" + acceptance.Tag("0001.1.1.1.1", "it gets their sum") + "}\n"
			}
			out, code := checked(t, acceptance.With(given(removed, named), map[string]string{
				"adder_test.go": gets("adder"),
				"pkgs/fitness-check-semantic-commit/requirements/0001_test.go": gets("requirements"),
			}))
			acceptance.Sees(t, out, code, 0)
		})
	})
}
