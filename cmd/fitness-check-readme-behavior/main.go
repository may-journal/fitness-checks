package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "readme-behavior", Scope: checkkit.ScopeFiles, Files: match.Base("README.md")},
		Run:      run,
	})
}

// Section headings a README keeps.
const (
	behaviorHeading = "## Behavior"
	exampleHeading  = "## Example"
)

// lone matches a line that is one Markdown link and nothing else, with its
// target in group 1.
var lone = regexp.MustCompile(`^\[[^\]]+\]\(([^)\s]+)\)$`)

// checkReadme matches the README of a check program's folder, with the
// check's name in group 1.
var checkReadme = regexp.MustCompile(`(?:^|/)cmd/fitness-check-([a-z0-9-]+)/README\.md$`)

// describeTimeout is how long a check program may take to describe itself.
const describeTimeout = 20 * time.Second

func run(root string, _ []string) (checkkit.Result, error) {
	var errs, info []string
	files := 0
	for _, rel := range walkfs.InScope(walkfs.FilesByExt(root, "README.md")) {
		if path.Base(rel) != "README.md" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			info = append(info, walkfs.Unreadable(rel, err))
			continue
		}
		files++
		found, note := judge(root, rel, string(raw))
		errs, info = append(errs, found...), append(info, note...)
	}
	return checkkit.Result{Ok: len(errs) == 0, Errors: append([]string{}, errs...), FilesChecked: files, Info: info}, nil
}

// judge lists the errors of one README, and the notes about what it could
// not judge.
func judge(root, rel, text string) (errs, info []string) {
	lines := strings.Split(text, "\n")
	errs = behaviorErrors(root, rel, lines)
	errs = append(errs, stepErrors(rel, lines)...)
	m := checkReadme.FindStringSubmatch(rel)
	if m == nil {
		return errs, nil
	}
	more, note := exampleErrors(rel, m[1], lines)
	return append(errs, more...), note
}

// section is a heading's line number, counted from 1, and the lines under it
// up to the next level 2 heading. It is 0 when the README has no such heading.
func section(lines []string, heading string) (int, []string) {
	for i, line := range lines {
		if strings.TrimRight(line, " \t\r") == heading {
			return i + 1, until(lines[i+1:])
		}
	}
	return 0, nil
}

// until is the lines before the first level 2 heading.
func until(lines []string) []string {
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			return lines[:i]
		}
	}
	return lines
}

// behaviorErrors judges the Behavior section: one link to the README's own
// requirements and nothing else.
func behaviorErrors(root, rel string, lines []string) []string {
	at, body := section(lines, behaviorHeading)
	if at == 0 {
		return nil
	}
	texts := textLines(body)
	if len(texts) == 0 {
		return []string{fmt.Sprintf("%s:%d: the Behavior section holds no link to this folder's docs/requirements", rel, at)}
	}
	if len(texts) > 1 {
		return []string{fmt.Sprintf("%s:%d: the Behavior section holds more than its one link", rel, at+1+texts[1])}
	}
	return linkErrors(root, rel, at+1+texts[0], strings.TrimSpace(body[texts[0]]))
}

// textLines lists the positions of the lines of body that hold text.
func textLines(body []string) []int {
	var texts []int
	for i, line := range body {
		if strings.TrimSpace(line) != "" {
			texts = append(texts, i)
		}
	}
	return texts
}

// linkErrors judges the one line of the Behavior section, at line number n.
func linkErrors(root, rel string, n int, line string) []string {
	m := lone.FindStringSubmatch(line)
	if m == nil {
		return []string{fmt.Sprintf("%s:%d: the Behavior section holds more than one Markdown link and nothing else", rel, n)}
	}
	target, _, _ := strings.Cut(m[1], "#")
	dir := path.Dir(rel)
	joined := path.Join(dir, target)
	own := path.Join(dir, "docs", "requirements")
	if joined != own && !strings.HasPrefix(joined, own+"/") {
		return []string{fmt.Sprintf("%s:%d: the link %s is not this folder's docs/requirements or a doc in it", rel, n, m[1])}
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(joined))); err != nil {
		return []string{fmt.Sprintf("%s:%d: the link %s names nothing that exists", rel, n, m[1])}
	}
	return nil
}

// stepItem matches a list item, bulleted or numbered, that opens with a Gherkin
// keyword in any case, with the keyword in group 1.
var stepItem = regexp.MustCompile(`(?i)^\s*(?:[-*+]|\d+[.)])\s+(given|when|then)\b`)

// fenceOf is the run of three or more backticks or tildes that makes line a
// code fence, or "" when line is no fence.
func fenceOf(line string) string {
	t := strings.TrimSpace(line)
	for _, mark := range []string{"`", "~"} {
		if run := t[:len(t)-len(strings.TrimLeft(t, mark))]; len(run) >= 3 {
			return run
		}
	}
	return ""
}

