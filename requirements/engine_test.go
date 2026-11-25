// Helpers the moved requirement tests share with the engine's own tests,
// which stay in fitness-runner.
package requirements

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// allPassedWith is the summary of an all-passing run with extra checks beyond
// the bundled catalog (negative when some are off).
func allPassedWith(extra int) string {
	return fmt.Sprintf("All %d checks passed", len(strings.Fields(catalogRows))+extra)
}

// branchAnswer is the saved result of the filtered open-PR request.
func branchAnswer(t *testing.T, repo, title, body string) response {
	t.Helper()
	base, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
	mustDo(t, err)
	data, err := json.Marshal(map[string]any{
		"number": 1, "title": title, "body": body, "head": "work", "source": exampleRepo,
		"base": "main", "baseSHA": strings.TrimSpace(string(base)), "url": "https://github.com/" + exampleRepo + "/pull/1", "draft": true,
	})
	mustDo(t, err)
	return response{Match: []string{"pulls?state=open&head=" + url.QueryEscape("may-journal:work")}, Stdout: string(data) + "\n"}
}

// branchRepo supplies a local branch and fetched base without network writes.
func branchRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	repo := example(t, "happyRepo", files)
	git(t, repo, "branch", "-M", "work")
	git(t, repo, "remote", "add", "origin", "https://github.com/"+exampleRepo+".git")
	git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	return repo
}

// epicBody is an Epic that follows the template, listing items.
func epicBody(items string) string {
	return "> Every fitness run finishes in a fraction of the time.\n\n## Background\n\nThis epic is a test fixture.\n\n## What needs to happen\n\n" + items + "\n"
}

// hooked copies happyRepo with files over it, runs `fitness init` with the
// installer placed by installerEnv, and returns the repo and that environment.
func hooked(t *testing.T, files map[string]string) (string, []string) {
	t.Helper()
	repo := example(t, "happyRepo", files)
	env := installerEnv(t)
	got, code := fitness(t, repo, env, "init")
	sees(t, got, code, 0, "fitness init: installed hooks")
	return repo, env
}

// hookEvent pipes one Claude Code hook event, as JSON, into the installed
// `fitness claude-hook` run from repo with a cache folder of its own. It
// asserts stdout holds only the JSON Claude Code parses, and returns the
// output and exit code.
func hookEvent(t *testing.T, repo string, event map[string]any) (string, int) {
	t.Helper()
	return hookEventIn(t, repo, t.TempDir(), event)
}

// houseReadme is a folder check's README: a title, the rule, and the fix.
const houseReadme = "# House rule\n\nKeeps the house in order. Each room stays tidy.\n\n## Fix\n\nTidy the room the error names.\n"

// installerEnv puts the installer where a user's setup does, without fitness on
// the PATH, and returns the environment the hooks need.
func installerEnv(t *testing.T) []string {
	t.Helper()
	dir := t.TempDir()
	mustDo(t, os.Symlink(installer(t), filepath.Join(dir, "fitness-install")))
	lookups := rlReleases(t, http.StatusOK, rlCandidate(t))
	return slices.Concat([]string{"FITNESS_INSTALL_DIR=" + dir}, lookups.env, fiPrivateCache(t))
}

// ghIssue is what gh prints for issue n read through fitness's filter:
// its body and label names.
func ghIssue(t *testing.T, n int, body string, labels ...string) response {
	t.Helper()
	out, err := json.Marshal(map[string]any{"body": body, "labels": append([]string{}, labels...)})
	mustDo(t, err)
	return response{Match: []string{"repos/" + exampleRepo + "/issues/" + itoa(n)}, Stdout: string(out) + "\n"}
}

func itoa(n int) string { return strconv.Itoa(n) }

// judgeBody runs pr-check on a pull_request event with title and body,
// GitHub answering from issues and the PR leaving the changelog alone.
func judgeBody(t *testing.T, title, body string, issues ...response) (string, int) {
	t.Helper()
	out, code, _ := rerunBody(t, title, body, issues...)
	return out, code
}

// judgeKind runs issue-check on issue 21 with body, as an issue edit
// triggers it, in happyRepo with files written over it.
func judgeKind(t *testing.T, body string, files map[string]string, saved []response) (string, int, replays) {
	t.Helper()
	issue := map[string]any{"number": 21, "title": "t", "body": body}
	return ghRun(t, "issue-check", "issues", bodyEdit(issue), files, saved)
}

