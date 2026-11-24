// Command fitness-check-eslint runs ESLint over the repo with the shared flat
// config, by executing the real ESLint binary. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