// fenceState is the fence open after line, given the fence open before it,
// and whether line is a fence. A fence closes on the same character at least
// as long as its opening.
func fenceState(open, line string) (string, bool) {
	f := fenceOf(line)
	switch {
	case f == "":
		return open, false
	case open == "":
		return f, true
	case f[0] == open[0] && len(f) >= len(open):
		return "", true
	}
	return open, false
}

// stepErrors fails each GIVEN, WHEN or THEN list item outside a code fence. A
// README keeps its requirements in docs/requirements, so such an item hides
// them where no check reads them.
func stepErrors(rel string, lines []string) []string {
	var errs []string
	open := ""
	for i, line := range lines {
		var fence bool
		if open, fence = fenceState(open, line); fence {
			continue
		}
		if m := stepItem.FindStringSubmatch(line); open == "" && m != nil {
			errs = append(errs, fmt.Sprintf("%s:%d: the README holds a %s step; move it to this folder's docs/requirements as a numbered acceptance", rel, i+1, strings.ToUpper(m[1])))
		}
	}
	return errs
}

// configKey is one key of a check's --describe config list.
type configKey struct {
	Key     string      `json:"key"`
	Default any         `json:"default"`
	Keys    []configKey `json:"keys"`
}

// described is the part of a check's --describe output the Example judges.
type described struct {
	Blob             bool        `json:"blob"`
	ContextInlineArg string      `json:"contextInlineArg"`
	Config           []configKey `json:"config"`
}

// exampleErrors judges the Example section of the README of check name. The
// notes say what it could not judge.
func exampleErrors(rel, name string, lines []string) (errs, info []string) {
	at, body := section(lines, exampleHeading)
	if at == 0 {
		return []string{rel + ":1: the README has no ## Example section"}, nil
	}
	text := strings.Join(body, "\n")
	var missing []string
	if !strings.Contains(text, "fitness "+name) {
		missing = append(missing, "the command `fitness "+name+"`")
	}
	d, err := describe(name)
	if err != nil {
		note := fmt.Sprintf("%s: cannot read `fitness-check-%s --describe`, so its flags and keys are not judged: %v", rel, name, err)
		return errorLine(rel, at, missing), []string{note}
	}
	missing = append(missing, missingFrom(text, d)...)
	return errorLine(rel, at, missing), nil
}

// errorLine is the error line for the parts the Example at line at misses, or
// none.
func errorLine(rel string, at int, missing []string) []string {
	if len(missing) == 0 {
		return nil
	}
	return []string{fmt.Sprintf("%s:%d: the Example misses %s", rel, at, strings.Join(missing, ", "))}
}

// missingFrom lists the flags and config keys of d that text leaves out.
func missingFrom(text string, d described) []string {
	missing := missingFlags(text, d)
	squeezed := strings.Join(strings.Fields(text), "")
	for _, key := range leaves(d.Config) {
		if !strings.Contains(squeezed, snippet(key)) {
			missing = append(missing, "the key `"+key.Key+"` set to its default in `.fitnessrc.json`")
		}
	}
	return missing
}

// missingFlags lists the flags d declares that text leaves out.
func missingFlags(text string, d described) []string {
	var missing []string
	for _, flag := range declaredFlags(d) {
		if !strings.Contains(text, flag) {
			missing = append(missing, "the flag `"+flag+"`")
		}
	}
	return missing
}

// declaredFlags lists the flags d says the program accepts.
func declaredFlags(d described) []string {
	var flags []string
	if d.Blob {
		flags = append(flags, "--body-file")
	}
	if d.ContextInlineArg != "" {
		flags = append(flags, d.ContextInlineArg)
	}
	return flags
}

// snippet is how a key reads in a squeezed Example: its last name, quoted,
// then its default when it has one.
func snippet(key configKey) string {
	last := key.Key[strings.LastIndex(key.Key, ".")+1:]
	out := `"` + last + `":`
	if key.Default == nil {
		return out
	}
	raw, _ := json.Marshal(key.Default)
	return out + string(raw)
}

// leaves lists the keys that hold no keys of their own, at any depth.
func leaves(keys []configKey) []configKey {
	var out []configKey
	for _, key := range keys {
		if len(key.Keys) == 0 {
			out = append(out, key)
			continue
		}
		out = append(out, leaves(key.Keys)...)
	}
	return out
}

// describe runs the program of check name that sits beside this one and
// decodes what its --describe prints.
func describe(name string) (described, error) {
	self, err := os.Executable()
	if err != nil {
		return described{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), describeTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, filepath.Join(filepath.Dir(self), "fitness-check-"+name), "--describe").Output()
	if err != nil {
		return described{}, err
	}
	var d described
	return d, json.Unmarshal(out, &d)
}
