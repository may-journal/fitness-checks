package requirements

import (
	"encoding/json"
	"maps"
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

func Test0160_2(t *testing.T) {
	t.Parallel()
	ours := goBuild0160(t, "./cmd/...")
	engine := goBuild0160(t, "github.com/may-journal/fitness-runner/go/cmd/fitness")
	for _, folder := range checkFolders0160(t, "..") {
		name := filepath.Base(folder)
		mustDo(t, os.Symlink(filepath.Join(engine, "fitness"), filepath.Join(engine, name)))
		got, want := describe0160(t, filepath.Join(ours, name)), describe0160(t, filepath.Join(engine, name))
		if got["protocol"] != float64(1) {
			t.Errorf("%s --describe has protocol %v, want 1", name, got["protocol"])
		}
		if !maps.Equal(got, want) {
			t.Errorf("%s --describe = %v, want the bundled check's %v", name, got, want)
		}
	}
}

// ranOurs0160 reports whether out says the engine started this repo's program
// for check. For local-branch it says the engine took the program as the
// declaration, since the judgment lives in the engine.
func ranOurs0160(out, check string) bool {
	return strings.Contains(out, "fitness-check-"+check+" replaces the built-in "+check+" check") ||
		strings.Contains(out, "fitness-check-"+check+" declares the "+check+" check, which fitness judges itself")
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
			t.Errorf("the engine ran its built-in %s, not this repo's program", check)
		}
	}
}
