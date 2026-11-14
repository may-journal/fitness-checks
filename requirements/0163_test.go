package requirements

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// remote0163 gives clone an empty bare remote of its own as origin, so a
// script's push never reaches this checkout.
func remote0163(t *testing.T, clone string) string {
	t.Helper()
	bare := filepath.Join(t.TempDir(), "remote.git")
	git(t, clone, "init", "-q", "--bare", bare)
	git(t, bare, "config", "receive.shallowUpdate", "true")
	git(t, clone, "remote", "set-url", "origin", bare)
	return bare
}

// script0163 runs this checkout's script in dir with args and returns its
// standard output.
func script0163(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	return scriptWith0163(t, dir, "", name, args...)
}

// scriptWith0163 is script0163 with gh standing in for GitHub over the bare
// remote, when there is one.
func scriptWith0163(t *testing.T, dir, remote, name string, args ...string) string {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("..", "scripts", name))
	mustDo(t, err)
	cmd := exec.Command(script, args...)
	cmd.Dir = dir
	if remote != "" {
		cmd.Env = append(os.Environ(), fakeGh0163(t, remote)...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", name, err, out)
	}
	return string(out)
}

// commitOf0163 returns the commit a tag names in repo, empty when it is missing.
func commitOf0163(repo, tag string) string {
	out, _ := exec.Command("git", "-C", repo, "rev-parse", "-q", "--verify", "refs/tags/"+tag+"^{commit}").Output()
	return strings.TrimSpace(string(out))
}

// release0163 tags an empty commit as the check's release of version, and
// returns its commit.
func release0163(t *testing.T, clone, version string) string {
	t.Helper()
	git(t, clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "--allow-empty", "-m", "chore: release "+version)
	git(t, clone, "tag", "fitness-check-gofmt-v"+version)
	return commitOf0163(clone, "fitness-check-gofmt-v"+version)
}

// fakeGh0163 puts the stand-in gh first on PATH, applying tag calls to the bare
// remote given in FAKE_REMOTE.
func fakeGh0163(t *testing.T, remote string) []string {
	t.Helper()
	bin := t.TempDir()
	shim, err := filepath.Abs(filepath.Join("..", "tests", "fake-gh-refs.sh"))
	mustDo(t, err)
	mustDo(t, os.Symlink(shim, filepath.Join(bin, "gh")))
	return []string{"PATH=" + bin + string(os.PathListSeparator) + os.Getenv("PATH"), "GITHUB_REPOSITORY=acme/checks", "FAKE_REMOTE=" + remote}
}

func Test0163_1(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
	bare := remote0163(t, clone)
	release0163(t, clone, "1.0.0")
	one := release0163(t, clone, "1.2.0")
	two := release0163(t, clone, "2.0.0")
	git(t, clone, "tag", "fitness-check-gofmt-v1", "HEAD~2")
	git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
	scriptWith0163(t, clone, bare, "release-major.sh")
	got := []string{commitOf0163(bare, "fitness-check-gofmt-v1"), commitOf0163(bare, "fitness-check-gofmt-v2")}
	if !strings.Contains(calls0163(t, bare), "-X POST") {
		t.Errorf("the missing v2 tag was not created with a POST: %s", calls0163(t, bare))
	}
	if !slices.Equal(got, []string{one, two}) {
		t.Errorf("remote major tags = %v, want %v", got, []string{one, two})
	}
	before := calls0163(t, bare)
	scriptWith0163(t, clone, bare, "release-major.sh")
	if after := calls0163(t, bare); after != before {
		t.Errorf("the second run called gh again: %s", after)
	}
}

func Test0163_2(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "fix(gofmt): change gofmt", true, "cmd/fitness-check-gofmt/main.go")
	git(t, clone, "tag", "fitness-check-gofmt-v1", "HEAD")
	if got := releasePlan(t, clone); !slices.Equal(got, []string{"fitness-check-gofmt-v1.0.1"}) {
		t.Errorf("release-plan = %v, want only gofmt at 1.0.1", got)
	}
}

