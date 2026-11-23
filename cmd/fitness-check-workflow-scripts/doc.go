// Command fitness-check-workflow-scripts fails a `run:` step in a workflow
// or action that holds more than a trivial command ([#946]). Shell inside YAML
// cannot be linted or tested, so real work belongs in a checked-in script.
//
// # Behavior
//
//   - It reads the `.yml` and `.yaml` files directly under `.github/workflows`, and every `action.yml` and `action.yaml`.
//   - Pass: a one-line step with no control flow, heredoc, pipe, `;`, `&&` or `||`, such as `run: scripts/release.sh`.
//   - Fail: a block scalar with more than one line of shell, not counting blank lines and comments.
//   - Fail: one line that opens with `if`, `for`, `while`, `until`, `case`, `select` or `function`.
//   - Fail: one line with a heredoc, a pipe, a `;`, an `&&` or `||` chain, or a trailing `\`.
//   - Text inside quotes, `${{ }}` expressions and `$(( ))`, a trailing comment, `\;` and `<<<` never counts, so `printf '%s | %s'` passes.
//   - A step may span lines as a plain or quoted scalar too, and counts the same.
//
// A repo with no workflow or action files passes having checked none. In a
// pull request, or a staged commit, it reads only the changed files, and it
// never runs a step.
//
// # Errors
//
//	.github/workflows/ci.yml:12: this run step holds 4 lines of shell; move it into a checked-in script such as scripts/<name>.sh, and call that script
//	action.yml:30: this run step holds a pipe; move it into a checked-in script such as scripts/<name>.sh, and call that script
//
// # Fixing a step
//
// Put the shell in `scripts/<name>.sh`, pass values through `env:`, and call
// the script from the step:
//
//	env:
//	  TAG: ${{ steps.release.outputs.tag }}
//	run: scripts/move-major-tag.sh
//
// [#946]: https://github.com/may-journal/fitness-runner/issues/946
package main

import _ "embed"

// docSource is this file, whose comment --describe summarizes.
//
//go:embed doc.go
var docSource string
