// Command fitness-check-gitignore-why requires every .gitignore pattern to
// follow a `#` comment explaining why it exists. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
