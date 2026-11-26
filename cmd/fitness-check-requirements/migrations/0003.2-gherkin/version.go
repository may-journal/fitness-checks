package gherkin

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// entryVersion matches a changelog entry's version heading.
var entryVersion = regexp.MustCompile(`(?m)^### (\d{4}\.\d{2}\.\d{2}\.\d{4})\s*$`)

// runVersions holds removedVersion's answer for each root, so every doc of
// one run takes the same version.
var runVersions = map[string]string{}

// removedVersion is the version of up's Removed lines: the newest changelog
// entry the change adds, else the time of the run.
func removedVersion(root string) string {
	if v, ok := runVersions[root]; ok {
		return v
	}
	v := time.Now().UTC().Format(requirement.VersionLayout)
	if versions := addedVersions(root); len(versions) > 0 {
		v = slices.Max(versions)
	}
	runVersions[root] = v
	return v
}

// addedVersions lists the versions of root's new changelog entries, which
// belong to the change being made.
func addedVersions(root string) []string {
	var versions []string
	for _, f := range addedEntries(root) {
		b, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			continue
		}
		for _, m := range entryVersion.FindAllStringSubmatch(string(b), -1) {
			versions = append(versions, m[1])
		}
	}
	return versions
}

// addedEntries lists the changelog files git sees as untracked or staged as
// added.
func addedEntries(root string) []string {
	untracked, _ := exec.Command("git", "-C", root, "ls-files", "-z", "--others", "--exclude-standard", "--", "changelog").Output()
	staged, _ := exec.Command("git", "-C", root, "diff", "--cached", "--relative", "--name-only", "--diff-filter=A", "-z", "--", "changelog").Output()
	return slices.DeleteFunc(strings.Split(string(untracked)+"\x00"+string(staged), "\x00"), func(f string) bool { return f == "" })
}
