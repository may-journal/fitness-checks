package requirements

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/may-journal/fitness-checks/internal/program"
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"engine-free go dependencies for every check", func(t *testing.T) {
			t.Parallel()
			for _, folder := range checkFolders(t, "..") {
				for _, dep := range strayDeps(t, folder) {
					t.Errorf("%s depends on %s", filepath.Base(folder), dep)
				}
			}
		}},
		{"checks read no config key the engine declares", func(t *testing.T) {
			t.Parallel()
			for _, use := range engineConfigReads(t, "..") {
				t.Errorf("%s reads engine config; declare the key in the check", use)
			}
		}},
		{"describe prints name summary scope and the keys with their defaults", func(t *testing.T) {
			t.Parallel()
			ours := goBuild(t, "./cmd/...")
			for _, folder := range checkFolders(t, "..") {
				name := filepath.Base(folder)
				got := describeMap(t, filepath.Join(ours, name))
				for _, problem := range append(describeProblems(name, got), keyProblems(name, got)...) {
					t.Error(problem)
				}
			}
		}},
		{"config schema and explain show each check's keys", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			var schema map[string]any
			stdoutJSON(t, repo, &schema, "config", "schema")
			for check, keys := range checkKeys {
				schemaLists(t, schema, keys)
				if check != byNameOnly {
					explainLists(t, repo, check, keys)
				}
			}
		}},
		{"installed engine starts each check program", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			all, code := fitness(t, repo, nil)
			sees(t, all, code, 0)
			for _, folder := range checkFolders(t, "..") {
				check := strings.TrimPrefix(filepath.Base(folder), "fitness-check-")
				out := all
				if !strings.Contains(all, "→ "+check+"\n") {
					out, _ = fitness(t, repo, nil, check)
				}
				if !ranOurs(out, check) {
					t.Errorf("the engine did not run this repo's program for %s", check)
				}
			}
		}},
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"one file in one check names only that check", func(t *testing.T) {
			t.Parallel()
			clone, base := changedClone(t, "cmd/fitness-check-gofmt/main.go")
			if got := changedChecks(t, clone, base); !slices.Equal(got, []string{"gofmt"}) {
				t.Errorf("changed-checks = %v, want only gofmt", got)
			}
		}},
		{"shared helper names dependent checks", func(t *testing.T) {
			t.Parallel()
			for _, file := range []string{"internal/spell/mask.go", "internal/spell/dict/git.txt"} {
				clone, base := changedClone(t, file)
				if got := changedChecks(t, clone, base); !slices.Equal(got, []string{"cspell"}) {
					t.Errorf("after %s, changed-checks = %v, want only cspell, the one check that imports spell", file, got)
				}
			}
		}},
		{"go mod change or no base names all checks", func(t *testing.T) {
			t.Parallel()
			clone, base := changedClone(t, "go.mod")
			var every []string
			for _, folder := range checkFolders(t, clone) {
				every = append(every, strings.TrimPrefix(filepath.Base(folder), "fitness-check-"))
			}
			for _, got := range [][]string{changedChecks(t, clone, base), changedChecks(t, clone)} {
				if !slices.Equal(got, every) {
					t.Errorf("changed-checks = %v, want every check, %v", got, every)
				}
			}
		}},
		{"test script builds vets and tests one check", func(t *testing.T) {
			t.Parallel()
			script, err := filepath.Abs(filepath.Join("..", "scripts", "test-check.sh"))
			mustDo(t, err)
			out, err := exec.Command(script, "gofmt").CombinedOutput()
			if err != nil || !strings.Contains(string(out), `{"name":"gofmt"`) {
				t.Errorf("test-check.sh gofmt: %v\n%s", err, out)
			}
		}},
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"major tag is never the last release", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "fix(gofmt): change gofmt", true, "cmd/fitness-check-gofmt/main.go")
			git(t, clone, "tag", "fitness-check-gofmt-v1", "HEAD")
			if got := releasePlan(t, clone); !slices.Equal(got, []string{"fitness-check-gofmt-v1.0.1"}) {
				t.Errorf("release-plan = %v, want only gofmt at 1.0.1", got)
			}
		}},
		{"fix commit plans patch for one check", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "fix(gofmt): change gofmt", true, "cmd/fitness-check-gofmt/main.go")
			if got := releasePlan(t, clone); !slices.Equal(got, []string{"fitness-check-gofmt-v1.0.1"}) {
				t.Errorf("release-plan = %v, want only gofmt at 1.0.1", got)
			}
		}},
		{"feat on helper plans minor for dependents", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "feat(spell): change spell", true, "internal/spell/mask.go")
			if got := releasePlan(t, clone); !slices.Equal(got, []string{"fitness-check-cspell-v1.1.0"}) {
				t.Errorf("release-plan = %v, want only cspell, the one check that imports spell, at 1.1.0", got)
			}
		}},
		{"deps commit plans patch for all checks", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "build(deps): change go.mod", true, "go.mod")
			if got, want := releasePlan(t, clone), tagsFor(t, clone, "1.0.1"); !slices.Equal(got, want) {
				t.Errorf("release-plan = %v, want every check at 1.0.1, %v", got, want)
			}
		}},
		{"docs breaking and revert commit kinds", func(t *testing.T) {
			t.Parallel()
			cases := map[string][]string{
				"docs(gofmt): change gofmt":          {},
				"refactor(gofmt)!: change gofmt":     {"fitness-check-gofmt-v2.0.0"},
				`Revert "feat(gofmt): change gofmt"`: {"fitness-check-gofmt-v1.0.1"},
			}
			for message, want := range cases {
				clone := releaseClone(t, message, true, "cmd/fitness-check-gofmt/main.go")
				if got := releasePlan(t, clone); !slices.Equal(got, want) {
					t.Errorf("after %q, release-plan = %v, want %v", message, got, want)
				}
			}
		}},
		{"unreleased check plans first release", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "docs(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
			if got, want := releasePlan(t, clone), tagsFor(t, clone, "1.0.0"); !slices.Equal(got, want) {
				t.Errorf("release-plan = %v, want every check's first release, %v", got, want)
			}
		}},
	})
}

