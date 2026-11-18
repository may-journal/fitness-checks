// Command describe-patterns reads a check's --describe line on stdin, and
// exits 1 when its Files, Exclude or Needs pattern does not compile.
// scripts/release-build.sh runs it, so no release ships a pattern the
// engine rejects.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	var d checkkit.Describe
	err := json.NewDecoder(os.Stdin).Decode(&d)
	if err == nil {
		_, err = checkkit.PatternsOf(d)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "describe-patterns: %s: %v\n", d.Name, err)
		os.Exit(1)
	}
}
