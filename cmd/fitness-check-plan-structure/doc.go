// Command fitness-check-plan-structure validates that a plan, a GitHub Issue
// body under the `Plan` label, follows the Plan issue template. Its rule is
// in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
