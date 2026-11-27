// Package requirements proves the requirements-changed check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

// pairedDoc is the one requirement doc of acceptance.GoRepo.
const pairedDoc = "docs/requirements/0001-adds-numbers.md"

// passed is what a passing single-check run prints.
const passed = "All 1 checks passed"

const readme = "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n"

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 changed product code without its test or doc fails naming the file and the missing half", func(t *testing.T) {
		t.Parallel()
		adder := map[string]string{"adder.go": adderPlus()}
		req.Then(t, "0001.1.2.1.1", "f7956a3e", func(_ req.A) {
			out, code := staged(t, acceptance.GoRepo(), adder)
			acceptance.Sees(t, out, code, 1, "adder.go changed without a test and a requirement doc; change both with it")
			out, code = pull(t, acceptance.GoRepo(), adder, "feat(adder): double numbers")
			acceptance.Sees(t, out, code, 1, "adder.go changed without a test and a requirement doc")
			codeWithComments(t)
		})
		req.Then(t, "0001.1.2.1.2", "74d28512", func(_ req.A) {
			out, code := staged(t, acceptance.GoRepo(), acceptance.With(adder, map[string]string{"adder_test.go": edited(acceptance.GoRepo()["adder_test.go"])}))
			acceptance.Sees(t, out, code, 1, "adder.go changed without a requirement doc; change one with it")
		})
	})
	t.Run("0001.1.2 product code beside a changed test plus a changed doc passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "07b8e753", func(_ req.A) {
			out, code := staged(t, acceptance.GoRepo(), map[string]string{
				"adder.go":      adderPlus(),
				"adder_test.go": edited(acceptance.GoRepo()["adder_test.go"]),
				pairedDoc:       strings.Replace(acceptance.GoRepo()[pairedDoc], "by hand.", "by hand at all.", 1),
			})
			acceptance.Sees(t, out, code, 0, passed)
		})
	})
	t.Run("0001.1.3 docs or config or a rename or a comment change passes since nothing there needs a pair", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "1e1fd798", func(_ req.A) {
			adder := acceptance.GoRepo()["adder.go"]
			for _, files := range []map[string]string{
				{"README.md": readme + "\nMore words.\n"},
				{".fitnessrc.json": `{"proseBudget": {}}` + "\n"},
				{"CHANGELOG.md": "# Changelog\n\nNothing yet.\n"},
				{"adder.go": "", "sum.go": adder},
				{"adder.go": strings.Replace(adder, "\treturn a + b\n", "\n\treturn  a + b  \n", 1)},
			} {
				out, code := staged(t, acceptance.GoRepo(), files)
				acceptance.Sees(t, out, code, 0, passed)
			}
			commentOnly(t)
		})
	})
	t.Run("0001.1.4 a merge replaying product code is not judged again", func(t *testing.T) {
		t.Parallel()
		replayed := func(merge []string, base string) (string, int) {
			repo := acceptance.Example(t, "happyRepo", acceptance.GoRepo())
			if base != "" {
				acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
			}
			acceptance.Git(t, repo, "checkout", "-q", "-b", "side")
			acceptance.Commit(t, repo, map[string]string{"adder.go": adderPlus()}, "feat(adder): double numbers")
			acceptance.Git(t, repo, "checkout", "-q", "-")
			acceptance.Commit(t, repo, map[string]string{"README.md": readme + "\nMore words.\n"}, "docs(readme): add words")
			acceptance.Git(t, repo, append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "merge", "-q", "--no-ff"}, merge...)...)
			if base != "" {
				return acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "requirements-changed")
			}
			return acceptance.Fitness(t, repo, nil, "requirements-changed")
		}
		req.Then(t, "0001.1.2.4.1", "a70e515b", func(_ req.A) {
			out, code := replayed([]string{"--no-commit", "side"}, "")
			acceptance.Sees(t, out, code, 0, passed)
		})
		req.Then(t, "0001.1.2.4.2", "74878220", func(_ req.A) {
			out, code := replayed([]string{"-m", "Merge side", "side"}, "main")
			acceptance.Sees(t, out, code, 1, "adder.go changed without a test and a requirement doc")
		})
	})
	t.Run("0001.1.5 a deletion passes unless a doc names the file then it names the file with the doc and its tests", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "63b7b875", func(_ req.A) {
			given := acceptance.With(acceptance.GoRepo(), map[string]string{"three.go": "package adder\n\n// Three is three.\nconst Three = 3\n"})
			out, code := staged(t, given, map[string]string{"three.go": ""})
			acceptance.Sees(t, out, code, 0, passed)
		})
		named := strings.Replace(acceptance.GoRepo()[pairedDoc], "by hand.", "by hand, with `adder.go`.", 1)
		deleted := func() (string, int) {
			return staged(t, acceptance.With(acceptance.GoRepo(), map[string]string{pairedDoc: named}), map[string]string{"adder.go": ""})
		}
		req.Then(t, "0001.1.2.5.2", "2df5c627", func(_ req.A) {
			out, code := deleted()
			acceptance.Sees(t, out, code, 1, "adder.go was deleted, leaving requirement doc "+pairedDoc)
		})
		req.Then(t, "0001.1.2.5.3", "f6d371f8", func(_ req.A) {
			out, code := deleted()
			acceptance.Sees(t, out, code, 1, "and test adder_test.go unchanged; change or delete them with it")
		})
	})
	t.Run("0001.1.6 a repo with no requirement doc passes since there is no doc to pair with", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "da870fe4", func(_ req.A) {
			given := acceptance.With(acceptance.GoRepo(), map[string]string{pairedDoc: ""})
			out, code := staged(t, given, map[string]string{"adder.go": adderPlus()})
			acceptance.Sees(t, out, code, 0, passed)
		})
	})
	t.Run("0001.1.7 a commit message cannot turn the check off", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.7.1", "25cf7d7e", func(_ req.A) {
			message := "feat(adder): double numbers\n\nRequirements: none, a helper no user calls yet"
			out, code := pull(t, acceptance.GoRepo(), map[string]string{"adder.go": adderPlus()}, message)
			acceptance.Sees(t, out, code, 1, "adder.go changed without a test and a requirement doc")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 vendored or fixture or generated files are not product code", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.2.1.1", "0b8a4d41", func(_ req.A) {
			given := acceptance.With(acceptance.GoRepo(), map[string]string{".gitattributes": "*.gen linguist-generated\n"})
			for _, files := range []map[string]string{
				{"vendor/x/x.go": "package x\n\nconst X = 1\n"},
				{"testdata/x.go": "package x\n\nconst X = 1\n"},
				{"dist/app.js": "console.log(1)\n"},
				{"table.gen": "generated\n"},
				{"zz.go": "// Code generated by hand. DO NOT EDIT.\n\npackage adder\n\nconst Z = 1\n"},
			} {
				out, code := staged(t, given, files)
				acceptance.Sees(t, out, code, 0, passed)
			}
		})
	})
	t.Run("0001.2.2 a full scan judges nothing", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", acceptance.GoRepo())
		acceptance.Commit(t, repo, map[string]string{"adder.go": adderPlus()}, "feat(adder): double numbers")
		out, code := acceptance.Fitness(t, repo, nil, "requirements-changed")
		req.Then(t, "0001.2.2.2.1", "8cc80e9c", func(_ req.A) {
			acceptance.PassedWithNoFiles(t, out, "requirements-changed")
		})
		req.Then(t, "0001.2.2.2.2", "07b8e753", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, passed)
		})
	})
}

