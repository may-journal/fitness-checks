// Command describe-patterns reads a check's --describe line on stdin, and
// exits 1 when a pattern does not compile or a requirement is unknown.
// scripts/release-build.sh runs it with the release tag as its one argument.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	if err := validate(); err != nil {
		fmt.Fprintf(os.Stderr, "describe-patterns: %s: %v\n", label(), err)
		os.Exit(1)
	}
}

// validate checks the describe line on stdin as the engine would.
func validate() error {
	var d checkkit.Describe
	if err := json.NewDecoder(os.Stdin).Decode(&d); err != nil {
		return err
	}
	if _, err := checkkit.PatternsOf(d); err != nil {
		return err
	}
	if r := checkkit.UnknownRequirement(d); r != "" {
		return errors.New("requires " + string(r) + ", which fitness-runner does not define")
	}
	return nil
}

// label names the release in a failure, or stdin when no tag was given.
func label() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "stdin"
}
