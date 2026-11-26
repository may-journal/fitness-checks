package acceptance

import "strings"

// LongSentence has 30 words, past the 23-word default.
var LongSentence = strings.Repeat("word ", 29) + "end.\n"

// InComments writes a doc comment of lines on line 3 of a Go, a Swift and a TypeScript file.
func InComments(lines ...string) map[string]string {
	mark := func(marker, code string, lines []string) string {
		var b strings.Builder
		for _, l := range lines {
			b.WriteString(strings.TrimRight(marker+" "+l, " ") + "\n")
		}
		return b.String() + code
	}
	return map[string]string{
		"app/main.go":   "package main\n\n" + mark("//", "func main() {}\n", lines),
		"app/App.swift": "import Foundation\n\n" + mark("///", "func run() {}\n", lines),
		"app/index.ts":  "export {};\n\n" + mark("//", "export const x = 1;\n", lines),
	}
}

// AtEachFile names line 3 of each InComments file, followed by message.
func AtEachFile(message string) []string {
	return []string{"app/App.swift:3: " + message, "app/index.ts:3: " + message, "app/main.go:3: " + message}
}
