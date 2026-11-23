// Command fitness-check-requirements enforces docs/requirements in every repo
// a language module covers. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
