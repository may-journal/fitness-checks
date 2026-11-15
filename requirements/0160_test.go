package requirements

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// engineOnly0160 reports whether dep is, or lies under, an engine package no
// check may depend on.
func engineOnly0160(dep string) bool {
	if strings.HasSuffix(dep, "/cmd/fitness") {
		return true
	}
	return slices.ContainsFunc([]string{"metadata", "aftereffect", "dupe", "report"}, func(name string) bool {
		return strings.HasSuffix(dep, "/internal/"+name) || strings.Contains(dep, "/internal/"+name+"/")
	})
}

// checkFolders0160 lists the check folders of the repo at root, such as
// cmd/fitness-check-gofmt.
func checkFolders0160(t *testing.T, root string) []string {
	t.Helper()
	folders, err := filepath.Glob(filepath.Join(root, "cmd", "fitness-check-*"))
	mustDo(t, err)
	if len(folders) == 0 {
		t.Fatal("no check folders under cmd")
	}
	return folders
}

// strayDeps lists what the check in folder depends on, tests included, that
// is the engine or another check.
func strayDeps(t *testing.T, folder string) []string {
	t.Helper()
	own := "/cmd/" + filepath.Base(folder)
	cmd := exec.Command("go", "list", "-deps", "-test", "-f", "{{.ImportPath}}", "."+own)
	cmd.Dir = ".."
	out, err := cmd.Output()
	mustDo(t, err)
	return slices.DeleteFunc(strings.Fields(string(out)), func(dep string) bool {
		pkg := strings.TrimSuffix(strings.Trim(dep, "[]"), ".test")
		other := strings.Contains(pkg, "/cmd/fitness-check-") && !strings.HasSuffix(pkg, own)
		return !other && !engineOnly0160(pkg)
	})
}

func Test0160_1(t *testing.T) {
	t.Parallel()
	for _, folder := range checkFolders0160(t, "..") {
		for _, dep := range strayDeps(t, folder) {
			t.Errorf("%s depends on %s", filepath.Base(folder), dep)
		}
	}
}

// goBuild0160 builds the packages into a fresh folder and returns it.
func goBuild0160(t *testing.T, packages ...string) string {
	t.Helper()
	bin := t.TempDir()
	build := exec.Command("go", append([]string{"build", "-o", bin + string(os.PathSeparator)}, packages...)...)
	build.Dir = ".."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

// describe0160 runs program with --describe and decodes what it prints.
func describe0160(t *testing.T, program string) map[string]any {
	t.Helper()
	out, err := exec.Command(program, "--describe").Output()
	mustDo(t, err)
	var got map[string]any
	mustDo(t, json.Unmarshal(out, &got))
	return got
}

// describeProblems0160 lists what is wrong with a check's --describe output.
func describeProblems0160(name string, got map[string]any) []string {
	var problems []string
	if got["protocol"] != float64(1) {
		problems = append(problems, fmt.Sprintf("%s --describe has protocol %v, want 1", name, got["protocol"]))
	}
	if got["name"] != strings.TrimPrefix(name, "fitness-check-") {
		problems = append(problems, fmt.Sprintf("%s --describe names the check %v", name, got["name"]))
	}
	if summary, _ := got["summary"].(string); summary == "" {
		problems = append(problems, name+" --describe has no summary")
	}
	if !slices.Contains([]any{"files", "blob", "git", "github"}, got["scope"]) {
		problems = append(problems, fmt.Sprintf("%s --describe has scope %v, want files, blob, git or github", name, got["scope"]))
	}
	return problems
}

func Test0160_2(t *testing.T) {
	t.Parallel()
	ours := goBuild0160(t, "./cmd/...")
	for _, folder := range checkFolders0160(t, "..") {
		name := filepath.Base(folder)
		for _, problem := range describeProblems0160(name, describe0160(t, filepath.Join(ours, name))) {
			t.Error(problem)
		}
	}
}

// ranOurs0160 reports whether out shows this repo's program for check ran,
// or for local-branch that the engine took the program as the declaration.
func ranOurs0160(out, check string) bool {
	if check == "local-branch" {
		return strings.Contains(out, "fitness-check-"+check+" declares the "+check+" check, which fitness judges itself")
	}
	return strings.Contains(out, "→ "+check+"\n")
}

func Test0160_3(t *testing.T) {
	t.Parallel()
	repo := example(t, "happyRepo", nil)
	all, code := fitness(t, repo, nil)
	sees(t, all, code, 0)
	for _, folder := range checkFolders0160(t, "..") {
		check := strings.TrimPrefix(filepath.Base(folder), "fitness-check-")
		out := all
		if !strings.Contains(all, "→ "+check+"\n") {
			out, _ = fitness(t, repo, nil, check)
		}
		if !ranOurs0160(out, check) {
			t.Errorf("the engine did not run this repo's program for %s", check)
		}
	}
}
