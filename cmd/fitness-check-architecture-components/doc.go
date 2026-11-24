// Command fitness-check-architecture-components gives every product source
// one owner among the components the components view draws. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
