// Package requirements proves the markdown-links check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"strings"
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
	t.Run("0001.1.1.1.1 it fails as `path:line: broken relative link: target`", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith(`See [the guide](gone.md "Guide") and [home](README.md).`+"\n"))
		acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: gone.md")
	})
	t.Run("0001.1.1.2.1 it fails resolving the target from that subfolder", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, map[string]string{"docs/guide.md": "# Guide\n\nSee [setup](setup.md#install).\n"})
		acceptance.Sees(t, out, code, 1, "docs/guide.md:3: broken relative link: setup.md#install")
	})
	t.Run("0001.1.1.3.1 it fails naming the image", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith("![Logo](logo.png)\n"))
		acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: logo.png")
	})
	t.Run("0001.1.1.3.2 a reference definition fails alike", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith("[spec]: spec.md\n"))
		acceptance.Sees(t, out, code, 1, "README.md:9: broken relative link: spec.md")
	})
	t.Run("0001.1.1.4.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, map[string]string{"docs/guide.markdown": "# Guide\n\nSee [setup](setup.md).\n"})
		acceptance.Sees(t, out, code, 1, "docs/guide.markdown:3: broken relative link: setup.md")
	})
	t.Run("0001.1.1.4.2 an `.mdx` file fails alike", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, map[string]string{"docs/intro.mdx": "# Intro\n\nSee [the guide](gone.mdx).\n"})
		acceptance.Sees(t, out, code, 1, "docs/intro.mdx:3: broken relative link: gone.mdx")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	exists := func(t *testing.T, body string) (string, int) {
		files := readmeWith(body)
		files["docs/guide.md"] = "# Guide\n\n## Setup\n\nBack to [the readme](../README.md).\n"
		return markdownLinks(t, files)
	}
	t.Run("0001.2.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := exists(t, "See [guide](./docs/guide.md) and [home](docs/../README.md).\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.1.2 a folder passes", func(t *testing.T) {
		t.Parallel()
		out, code := exists(t, "See [docs](docs/).\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.1.3 a heading passes since a `#fragment` is stripped", func(t *testing.T) {
		t.Parallel()
		out, code := exists(t, "See [setup](./docs/guide.md#setup).\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.1.4 a `#fragment` is never validated", func(t *testing.T) {
		t.Parallel()
		out, code := exists(t, "See [nothing](docs/guide.md#no-such-heading).\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	web := func(t *testing.T, link string) (string, int) {
		return markdownLinks(t, readmeWith("See "+link+".\n"))
	}
	t.Run("0001.2.1.2.1 it passes untouched", func(t *testing.T) {
		t.Parallel()
		out, code := web(t, "[site](https://example.com/gone)")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.2.2 a mail link passes alike", func(t *testing.T) {
		t.Parallel()
		out, code := web(t, "[mail](mailto:a@example.com)")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.2.3 a protocol-relative link passes alike", func(t *testing.T) {
		t.Parallel()
		out, code := web(t, "[cdn](//example.com/x)")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.2.4 a same-page link passes alike", func(t *testing.T) {
		t.Parallel()
		out, code := web(t, "[top](#app)")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.3.1 it is ignored", func(t *testing.T) {
		t.Parallel()
		out, code := markdownLinks(t, readmeWith("Write `[text](gone.md)` like this.\n\n```\n[text](gone.md)\n```\n"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.4.1 it is decoded", func(t *testing.T) {
		t.Parallel()
		files := readmeWith("See [notes](docs/release%20notes.md).\n")
		files["docs/release notes.md"] = "# Release Notes\n"
		out, code := markdownLinks(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "%20") {
			t.Errorf("the encoded space must be decoded before the link is looked up:\n%s", out)
		}
	})
	t.Run("0001.2.1.4.2 it resolves", func(t *testing.T) {
		t.Parallel()
		files := readmeWith("See [notes](docs/release%20notes.md).\n")
		files["docs/release notes.md"] = "# Release Notes\n"
		out, code := markdownLinks(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}
