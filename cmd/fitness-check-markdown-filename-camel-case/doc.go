// Command fitness-check-markdown-filename-camel-case requires every markdown
// file name to be camelCase, standard root docs exempt. The validation core
// is shared with the kebab-case flavor in `internal/mdfilename`; see
// [markdown-filename-kebab-case] for behavior and conventions covering both
// flavors.
//
// This flavor pins the camelCase convention: every Markdown basename must match
// `^[a-z][a-zA-Z0-9]*\.md$`, standard root docs and names GitHub requires
// exempt.
//
// [markdown-filename-kebab-case]: https://pkg.go.dev/github.com/may-journal/fitness-checks/cmd/fitness-check-markdown-filename-kebab-case
package main

import _ "embed"

// docSource is this file, whose comment --describe summarizes.
//
//go:embed doc.go
var docSource string