// staged writes files over a fresh copy of given, stages them, and runs the
// check on the staged commit.
func staged(t *testing.T, given, files map[string]string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", given)
	acceptance.Write(t, repo, files)
	acceptance.Git(t, repo, "add", "-A")
	return acceptance.Fitness(t, repo, nil, "requirements-changed")
}

// pull commits files with message on a pull request branch off main in a
// copy of given. It runs the check there.
func pull(t *testing.T, given, files map[string]string, message string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", given)
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	acceptance.Commit(t, repo, files, message)
	return acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "requirements-changed")
}

// edited is body with one more function, a code change that keeps Go valid.
func edited(body string) string {
	return body + "\nfunc oneMore() {}\n"
}

// withComment is body with one more comment line, a change to no code.
func withComment(body string) string {
	return body + "\n// One more line.\n"
}

// adderPlus is GoRepo's adder.go with one more function, a product change.
func adderPlus() string {
	return acceptance.GoRepo()["adder.go"] + "\n// Double returns twice a.\nfunc Double(a int) int {\n\treturn a + a\n}\n"
}

// commented is GoRepo's adder.go with every comment rewritten and one added
// after code. Its body is indented anew, as a comment edit may leave it.
func commented() string {
	return "// Package adder sums numbers.\npackage adder\n\n// Add returns a plus b,\n// in either order.\nfunc Add(a, b int) int {\n    return a + b // the sum\n}\n"
}

