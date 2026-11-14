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
	script, err := filepath.Abs(filepath.Join("..", "scripts", name))
	mustDo(t, err)
	cmd := exec.Command(script, args...)
	cmd.Dir = dir
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

func Test0163_1(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "fix(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
	bare := remote0163(t, clone)
	release0163(t, clone, "1.0.0")
	one := release0163(t, clone, "1.2.0")
	two := release0163(t, clone, "2.0.0")
	git(t, clone, "tag", "fitness-check-gofmt-v1", "HEAD~2")
	script0163(t, clone, "release-major.sh")
	got := []string{commitOf0163(bare, "fitness-check-gofmt-v1"), commitOf0163(bare, "fitness-check-gofmt-v2")}
	if !slices.Equal(got, []string{one, two}) {
		t.Errorf("remote major tags = %v, want %v", got, []string{one, two})
	}
	if out := script0163(t, clone, "release-major.sh"); strings.Contains(out, "new tag") || strings.Contains(out, "forced update") {
		t.Errorf("the second run moved a tag: %s", out)
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
