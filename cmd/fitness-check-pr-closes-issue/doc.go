// Command fitness-check-pr-closes-issue validates that a pull request
// description closes at least one issue on merge (#69). Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
