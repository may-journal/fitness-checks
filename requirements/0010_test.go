package requirements

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

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
		for _, check := range everyCheck0010(t, clone) {
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

// everyCheck0010 names each check folder's check in the repo at root.
func everyCheck0010(t *testing.T, root string) []string {
	t.Helper()
	var every []string
	for _, folder := range checkFolders0008(t, root) {
		every = append(every, strings.TrimPrefix(filepath.Base(folder), "fitness-check-"))
	}
	return every
}

// tagsFor0010 names the tag each check of the repo at root gets at version.
func tagsFor0010(t *testing.T, root, version string) []string {
	t.Helper()
	var tags []string
	for _, check := range everyCheck0010(t, root) {
		tags = append(tags, "fitness-check-"+check+"-v"+version)
	}
	return tags
}

func Test0010_1(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "fix(gofmt): change gofmt", true, "cmd/fitness-check-gofmt/main.go")
	if got := releasePlan(t, clone); !slices.Equal(got, []string{"fitness-check-gofmt-v1.0.1"}) {
		t.Errorf("release-plan = %v, want only gofmt at 1.0.1", got)
	}
}

func Test0010_2(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "feat(spell): change spell", true, "internal/spell/mask.go")
	if got := releasePlan(t, clone); !slices.Equal(got, []string{"fitness-check-cspell-v1.1.0"}) {
		t.Errorf("release-plan = %v, want only cspell, the one check that imports spell, at 1.1.0", got)
	}
}

func Test0010_3(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "build(deps): change go.mod", true, "go.mod")
	if got, want := releasePlan(t, clone), tagsFor0010(t, clone, "1.0.1"); !slices.Equal(got, want) {
		t.Errorf("release-plan = %v, want every check at 1.0.1, %v", got, want)
	}
}

func Test0010_4(t *testing.T) {
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
}

func Test0010_5(t *testing.T) {
	t.Parallel()
	clone := releaseClone(t, "docs(gofmt): change gofmt", false, "cmd/fitness-check-gofmt/main.go")
	if got, want := releasePlan(t, clone), tagsFor0010(t, clone, "1.0.0"); !slices.Equal(got, want) {
		t.Errorf("release-plan = %v, want every check's first release, %v", got, want)
	}
}

func Test0010_6(t *testing.T) {
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
}
