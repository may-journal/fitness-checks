// Package requirements proves the markdown-filename-kebab-case check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"maps"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

const check = "markdown-filename-kebab-case"

// docsRepo is a repo whose docs follow kebab-case, with files added.
func docsRepo(files map[string]string) map[string]string {
	repo := map[string]string{
		"docs/api-design.md": "# API\n",
		"docs/user-guide.md": "# Guide\n",
	}
	maps.Copy(repo, files)
	return repo
}

// docs runs `fitness-install -- markdown-filename-kebab-case` on happyRepo with files written
// over it.
func docs(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, check)
}

// exempt is docs whose names no convention judges.
var exempt = map[string]string{
	"CODE_OF_CONDUCT.md":               "# Code of conduct\n",
	"docs/AGENTS.md":                   "# Agents\n",
	".github/pull_request_template.md": "# Pull request\n",
	"docs/issue_template.md":           "# Issue\n",
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 it fails as `path: filename must be kebab-case`", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{"docs/releaseNotes.md": "# Notes\n"}))
		acceptance.Sees(t, out, code, 1, "docs/releaseNotes.md: filename must be kebab-case")
	})
	t.Run("0001.1.1.2.1 each fails as `path: filename must be kebab-case`", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{
			"docs/api_notes.md": "# Notes\n",
			"docs/ApiNotes.md":  "# Notes\n",
		}))
		acceptance.Sees(t, out, code, 1,
			"docs/api_notes.md: filename must be kebab-case",
			"docs/ApiNotes.md: filename must be kebab-case")
	})
	t.Run("0001.1.1.2.2 a name with a space fails alike", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{"docs/api notes.md": "# Notes\n"}))
		acceptance.Sees(t, out, code, 1, "docs/api notes.md: filename must be kebab-case")
	})
	t.Run("0001.1.1.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(nil))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.3.2 digits in a name pass", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{
			"architecture/02-containers.md": "# Containers\n",
			"docs/adr001.md":                "# ADR\n",
		}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.4.1 it passes since only the basename is judged", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{"Guides/Api_Docs/api-design.md": "# API\n"}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.5.1 it fails that doc alone", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", docsRepo(map[string]string{
			"docs/oldNotes.md":    "# Old\n",
			"docs/setup-steps.md": "# Setup\n",
		}))
		acceptance.Write(t, repo, map[string]string{"docs/newNotes.md": "# New\n"})
		acceptance.Git(t, repo, "add", "-A")
		out, code := acceptance.Fitness(t, repo, nil, check)
		acceptance.Sees(t, out, code, 1, "docs/newNotes.md: filename must be kebab-case")
		if strings.Contains(out, "oldNotes.md") {
			t.Errorf("only the staged file may be judged:\n%s", out)
		}
	})
	t.Run("0001.1.1.6.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{"CODE_OF_CONDUCT.md": exempt["CODE_OF_CONDUCT.md"], "docs/AGENTS.md": exempt["docs/AGENTS.md"]}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.6.2 a GitHub template name passes", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, docsRepo(map[string]string{".github/pull_request_template.md": exempt[".github/pull_request_template.md"], "docs/issue_template.md": exempt["docs/issue_template.md"]}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1.1.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, map[string]string{
			"docs/releaseNotes.md": "# Notes\n",
			"docs/installGuide.md": "# Guide\n",
			"docs/api-design.md":   "# API\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, check)
	})
	t.Run("0001.2.1.2.1 it fails the camelCase file", func(t *testing.T) {
		t.Parallel()
		out, code := docs(t, map[string]string{
			"docs/api-design.md":   "# API\n",
			"docs/releaseNotes.md": "# Notes\n",
		})
		acceptance.Sees(t, out, code, 1, "docs/releaseNotes.md: filename must be kebab-case")
	})
}