func Test0001_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"archive and checksum per platform", func(t *testing.T) {
			t.Parallel()
			out := filepath.Join(t.TempDir(), "release")
			build := exec.Command("./scripts/release-build.sh", "fitness-check-adr-checks-v1.2.3", out)
			build.Dir = ".."
			if got, err := build.CombinedOutput(); err != nil {
				t.Fatalf("release-build.sh: %v\n%s", err, got)
			}
			for _, platform := range []string{"linux-amd64", "linux-arm64", "darwin-amd64", "darwin-arm64"} {
				for _, suffix := range []string{".tar", ".tar.sha256", ".tar.gz", ".tar.gz.sha256"} {
					name := "fitness-check-adr-checks-v1.2.3-" + platform + suffix
					if _, err := os.Stat(filepath.Join(out, name)); err != nil {
						t.Errorf("the release lacks %s: %v", name, err)
					}
				}
			}
		}},
		{"stripped program is smaller in gzip", func(t *testing.T) {
			t.Parallel()
			out := filepath.Join(t.TempDir(), "release")
			got, code := runSize(t, nil, "release-build.sh", "fitness-check-gofmt-v1.2.3", out)
			if code != 0 {
				t.Fatalf("release-build.sh exited %d: %s", code, got)
			}
			unpacked := t.TempDir()
			mustDo(t, exec.Command("tar", "-xzf", filepath.Join(out, "fitness-check-gofmt-v1.2.3-linux-arm64.tar.gz"), "-C", unpacked).Run())
			plain := filepath.Join(t.TempDir(), "plain")
			build := exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-o", plain, "./cmd/fitness-check-gofmt")
			build.Dir = ".."
			build.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=arm64")
			if b, err := build.CombinedOutput(); err != nil {
				t.Fatalf("go build: %v\n%s", err, b)
			}
			stripped, err := os.Stat(filepath.Join(unpacked, "fitness-check-gofmt"))
			mustDo(t, err)
			unstripped, err := os.Stat(plain)
			mustDo(t, err)
			if stripped.Size() >= unstripped.Size() {
				t.Errorf("the released program is %d bytes, want fewer than the unstripped %d", stripped.Size(), unstripped.Size())
			}
		}},
	})
}

func Test0001_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"major tags move to newest release on every rerun", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
			bare := remoteOf(t, clone)
			cutRelease(t, clone, "1.0.0")
			one := cutRelease(t, clone, "1.2.0")
			two := cutRelease(t, clone, "2.0.0")
			git(t, clone, "tag", "fitness-check-gofmt-v1", "HEAD~2")
			git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
			runScriptWith(t, clone, bare, "release-major.sh")
			got := []string{tagCommit(bare, "fitness-check-gofmt-v1"), tagCommit(bare, "fitness-check-gofmt-v2")}
			if calls := remoteCalls(t, bare); !strings.Contains(calls, "-X POST repos/acme/checks/git/refs -f ref=refs/tags/fitness-check-gofmt-v2") {
				t.Errorf("the missing v2 tag was not created with a POST: %s", calls)
			}
			if !slices.Equal(got, []string{one, two}) {
				t.Errorf("remote major tags = %v, want %v", got, []string{one, two})
			}
			before := remoteCalls(t, bare)
			runScriptWith(t, clone, bare, "release-major.sh")
			if after := remoteCalls(t, bare); after != before {
				t.Errorf("the second run called gh again: %s", after)
			}
		}},
		{"forbidden create exits with github message", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
			bare := remoteOf(t, clone)
			cutRelease(t, clone, "1.0.0")
			git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
			out, err := forbiddenMajor(t, clone, bare)
			if err == nil || !strings.Contains(string(out), "HTTP 403") || tagCommit(bare, "fitness-check-gofmt-v1") != "" {
				t.Errorf("run = %v with %q, want an error naming HTTP 403 and no tag created", err, out)
			}
		}},
		{"non-ancestor major tag is forced forward", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
			bare := remoteOf(t, clone)
			release := cutRelease(t, clone, "1.0.0")
			git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
			// An orphan commit in the remote is no ancestor of the release.
			tree := gitOut(t, bare, "rev-parse", "main^{tree}")
			other := gitOut(t, bare, "-c", "user.name=t", "-c", "user.email=t@t", "commit-tree", "-m", "other line", tree)
			git(t, bare, "update-ref", "refs/tags/fitness-check-gofmt-v1", other)
			runScriptWith(t, clone, bare, "release-major.sh")
			if got := tagCommit(bare, "fitness-check-gofmt-v1"); got != release {
				t.Errorf("remote major tag = %q, want the release %q, moved by force", got, release)
			}
		}},
		{"refused move names tag and never creates", func(t *testing.T) {
			t.Parallel()
			clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
			bare := remoteOf(t, clone)
			first := cutRelease(t, clone, "1.0.0")
			git(t, clone, "tag", "fitness-check-gofmt-v1", "HEAD~1")
			git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
			git(t, bare, "update-ref", "refs/tags/fitness-check-gofmt-v1", first)
			git(t, clone, "tag", "-d", "fitness-check-gofmt-v1")
			out, err := forbiddenMajor(t, clone, bare)
			if err == nil || !strings.Contains(string(out), "could not move fitness-check-gofmt-v1") || strings.Contains(remoteCalls(t, bare), "-X POST") {
				t.Errorf("run = %v with %q, want an error naming the tag and no POST for an existing tag", err, out)
			}
		}},
	})
}

