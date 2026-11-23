// Command fitness-check-issue-checklist fails an issue body whose task list
// still has unchecked items, so closing it would close unfinished work. Its
// rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
