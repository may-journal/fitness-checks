// Command fitness-check-java-version checks a Java repo pins its JDK and that
// the machine's JDKs satisfy the pin, through the Java language module. Its
// rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
