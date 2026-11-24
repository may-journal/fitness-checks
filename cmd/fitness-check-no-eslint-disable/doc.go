// Command fitness-check-no-eslint-disable fails any JS or TS source file that
// holds an ESLint disable directive, in file, block, `-line`, or `-next-line`
// form. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