// kindGitHub is the saved GitHub for issue 21 of kind, with no verdict yet,
// carrying labels.
func kindGitHub(kind string, labels ...string) []response {
	return issueGitHub(kindIssue, "fitness:issue-check:"+kind+":", "", labels...)
}

// passing is a local check's verdict that judged one file and found nothing.
const passing = `{"ok": true, "errors": [], "filesChecked": 1}`

// planTemplate is a Plan that follows the template, with one task.
func planTemplate(task string) string {
	return "> One clear pitch for this Plan.\n\n## Background\n\nThis Plan is a test fixture.\n\n## What needs to happen\n\n" + task + "\n"
}

// prClosing is a PR description that follows the template, links the entry
// stamped entryStamp, and closes issue n.
func prClosing(n int) string {
	return prBody(changelogLink("CHANGELOG.md#"+strings.ReplaceAll(entryStamp, ".", "")), n)
}

// readme is happyRepo's README, for Givens that add to it.
const readme = "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n"

// remote gives repo an empty bare origin to push to, so no test needs GitHub.
func remote(t *testing.T, repo string) {
	t.Helper()
	bare := t.TempDir()
	git(t, bare, "init", "-q", "--bare")
	git(t, repo, "remote", "add", "origin", bare)
}

// exampleRepo is the repository the Actions event comes from.
const exampleRepo = "may-journal/example"

// shellCheck is a repo with one shell check, name, and a README.
// The check answers --describe with describe, else ignores the handshake, and prints verdict; an empty verdict hangs.
func shellCheck(name, describe, verdict string) map[string]string {
	script := verdictScript(verdict)
	if describe != "" {
		script = strings.Replace(script, "#!/bin/sh\n", "#!/bin/sh\n[ \"$1\" = --describe ] && { echo '"+describe+"'; exit 0; }\n", 1)
	}
	return map[string]string{"fitness/" + name + "/README.md": houseReadme, "fitness/" + name + "/check": script}
}

// silentPass asserts the hook let the event pass without a word.
func silentPass(t *testing.T, out string, code int) {
	t.Helper()
	if code != 0 || out != "" {
		t.Errorf("want a silent pass, got exit %d:\n%s", code, out)
	}
}

// timedCheck is one slow check that logs its step's start and end around a
// pause, answering --describe with its stage.
func timedCheck(name string, args []string) int {
	check := strings.TrimPrefix(name, "fitness-check-")
	dir := os.Getenv("FITNESS_TIMED")
	step, pause, answer := "run", time.Second, `{"ok":true,"errors":[],"filesChecked":0}`
	if slices.Contains(args, "--describe") {
		stage, err := os.ReadFile(filepath.Join(dir, check+".stage"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "slow check has no stage:", err)
			return 3
		}
		step, pause, answer = "describe", 300*time.Millisecond, `{"name":"`+check+`","stage":"`+string(stage)+`"}`
	}
	logCall(filepath.Join(dir, check+".log"), step+" start "+strconv.FormatInt(time.Now().UnixNano(), 10))
	time.Sleep(pause)
	logCall(filepath.Join(dir, check+".log"), step+" end "+strconv.FormatInt(time.Now().UnixNano(), 10))
	fmt.Println(answer)
	return 0
}

// timedPrefix names the slow checks Test0001_13 installs. Run under such a
// name, this test binary is that check, and logs when each step starts and
// ends.
const timedPrefix = "fitness-check-slow-"

// user runs git in repo as a person at the terminal does, with hooks, returning
// its output and exit code.
func user(t *testing.T, repo string, env []string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t"}, args...)...)
	cmd.Dir, cmd.Env = repo, append(userEnv(), env...)
	out, _ := cmd.CombinedOutput()
	return ansi.ReplaceAllString(string(out), ""), cmd.ProcessState.ExitCode()
}

// verdictScript is a check program that prints verdict, or hangs when it
// is empty.
func verdictScript(verdict string) string {
	if verdict == "" {
		return "#!/bin/sh\nexec /bin/sleep 30\n"
	}
	return "#!/bin/sh\necho '" + verdict + "'\n"
}

// bodyEdit is the event GitHub delivers when issue's body is edited: its
// changes name the body's earlier text.
func bodyEdit(issue map[string]any) map[string]any {
	return map[string]any{"action": "edited", "issue": issue, "changes": map[string]any{"body": map[string]any{"from": "An earlier body."}}}
}