// script is a shell script, product code, with comments to rewrite.
const script = "#!/bin/sh\n# greet says hello.\necho \"# not a comment\" # says it\necho done\n"

// given is GoRepo with a shell script, a file behind a build tag, and a cgo
// preamble, which the toolchain reads as code.
func given() map[string]string {
	return acceptance.With(acceptance.GoRepo(), map[string]string{
		"greet.sh": script,
		"tag.go":   "//go:build darwin\n\npackage adder\n",
		"cgo.go":   "package adder\n\n// int three(void) { return 3; }\nimport \"C\"\n",
	})
}

// commentOnly proves that edits to comments alone, in Go and shell, pass
// staged and on a pull request.
func commentOnly(t *testing.T) {
	t.Helper()
	files := map[string]string{
		"adder.go": commented(),
		"greet.sh": "#!/bin/sh\n# greet prints a greeting.\n\necho \"# not a comment\" # prints it\n  echo done\n",
	}
	out, code := staged(t, given(), files)
	acceptance.Sees(t, out, code, 0, passed)
	out, code = pull(t, given(), files, "feat(adder): double numbers")
	acceptance.Sees(t, out, code, 0, passed)
}

// codeWithComments proves that a change to code, a directive, or a cgo
// preamble beside comment edits still fails, staged and in a pull. Joined
// lines and a changed string count as code, and so does a test or doc changed
// only in its comments.
func codeWithComments(t *testing.T) {
	t.Helper()
	greeting := strings.Replace(script, "says hello", "prints a greeting", 1)
	for _, c := range []struct{ file, body string }{
		{"adder.go", strings.Replace(commented(), "a + b", "b + a", 1)},
		{"greet.sh", strings.Replace(greeting, "# not a comment\"", "# still not a comment\"", 1)},
		{"greet.sh", strings.Replace(greeting, "# says it\necho", "echo", 1)},
		{"tag.go", "//go:build linux\n\n// Package adder builds on Linux.\npackage adder\n"},
		{"cgo.go", "package adder\n\n// int three(void) { return 4; }\nimport \"C\"\n"},
	} {
		out, code := staged(t, given(), map[string]string{c.file: c.body})
		acceptance.Sees(t, out, code, 1, c.file+" changed without a test and a requirement doc; change both with it")
		out, code = pull(t, given(), map[string]string{c.file: c.body}, "feat(adder): double numbers")
		acceptance.Sees(t, out, code, 1, c.file+" changed without a test and a requirement doc")
	}
	partners := map[string]string{
		"adder.go":      adderPlus(),
		"adder_test.go": withComment(acceptance.GoRepo()["adder_test.go"]),
		pairedDoc:       acceptance.GoRepo()[pairedDoc] + "\n<!-- Sums stay small. -->\n",
	}
	out, code := staged(t, given(), partners)
	acceptance.Sees(t, out, code, 1, "adder.go changed without a test and a requirement doc; change both with it")
}
