// cspell:ignore trimpath buildvcs
package requirements

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// run0167 runs this checkout's script from its root with env added, and
// returns the combined output and exit code.
func run0167(t *testing.T, env []string, name string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(filepath.Join("scripts", name), args...)
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); ok {
		return string(out), exit.ExitCode()
	}
	mustDo(t, err)
	return string(out), 0
}

func Test0167_1(t *testing.T) {
	t.Parallel()
	out := filepath.Join(t.TempDir(), "release")
	got, code := run0167(t, nil, "release-build.sh", "fitness-check-gofmt-v1.2.3", out)
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
}

func Test0167_2(t *testing.T) {
	t.Parallel()
	baseline := filepath.Join(t.TempDir(), "baseline.tsv")
	env := []string{"SIZE_BASELINE=" + baseline, "SIZE_BUILD_DIR=" + t.TempDir()}
	if out, code := run0167(t, env, "check-binary-sizes.sh", "--update"); code != 0 {
		t.Fatalf("--update exited %d: %s", code, out)
	}
	if out, code := run0167(t, env, "check-binary-sizes.sh"); code != 0 {
		t.Errorf("current sizes exited %d, want 0: %s", code, out)
	}
	b, err := os.ReadFile(baseline)
	mustDo(t, err)
	if rows := strings.Count(string(b), "\nfitness-check-"); rows != len(checkFolders0160(t, "..")) {
		t.Errorf("the baseline lists %d programs, want one for each check folder", rows)
	}
	grown := regexp.MustCompile(`(?m)^fitness-check-gofmt\t\d+$`).ReplaceAllString(string(b), "fitness-check-gofmt\t1000")
	shrunk := regexp.MustCompile(`(?m)^fitness-check-gofmt\t\d+$`).ReplaceAllString(string(b), "fitness-check-gofmt\t900000000")
	stale := string(b) + "fitness-check-gone\t1000\n"
	for want, text := range map[string]string{"SIZE: fitness-check-gofmt is": grown, "under its baseline": shrunk, "SIZE: fitness-check-gone": stale} {
		fails0167(t, env, baseline, text, want)
	}
}

// fails0167 writes text as the baseline and expects the size check to exit 1
// with want in its output.
func fails0167(t *testing.T, env []string, baseline, text, want string) {
	t.Helper()
	mustDo(t, os.WriteFile(baseline, []byte(text), 0o644))
	if out, code := run0167(t, env, "check-binary-sizes.sh"); code != 1 || !strings.Contains(out, want) {
		t.Errorf("baseline %.40q exited %d with %q, want 1 containing %q", text, code, out, want)
	}
}

func Test0167_3(t *testing.T) {
	t.Parallel()
	repo := filepath.Join(t.TempDir(), "repo")
	mustDo(t, os.CopyFS(repo, os.DirFS(filepath.Join("..", "tests", "examples", "happyRepo"))))
	git(t, repo, "init", "-q")
	git(t, repo, "add", "-A")
	git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "docs(readme): add readme")
	out := filepath.Join(t.TempDir(), "measure")
	if got, code := run0167(t, []string{"MEASURE_REPO=" + repo, "MEASURE_RUNS=1"}, "measure-checks.sh", out); code != 0 {
		t.Fatalf("measure-checks.sh exited %d: %s", code, got)
	}
	sizes, err := os.ReadFile(filepath.Join(out, "sizes.tsv"))
	mustDo(t, err)
	timing, err := os.ReadFile(filepath.Join(out, "timing.tsv"))
	mustDo(t, err)
	rows := strings.Split(strings.TrimSpace(string(sizes)), "\n")
	if want := len(checkFolders0160(t, "..")) + 1; len(rows) != want || strings.Count(string(timing), "\n") != 2 || strings.Contains(string(sizes), "\t0\t") {
		t.Errorf("sizes has %d rows and timing %q, want %d rows with sizes, and one timed run\n%s", len(rows), timing, want, sizes)
	}
}

// releases0167 lists a release shipped before archives were compressed, and
// one that ships a tar and its gzip.
const releases0167 = `[
 {"tag_name":"fitness-check-gofmt-v1.0.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.0.0-linux-arm64.tar","browser_download_url":"https://example.test/a.tar","digest":"sha256:aaa"}]},
 {"tag_name":"fitness-check-gofmt-v1.1.0","draft":false,"prerelease":false,"assets":[
  {"name":"fitness-check-gofmt-v1.1.0-linux-arm64.tar","browser_download_url":"https://example.test/b.tar","digest":"sha256:bbb"},
  {"name":"fitness-check-gofmt-v1.1.0-linux-arm64.tar.gz","browser_download_url":"https://example.test/b.tar.gz","digest":"sha256:ccc"}]}
]`

// index0167 is the index those releases make. Format 1 keeps url and sha256
// on the tar, with gz only beside the release that has a gzip.
const index0167 = `{"format":1,"repository":"may-journal/fitness-checks","checks":{"gofmt":{
 "latest":"1.1.0","majors":{"1":"1.1.0"},"versions":{
  "1.0.0":{"tag":"fitness-check-gofmt-v1.0.0","artifacts":{"linux-arm64":{"url":"https://example.test/a.tar","sha256":"aaa"}}},
  "1.1.0":{"tag":"fitness-check-gofmt-v1.1.0","artifacts":{"linux-arm64":{"url":"https://example.test/b.tar","sha256":"bbb","gz":{"url":"https://example.test/b.tar.gz","sha256":"ccc"}}}}}}}}`

func Test0167_4(t *testing.T) {
	t.Parallel()
	file := filepath.Join(t.TempDir(), "releases.json")
	mustDo(t, os.WriteFile(file, []byte(releases0167), 0o644))
	var got, want any
	mustDo(t, json.Unmarshal([]byte(script0163(t, "..", "release-index.sh", file)), &got))
	mustDo(t, json.Unmarshal([]byte(index0167), &want))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("index = %v, want %v", got, want)
	}
}
