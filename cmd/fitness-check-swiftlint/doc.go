// Command fitness-check-swiftlint runs SwiftLint over the repo through the
// Swift language module, which owns the tool run and its judgment, via
// `swiftlint lint --strict`. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
