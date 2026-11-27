// Package requirements proves the go-test-coverage check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo holding a small Go module.
package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo with no Go modules passes without files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.1.1", "857dd667", func(_ req.A) {
			out, _ := goTestCoverage(t, nil)
			acceptance.PassedWithNoFiles(t, out, "go-test-coverage")
		})
	})
	t.Run("0001.1.2 a missing Go shows the download hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "188437da", func(_ req.A) {
			repo := acceptance.Example(t, "happyRepo", command(entryTest))
			out, code := acceptance.Fitness(t, repo, []string{acceptance.PathWithout(t, "go")}, "go-test-coverage")
			acceptance.ShowsInfo(t, out, code, "go-test-coverage", "Go not installed: see https://go.dev/dl")
		})
	})
	t.Run("0001.1.3 a library with no declared entry fails asking to declare entries", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.3.1", "c47c911d", func(_ req.A) {
			coverageSees(t, 1, "no entry packages; declare library entries with --entry=import/path", library("{}\n"))
		})
	})
	t.Run("0001.1.4 an unknown entry fails naming it", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "e721e89c", func(_ req.A) {
			coverageSees(t, 1, `unknown entry package "example.test/nope"`, command(entryTest), "--entry=example.test/nope")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 untested statements in each entry fail", func(t *testing.T) {
		t.Parallel()
		lib := library(`{"goTestCoverage": {"entries": ["example.test/cli/lib"]}}` + "\n")
		files := command("")
		files[".fitnessrc.json"] = lib[".fitnessrc.json"]
		for _, name := range []string{"lib.go", "lib_test.go"} {
			files["lib/"+name] = lib[name]
		}
		out, code := goTestCoverage(t, files)
		req.Then(t, "0001.2.1.1.1", "4c82cec8", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "entry example.test/cli:", "entry example.test/cli/lib:")
		})
		req.Then(t, "0001.2.1.1.2", "275cd92d", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "entry example.test/cli: 0/", "entry example.test/cli/lib: 1/")
		})
		req.Then(t, "0001.2.1.1.3", "5fd395cb", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "0/1 statements; required 100%", "1/2 statements; required 100%")
		})
	})
	t.Run("0001.2.2 a named helper passed to t.Run fails", func(t *testing.T) {
		t.Parallel()
		wrapped := "package main\n\nimport \"testing\"\n\nfunc TestAll(t *testing.T) { t.Run(\"entry\", runsMain) }\n\nfunc runsMain(t *testing.T) { main() }\n"
		out, code := goTestCoverage(t, command(wrapped))
		req.Then(t, "0001.2.1.2.1", "963cecf1", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "main_test.go:5:")
		})
		req.Then(t, "0001.2.1.2.2", "988d6b01", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "main_test.go:5: TestAll passes named helper runsMain to t.Run")
		})
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 the audit map is written outside the repo", func(t *testing.T) {
		t.Parallel()
		claims := filepath.Join(t.TempDir(), "claims.json")
		out, code := goTestCoverage(t, command(entryTest), "--audit-map="+claims)
		req.Then(t, "0001.3.1.1.1", "07b8e753", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.3.1.1.2", "c4a15532", func(a req.A) {
			_, err := os.Stat(claims)
			a.True(err == nil, "source map was not written: %v", err)
		})
	})
	t.Run("0001.3.2 unstable test coverage is reported", func(t *testing.T) {
		t.Parallel()
		files := command(flipTest)
		files["flip.go"] = flip
		out, code := goTestCoverage(t, files)
		req.Then(t, "0001.3.1.2.1", "1b963b89", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "example.test/cli/TestFlip (", "unstable coverage")
		})
		req.Then(t, "0001.3.1.2.2", "034c6558", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "excluded from overlap advice")
		})
	})
	t.Run("0001.3.3 parallel plus serial tests overlap correctly", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", slowTests())
		out, code := acceptance.Fitness(t, repo, nil, "go-test-coverage")
		runs := auditRuns(t, repo)
		if len(runs) != 32 {
			t.Fatalf("audit ran the sixteen tests %d times, want 32", len(runs))
		}
		counts := overlaps(runs)
		req.Then(t, "0001.3.1.3.1", "81beebb7", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			a.True(counts["parallel-parallel"] != 0, "no two runs of the parallel tests overlapped; go test runs them side by side")
		})
		req.Then(t, "0001.3.1.3.2", "74c079d8", func(a req.A) {
			a.True(counts["serial-serial"] != 0, "no two serial packages' runs overlapped; go test runs packages side by side")
		})
		req.Then(t, "0001.3.1.3.3", "a722ea83", func(a req.A) {
			a.True(counts["forbidden"] == 0, "%d pairs of runs overlapped within a serial package or one test; go test runs them in turn", counts["forbidden"])
		})
	})
}

// goTestCoverage runs `fitness-install -- go-test-coverage <args>` on
// happyRepo with files written over it.
func goTestCoverage(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"go-test-coverage"}, args...)...)
}

// command is a Go module with one command package of one statement.
func command(test string) map[string]string {
	return map[string]string{
		"go.mod":       "module example.test/cli\n\ngo 1.24\n",
		"main.go":      "package main\n\nfunc main() { println(42) }\n",
		"main_test.go": test,
	}
}

