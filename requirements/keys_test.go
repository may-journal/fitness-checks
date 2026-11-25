package requirements

// checkKeys maps each check that reads .fitnessrc.json keys to each leaf key
// it declares and that key's default, nil for none. No other check declares
// a key; adr-checks reads adrAsIssue, which adr-as-issue declares.
var checkKeys = map[string]map[string]any{
	"repeated-string-literals": {"repeatedStringLiterals.allow": nil},
	"text-readability":         {"textReadability.maxGrade": 18.0, "textReadability.maxLix": 60.0, "textReadability.minWords": 100.0},
	"prose-budget": {
		"proseBudget.maxSentenceWords": 23.0, "proseBudget.maxParagraphSentences": 4.0, "proseBudget.maxSectionParagraphs": 3.0,
		"proseBudget.maxListItemWords": 23.0, "proseBudget.maxListItems": 8.0, "proseBudget.maxWords": 300.0,
		"prose.comments.maxWordsPerToken": 3.0, "prose.comments.minBodyComments": 2.0,
	},
	"comment-structure": {"comment.maxWords": 100.0, "comment.maxTables": 1.0},
	"adr-as-issue":      {"adrAsIssue.allow": nil},
	"go-test-coverage":  {"goTestCoverage.entries": nil},
	"go-complexity":     {"goComplexity.max": 5.0},
	"localized-text":    {"localizedText.reference": nil, "localizedText.allow": nil},
	"requirements":      {"requirements.maxAcceptances": 8.0},
}
