package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{Name: "workflow-scripts", Scope: checkkit.ScopeFiles},
		Run:      run,
	})
}

// ruleScript is the rule every finding names.
const ruleScript = "inline-script"

// runRe finds a `run:` key, with its list dash, indent and value.
var runRe = regexp.MustCompile(`^(\s*)(- +)?run:[ \t]*(.*?)\s*$`)

// blockRe is a block scalar header such as `|`, `>-` or `|2`.
var blockRe = regexp.MustCompile(`^[|>][+-]?\d?[+-]?(\s+#.*)?$`)

// maskRe finds what hides an operator: an expression, a quoted string, `\;` or `<<<`.
var maskRe = regexp.MustCompile(`\$\{\{.*?\}\}|\$\(\(.*?\)\)|'[^']*'|"(?:[^"\\]|\\.)*"|\\;|<<<`)

// controlRe finds a line that opens shell control flow.
var controlRe = regexp.MustCompile(`^(if|for|while|until|case|select|function)\b`)

// operators are the shell operators one trivial line never holds.
var operators = []struct{ text, name string }{
	{"&&", "an && chain"}, {"||", "an || chain"}, {"<<", "a heredoc"},
	{";", "a ; chain"}, {"|", "a pipe"},
}

func run(root string, _ []string) (checkkit.Result, error) {
	files := scriptFiles(root)
	var found []checkkit.Finding
	var skipped []string
	for _, file := range files {
		got, err := judgeFile(root, file)
		found = append(found, got...)
		if err != nil {
			skipped = append(skipped, walkfs.Unreadable(file, err))
		}
	}
	res := checkkit.Found(len(files)-len(skipped), render, found...)
	res.Info = skipped
	return res, nil
}

// render is a finding's error line: where, what is wrong, and the fix.
func render(f checkkit.Finding) string {
	return fmt.Sprintf("%s:%d: %s; %s", f.File, f.Line, f.Message, f.Fix)
}

// scriptFiles is the in-scope workflow files and action.yml files.
func scriptFiles(root string) []string {
	var out []string
	for _, f := range walkfs.InScope(walkfs.FilesByExt(root, ".yml", ".yaml")) {
		if isWorkflow(f) || isAction(f) {
			out = append(out, f)
		}
	}
	return out
}

// isWorkflow reports whether f is directly under .github/workflows.
func isWorkflow(f string) bool {
	rest, ok := strings.CutPrefix(f, ".github/workflows/")
	return ok && !strings.Contains(rest, "/")
}

// isAction reports whether f is an action.yml or action.yaml file.
func isAction(f string) bool {
	base := filepath.Base(f)
	return base == "action.yml" || base == "action.yaml"
}

// blockKeyRe is a key whose value is a block scalar, so its body is no step.
var blockKeyRe = regexp.MustCompile(`^(\s*)(?:- +)?[\w.-]+:[ \t]*[|>][+-]?\d?[+-]?(\s+#.*)?$`)

// mappingRe is a line that opens a mapping entry, such as `shell: bash`.
var mappingRe = regexp.MustCompile(`^[\w.-]+:( |$)`)

// judgeFile finds each run step in one file that is no trivial command.
func judgeFile(root, file string) ([]checkkit.Finding, error) {
	data, err := os.ReadFile(filepath.Join(root, file))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var out []checkkit.Finding
	for i := 0; i < len(lines); i++ {
		why, skip := whyHeavy(lines, i)
		if why != "" {
			out = append(out, finding(file, i+1, why))
		}
		i += skip
	}
	return out, nil
}

// finding is the finding for a heavy run step at line.
func finding(file string, line int, why string) checkkit.Finding {
	return checkkit.Finding{File: file, Line: line, Rule: ruleScript,
		Message: "this run step holds " + why,
		Fix:     "move it into a checked-in script such as scripts/<name>.sh, and call that script"}
}

// whyHeavy says why the run step on lines[i] is no trivial command, or "".
// It also counts the body lines to skip, so text inside a block is no step.
func whyHeavy(lines []string, i int) (string, int) {
	m := runRe.FindStringSubmatch(lines[i])
	if m == nil {
		return "", blockBody(lines, i)
	}
	body := bodyLines(lines[i+1:], len(m[1])+len(m[2]))
	code := scriptLines(body)
	if !blockRe.MatchString(m[3]) {
		if nestedRun(m[3], code) {
			return "", 0
		}
		code = append(nonEmpty(m[3]), code...)
	}
	return blockWhy(code), len(body)
}

// nestedRun reports whether a run key holds a mapping, as defaults.run does, not a script.
func nestedRun(value string, code []string) bool {
	return value == "" && (len(code) == 0 || mappingRe.MatchString(code[0]))
}

// blockBody counts the body lines of a block scalar that line i opens, or 0.
func blockBody(lines []string, i int) int {
	m := blockKeyRe.FindStringSubmatch(lines[i])
	if m == nil {
		return 0
	}
	return len(bodyLines(lines[i+1:], len(m[1])))
}

// nonEmpty is a one-item list holding s, or none when s is empty.
func nonEmpty(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}

// blockWhy says why a step's code lines are no trivial command, or "".
func blockWhy(code []string) string {
	if len(code) > 1 {
		return fmt.Sprintf("%d lines of shell", len(code))
	}
	if len(code) == 1 {
		return operatorsIn(code[0])
	}
	return ""
}

// bodyLines is the lines under a key at indent: those indented deeper, blank lines included.
func bodyLines(rest []string, indent int) []string {
	n := 0
	for i, l := range rest {
		if strings.TrimSpace(l) != "" && len(l)-len(strings.TrimLeft(l, " ")) <= indent {
			break
		}
		if strings.TrimSpace(l) != "" {
			n = i + 1
		}
	}
	return rest[:n]
}

// scriptLines is the code lines of a body: neither blank nor shell comments.
func scriptLines(body []string) []string {
	var out []string
	for _, l := range body {
		text := strings.TrimSpace(l)
		if text != "" && !strings.HasPrefix(text, "#") {
			out = append(out, text)
		}
	}
	return out
}

// operatorsIn names the control flow or operator a single command line holds, or "".
func operatorsIn(command string) string {
	code := maskRe.ReplaceAllString(unquote(command), "x")
	if before, _, found := strings.Cut(code, " #"); found {
		code = before
	}
	if controlRe.MatchString(code) {
		return "shell control flow"
	}
	return operatorName(code)
}

// operatorName names the operator in a masked command line, or "".
func operatorName(code string) string {
	if strings.HasSuffix(code, `\`) {
		return "a line continuation"
	}
	for _, op := range operators {
		if strings.Contains(code, op.text) {
			return op.name
		}
	}
	return ""
}

// unquote is a YAML scalar's text without the quotes around the whole value.
func unquote(value string) string {
	if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
		return value[1 : len(value)-1]
	}
	return value
}
