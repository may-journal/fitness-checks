// Package requirements proves the release-changelog check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo with no manifest passes with zero files whatever version file or changelog it holds", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.1.1", "cc0e4bfd", func(_ req.A) {
			out, code := releaseChangelog(t, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "release-changelog")
		})
	})
	t.Run("0001.1.2 an empty bootstrap manifest passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.2.1", "5858b32f", func(_ req.A) {
			changelogSees(t, map[string]string{".release-please-manifest.json": "{}\n"}, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.3 a malformed manifest or one without a root version fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.3.1", "95096028", func(_ req.A) {
			for _, c := range []struct {
				files map[string]string
				want  string
			}{
				{map[string]string{".release-please-manifest.json": "{\n"}, ".release-please-manifest.json is invalid JSON"},
				{map[string]string{".release-please-manifest.json": "null\n"}, ".release-please-manifest.json must be a JSON object"},
				{map[string]string{".release-please-manifest.json": `{"go":"1.2.3"}`}, `.release-please-manifest.json must contain a root "." release version`},
				{map[string]string{".release-please-manifest.json": `{".":"v1.2.3"}`}, `.release-please-manifest.json root version must be numeric major.minor.patch (got "v1.2.3")`},
			} {
				changelogSees(t, c.files, 1, c.want)
			}
		})
	})
	t.Run("0001.1.4 a version file that differs fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.4.1", "94b0893c", func(_ req.A) {
			files := releaseChangelogRelease(upstreamChangelog)
			files["version.txt"] = "1.2.4\n"
			changelogSees(t, files, 1, `version.txt version "1.2.4" must match .release-please-manifest.json version "1.2.3"`)
		})
	})
	t.Run("0001.1.5 a release with no version file fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.5.1", "0e11e2d4", func(_ req.A) {
			changelogSees(t, map[string]string{".release-please-manifest.json": `{".":"1.2.3"}` + "\n"}, 1, "version.txt is missing; it must match release 1.2.3")
		})
	})
	t.Run("0001.1.6 a matching version passes with or without a changelog", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.6.1", "bfa4ddfa", func(_ req.A) {
			for _, changelog := range []string{upstreamChangelog, ""} {
				repo := acceptance.Example(t, "happyRepo", releaseChangelogRelease(changelog))
				acceptance.Write(t, repo, map[string]string{".release-please-manifest.json": `{".": "1.2.3"}` + "\n"})
				acceptance.Git(t, repo, "add", ".release-please-manifest.json")
				out, code := acceptance.Fitness(t, repo, nil, "release-changelog")
				acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			}
		})
	})
}

// upstreamChangelog is a CHANGELOG.md as an upstream project keeps it:
// undated, typed sections of numbered items under a version heading.
const upstreamChangelog = "# Changelog\n\n## v2.0.0\n\n### Added\n\n1. Catalog-level publishing\n2. Plugin hooks\n\n### Removed\n1. Old lint tools\n"

// releaseChangelog runs `fitness-install -- release-changelog` on happyRepo
// with files written over it.
func releaseChangelog(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "release-changelog")
}

// releaseChangelogRelease is a repo releasing 1.2.3 with changelog as its
// CHANGELOG.md, or none when it is empty.
func releaseChangelogRelease(changelog string) map[string]string {
	return map[string]string{
		".release-please-manifest.json": `{".":"1.2.3"}` + "\n",
		"version.txt":                   "1.2.3\n",
		"CHANGELOG.md":                  changelog,
	}
}

// changelogSees runs release-changelog on happyRepo with files written over it and expects code and want.
func changelogSees(t *testing.T, files map[string]string, wantCode int, want string) {
	t.Helper()
	out, code := releaseChangelog(t, files)
	acceptance.Sees(t, out, code, wantCode, want)
}
