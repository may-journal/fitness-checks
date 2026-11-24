// Command fitness-check-changelog-updated checks that each staged commit
// writes a changelog entry about its change. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
