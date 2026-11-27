// Package digest hashes the words of a THEN into the short string its test
// tag holds. It imports no test package, so a check can use it.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Length is how many hex digits of the sha256 a tag keeps.
const Length = 8

// Of is the first 8 hex digits of the sha256 of text with whitespace collapsed.
func Of(text string) string {
	sum := sha256.Sum256([]byte(strings.Join(strings.Fields(text), " ")))
	return hex.EncodeToString(sum[:])[:Length]
}
