package main

import (
	"fmt"
	"regexp"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// nodeReportsCases is the Node switch, false until fitness-runner#1108 ships
// Vitest and Jest case names. While false, a Node test needs no case per THEN.
const nodeReportsCases = false

// casesReported names the language modules whose requirement tests report
// the cases they run. Swift reports XCTest activities, per requirement 0009.3
// of fitness-runner.
var casesReported = map[string]bool{"go": true, "java": true, "swift": true, "node": nodeReportsCases}

// owned is an acceptance, the path of the doc declaring it, and its tests.
type owned struct {
	requirement.IDLine
	path  string
	tests []reqTest
}

// caseID matches the dotted path a case name opens with, such as
// 0005.3.1.2.1, up to a space or the end. The name's number is any digits
// after the doc number, so a path with too few parts is still judged.
var caseID = regexp.MustCompile(`^(\d{4}(?:\.\d+)+)(?:\s|$)`)

// thenCaseErrors requires exactly one case for each THEN, named with its ID
// and its words. It also fails a case naming a path no doc holds or holding
// no assertion call.
func thenCaseErrors(root string, docs []requirement.Document, tests []reqTest) []string {
	var errs []string
	holders := thenHolders(docs)
	for _, a := range ownedAcceptances(docs, tests) {
		errs = append(errs, missingCases(a)...)
		for _, t := range a.tests {
			errs = append(errs, strayCases(a, t, holders)...)
		}
	}
	return append(errs, assertionErrors(root, tests)...)
}

// thenHolders maps each THEN ID of every doc to the acceptance holding it.
func thenHolders(docs []requirement.Document) map[string]string {
	out := map[string]string{}
	for _, d := range docs {
		for _, a := range d.Acceptances {
			for _, s := range a.Steps {
				for _, th := range s.Thens {
					out[th.ID] = a.ID
				}
			}
		}
	}
	return out
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

// caseName is the name a THEN's case carries: its ID, then its words.
func caseName(th requirement.Then) string {
	return th.ID + " " + th.Text
}

// missingCases fails each THEN of a without exactly one case in its test, or
// with a case of other words. Two owners already fail, so they are skipped.
func missingCases(a *owned) []string {
	t, ok := soleReporter(a)
	if !ok {
		return nil
	}
	byID := namesByID(t)
	var errs []string
	for _, s := range a.Steps {
		for _, th := range s.Thens {
			errs = append(errs, judgeThen(a, t, th, byID[th.ID])...)
		}
	}
	return errs
}

// namesByID groups the names of t's cases by the path each opens with.
func namesByID(t reqTest) map[string][]string {
	byID := map[string][]string{}
	for _, c := range t.Cases {
		if m := caseID.FindStringSubmatch(c.Name); m != nil {
			byID[m[1]] = append(byID[m[1]], c.Name)
		}
	}
	return byID
}

// judgeThen fails a THEN whose cases, the names opening with its ID, are not
// exactly one case named for it.
func judgeThen(a *owned, t reqTest, th requirement.Then, names []string) []string {
	want := caseName(th)
	switch {
	case len(names) == 0:
		return []string{requirement.At(a.path, th.Line, fmt.Sprintf("THEN %s has no test case; add a case to %s named %q", th.ID, t.Name, want))}
	case len(names) > 1:
		return []string{requirement.At(a.path, th.Line, fmt.Sprintf("THEN %s has %d test cases in %s; keep exactly one named %q", th.ID, len(names), t.Name, want))}
	case names[0] != want:
		return []string{requirement.At(a.path, th.Line, fmt.Sprintf("THEN %s has a case named %q in %s; rename it to %q, its ID then its words without THEN", th.ID, names[0], t.Name, want))}
	}
	return nil
}

// strayCases fails each case of t that opens with a path no THEN of a holds.
// It names the acceptance that holds the path when another does.
func strayCases(a *owned, t reqTest, holders map[string]string) []string {
	var errs []string
	for _, c := range t.Cases {
		m := caseID.FindStringSubmatch(c.Name)
		if m == nil {
			continue
		}
		switch home := holders[m[1]]; {
		case home == a.ID:
		case home == "":
			errs = append(errs, fmt.Sprintf("%s: case %q names %s, which no requirement doc holds; rename the case for a THEN of %s or delete it", c.Loc, c.Name, m[1], a.ID))
		default:
			errs = append(errs, fmt.Sprintf("%s: case %q names %s, which %s holds, not %s; move the case to its test", c.Loc, c.Name, m[1], home, a.ID))
		}
	}
	return errs
}