// releases0163 is what GitHub lists: two published versions of one check, a
// draft, a prerelease and a release that is no check.
const releases0163 = `[
 {"tag_name":"fitness-check-gofmt-v1.0.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.0.0-linux-arm64.tar","browser_download_url":"https://example.test/a.tar","digest":"sha256:aaa"},
  {"name":"fitness-check-gofmt-v1.0.0-linux-arm64.tar.sha256","browser_download_url":"https://example.test/a.sha","digest":"sha256:bbb"}]},
 {"tag_name":"fitness-check-gofmt-v1.10.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.10.0-linux-arm64.tar","browser_download_url":"https://example.test/b.tar","digest":"sha256:ccc"}]},
 {"tag_name":"fitness-check-gofmt-v1.11.0","draft":true,"prerelease":false,"assets":[]},
 {"tag_name":"fitness-check-vet-v2.0.0-rc1","draft":false,"prerelease":true,"assets":[]},
 {"tag_name":"fitness-check-index","draft":false,"prerelease":false,"assets":[]}
]`

// index0163 is the index the releases above make: the draft, the prerelease,
// the other release and each .sha256 file are left out.
const index0163 = `{"format":1,"repository":"may-journal/fitness-checks","checks":{"gofmt":{
 "latest":"1.10.0","majors":{"1":"1.10.0"},"versions":{
  "1.0.0":{"tag":"fitness-check-gofmt-v1.0.0","artifacts":{"linux-arm64":{"url":"https://example.test/a.tar","sha256":"aaa"}}},
  "1.10.0":{"tag":"fitness-check-gofmt-v1.10.0","artifacts":{"linux-arm64":{"url":"https://example.test/b.tar","sha256":"ccc"}}}}}}}`

func Test0163_3(t *testing.T) {
	t.Parallel()
	file := filepath.Join(t.TempDir(), "releases.json")
	mustDo(t, os.WriteFile(file, []byte(releases0163), 0o644))
	var got, want any
	mustDo(t, json.Unmarshal([]byte(script0163(t, "..", "release-index.sh", file)), &got))
	mustDo(t, json.Unmarshal([]byte(index0163), &want))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("index = %v, want %v", got, want)
	}
}

// calls0163 returns the gh calls the stand-in logged for the bare remote.
func calls0163(t *testing.T, bare string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(bare, "calls.log"))
	mustDo(t, err)
	return string(b)
}

func Test0163_4(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
	bare := remote0163(t, clone)
	release0163(t, clone, "1.0.0")
	git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
	script, err := filepath.Abs(filepath.Join("..", "scripts", "release-major.sh"))
	mustDo(t, err)
	cmd := exec.Command(script)
	cmd.Dir = clone
	cmd.Env = append(append(os.Environ(), fakeGh0163(t, bare)...), "FAKE_PATCH_ERROR=gh: Forbidden (HTTP 403)")
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "HTTP 403") || commitOf0163(bare, "fitness-check-gofmt-v1") != "" {
		t.Errorf("run = %v with %q, want an error naming HTTP 403 and no tag created", err, out)
	}
}

func Test0163_5(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
	bare := remote0163(t, clone)
	release := release0163(t, clone, "1.0.0")
	git(t, clone, "push", "-q", "origin", "HEAD:refs/heads/main")
	// An orphan commit in the remote is no ancestor of the release.
	tree := git0163(t, bare, "rev-parse", "main^{tree}")
	other := git0163(t, bare, "-c", "user.name=t", "-c", "user.email=t@t", "commit-tree", "-m", "other line", tree)
	git(t, bare, "update-ref", "refs/tags/fitness-check-gofmt-v1", other)
	scriptWith0163(t, clone, bare, "release-major.sh")
	if got := commitOf0163(bare, "fitness-check-gofmt-v1"); got != release {
		t.Errorf("remote major tag = %q, want the release %q, moved by force", got, release)
	}
}

// git0163 runs git in dir and returns its trimmed standard output.
func git0163(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	mustDo(t, err)
	return strings.TrimSpace(string(out))
}