// entryTest is a command test that runs every statement.
const entryTest = "package main\n\nimport \"testing\"\n\nfunc TestEntry(t *testing.T) { main() }\n"

// library is a Go module with a two-statement library package, a config
// written over happyRepo's, and a test running only One.
func library(config string) map[string]string {
	return map[string]string{
		".fitnessrc.json": config,
		"go.mod":          "module example.test/lib\n\ngo 1.24\n",
		"lib.go":          "package lib\n\nfunc One() int { return 1 }\n\nfunc Two() int { return 2 }\n",
		"lib_test.go":     "package lib\n\nimport \"testing\"\n\nfunc TestOne(t *testing.T) { One() }\n",
	}
}

// flip is a command source file whose flip takes one branch for even
// counts and the other for odd ones.
const flip = "package main\n\nfunc flip(n int) int {\n\tif n%2 == 0 {\n\t\treturn 1\n\t}\n\treturn 2\n}\n"

// flipTest covers the command in TestEntry, while TestFlip takes the other
// branch of flip on each run, counted in a file.
const flipTest = `package main

import (
	"os"
	"testing"
)

func TestEntry(t *testing.T) { main(); flip(0); flip(1) }

func TestFlip(t *testing.T) {
	runs, _ := os.ReadFile("runs")
	if err := os.WriteFile("runs", append(runs, 'x'), 0o600); err != nil {
		t.Fatal(err)
	}
	flip(len(runs))
}
`

// slowTests adds sixteen one-second tests: eight call t.Parallel() in one
// package, and eight run serially, two in each of four more packages. Each
// run appends its package, test, process, start and end to runs.log at the
// module root.
func slowTests() map[string]string {
	files := command(entryTest)
	slowPackage(files, "parallel", 8, "\tt.Parallel()\n")
	for i := range 4 {
		slowPackage(files, fmt.Sprintf("serial%d", i), 2, "")
	}
	return files
}

// slowPackage writes package pkg with count one-second tests, each
// starting with first.
func slowPackage(files map[string]string, pkg string, count int, first string) {
	files[pkg+"/wait.go"] = "package " + pkg + "\n\nimport \"time\"\n\nfunc Wait() { time.Sleep(time.Second) }\n"
	test := "package " + pkg + "\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"testing\"\n\t\"time\"\n)\n" + record
	for i := range count {
		test += fmt.Sprintf("\nfunc TestWait%d(t *testing.T) {\n%s\trecord(t, %q)\n}\n", i, first, pkg)
	}
	files[pkg+"/wait_test.go"] = test
}

// record is the slow tests' helper: it waits a second and logs the run.
const record = `
func record(t *testing.T, pkg string) {
	start := time.Now().UnixNano()
	Wait()
	f, err := os.OpenFile("../runs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "%s %s %d %d %d\n", pkg, t.Name(), os.Getpid(), start, time.Now().UnixNano())
}
`

// run is one logged test run.
type run struct {
	pkg, test  string
	start, end int64
}

// auditRuns reads runs.log and keeps the audit runs whose process ran one
// test, since the module-wide run logs several per process.
func auditRuns(t *testing.T, repo string) []run {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repo, "runs.log"))
	if err != nil {
		t.Fatalf("the slow tests logged no runs: %v", err)
	}
	perProcess := map[string][]run{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		pid, run := parseRun(t, line)
		perProcess[pid] = append(perProcess[pid], run)
	}
	var runs []run
	for _, group := range perProcess {
		runs = append(runs, single(group)...)
	}
	return runs
}

// parseRun reads one runs.log line as its process ID and run.
func parseRun(t *testing.T, line string) (string, run) {
	t.Helper()
	var run run
	var pid string
	if _, err := fmt.Sscan(line, &run.pkg, &run.test, &pid, &run.start, &run.end); err != nil {
		t.Fatalf("bad runs.log line %q: %v", line, err)
	}
	return pid, run
}

// single keeps a process's run only when it ran exactly one test.
func single(group []run) []run {
	if len(group) != 1 {
		return nil
	}
	return group
}

// overlaps counts overlapping pairs of runs by pairKind.
func overlaps(runs []run) map[string]int {
	counts := map[string]int{}
	for i, a := range runs {
		for _, b := range runs[i+1:] {
			if a.start < b.end && b.start < a.end {
				counts[pairKind(a, b)]++
			}
		}
	}
	return counts
}

// pairKind names what an overlap of runs a and b shows: "forbidden"
// within one serial package or test, else the runs' sides.
func pairKind(a, b run) string {
	if a.pkg == b.pkg && (a.test == b.test || a.pkg != "parallel") {
		return "forbidden"
	}
	return side(a) + "-" + side(b)
}

// side is "parallel" for the parallel package's runs, else "serial".
func side(run run) string {
	if run.pkg == "parallel" {
		return "parallel"
	}
	return "serial"
}

// coverageSees runs go-test-coverage on files and args, asserting the exit code and message.
func coverageSees(t *testing.T, want int, msg string, files map[string]string, args ...string) {
	t.Helper()
	out, code := goTestCoverage(t, files, args...)
	acceptance.Sees(t, out, code, want, msg)
}