func Test0001_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"index holds only the newest published version of each major", func(t *testing.T) {
			t.Parallel()
			file := filepath.Join(t.TempDir(), "releases.json")
			mustDo(t, os.WriteFile(file, []byte(publishedReleases), 0o644))
			var got, want any
			mustDo(t, json.Unmarshal([]byte(runScript(t, "..", "release-index.sh", file)), &got))
			mustDo(t, json.Unmarshal([]byte(publishedIndex), &want))
			if !reflect.DeepEqual(got, want) {
				t.Errorf("index = %v, want %v", got, want)
			}
		}},
		{"describe fields appear only on newer version", func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			releases, describes := filepath.Join(dir, "releases.json"), filepath.Join(dir, "describes.json")
			mustDo(t, os.WriteFile(releases, []byte(publishedReleases), 0o644))
			mustDo(t, os.WriteFile(describes, []byte(`{"fitness-check-gofmt-v1.10.0":
		{"name":"gofmt","scope":"files","files":"\\.(go)$","exclude":"(^|/)vendor/","needs":"","requires":["repo"]}}`), 0o644))
			var index struct {
				Checks map[string]struct {
					Versions map[string]map[string]any `json:"versions"`
				} `json:"checks"`
			}
			mustDo(t, json.Unmarshal([]byte(runScript(t, "..", "release-index.sh", releases, describes)), &index))
			versions := index.Checks["gofmt"].Versions
			got := map[string]any{}
			for _, key := range []string{"files", "exclude", "needs", "requires", "name", "scope"} {
				if v, ok := versions["1.10.0"][key]; ok {
					got[key] = v
				}
			}
			want := map[string]any{"files": `\.(go)$`, "exclude": "(^|/)vendor/", "requires": []any{"repo"}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("1.10.0 copies %v from its describe line, want only %v", got, want)
			}
			if _, ok := versions["2.0.0"]["files"]; ok {
				t.Errorf("2.0.0 has no describe line, so it must hold no files: %v", versions["2.0.0"])
			}
		}},
		{"index downloads only new versions", func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			published := filepath.Join(dir, "published.json")
			mustDo(t, os.WriteFile(published, []byte(`{"checks":{"gofmt":{"versions":{"1.1.0":{"tag":"fitness-check-gofmt-v1.1.0","files":"kept"}}}}}`), 0o644))
			// The 1.0.0 and 1.1.0 describe files are missing, so a download of either would fail the script.
			releases := reuseIndex(t, dir, false)
			got := indexFields(t, runScript(t, "..", "release-index.sh", releases, "", published))
			if want := map[string]any{"1.1.0": "kept", "2.0.0": `\.(2.0.0)$`}; !reflect.DeepEqual(got, want) {
				t.Errorf("fields with a published index = %v, want %v", got, want)
			}
			releases = reuseIndex(t, dir, true)
			unreadable := filepath.Join(dir, "unreadable.json")
			mustDo(t, os.WriteFile(unreadable, []byte("not json"), 0o644))
			empty := filepath.Join(dir, "empty.json")
			bare := filepath.Join(dir, "bare.json")
			mustDo(t, os.WriteFile(empty, nil, 0o644))
			mustDo(t, os.WriteFile(bare, []byte(`{"checks":{"gofmt":{"versions":{"1.1.0":{"tag":"fitness-check-gofmt-v1.1.0"}}}}}`), 0o644))
			for name, index := range map[string]string{"unreadable": unreadable, "missing (empty)": empty, "field-less entry": bare} {
				got = indexFields(t, runScript(t, "..", "release-index.sh", releases, "", index))
				if want := map[string]any{"1.1.0": `\.(1.1.0)$`, "2.0.0": `\.(2.0.0)$`}; !reflect.DeepEqual(got, want) {
					t.Errorf("fields with a %s index = %v, want %v", name, got, want)
				}
			}
		}},
		{"tar url stays with optional gz", func(t *testing.T) {
			t.Parallel()
			file := filepath.Join(t.TempDir(), "releases.json")
			mustDo(t, os.WriteFile(file, []byte(sizedReleases), 0o644))
			var got, want any
			mustDo(t, json.Unmarshal([]byte(runScript(t, "..", "release-index.sh", file)), &got))
			mustDo(t, json.Unmarshal([]byte(sizedIndex), &want))
			if !reflect.DeepEqual(got, want) {
				t.Errorf("index = %v, want %v", got, want)
			}
		}},
	})
}

