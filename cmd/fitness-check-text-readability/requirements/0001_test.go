// Package requirements proves the text-readability check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

// plain is 140 prose words of short, plain sentences.
var plain = strings.Repeat("The runner builds each check and then runs the full suite on every commit. ", 10)

// dense is 141 long words in one sentence.
var dense = strings.Repeat("extraordinary bureaucratic infrastructure considerations regarding administrative complexity ", 20) + "conclude.\n"

// hardToRead holds 104 long words in sentences within their limit.
var hardToRead = strings.Repeat("Configuration interdependencies necessitate comprehensive documentation regarding initialization sequencing, synchronization guarantees, and serialization compatibility. ", 8)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 long words in one endless sentence fail with each piece of guidance", func(t *testing.T) {
		t.Parallel()
		out, code := textReadability(t, map[string]string{"dense.md": dense})
		req.Then(t, "0001.1.1.1.1", "f7956a3e", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "dense.md: readability alarm")
		})
		req.Then(t, "0001.1.1.1.2", "63334658", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "dense.md: readability alarm on 141 prose words", "(grade bands, alarm at 18)", "(alarm at 60)")
		})
		req.Then(t, "0001.1.1.1.3", "214d9b00", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "to fix: split sentences over ~25 words")
		})
		req.Then(t, "0001.1.1.1.4", "a86aec30", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "methodology: go/cmd/fitness-check-text-readability/README.md")
		})
	})
	t.Run("0001.1.2 a plain or short file or one formula in alarm or dense text only in fenced code passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "07b8e753", func(_ req.A) {
			out, code := textReadability(t, map[string]string{"plain.md": plain})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			short := strings.Repeat("extraordinary bureaucratic infrastructure considerations ", 20) + "conclude.\n"
			out, code = textReadability(t, map[string]string{"short.md": short})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			// Half seven-letter words in 21-word sentences: only the long-word score is high.
			lix := strings.Repeat(strings.Repeat("journey to ", 10)+"end. ", 6)
			out, code = textReadability(t, map[string]string{"lix.md": lix}, "--report")
			acceptance.Sees(t, out, code, 0, "LIX 68.6", "All 1 checks passed")
			doc := plain + "\n\n```text\n" + dense + "```\n"
			out, code = textReadability(t, map[string]string{"notation.md": doc})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a dense issue or pull request body fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.1.1", "0276e931", func(_ req.A) {
			out, code := textReadability(t, map[string]string{"body.md": dense}, "--body-file", "body.md")
			acceptance.Sees(t, out, code, 1, "(description): readability alarm on 141 prose words")
		})
	})
	t.Run("0001.2.2 a repo lowers its readability bands", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.2.1", "f528a209", func(_ req.A) {
			out, code := textReadability(t, map[string]string{
				".fitnessrc.json": `{"textReadability": {"maxGrade": 5, "minWords": 10}}` + "\n",
				"plain.md":        plain,
			})
			acceptance.Sees(t, out, code, 1, "plain.md: readability alarm on 140 prose words", "(grade bands, alarm at 5)")
		})
	})
	t.Run("0001.2.3 the report flag prints word counts with scores and leaves the verdict", func(t *testing.T) {
		t.Parallel()
		out, code := textReadability(t, map[string]string{"plain.md": plain}, "--report")
		req.Then(t, "0001.2.1.3.1", "96381708", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "README.md 3 words (below 100 — not judged)", "plain.md 140 words")
		})
		req.Then(t, "0001.2.1.3.2", "ecf42b05", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	req.Then(t, "0001.3.1.1.1", "38296b83", func(_ req.A) {
		comment := strings.TrimRight(hardToRead, " ")
		files := map[string]string{
			"app/main.go":   "package main\n\n// " + comment + "\nfunc main() {}\n",
			"app/App.swift": "import Foundation\n\n/// " + comment + "\nfunc run() {}\n",
			"app/index.ts":  "export {};\n\n// " + comment + "\nexport const x = 1;\n",
		}
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "text-readability")
		acceptance.Sees(t, out, code, 1,
			"app/App.swift:3: readability alarm on 104 prose words",
			"app/index.ts:3: readability alarm on 104 prose words",
			"app/main.go:3: readability alarm on 104 prose words")
	})
}

// textReadability runs text-readability with args on happyRepo with files
// written over it.
func textReadability(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, append([]string{"text-readability"}, args...)...)
}
