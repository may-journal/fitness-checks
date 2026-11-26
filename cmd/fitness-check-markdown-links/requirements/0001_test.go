// Package requirements proves the markdown-links check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

// readme is happyRepo's README, which body is appended to.
const readme = "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n"

// markdownLinks runs `fitness-install -- markdown-links` on happyRepo with
// files written over it.
func markdownLinks(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "markdown-links")
}

// readmeWith is the happy README with body appended; body starts on line 9.
func readmeWith(body string) map[string]string {
	return map[string]string{"README.md": readme + "\n" + body}
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a link to a missing file fails naming file and line", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith(`See [the guide](gone.md "Guide") and [home](README.md).`+"\n"))
		acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: gone.md")
	})
	t.Run("0001.1.2 a missing sibling is resolved from its folder", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, map[string]string{"docs/guide.md": "# Guide\n\nSee [setup](setup.md#install).\n"})
		acceptance.Sees(t, out, code, 1, "docs/guide.md:3: broken relative link: setup.md#install")
	})
	t.Run("0001.1.3 an image and a reference definition to a missing file fail", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith("![Logo](logo.png)\n\n[spec]: spec.md\n"))
		acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: logo.png", "README.md:11: broken relative link: spec.md")
	})
	t.Run("0001.1.4 markdown and mdx endings count", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, map[string]string{
			"docs/guide.markdown": "# Guide\n\nSee [setup](setup.md).\n",
			"docs/intro.mdx":      "# Intro\n\nSee [the guide](gone.mdx).\n",
		})
		acceptance.Sees(t, out, code, 1, "docs/guide.markdown:3: broken relative link: setup.md", "docs/intro.mdx:3: broken relative link: gone.mdx")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 links to a file or a folder or a heading pass with the fragment stripped", func(t *testing.T) {
		t.Parallel()
		files := readmeWith("See [guide](./docs/guide.md#setup), [docs](docs/), and [home](docs/../README.md).\n")
		files["docs/guide.md"] = "# Guide\n\n## Setup\n\nBack to [the readme](../README.md).\n"
		out, code := markdownLinks(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2 web and mail and same-page links are never touched", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith(
			"See [site](https://example.com/gone), [mail](mailto:a@example.com), [cdn](//example.com/x), and [top](#app).\n"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.3 a broken link inside code is ignored", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith("Write `[text](gone.md)` like this.\n\n```\n[text](gone.md)\n```\n"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.4 a percent-encoded space is decoded then resolves", func(t *testing.T) {
		t.Parallel()
		files := readmeWith("See [notes](docs/release%20notes.md).\n")
		files["docs/release notes.md"] = "# Release Notes\n"
		out, code := markdownLinks(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}
