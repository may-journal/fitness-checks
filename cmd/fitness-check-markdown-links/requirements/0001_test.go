// Package requirements proves the markdown-links check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
		req.Then(t, "0001.1.1.1.1", "0a603ba3", func(_ req.A) {
			out, code := markdownLinks(t, readmeWith(`See [the guide](gone.md "Guide") and [home](README.md).`+"\n"))
			acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: gone.md")
		})
	})
	t.Run("0001.1.2 a missing sibling is resolved from its folder", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "2322d1a3", func(_ req.A) {
			out, code := markdownLinks(t, map[string]string{"docs/guide.md": "# Guide\n\nSee [setup](setup.md#install).\n"})
			acceptance.Sees(t, out, code, 1, "docs/guide.md:3: broken relative link: setup.md#install")
		})
	})
	t.Run("0001.1.3 an image and a reference definition to a missing file fail", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith("![Logo](logo.png)\n\n[spec]: spec.md\n"))
		req.Then(t, "0001.1.1.3.1", "9dd8fec6", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: logo.png")
		})
		req.Then(t, "0001.1.1.3.2", "ef59dc24", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "README.md:11: broken relative link: spec.md")
		})
	})
	t.Run("0001.1.4 markdown and mdx endings count", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, map[string]string{
			"docs/guide.markdown": "# Guide\n\nSee [setup](setup.md).\n",
			"docs/intro.mdx":      "# Intro\n\nSee [the guide](gone.mdx).\n",
		})
		req.Then(t, "0001.1.1.4.1", "07f1343a", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/guide.markdown:3: broken relative link: setup.md")
		})
		req.Then(t, "0001.1.1.4.2", "e42c3dea", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/intro.mdx:3: broken relative link: gone.mdx")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 links to a file or a folder or a heading pass with the fragment stripped", func(t *testing.T) {
		t.Parallel()
		linking := func(body string) (string, int) {
			files := readmeWith(body)
			files["docs/guide.md"] = "# Guide\n\n## Setup\n\nBack to [the readme](../README.md).\n"
			return markdownLinks(t, files)
		}
		file, fileCode := linking("See [guide](docs/guide.md) and [home](docs/../README.md).\n")
		folder, folderCode := linking("See [docs](docs/).\n")
		heading, headingCode := linking("See [guide](./docs/guide.md#setup).\n")
		fragment, fragmentCode := linking("See [guide](docs/guide.md#no-such-heading).\n")
		req.Then(t, "0001.2.1.1.1", "07b8e753", func(a req.A) {
			acceptance.Sees(t, file, fileCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.1.2", "3a2faab4", func(a req.A) {
			acceptance.Sees(t, folder, folderCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.1.3", "d506ade1", func(a req.A) {
			acceptance.Sees(t, heading, headingCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.1.4", "7956e904", func(a req.A) {
			acceptance.Sees(t, fragment, fragmentCode, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.2.2 web and mail and same-page links are never touched", func(t *testing.T) {
		t.Parallel()
		web, webCode := markdownLinks(t, readmeWith("See [site](https://example.com/gone).\n"))
		mail, mailCode := markdownLinks(t, readmeWith("See [mail](mailto:a@example.com).\n"))
		cdn, cdnCode := markdownLinks(t, readmeWith("See [cdn](//example.com/x).\n"))
		top, topCode := markdownLinks(t, readmeWith("See [top](#app).\n"))
		req.Then(t, "0001.2.1.2.1", "17e1456d", func(a req.A) {
			acceptance.Sees(t, web, webCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.2.2", "c016231c", func(a req.A) {
			acceptance.Sees(t, mail, mailCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.2.3", "9cf417db", func(a req.A) {
			acceptance.Sees(t, cdn, cdnCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.2.4", "9ea42560", func(a req.A) {
			acceptance.Sees(t, top, topCode, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.2.3 a broken link inside code is ignored", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.3.1", "d958af26", func(_ req.A) {
			out, code := markdownLinks(t, readmeWith("Write `[text](gone.md)` like this.\n\n```\n[text](gone.md)\n```\n"))
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.2.4 a percent-encoded space is decoded then resolves", func(t *testing.T) {
		t.Parallel()
		files := readmeWith("See [notes](docs/release%20notes.md).\n")
		files["docs/release notes.md"] = "# Release Notes\n"
		out, code := markdownLinks(t, files)
		req.Then(t, "0001.2.1.4.1", "47fab0ca", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.1.4.2", "10b8231e", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
}