func Test0001_7(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"size baseline drift over ten percent fails", func(t *testing.T) {
			t.Parallel()
			baseline := filepath.Join(t.TempDir(), "baseline.tsv")
			env := []string{"SIZE_BASELINE=" + baseline, "SIZE_BUILD_DIR=" + t.TempDir()}
			if out, code := runSize(t, env, "check-binary-sizes.sh", "--update"); code != 0 {
				t.Fatalf("--update exited %d: %s", code, out)
			}
			if out, code := runSize(t, env, "check-binary-sizes.sh"); code != 0 {
				t.Errorf("current sizes exited %d, want 0: %s", code, out)
			}
			b, err := os.ReadFile(baseline)
			mustDo(t, err)
			if rows := strings.Count(string(b), "\nfitness-check-"); rows != len(checkFolders(t, "..")) {
				t.Errorf("the baseline lists %d programs, want one for each check folder", rows)
			}
			grown := regexp.MustCompile(`(?m)^fitness-check-gofmt\t\d+$`).ReplaceAllString(string(b), "fitness-check-gofmt\t1000")
			shrunk := regexp.MustCompile(`(?m)^fitness-check-gofmt\t\d+$`).ReplaceAllString(string(b), "fitness-check-gofmt\t900000000")
			stale := string(b) + "fitness-check-gone\t1000\n"
			for want, text := range map[string]string{"SIZE: fitness-check-gofmt is": grown, "under its baseline": shrunk, "SIZE: fitness-check-gone": stale} {
				sizeFails(t, env, baseline, text, want)
			}
		}},
		{"measure script records sizes and timing", func(t *testing.T) {
			t.Parallel()
			repo := filepath.Join(t.TempDir(), "repo")
			mustDo(t, os.CopyFS(repo, os.DirFS(filepath.Join("..", "tests", "examples", "happyRepo"))))
			git(t, repo, "init", "-q")
			git(t, repo, "add", "-A")
			git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "docs(readme): add readme")
			out := filepath.Join(t.TempDir(), "measure")
			if got, code := runSize(t, []string{"MEASURE_REPO=" + repo, "MEASURE_RUNS=1"}, "measure-checks.sh", out); code != 0 {
				t.Fatalf("measure-checks.sh exited %d: %s", code, got)
			}
			sizes, err := os.ReadFile(filepath.Join(out, "sizes.tsv"))
			mustDo(t, err)
			timing, err := os.ReadFile(filepath.Join(out, "timing.tsv"))
			mustDo(t, err)
			rows := strings.Split(strings.TrimSpace(string(sizes)), "\n")
			if want := len(checkFolders(t, "..")) + 1; len(rows) != want || strings.Count(string(timing), "\n") != 2 || strings.Contains(string(sizes), "\t0\t") {
				t.Errorf("sizes has %d rows and timing %q, want %d rows with sizes, and one timed run\n%s", len(rows), timing, want, sizes)
			}
		}},
	})
}

func Test0001_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"README names the check and opens with the described summary", func(t *testing.T) {
			t.Parallel()
			var ours string
			eachReadme(t, func(folder, name, readme string) {
				if ours == "" {
					ours = goBuild(t, "./cmd/...")
				}
				checkDocumentedReadme(t, ours, folder, name, readme)
			})
		}},
		{"relative links resolve", func(t *testing.T) {
			t.Parallel()
			eachReadme(t, func(folder, _, readme string) {
				for _, target := range readmeLinks(readme) {
					if _, err := os.Stat(filepath.Join(folder, target)); err != nil {
						t.Errorf("%s/README.md links %s, which does not resolve: %v", folder, target, err)
					}
				}
			})
		}},
	})
}

// checkDocumentedReadme fails t for each way a check folder's README differs
// from what its binary says.
func checkDocumentedReadme(t *testing.T, ours, folder, name, readme string) {
	t.Helper()
	summary := program.Summary(readme)
	if !plainSentence(summary) {
		t.Errorf("%s/README.md has no plain sentence in the paragraph under its title; --describe prints %q", folder, summary)
	}
	if described := describeMap(t, filepath.Join(ours, "fitness-check-"+name))["summary"]; described != summary {
		t.Errorf("%s --describe summary is %q, want the README's %q", name, described, summary)
	}
	if title := "# " + name + "\n"; !strings.HasPrefix(readme, title) {
		t.Errorf("%s/README.md does not open with the title %q", folder, title)
	}
}

// plainSentence reports whether summary starts with a letter and holds no
// Markdown link, as a sentence of prose does.
func plainSentence(summary string) bool {
	return summary != "" && unicode.IsLetter([]rune(summary)[0]) && !strings.Contains(summary, "](")
}

// eachReadme calls visit with each check folder that has a README.md, its
// name, and the README's text.
func eachReadme(t *testing.T, visit func(folder, name, readme string)) {
	t.Helper()
	for _, folder := range checkFolders(t, "..") {
		if readme, ok := readFolder(t, folder, "README.md"); ok {
			visit(folder, strings.TrimPrefix(filepath.Base(folder), "fitness-check-"), readme)
		}
	}
}

// relativeLink matches a Markdown link whose target is a file path, without
// its #fragment.
var relativeLink = regexp.MustCompile(`\]\(([^)#:\s]+)(?:#[^)]*)?(?:\s+"[^"]*")?\)`)

// referenceLink matches a reference-style link definition, "[x]: path", whose
// target is a file path, without its #fragment.
var referenceLink = regexp.MustCompile(`(?m)^ {0,3}\[[^\]]+\]:\s*<?([^\s>#:]+)(?:#\S*)?>?(?:\s+.*)?$`)

// codeFence matches a fenced code block, which holds example text and no links.
var codeFence = regexp.MustCompile("(?ms)^ {0,3}```.*?^ {0,3}```[^\\n]*$|^ {0,3}~~~.*?^ {0,3}~~~[^\\n]*$")

// readmeLinks lists the relative file targets the Markdown text links, inline
// or by reference, outside fenced code.
func readmeLinks(readme string) []string {
	readme = codeFence.ReplaceAllString(readme, "")
	var targets []string
	for _, m := range relativeLink.FindAllStringSubmatch(readme, -1) {
		targets = append(targets, m[1])
	}
	for _, m := range referenceLink.FindAllStringSubmatch(readme, -1) {
		targets = append(targets, m[1])
	}
	return targets
}

// engineOnly reports whether dep is, or lies under, an engine package no
// check may depend on.
func engineOnly(dep string) bool {
	if strings.HasSuffix(dep, "/cmd/fitness") {
		return true
	}
	return slices.ContainsFunc([]string{"metadata", "aftereffect", "dupe", "report"}, func(name string) bool {
		return strings.HasSuffix(dep, "/internal/"+name) || strings.Contains(dep, "/internal/"+name+"/")
	})
}

