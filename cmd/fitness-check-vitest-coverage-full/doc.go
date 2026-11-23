// Command fitness-check-vitest-coverage-full requires Vitest coverage
// thresholds of 100 in the repo and the shared config, then runs `vitest run
// --coverage`. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
