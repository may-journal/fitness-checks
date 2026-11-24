// Command fitness-check-prettier runs Prettier `--check` over the repo's
// files, to ensure they are formatted. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
