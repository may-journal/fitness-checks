// Command fitness-check-architecture-actors requires each requirement's Given
// to start with an actor the system context declares. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
