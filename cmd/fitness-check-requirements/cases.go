package main

import (
	"fmt"

	"github.com/may-journal/fitness-checks/internal/marks"
	"github.com/may-journal/fitness-checks/internal/requirement"
)

// casesReported names the language modules whose requirement tests report
// the cases they run. Node reports none until fitness-runner issue 1099, and
// Swift only XCTest activities, which a Swift Testing test may not use.
var casesReported = map[string]bool{"go": true, "java": true}

// owned is an acceptance, the path of the doc declaring it, and its tests.
type owned struct {
	requirement.IDLine
	path  string
	tests []reqTest
}

// stepCaseErrors requires a case named for each chained step, and fails a
// case naming a step its acceptance lacks.
func stepCaseErrors(docs []requirement.Document, tests []reqTest, tags []marks.Tag) []string {
	var errs []string
	tagged := taggedThens(tags)
	for _, a := range ownedAcceptances(docs, tests) {
		errs = append(errs, missingCases(a, tagged)...)
		for _, t := range a.tests {
			errs = append(errs, strayCases(a, t)...)
		}
	}
	return errs
}

// ownedAcceptances lists each well-formed acceptance in docs with the tests
// named for it.
func ownedAcceptances(docs []requirement.Document, tests []reqTest) []*owned {
	out, byID := wellFormed(docs)
	for _, t := range tests {
		if a := byID[t.ID]; a != nil {
			a.tests = append(a.tests, t)
		}
	}
	return out
}

// wellFormed lists each acceptance in docs that holds steps, and indexes
// them by ID. A malformed chain holds none, and its doc already fails.
func wellFormed(docs []requirement.Document) ([]*owned, map[string]*owned) {
	var out []*owned
	byID := map[string]*owned{}
	for _, d := range docs {
		for _, a := range d.Acceptances {
			if len(a.Steps) > 0 {
				byID[a.ID] = &owned{IDLine: a, path: d.Path}
				out = append(out, byID[a.ID])
			}
		}
	}
	return out, byID
}

// soleReporter is a's one test, when its language reports cases.
func soleReporter(a *owned) (reqTest, bool) {
	if len(a.tests) != 1 {
		return reqTest{}, false
	}
	return a.tests[0], casesReported[a.tests[0].module.Name]
}

// stepIDs is the set of a's numbered step IDs.
func stepIDs(a *owned) map[string]bool {
	ids := map[string]bool{}
	for _, s := range a.Steps {
		ids[s.ID] = s.ID != ""
	}
	return ids
}

// caseSteps is the set of step IDs t's case names open with.
func caseSteps(t reqTest) map[string]bool {
	steps := map[string]bool{}
	for _, c := range t.Cases {
		steps[c.Step] = c.Step != ""
	}
	return steps
}

// missingCases fails each numbered step of a that no case of its one test
// names. Two owners already fail, so they are not judged here.
func missingCases(a *owned, tagged map[string]bool) []string {
	t, ok := soleReporter(a)
	if !ok {
		return nil
	}
	named := caseSteps(t)
	var errs []string
	for _, s := range a.Steps {
		if uncovered(s, named, tagged) {
			errs = append(errs, requirement.At(a.path, s.Line, fmt.Sprintf("step %s has no test case; add a case to %s whose name opens with %s", s.ID, t.Name, s.ID)))
		}
	}
	return errs
}

// uncovered reports whether step s has no case named for it and no tag on each THEN.
func uncovered(s requirement.Step, named, tagged map[string]bool) bool {
	return s.ID != "" && !named[s.ID] && !stepTagged(s, tagged)
}

// strayCases fails each case of t whose name opens with a step ID that a
// does not chain.
func strayCases(a *owned, t reqTest) []string {
	chained := stepIDs(a)
	var errs []string
	for _, c := range t.Cases {
		if c.Step != "" && !chained[c.Step] {
			errs = append(errs, fmt.Sprintf("%s: case %q names step %s, which %s does not chain; rename the case for one of its steps or delete it", c.Loc, c.Name, c.Step, a.ID))
		}
	}
	return errs
}