// engineConfig names, by import path, the engine identifiers that read a
// check's config key or file list. A check declares its own instead.
var engineConfig = map[string][]string{
	"github.com/may-journal/fitness-runner/go/pkg/conf":        {"Load", "Defaults", "Config", "MaxListItems"},
	"github.com/may-journal/fitness-runner/go/pkg/adr":         {"LoadAllow", "Allow"},
	"github.com/may-journal/fitness-runner/go/pkg/lang/golang": {"LoadConfig", "Defaults", "Config"},
	"github.com/may-journal/fitness-runner/go/pkg/lang":        {"ScanExts"},
}

// engineConfigReads lists each use, as file:line name, of an engineConfig
// identifier in the Go files under root's cmd and internal folders.
func engineConfigReads(t *testing.T, root string) []string {
	t.Helper()
	var uses []string
	fset := token.NewFileSet()
	for _, dir := range []string{"cmd", "internal"} {
		mustDo(t, filepath.WalkDir(filepath.Join(root, dir), func(name string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(name) != ".go" {
				return err
			}
			file, err := parser.ParseFile(fset, name, nil, parser.SkipObjectResolution)
			if err != nil {
				return err
			}
			uses = append(uses, configSelectors(fset, file)...)
			return nil
		}))
	}
	return uses
}

// configSelectors lists each selector in file that names an engineConfig
// identifier through its import.
func configSelectors(fset *token.FileSet, file *ast.File) []string {
	names := map[string][]string{}
	for _, spec := range file.Imports {
		imported, _ := strconv.Unquote(spec.Path.Value)
		if banned, ok := engineConfig[imported]; ok {
			local := path.Base(imported)
			if spec.Name != nil {
				local = spec.Name.Name
			}
			names[local] = banned
		}
	}
	var uses []string
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); ok && slices.Contains(names[pkg.Name], sel.Sel.Name) {
			uses = append(uses, fmt.Sprintf("%s %s.%s", fset.Position(sel.Pos()), pkg.Name, sel.Sel.Name))
		}
		return true
	})
	return uses
}

// checkFolders lists the check folders of the repo at root, such as
// cmd/fitness-check-gofmt.
func checkFolders(t *testing.T, root string) []string {
	t.Helper()
	folders, err := filepath.Glob(filepath.Join(root, "cmd", "fitness-check-*"))
	mustDo(t, err)
	if len(folders) == 0 {
		t.Fatal("no check folders under cmd")
	}
	return folders
}

// strayDeps lists what the check in folder depends on, tests included, that
// is the engine or another check.
func strayDeps(t *testing.T, folder string) []string {
	t.Helper()
	own := "/cmd/" + filepath.Base(folder)
	cmd := exec.Command("go", "list", "-deps", "-test", "-f", "{{.ImportPath}}", "."+own)
	cmd.Dir = ".."
	out, err := cmd.Output()
	mustDo(t, err)
	return slices.DeleteFunc(strings.Fields(string(out)), func(dep string) bool {
		pkg := strings.TrimSuffix(strings.Trim(dep, "[]"), ".test")
		other := strings.Contains(pkg, "/cmd/fitness-check-") && !strings.HasSuffix(pkg, own)
		return !other && !engineOnly(pkg)
	})
}

