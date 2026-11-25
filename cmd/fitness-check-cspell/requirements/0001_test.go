// Package requirements proves the cspell check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 an unknown word fails naming path line and column", func(t *testing.T) {
		t.Parallel()
		out, code := cspellRun(t, map[string]string{"notes.md": "good words\nthe " + cspellTypo + " fox\n"})
		acceptance.Sees(t, out, code, 1, "notes.md:2:5 - Unknown word ("+cspellTypo+")")
	})
	t.Run("0001.1.2 a file over 64 KB is checked in chunks with its directives", func(t *testing.T) {
		t.Parallel()
		// 2,500 filler lines are 65,000 bytes, so the block opens just before
		// cspell's 64 KB chunk boundary and closes after it.
		filler := strings.Repeat("good words in a long file\n", 2500)
		long := "the " + cspellTypo + " fox\n" + filler +
			"cspell:disable\n" + filler[:2600] + "the " + cspellBodyTypo + " fox\ncspell:enable\n" +
			filler + "the " + cspellRepoWord + " fox\n"
		out, code := cspellRun(t, map[string]string{"long.md": long})
		first := "long.md:1:5 - Unknown word (" + cspellTypo + ")"
		last := "long.md:5105:5 - Unknown word (" + cspellRepoWord + ")"
		acceptance.Sees(t, out, code, 1, first, last)
		if strings.Contains(out, cspellBodyTypo) || strings.Index(out, first) > strings.Index(out, last) {
			t.Errorf("only the words outside the block may be named, in line order:\n%s", out)
		}
	})
	t.Run("0001.1.3 the repo's own word list wins and the owner and agent names pass", func(t *testing.T) {
		t.Parallel()
		out, code := cspellRun(t, map[string]string{
			"cspell.json":   `{"words":["fitnessrc","` + cspellRepoWord + `"]}` + "\n",
			cspellInstalled: `{"words":["fitnessrc","` + cspellSharedWord + `"]}` + "\n",
			"notes.md":      cspellRepoWord + " and " + cspellSharedWord + "\n@jasonvanderslice from Vanderbot\n",
		})
		acceptance.Sees(t, out, code, 1, "notes.md:1:11 - Unknown word ("+cspellSharedWord+")")
		if strings.Contains(out, "Unknown word ("+cspellRepoWord+")") {
			t.Errorf("the repo's own word must be accepted:\n%s", out)
		}
		if strings.Contains(out, "Unknown word (jasonvanderslice)") || strings.Contains(out, "Unknown word (Vanderbot)") {
			t.Errorf("the owner and agent names must be accepted:\n%s", out)
		}
	})
	t.Run("0001.1.4 an ignorePaths setting fails before any scan", func(t *testing.T) {
		t.Parallel()
		out, code := cspellRun(t, map[string]string{"cspell.json": `{"ignorePaths":["docs"]}` + "\n"})
		acceptance.Sees(t, out, code, 1, "cspell.json sets ignorePaths; cspell checks every tracked file, so remove ignorePaths and fix the findings instead")
	})
	t.Run("0001.1.5 binary and generated and lock and wrapper files are left out", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitattributes":       "*.gen linguist-generated\n",
			"deps/pkg.gen":         cspellSharedWord + "\n",
			"blob.bin":             cspellBinaryWord + "\x00" + cspellBinaryOther,
			"package-lock.json":    `{"packages":{"node_modules/` + cspellTypo + `":{}}}` + "\n",
			"web/yarn.lock":        cspellRepoWord + "@1.0.0\n",
			"pnpm-lock.yaml":       cspellBodyTypo + ": 1.0.0\n",
			"go.sum":               "example.com/" + cspellBinaryOther + " v1.0.0 h1:abc=\n",
			"App/Package.resolved": `{"identity":"` + cspellBinaryWord + `"}` + "\n",
			"lib/gradle.lockfile":  "com.example:" + cspellTypo + ":1.0=compileClasspath\n",
			"mvnw":                 "url=" + cspellTypo + "\n",
			"app/gradlew.bat":      "set " + cspellRepoWord + "=1\n",
			".mvn/wrapper/MavenWrapperDownloader.java": "class Downloader { String " + cspellRepoWord + "; }\n",
			"gradle/wrapper/gradle-wrapper.properties": "distributionUrl=" + cspellBodyTypo + "\n",
		})
		acceptance.Write(t, repo, map[string]string{"untracked.md": "the " + cspellRepoWord + " fox\n"})
		out, code := acceptance.Fitness(t, repo, nil, "cspell")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.6 a scoped run checks only the changed paths on disk", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{"old.md": "the " + cspellTypo + " fox\n", "gone.md": "clean words\n"})
		acceptance.Write(t, repo, map[string]string{"new.md": "the " + cspellRepoWord + " fox\n"})
		acceptance.Git(t, repo, "add", "new.md")
		acceptance.Git(t, repo, "rm", "-q", "gone.md")
		out, code := acceptance.Fitness(t, repo, nil, "cspell")
		acceptance.Sees(t, out, code, 1, "new.md:1:5 - Unknown word ("+cspellRepoWord+")")
		if strings.Contains(out, cspellTypo) || strings.Contains(out, "gone.md") {
			t.Errorf("only the staged file on disk may be checked:\n%s", out)
		}
	})
	t.Run("0001.1.7 an issue or PR body is spelled as description", func(t *testing.T) {
		t.Parallel()
		out, code := cspellRun(t, map[string]string{"body.md": "a " + cspellBodyTypo + " word\n"}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, "(description):1:3 - Unknown word ("+cspellBodyTypo+")")
	})
	t.Run("0001.1.8 no repo word list falls back to the shared and embedded words", func(t *testing.T) {
		t.Parallel()
		out, code := cspellRun(t, map[string]string{
			cspellInstalled: `{"words":["fitnessrc","` + cspellSharedWord + `"]}` + "\n",
			"notes.md":      cspellSharedWord + " is a shared word\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = cspellRun(t, map[string]string{"notes.md": "mayjournal is a shared word, and addrinfo a language one\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = cspellRun(t, map[string]string{"body.md": "sync `fastlane/Appfile`, `fastlane/Fastfile`, `fastlane/Gymfile`\n"}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = cspellRun(t, map[string]string{"static-preview.yml": "run: touch \"$PAGES/.nojekyll\"\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.9 a forbidden word or phrase fails except in cspell.json", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{
			"cspell.json": `{"words":["fitnessrc"],"flagWords":["green apple"]}` + "\n",
			"notes.md":    "a Green  Apple and green" + "apple\n",
			"body.md":     "no green\napple here\n",
		}
		out, code := cspellRun(t, files)
		acceptance.Sees(t, out, code, 1, "notes.md:1:3 - Forbidden word (green apple)", "body.md:1:4 - Forbidden word (green apple)")
		if strings.Contains(out, "cspell.json:") {
			t.Errorf("the cspell.json listing flagWords must not be flagged:\n%s", out)
		}
		out, code = cspellRun(t, files, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, "(description):1:4 - Forbidden word (green apple)")
	})
	t.Run("0001.1.10 Renovate's hidden debug comment is skipped", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		out, code := acceptance.Piped(t, repo, renovateBody, nil, "--blob", "-", "--checks=cspell")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = acceptance.Piped(t, repo, misspelledNote+renovateBody, nil, "--blob", "-", "--checks=cspell")
		acceptance.Sees(t, out, code, 1, "Unknown word (mispeled)")
	})
}

// Words no dictionary holds, split so cspell never reads them here.
const (
	cspellTypo        = "qu" + "ik"
	cspellRepoWord    = "bor" + "wn"
	cspellSharedWord  = "zz" + "zq" + "qqv"
	cspellBinaryWord  = "xq" + "zzt"
	cspellBinaryOther = "vb" + "nmm"
	cspellBodyTypo    = "mis" + "pe" + "led"
)

// cspellInstalled is where an installed shared package keeps its cspell.json.
const cspellInstalled = "node_modules/@mayjournal/fitness-shared/config/cspell.json"

// cspell:ignore mispeled

// misspelledNote is a note cspell fails.
const misspelledNote = "A mispeled wrd.\n"

// renovateBody is the end of a Renovate PR body, with the hidden debug block
// that holds its state encoded as Renovate encodes it.
var renovateBody = "This PR updates the tooling.\n\n<!--renovate-debug:" +
	base64.StdEncoding.EncodeToString([]byte(`{"createdInVer":"44.222.2","updatedInVer":"44.222.2","targetBranch":"main","labels":[]}`)) +
	"-->\n"

// cspellRun runs `fitness-install -- cspell <args>` on happyRepo with files
// written over it and committed.
func cspellRun(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, append([]string{"cspell"}, args...)...)
}
