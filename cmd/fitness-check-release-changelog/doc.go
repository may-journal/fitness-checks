// Command fitness-check-release-changelog binds Release Please's manifest to
// the version file a reviewer can inspect. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
