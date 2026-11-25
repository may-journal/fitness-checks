package main

import (
	"fmt"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// casesReported names the language modules whose scan reports the cases a
// test runs. Node reports none yet, so its chains need no case per step
// until fitness-runner issue 1099 lands.
var casesReported = map[string]bool{"go": true, "java": true, "swift": true}

// owned is an acceptance and the path of the doc declaring it.
type owned struct {
	requirement.IDLine
	path string
}

// stepCaseErrors requires a case named for each chained step, and fails a
// case naming a step its acceptance lacks.
func stepCaseErrors(docs []requirement.Document, tests []reqTest) []string {
	acceptances := byID(docs)
	var errs []string
	for _, t := range tests {
		if a, ok := acceptances[t.ID]; ok && casesReported[t.module.Name] {
			errs = append(errs, missingCases(a, t)...)
			errs = append(errs, strayCases(a.IDLine, t)...)
		}
	}
	return errs
}

// missingCases fails each numbered step of a that no case of t names.
func missingCases(a owned, t reqTest) []string {
	named := map[string]bool{}
	for _, c := range t.Cases {
		named[c.Step] = true
	}
	var errs []string
	for _, s := range a.Steps {
		if s.ID != "" && !named[s.ID] {
			errs = append(errs, requirement.At(a.path, s.Line, fmt.Sprintf("step %s has no test case; add a case to %s whose name opens with %s", s.ID, t.Name, s.ID)))
		}
	}
	return errs
}

// strayCases fails each case of t whose name opens with a step ID that a
// does not chain.
func strayCases(a requirement.IDLine, t reqTest) []string {
	chained := map[string]bool{}
	for _, s := range a.Steps {
		chained[s.ID] = s.ID != ""
	}
	var errs []string
	for _, c := range t.Cases {
		if c.Step != "" && !chained[c.Step] {
			errs = append(errs, fmt.Sprintf("%s: case %q names step %s, which %s does not chain; rename the case for one of its steps or delete it", c.Loc, c.Name, c.Step, a.ID))
		}
	}
	return errs
}

// byID maps each acceptance ID in docs to its acceptance and doc.
func byID(docs []requirement.Document) map[string]owned {
	out := map[string]owned{}
	for _, d := range docs {
		for _, a := range d.Acceptances {
			out[a.ID] = owned{a, d.Path}
		}
	}
	return out
}