// fiPrivateCache makes a cache holding only the verified release archive,
// copied from the cache the environment already filled, so nothing downloads.
func fiPrivateCache(t *testing.T) []string {
	t.Helper()
	out, code := fiRun(t, t.TempDir(), nil, "--install-only")
	if code != 0 {
		t.Fatalf("installing fitness:\n%s", out)
	}
	key := filepath.Dir(fiLastLine(out))
	cache := t.TempDir()
	copied := filepath.Join(cache, filepath.Base(key))
	mustDo(t, os.MkdirAll(copied, 0o700))
	mustDo(t, os.WriteFile(filepath.Join(copied, "archive.tar.gz"), []byte(fiFile(t, filepath.Join(key, "archive.tar.gz"))), 0o600))
	return []string{"FITNESS_CACHE_DIR=" + cache}
}

// ghRun runs a GitHub command as its Actions step does, in happyRepo with files
// over it, on GitHub's event payload and saved replies.
func ghRun(t *testing.T, command, event string, payload map[string]any, files map[string]string, saved []response) (string, int, replays) {
	t.Helper()
	repo := example(t, "happyRepo", files)
	path := filepath.Join(t.TempDir(), "event.json")
	data, err := json.Marshal(payload)
	mustDo(t, err)
	mustDo(t, os.WriteFile(path, data, 0o644))
	rp := standIn(t, map[string][]response{"gh": saved})
	env := append(rp.env, "GITHUB_EVENT_NAME="+event, "GITHUB_EVENT_PATH="+path, "GITHUB_REPOSITORY=may-journal/example")
	out, code := fitness(t, repo, env, command)
	return out, code, rp
}

// hookEventIn is hookEvent with the fitness cache folder cache, where the
// hook writes its log.
func hookEventIn(t *testing.T, repo, cache string, event map[string]any) (string, int) {
	t.Helper()
	out, code, _ := timedHook(t, filepath.Join(installOnly(t), "fitness"), repo, cache, event)
	return out, code
}

// issueGitHub is the saved GitHub holding issue n's labels and verdict comments,
// found by a query naming marker, where every write succeeds.
func issueGitHub(n, marker, verdicts string, labels ...string) []response {
	issue, _ := json.Marshal(map[string]any{"body": "", "labels": labels})
	return []response{
		{Match: []string{marker}, Stdout: verdicts},
		{Match: []string{"minimizeComment"}},
		{Match: []string{"labels: [.labels[].name]", "issues/" + n}, Stdout: string(issue) + "\n"},
		{Match: []string{"-X POST"}},
		{Match: []string{"-X PATCH", "issues/comments/"}},
		{Match: []string{"-X DELETE", "issues/" + n + "/labels/"}},
		{Match: []string{"--silent", "/labels/fitness"}},
	}
}

// kindIssue is the issue number every 0082 event and saved response names.
const kindIssue = "21"

// changelogLink is a `## Changelog` item linking target on the PR's branch.
func changelogLink(target string) string {
	return "- [CHANGELOG entry](https://github.com/" + exampleRepo + "/blob/feature/" + target + ")"
}

// prBody is a PR description with changelog as its `## Changelog`
// section, closing issue n.
func prBody(changelog string, n int) string {
	return fmt.Sprintf("> One clear summary of the change.\n\n## Background\n\nThis change is a test fixture.\n\n## Changelog\n\n%s\n\nCloses #%d\n", changelog, n)
}

// rerunBody runs pr-check on a saved pull_request event carrying title and
// body, as "Re-run job" replays it, GitHub answering from saved first.
func rerunBody(t *testing.T, title, body string, saved ...response) (string, int, replays) {
	t.Helper()
	rp := standIn(t, map[string][]response{"gh": append(saved, noChangelog)})
	pr := map[string]any{"number": 1, "title": title, "body": body}
	env := actionsEvent(t, "pull_request", map[string]any{"action": "edited", "pull_request": pr})
	out, code := fitness(t, example(t, "happyRepo", nil), append(rp.env, env...), "pr-check")
	return out, code, rp
}

