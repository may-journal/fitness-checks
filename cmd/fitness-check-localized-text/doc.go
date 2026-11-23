// Command fitness-check-localized-text fails SwiftUI text and requirement
// test literals that skip the app's String Catalog (#527). Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
