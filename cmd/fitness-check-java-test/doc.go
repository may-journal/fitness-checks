// Command fitness-check-java-test runs Java tests once per build root and
// judges their JUnit XML results and JaCoCo coverage. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