// goBuild builds the packages into a fresh folder and returns it.
func goBuild(t *testing.T, packages ...string) string {
	t.Helper()
	bin := t.TempDir()
	build := exec.Command("go", append([]string{"build", "-o", bin + string(os.PathSeparator)}, packages...)...)
	build.Dir = ".."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

// describeMap runs program with --describe and decodes what it prints.
func describeMap(t *testing.T, program string) map[string]any {
	t.Helper()
	out, err := exec.Command(program, "--describe").Output()
	mustDo(t, err)
	var got map[string]any
	mustDo(t, json.Unmarshal(out, &got))
	return got
}

// describeProblems lists what is wrong with a check's --describe output.
func describeProblems(name string, got map[string]any) []string {
	var problems []string
	if got["protocol"] != float64(1) {
		problems = append(problems, fmt.Sprintf("%s --describe has protocol %v, want 1", name, got["protocol"]))
	}
	if got["name"] != strings.TrimPrefix(name, "fitness-check-") {
		problems = append(problems, fmt.Sprintf("%s --describe names the check %v", name, got["name"]))
	}
	if summary, _ := got["summary"].(string); summary == "" {
		problems = append(problems, name+" --describe has no summary")
	}
	if !slices.Contains([]any{"files", "blob", "git", "github"}, got["scope"]) {
		problems = append(problems, fmt.Sprintf("%s --describe has scope %v, want files, blob, git or github", name, got["scope"]))
	}
	return problems
}

// keyProblems lists how the keys a check's --describe declares differ from
// checkKeys.
func keyProblems(name string, got map[string]any) []string {
	declared := declaredLeaves(got["config"])
	if want := checkKeys[strings.TrimPrefix(name, "fitness-check-")]; !reflect.DeepEqual(declared, want) {
		return []string{fmt.Sprintf("%s --describe declares %v, want %v", name, declared, want)}
	}
	return nil
}

// checkKeys maps each check that reads .fitnessrc.json keys to each leaf key
// it declares and that key's default, nil for none. No other check declares
// a key; adr-checks reads adrAsIssue, which adr-as-issue declares.
var checkKeys = map[string]map[string]any{
	"repeated-string-literals": {"repeatedStringLiterals.allow": nil},
	"text-readability":         {"textReadability.maxGrade": 18.0, "textReadability.maxLix": 60.0, "textReadability.minWords": 100.0},
	"prose-budget": {
		"proseBudget.maxSentenceWords": 23.0, "proseBudget.maxParagraphSentences": 4.0, "proseBudget.maxSectionParagraphs": 3.0,
		"proseBudget.maxListItemWords": 23.0, "proseBudget.maxListItems": 8.0, "proseBudget.maxWords": 300.0,
		"prose.comments.maxWordsPerToken": 3.0, "prose.comments.minBodyComments": 2.0,
	},
	"comment-structure": {"comment.maxWords": 100.0, "comment.maxTables": 1.0},
	"adr-as-issue":      {"adrAsIssue.allow": nil},
	"go-test-coverage":  {"goTestCoverage.entries": nil},
	"go-complexity":     {"goComplexity.max": 5.0},
	"localized-text":    {"localizedText.reference": nil, "localizedText.allow": nil},
	"requirements":      {"requirements.maxAcceptances": 8.0},
}

// byNameOnly is the check the engine's catalog leaves out, since it runs only
// when named, so `fitness explain` does not answer for it.
const byNameOnly = "go-test-coverage"

// declaredLeaves maps each leaf key in a decoded config list, as --describe
// and `fitness explain --json` print it, to its default; nil for no keys.
func declaredLeaves(config any) map[string]any {
	leaves := map[string]any{}
	list, _ := config.([]any)
	for _, item := range list {
		key, _ := item.(map[string]any)
		if inner, nested := key["keys"]; nested {
			maps.Copy(leaves, declaredLeaves(inner))
		} else if path, _ := key["key"].(string); path != "" {
			leaves[path] = key["default"]
		}
	}
	if len(leaves) == 0 {
		return nil
	}
	return leaves
}

// schemaLists fails each of keys the schema leaves out or gives another
// default.
func schemaLists(t *testing.T, schema map[string]any, keys map[string]any) {
	t.Helper()
	for path, want := range keys {
		if got, ok := schemaDefault(schema, path); !ok || got != want {
			t.Errorf("the schema lists %s with default %v (found %v), want %v", path, got, ok, want)
		}
	}
}

// explainLists fails when `fitness explain --json` lists other keys for
// check than keys.
func explainLists(t *testing.T, repo, check string, keys map[string]any) {
	t.Helper()
	var explained map[string]any
	stdoutJSON(t, repo, &explained, "explain", check, "--json")
	if got := declaredLeaves(explained["config"]); !reflect.DeepEqual(got, keys) {
		t.Errorf("fitness explain %s lists %v, want %v", check, got, keys)
	}
}

// schemaDefault is the default the JSON Schema gives the dotted key path,
// and whether the schema lists the key.
func schemaDefault(schema map[string]any, path string) (any, bool) {
	node := schema
	for _, name := range strings.Split(path, ".") {
		properties, _ := node["properties"].(map[string]any)
		next, ok := properties[name].(map[string]any)
		if !ok {
			return nil, false
		}
		node = next
	}
	return node["default"], true
}

// stdoutJSON runs fitness with args in repo and decodes its stdout into into.
func stdoutJSON(t *testing.T, repo string, into any, args ...string) {
	t.Helper()
	var stdout, stderr strings.Builder
	cmd := started(t, func() *exec.Cmd {
		stdout.Reset()
		stderr.Reset()
		cmd := exec.Command(installer(t), append([]string{"--"}, args...)...)
		cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = repo, userEnv(), &stdout, &stderr
		return cmd
	}, (*exec.Cmd).Run)
	if code := cmd.ProcessState.ExitCode(); code != 0 {
		t.Fatalf("fitness %s exited %d:\n%s%s", strings.Join(args, " "), code, stdout.String(), stderr.String())
	}
	mustDo(t, json.Unmarshal([]byte(stdout.String()), into))
}

// ranOurs reports whether out shows this repo's program for check ran,
// or skipped by its own --describe. For local-branch, the program declares.
func ranOurs(out, check string) bool {
	if check == "local-branch" {
		return strings.Contains(out, "fitness-check-"+check+" declares the "+check+" check, which fitness judges itself")
	}
	return strings.Contains(out, "→ "+check+"\n") || strings.Contains(row(out, check), "skipped")
}

// changedClone clones this repo, appends a line to each file, commits, and
// returns the clone with the commit before the change.
func changedClone(t *testing.T, files ...string) (string, string) {
	t.Helper()
	clone := filepath.Join(t.TempDir(), "repo")
	if out, err := exec.Command("git", "clone", "-q", "--local", "..", clone).CombinedOutput(); err != nil {
		t.Fatalf("git clone: %v\n%s", err, out)
	}
	base, err := exec.Command("git", "-C", clone, "rev-parse", "HEAD").Output()
	mustDo(t, err)
	for _, file := range files {
		f, err := os.OpenFile(filepath.Join(clone, file), os.O_APPEND|os.O_WRONLY, 0)
		mustDo(t, err)
		_, err = f.WriteString("\n// changed\n")
		mustDo(t, errors.Join(err, f.Close()))
	}
	git(t, clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qam", "test(change): change files")
	return clone, string(base[:len(base)-1])
}

// changedChecks runs this checkout's scripts/changed-checks.sh in clone with
// args and decodes the checks it names.
func changedChecks(t *testing.T, clone string, args ...string) []string {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("..", "scripts", "changed-checks.sh"))
	mustDo(t, err)
	cmd := exec.Command(script, args...)
	cmd.Dir = clone
	out, err := cmd.Output()
	mustDo(t, err)
	var checks []string
	mustDo(t, json.Unmarshal(out, &checks))
	return checks
}

// releaseClone changes files in a clone of this repo as one commit with the
// message given. When tagged, every check gets a v1.0.0 tag on the commit
// before the change, as if it had released there.
func releaseClone(t *testing.T, message string, tagged bool, files ...string) string {
	t.Helper()
	clone, base := changedClone(t, files...)
	git(t, clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "--amend", "-m", message)
	// The clone copies this checkout's tags, real releases included.
	tags, err := exec.Command("git", "-C", clone, "tag", "-l").Output()
	mustDo(t, err)
	if names := strings.Fields(string(tags)); len(names) > 0 {
		git(t, clone, append([]string{"tag", "-d"}, names...)...)
	}
	if tagged {
		for _, check := range everyCheck(t, clone) {
			git(t, clone, "tag", "fitness-check-"+check+"-v1.0.0", base)
		}
	}
	return clone
}

// releasePlan runs this checkout's scripts/release-plan.sh in clone and
// decodes the tags it plans.
func releasePlan(t *testing.T, clone string) []string {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("..", "scripts", "release-plan.sh"))
	mustDo(t, err)
	cmd := exec.Command(script)
	cmd.Dir = clone
	out, err := cmd.Output()
	mustDo(t, err)
	var tags []string
	mustDo(t, json.Unmarshal(out, &tags))
	return tags
}

