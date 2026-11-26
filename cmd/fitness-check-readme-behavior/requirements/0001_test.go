// Package requirements proves the readme-behavior check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// link is a Behavior section that names the doc the fixtures write.
const link = "[Requirement 0001](docs/requirements/0001-it-works.md)\n"

// readmeOf is the README text of a check with the given Behavior and Example.
func readmeOf(behavior, example string) string {
	return "# demo\n\nDoes a thing.\n\n## Behavior\n\n" + behavior + "\n## Example\n\n" + example
}

// judged runs the check on a repo that holds readme as the README of check
// name, with a requirement doc beside it.
func judged(t *testing.T, name, readme string) (string, int) {
	t.Helper()
	folder := "cmd/fitness-check-" + name + "/"
	repo := acceptance.Example(t, "happyRepo", map[string]string{
		folder + "README.md":                          readme,
		folder + "docs/requirements/0001-it-works.md": "# 0001 It Works\n",
		"other/README.md":                             "# other\n",
	})
	return acceptance.Fitness(t, repo, nil, "readme-behavior")
}

// command is an Example that shows just the command of check name.
func command(name string) string {
	return "```\nfitness " + name + "\n```\n"
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 one link passes", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf(link, command("gofmt")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2 more than the link fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf(link+"\nIt formats code.\n", command("gofmt")))
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:9: the Behavior section holds more than its one link")
	})
	t.Run("0001.1.3 no link fails naming the heading", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf("", command("gofmt")))
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:5: the Behavior section holds no link")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a doc in its own requirements passes", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf("[Requirements](docs/requirements)\n", command("gofmt")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2 another folder fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf("[Requirements](../fitness-check-gofmt/README.md)\n", command("gofmt")))
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:7: the link ../fitness-check-gofmt/README.md is not this folder's docs/requirements")
	})
	t.Run("0001.2.3 a link to nothing fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf("[Requirement 0002](docs/requirements/0002-gone.md)\n", command("gofmt")))
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:7: the link docs/requirements/0002-gone.md names nothing that exists")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 no Example fails naming the README", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", "# demo\n\nDoes a thing.\n\n## Behavior\n\n"+link)
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:1: the README has no ## Example section")
	})
	t.Run("0001.3.2 an Example without the command fails naming it", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf(link, "Run it.\n"))
		acceptance.Sees(t, out, code, 1, "the Example misses the command `fitness gofmt`")
	})
}

func Test0001_4(t *testing.T) {
	t.Parallel()
	t.Run("0001.4.1 a missing flag fails naming it", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "markdown-no-bold-italic", readmeOf(link, command("markdown-no-bold-italic")))
		acceptance.Sees(t, out, code, 1, "the Example misses the flag `--body-file`")
		out, code = judged(t, "commit-attribution", readmeOf(link, command("commit-attribution")))
		acceptance.Sees(t, out, code, 1, "the Example misses the flag `--message`")
	})
	t.Run("0001.4.2 a shown flag passes", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "markdown-no-bold-italic", readmeOf(link, command("markdown-no-bold-italic")+"\n```\nfitness markdown-no-bold-italic --body-file body.md\n```\n"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_5(t *testing.T) {
	t.Parallel()
	const keys = "```json\n{\n  \"textReadability\": { \"maxGrade\": 18, \"maxLix\": 60, \"minWords\": 100 }\n}\n```\n"
	flagged := command("text-readability") + "\n```\nfitness text-readability --body-file body.md\n```\n"
	t.Run("0001.5.1 a missing key fails naming it", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "text-readability", readmeOf(link, flagged+"\n```json\n{ \"textReadability\": { \"maxGrade\": 18, \"maxLix\": 60 } }\n```\n"))
		acceptance.Sees(t, out, code, 1, "the Example misses the key `textReadability.minWords` set to its default")
	})
	t.Run("0001.5.2 a snippet with every default passes", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "text-readability", readmeOf(link, flagged+"\n"+keys))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_6(t *testing.T) {
	t.Parallel()
	t.Run("0001.6.1 a step item in any case fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf(link, command("gofmt"))+"\n## Notes\n\n- When it runs\n")
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:17: the README holds a WHEN step")
		out, code = judged(t, "gofmt", readmeOf(link, command("gofmt"))+"\n## Notes\n\n1. then it passes\n")
		acceptance.Sees(t, out, code, 1, "cmd/fitness-check-gofmt/README.md:17: the README holds a THEN step")
	})
	t.Run("0001.6.2 the words in a fence or a plain sentence pass", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, "gofmt", readmeOf(link, command("gofmt"))+"\n## Notes\n\nWhen it runs it formats.\n\n```\n- When it runs\n```\n\n~~~\n- When it runs\n~~~\n\n````\n```\n- When it runs\n```\n- Then it passes\n````\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}
