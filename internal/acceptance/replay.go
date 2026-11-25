package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// Response is one saved answer: when every Match string appears in the call's
// arguments, the stand-in writes Files, prints Stdout and exits.
type Response struct {
	Match  []string          `json:"match"`
	Files  map[string]string `json:"files"`
	Stdout string            `json:"stdout"`
	Stderr string            `json:"stderr"`
	Exit   int               `json:"exit"`
	// Wait, when set, is how many seconds the stand-in hangs first, as a
	// stalled service does.
	Wait int `json:"wait"`
	// Token, when set, also requires the call's GH_TOKEN to equal it.
	Token string `json:"token"`
}

// ServeStandIns lets a test binary run under one of names answer from saved
// responses and exit, so a folder's TestMain calls it first.
func ServeStandIns(names []string) {
	if name := filepath.Base(os.Args[0]); slices.Contains(names, name) {
		os.Exit(replay(name, os.Args[1:]))
	}
}

// replay answers one call from $FITNESS_REPLAY/<name>.json, the first
// response that matches, and logs the call to $FITNESS_REPLAY/<name>.log.
// A call that reads its input from stdin logs that input after its arguments.
func replay(name string, args []string) int {
	dir := os.Getenv("FITNESS_REPLAY")
	call := strings.Join(args, " ")
	if readsStdin(call) {
		sent, _ := io.ReadAll(os.Stdin)
		call += " " + string(sent)
	}
	LogCall(filepath.Join(dir, name+".log"), call)
	var saved []Response
	raw, _ := os.ReadFile(filepath.Join(dir, name+".json"))
	_ = json.Unmarshal(raw, &saved)
	for _, r := range saved {
		if r.answers(call) {
			time.Sleep(time.Duration(r.Wait) * time.Second)
			writeFiles(r.Files)
			fmt.Print(r.Stdout)
			fmt.Fprint(os.Stderr, r.Stderr)
			return r.Exit
		}
	}
	if pr, ok := eventPR(name, call); ok {
		fmt.Println(pr)
		return 0
	}
	fmt.Fprintf(os.Stderr, "%s stand-in: no saved response for: %s\n", name, call)
	return 97
}

// readsStdin reports whether a call reads its input or body from stdin.
func readsStdin(call string) bool {
	return strings.Contains(call, "--input -") || strings.Contains(call, "--body-file -")
}

// eventPR answers a gh read of the event's pull request with the event's
// number, title and body unless a saved response overrides it.
func eventPR(name, call string) (string, bool) {
	pr := savedEventPR()
	if name != "gh" || pr == nil {
		return "", false
	}
	path := fmt.Sprintf("repos/%s/pulls/%d", os.Getenv("GITHUB_REPOSITORY"), pr.Number)
	if !strings.HasPrefix(call, "api ") || !strings.HasSuffix(call, " "+path) {
		return "", false
	}
	out, err := json.Marshal(pr)
	return string(out), err == nil
}

// eventPRFields are the pull request fields a gh read of one PR returns.
type eventPRFields struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// savedEventPR is the pull request the saved Actions event carries, or nil
// when there is no event or it carries none.
func savedEventPR() *eventPRFields {
	var ev struct {
		PR *eventPRFields `json:"pull_request"`
	}
	raw, err := os.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil || json.Unmarshal(raw, &ev) != nil {
		return nil
	}
	return ev.PR
}

// writeFiles writes each file a saved answer leaves behind, such as a test
// report, relative to the folder the stand-in runs in.
func writeFiles(files map[string]string) {
	for name, body := range files {
		_ = os.MkdirAll(filepath.Dir(name), 0o755)
		_ = os.WriteFile(name, []byte(body), 0o644)
	}
}

// answers reports whether r answers call: every Match string appears in it,
// and it runs with r's Token, when set, as GH_TOKEN.
func (r Response) answers(call string) bool {
	return matches(call, r.Match) && (r.Token == "" || os.Getenv("GH_TOKEN") == r.Token)
}

// matches reports whether every want appears in call.
func matches(call string, want []string) bool {
	for _, w := range want {
		if !strings.Contains(call, w) {
			return false
		}
	}
	return true
}

// LogCall appends one call, newline-separated, to path.
func LogCall(path, call string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintln(f, strings.ReplaceAll(call, "\n", `\n`))
}

// Replays is a test's stand-ins: the environment putting them first on PATH and
// the folder holding their saved responses and call logs.
type Replays struct {
	Env []string
	Dir string
}

// StandIn puts a stand-in for each named program first on PATH, answering from
// saved responses, and returns the environment to run fitness with.
func StandIn(t *testing.T, saved map[string][]Response) Replays {
	t.Helper()
	self, err := os.Executable()
	MustDo(t, err)
	dir, bin := t.TempDir(), t.TempDir()
	for name, responses := range saved {
		data, err := json.Marshal(responses)
		MustDo(t, err)
		MustDo(t, os.WriteFile(filepath.Join(dir, name+".json"), data, 0o644))
		MustDo(t, os.Symlink(self, filepath.Join(bin, name)))
	}
	path := bin + string(os.PathListSeparator) + os.Getenv("PATH")
	return Replays{Env: []string{"PATH=" + path, "FITNESS_REPLAY=" + dir}, Dir: dir}
}

// Calls returns every call the named stand-in received, in order.
func (r Replays) Calls(name string) []string {
	raw, _ := os.ReadFile(filepath.Join(r.Dir, name+".log"))
	return strings.FieldsFunc(string(raw), func(c rune) bool { return c == '\n' })
}

// Called reports whether some call to name holds every want.
func (r Replays) Called(name string, want ...string) bool {
	return r.Count(name, want...) > 0
}

// Count returns how many calls to name hold every want.
func (r Replays) Count(name string, want ...string) int {
	n := 0
	for _, c := range r.Calls(name) {
		if matches(c, want) {
			n++
		}
	}
	return n
}

// PullRequestEnv is the GitHub Actions environment of a pull request into main.
var PullRequestEnv = []string{"GITHUB_ACTIONS=true", "GITHUB_BASE_REF=main"}

// PathWithout returns a PATH holding every real tool on the caller's PATH
// except names, since Linux runners keep Go and Node in /usr/bin.
func PathWithout(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, d := range filepath.SplitList(os.Getenv("PATH")) {
		linkTools(dir, d, names)
	}
	return "PATH=" + dir
}

// linkTools links each entry of src into dir, except names and entries an
// earlier PATH folder already supplied.
func linkTools(dir, src string, names []string) {
	entries, _ := os.ReadDir(src)
	for _, e := range entries {
		link := filepath.Join(dir, e.Name())
		if _, err := os.Lstat(link); err == nil || slices.Contains(names, e.Name()) {
			continue
		}
		_ = os.Symlink(filepath.Join(src, e.Name()), link)
	}
}