// everyCheck names each check folder's check in the repo at root.
func everyCheck(t *testing.T, root string) []string {
	t.Helper()
	var every []string
	for _, folder := range checkFolders(t, root) {
		every = append(every, strings.TrimPrefix(filepath.Base(folder), "fitness-check-"))
	}
	return every
}

// tagsFor names the tag each check of the repo at root gets at version.
func tagsFor(t *testing.T, root, version string) []string {
	t.Helper()
	var tags []string
	for _, check := range everyCheck(t, root) {
		tags = append(tags, "fitness-check-"+check+"-v"+version)
	}
	return tags
}

// remoteOf gives clone an empty bare remote of its own as origin, so a
// script's push never reaches this checkout.
func remoteOf(t *testing.T, clone string) string {
	t.Helper()
	bare := filepath.Join(t.TempDir(), "remote.git")
	git(t, clone, "init", "-q", "--bare", bare)
	git(t, bare, "config", "receive.shallowUpdate", "true")
	git(t, clone, "remote", "set-url", "origin", bare)
	return bare
}

// runScript runs this checkout's script in dir with args and returns its
// standard output.
func runScript(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	return runScriptWith(t, dir, "", name, args...)
}

// runScriptWith is runScript with gh standing in for GitHub over the bare
// remote, when there is one.
func runScriptWith(t *testing.T, dir, remote, name string, args ...string) string {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("..", "scripts", name))
	mustDo(t, err)
	cmd := exec.Command(script, args...)
	cmd.Dir = dir
	if remote != "" {
		cmd.Env = append(os.Environ(), fakeGh(t, remote)...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", name, err, out)
	}
	return string(out)
}

// forbiddenMajor runs release-major.sh in clone with GitHub refusing every
// tag write, and returns what it printed and its error.
func forbiddenMajor(t *testing.T, clone, bare string) ([]byte, error) {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("..", "scripts", "release-major.sh"))
	mustDo(t, err)
	cmd := exec.Command(script)
	cmd.Dir = clone
	cmd.Env = append(append(os.Environ(), fakeGh(t, bare)...), "FAKE_PATCH_ERROR=gh: Forbidden (HTTP 403)")
	return cmd.CombinedOutput()
}

// tagCommit returns the commit a tag names in repo, empty when it is missing.
func tagCommit(repo, tag string) string {
	out, _ := exec.Command("git", "-C", repo, "rev-parse", "-q", "--verify", "refs/tags/"+tag+"^{commit}").Output()
	return strings.TrimSpace(string(out))
}

// cutRelease tags an empty commit as the check's release of version, and
// returns its commit.
func cutRelease(t *testing.T, clone, version string) string {
	t.Helper()
	git(t, clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "--allow-empty", "-m", "chore: release "+version)
	git(t, clone, "tag", "fitness-check-gofmt-v"+version)
	return tagCommit(clone, "fitness-check-gofmt-v"+version)
}

// fakeGh puts the stand-in gh first on PATH, applying tag calls to the bare
// remote given in FAKE_REMOTE.
func fakeGh(t *testing.T, remote string) []string {
	t.Helper()
	bin := t.TempDir()
	shim, err := filepath.Abs(filepath.Join("..", "tests", "fake-gh-refs.sh"))
	mustDo(t, err)
	mustDo(t, os.Symlink(shim, filepath.Join(bin, "gh")))
	return []string{"PATH=" + bin + string(os.PathListSeparator) + os.Getenv("PATH"), "GITHUB_REPOSITORY=acme/checks", "FAKE_REMOTE=" + remote}
}

// publishedReleases is what GitHub lists for one check: two published
// versions of major 1 and one of major 2. A draft, a prerelease and a release
// that is no check come too.
const publishedReleases = `[
 {"tag_name":"fitness-check-gofmt-v1.0.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.0.0-linux-arm64.tar","browser_download_url":"https://example.test/a.tar","digest":"sha256:aaa"},
  {"name":"fitness-check-gofmt-v1.0.0-linux-arm64.tar.sha256","browser_download_url":"https://example.test/a.sha","digest":"sha256:bbb"}]},
 {"tag_name":"fitness-check-gofmt-v1.10.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.10.0-linux-arm64.tar","browser_download_url":"https://example.test/b.tar","digest":"sha256:ccc"}]},
 {"tag_name":"fitness-check-gofmt-v2.0.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v2.0.0-linux-arm64.tar","browser_download_url":"https://example.test/c.tar","digest":"sha256:ddd"}]},
 {"tag_name":"fitness-check-gofmt-v1.11.0","draft":true,"prerelease":false,"assets":[]},
 {"tag_name":"fitness-check-vet-v2.0.0-rc1","draft":false,"prerelease":true,"assets":[]},
 {"tag_name":"fitness-check-index","draft":false,"prerelease":false,"assets":[]}
]`

