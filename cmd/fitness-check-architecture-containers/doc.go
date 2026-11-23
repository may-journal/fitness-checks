// Command fitness-check-architecture-containers matches the containers view
// against the runtime products the repository builds and deploys. Its rule is
// in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