// rlCandidate is the release version the candidate installer carries.
func rlCandidate(t *testing.T) string {
	t.Helper()
	out, _ := fiRun(t, t.TempDir(), nil, "--help")
	match := rlDefault.FindStringSubmatch(out)
	if match == nil {
		t.Fatalf("no default version in --help:\n%s", out)
	}
	return match[1]
}

// rlReleases serves lookups answering status with release version, and
// returns them with the environment pointing the installer at them. The
// environment sets no token variable.
func rlReleases(t *testing.T, status int, version string) *rlLookups {
	t.Helper()
	lookups := &rlLookups{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lookups.mu.Lock()
		lookups.auth = append(lookups.auth, r.Header.Get("Authorization"))
		lookups.mu.Unlock()
		w.WriteHeader(status)
		if status != http.StatusOK {
			fmt.Fprint(w, `{"message":"API rate limit exceeded"}`)
			return
		}
		fmt.Fprintf(w, `{"tag_name":"v%s"}`, version)
	}))
	t.Cleanup(srv.Close)
	lookups.env = []string{"GITHUB_API_URL=" + srv.URL, "GH_TOKEN=", "GITHUB_TOKEN="}
	return lookups
}

// entryStamp is the dated changelog entry a fixture PR adds and links.
const entryStamp = "2026.10.10.0900"

// actionsEvent saves an Actions event payload and returns the environment
// a workflow step sees: event name, payload path, and repository.
func actionsEvent(t *testing.T, event string, payload any) []string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "event.json")
	data, err := json.Marshal(payload)
	mustDo(t, err)
	mustDo(t, os.WriteFile(path, data, 0o644))
	return []string{"GITHUB_EVENT_NAME=" + event, "GITHUB_EVENT_PATH=" + path, "GITHUB_REPOSITORY=" + exampleRepo}
}

// fiFile returns the text of a file the installer wrote.
func fiFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	mustDo(t, err)
	return string(data)
}

// fiLastLine is the last printed line, where install-only names its folder.
func fiLastLine(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

// fiRun runs the candidate `fitness-install <args>` in dir with env added
// and returns its output and exit code.
func fiRun(t *testing.T, dir string, env []string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(installer(t), args...)
	cmd.Dir, cmd.Env = dir, append(userEnv(), env...)
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		t.Fatalf("running fitness-install: %v", err)
	}
	return ansi.ReplaceAllString(string(out), ""), cmd.ProcessState.ExitCode()
}

// installOnly runs `fitness-install --install-only`, which prints the
// directory holding the runner and every check, and returns that directory.
func installOnly(t *testing.T) string {
	t.Helper()
	out, err := exec.Command(installer(t), "--install-only").Output()
	mustDo(t, err)
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

// noChangelog is GitHub listing no CHANGELOG.md change in a PR.
var noChangelog = response{Match: []string{"/files"}}

// rlDefault matches the version --help names as the installer's default.
var rlDefault = regexp.MustCompile(`exact release version, or explicit latest \(default "([^"]+)"\)`)

// rlLookups is a saved releases API answering latest release lookups with
// one status, naming a release on 200, and keeping each Authorization header.
type rlLookups struct {
	env  []string
	mu   sync.Mutex
	auth []string
}

// timedHook runs the hook event through runner and also returns how long
// the hook alone took.
func timedHook(t *testing.T, runner, repo, cache string, event map[string]any) (string, int, time.Duration) {
	t.Helper()
	event["cwd"] = repo
	raw, err := json.Marshal(event)
	mustDo(t, err)
	cmd := exec.Command(runner, "claude-hook")
	cmd.Dir, cmd.Env, cmd.Stdin = repo, hookEnv(t, cache), strings.NewReader(string(raw))
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	start := time.Now()
	_ = cmd.Run()
	took := time.Since(start)
	onlyJSON(t, stdout.String())
	return stdout.String() + stderr.String(), cmd.ProcessState.ExitCode(), took
}

// hookEnv is the user's environment with the fitness cache folder cache.
func hookEnv(t *testing.T, cache string) []string {
	t.Helper()
	return append(userEnv(), "FITNESS_CACHE_DIR="+cache)
}

// onlyJSON asserts out is empty or one JSON object and nothing else.
func onlyJSON(t *testing.T, out string) {
	t.Helper()
	var verdict map[string]any
	if out != "" && json.Unmarshal([]byte(out), &verdict) != nil {
		t.Errorf("stdout must hold only the hook's JSON:\n%s", out)
	}
}
