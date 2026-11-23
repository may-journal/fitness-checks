// Command fitness-check-swift-storybook-requirements fails each Swift
// Storybook page not named by a requirement acceptance, and each acceptance
// with no page (#79). Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
