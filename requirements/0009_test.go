package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// dispatchStep is the Release workflow's step that announces a release.
func dispatchStep(t *testing.T) string {
	t.Helper()
	text, err := os.ReadFile(filepath.Join("..", ".github", "workflows", "release.yml"))
	mustDo(t, err)
	_, step, ok := strings.Cut(string(text), "- name: Announce the release\n")
	if !ok {
		t.Fatal("release.yml has no step that announces the release")
	}
	step, _, _ = strings.Cut(step, "      - name:")
	return step
}

func Test0009_1(t *testing.T) {
	t.Parallel()
	bin := t.TempDir()
	log := filepath.Join(bin, "calls.log")
	mustDo(t, os.WriteFile(filepath.Join(bin, "gh"), []byte("#!/bin/sh\nprintf '%s\\n' \"$*\" >> \""+log+"\"\n"), 0o755))
	script, err := filepath.Abs(filepath.Join("..", "scripts", "release-dispatch.sh"))
	mustDo(t, err)
	cmd := exec.Command(script)
	cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "DISPATCH_REPOSITORY=acme/.github")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("release-dispatch.sh: %v\n%s", err, out)
	}
	calls, err := os.ReadFile(log)
	mustDo(t, err)
	for _, want := range []string{"api -X POST repos/acme/.github/dispatches", "event_type=fitness-check-released", "client_payload[released_at]="} {
		if !strings.Contains(string(calls), want) {
			t.Errorf("gh was called with %q, want %q", calls, want)
		}
	}
	if step := dispatchStep(t); !strings.Contains(step, "scripts/release-dispatch.sh") {
		t.Errorf("the announce step does not run the script:\n%s", step)
	}
}

func Test0009_2(t *testing.T) {
	t.Parallel()
	step := dispatchStep(t)
	for _, want := range []string{"needs.publish.result == 'success'", "needs.plan.outputs.tags != '[]'"} {
		if !strings.Contains(step, want) {
			t.Errorf("the announce step does not wait for %q:\n%s", want, step)
		}
	}
}