// publishedIndex is the index the releases above make, keeping the newest
// version of each major by number. It leaves out the draft, the prerelease,
// the other release and each .sha256 file.
const publishedIndex = `{"format":1,"repository":"may-journal/fitness-checks","checks":{"gofmt":{
 "latest":"2.0.0","majors":{"1":"1.10.0","2":"2.0.0"},"versions":{
  "1.10.0":{"tag":"fitness-check-gofmt-v1.10.0","artifacts":{"linux-arm64":{"url":"https://example.test/b.tar","sha256":"ccc"}}},
  "2.0.0":{"tag":"fitness-check-gofmt-v2.0.0","artifacts":{"linux-arm64":{"url":"https://example.test/c.tar","sha256":"ddd"}}}}}}}`

// remoteCalls returns the gh calls the stand-in logged for the bare remote.
func remoteCalls(t *testing.T, bare string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(bare, "calls.log"))
	mustDo(t, err)
	return string(b)
}

// gitOut runs git in dir and returns its trimmed standard output.
func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	mustDo(t, err)
	return strings.TrimSpace(string(out))
}

// reuseIndex lists gofmt 1.0.0, 1.1.0 and 2.0.0 with describe files in dir,
// the 1.0.0 one always missing. The 1.1.0 file exists when old is set, and
// fetching a missing file fails the script.
func reuseIndex(t *testing.T, dir string, old bool) string {
	t.Helper()
	asset := func(version string, exists bool) string {
		path := filepath.Join(dir, version+".describe.json")
		if exists {
			mustDo(t, os.WriteFile(path, []byte(`{"files":"\\.(`+version+`)$"}`), 0o644))
		}
		tag := "fitness-check-gofmt-v" + version
		return `{"tag_name":"` + tag + `","draft":false,"prerelease":false,"assets":[{"name":"` + tag + `.describe.json","browser_download_url":"file://` + path + `"}]}`
	}
	file := filepath.Join(dir, "releases.json")
	mustDo(t, os.WriteFile(file, []byte("["+asset("1.0.0", false)+","+asset("1.1.0", old)+","+asset("2.0.0", true)+"]"), 0o644))
	return file
}

func indexFields(t *testing.T, out string) map[string]any {
	t.Helper()
	var index struct {
		Checks map[string]struct {
			Versions map[string]map[string]any `json:"versions"`
		} `json:"checks"`
	}
	mustDo(t, json.Unmarshal([]byte(out), &index))
	got := map[string]any{}
	for version, fields := range index.Checks["gofmt"].Versions {
		got[version] = fields["files"]
	}
	return got
}

// readFolder reads file in folder, or reports why it could not.
func readFolder(t *testing.T, folder, file string) (string, bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(folder, file))
	if err != nil {
		t.Errorf("%s: %v", folder, err)
		return "", false
	}
	return string(data), true
}

// cspell:ignore trimpath buildvcs

// runSize runs this checkout's script from its root with env added, and
// returns the combined output and exit code.
func runSize(t *testing.T, env []string, name string, args ...string) (string, int) {
	t.Helper()
	return execute(t, func() *exec.Cmd {
		cmd := exec.Command(filepath.Join("scripts", name), args...)
		cmd.Dir = ".."
		cmd.Env = append(os.Environ(), env...)
		return cmd
	})
}

// sizeFails writes text as the baseline and expects the size check to exit 1
// with want in its output.
func sizeFails(t *testing.T, env []string, baseline, text, want string) {
	t.Helper()
	mustDo(t, os.WriteFile(baseline, []byte(text), 0o644))
	if out, code := runSize(t, env, "check-binary-sizes.sh"); code != 1 || !strings.Contains(out, want) {
		t.Errorf("baseline %.40q exited %d with %q, want 1 containing %q", text, code, out, want)
	}
}

// sizedReleases lists a release shipped before archives were compressed, and
// one of the next major that ships a tar and its gzip.
const sizedReleases = `[
 {"tag_name":"fitness-check-gofmt-v1.0.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.0.0-linux-arm64.tar","browser_download_url":"https://example.test/a.tar","digest":"sha256:aaa"}]},
 {"tag_name":"fitness-check-gofmt-v2.0.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v2.0.0-linux-arm64.tar","browser_download_url":"https://example.test/b.tar","digest":"sha256:bbb"},
  {"name":"fitness-check-gofmt-v2.0.0-linux-arm64.tar.gz","browser_download_url":"https://example.test/b.tar.gz","digest":"sha256:ccc"}]}
]`

// sizedIndex is the index those releases make. Format 1 keeps url and sha256
// on the tar, with gz only beside the release that has a gzip.
const sizedIndex = `{"format":1,"repository":"may-journal/fitness-checks","checks":{"gofmt":{
 "latest":"2.0.0","majors":{"1":"1.0.0","2":"2.0.0"},"versions":{
  "1.0.0":{"tag":"fitness-check-gofmt-v1.0.0","artifacts":{"linux-arm64":{"url":"https://example.test/a.tar","sha256":"aaa"}}},
  "2.0.0":{"tag":"fitness-check-gofmt-v2.0.0","artifacts":{"linux-arm64":{"url":"https://example.test/b.tar","sha256":"bbb","gz":{"url":"https://example.test/b.tar.gz","sha256":"ccc"}}}}}}}}`
