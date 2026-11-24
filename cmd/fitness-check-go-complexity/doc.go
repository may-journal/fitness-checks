// Command fitness-check-go-complexity enforces a cyclomatic complexity
// ceiling on Go functions, the Go-native counterpart of the house eslint rule
// (`complexity: max 5`). Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
