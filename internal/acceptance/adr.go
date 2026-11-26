package acceptance

// The local check fixtures the adr-as-issue and adr-checks requirement tests
// share with the repo's own.

// Passing is a local check's verdict that judged one file and found nothing.
const Passing = `{"ok": true, "errors": [], "filesChecked": 1}`

// VerdictScript is a check program that prints verdict, or hangs when it
// is empty.
func VerdictScript(verdict string) string {
	if verdict == "" {
		return "#!/bin/sh\nexec /bin/sleep 30\n"
	}
	return "#!/bin/sh\necho '" + verdict + "'\n"
}

// HeldDecision is a check README that still holds its decision.
const HeldDecision = "# 0001 — Small files\n\nSmall files read fast.\n\n## Context\n\nLong files hide their point.\n\n" +
	"## Decision\n\n1. Every Markdown file stays small.\n\n## Consequences\n\n1. Readers find the point fast.\n\n" +
	"## Check\n\nFails a long Markdown file.\n"

// HeldCheck is ADR 0001's check folder, its README holding the decision.
func HeldCheck() map[string]string {
	return map[string]string{"fitness/adr-0001/README.md": HeldDecision, "fitness/adr-0001/check": VerdictScript(Passing)}
}

// AllowADRs is a config whose adrAsIssue.allow holds entries.
func AllowADRs(entries string) string {
	return `{"adrAsIssue": {"allow": [` + entries + `]}}` + "\n"
}
