package main

import "github.com/may-journal/fitness-checks/internal/requirement"

// owned is an acceptance, the path of the doc declaring it, and its tests.
type owned struct {
	requirement.IDLine
	path  string
	tests []reqTest
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
