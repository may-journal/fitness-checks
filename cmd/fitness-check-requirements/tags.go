package main

import (
	"fmt"
	"strings"

	"github.com/may-journal/fitness-checks/internal/marks"
	"github.com/may-journal/fitness-checks/internal/requirement"
)

// heldThen is a THEN with the doc that holds it.
type heldThen struct {
	requirement.Then
	path string
}

// allThens lists the THENs of docs in doc order.
func allThens(docs []requirement.Document) []heldThen {
	var out []heldThen
	for _, d := range docs {
		for _, a := range d.Acceptances {
			for _, s := range a.Steps {
				for _, th := range s.Thens {
					out = append(out, heldThen{Then: th, path: d.Path})
				}
			}
		}
	}
	return out
}

// thensByID maps each THEN ID of docs to its THEN.
func thensByID(docs []requirement.Document) map[string]heldThen {
	out := map[string]heldThen{}
	for _, th := range allThens(docs) {
		out[th.ID] = th
	}
	return out
}

// tagErrors fails each tag that names no THEN, holds a stale hash or repeats a tag
// of its THEN.
func tagErrors(docs []requirement.Document, tags []marks.Tag) []string {
	thens := thensByID(docs)
	first := map[string]marks.Tag{}
	var errs []string
	for _, tag := range tags {
		if msg := tagProblem(tag, thens, first); msg != "" {
			errs = append(errs, msg)
		} else {
			first[tag.ID] = tag
		}
	}
	return errs
}

// tagProblem is what is wrong with a tag, or "" for a tag that names its THEN well.
func tagProblem(tag marks.Tag, thens map[string]heldThen, first map[string]marks.Tag) string {
	th, ok := thens[tag.ID]
	switch {
	case !ok:
		return fmt.Sprintf("%s: tag %s names no THEN; delete it or tag the case with the ID of its THEN", tag.Loc(), tag.Value())
	case !tag.Placed():
		return fmt.Sprintf("%s: tag %s sits in %s, outside the test of its acceptance; move the case into Test%s", tag.Loc(), tag.Value(), tag.Func, strings.Replace(acceptanceOf(tag.ID), ".", "_", 1))
	case tag.Hash != th.Hash():
		return fmt.Sprintf("%s: tag %s is stale; reread THEN %s at %s:%d, then tag the case %s#%s", tag.Loc(), tag.Value(), tag.ID, th.path, th.Line, tag.ID, th.Hash())
	case first[tag.ID].File != "":
		return fmt.Sprintf("%s: THEN %s is tagged twice, also at %s; keep one case for it", tag.Loc(), tag.ID, first[tag.ID].Loc())
	}
	return ""
}

// acceptanceOf is the acceptance ID a THEN ID opens with, such as 0001.1.
func acceptanceOf(thenID string) string {
	parts := strings.SplitN(thenID, ".", 3)
	return parts[0] + "." + parts[1]
}

// taggedThens is the set of THEN IDs a tag names.
func taggedThens(tags []marks.Tag) map[string]bool {
	out := map[string]bool{}
	for _, t := range tags {
		out[t.ID] = out[t.ID] || t.Placed()
	}
	return out
}

// stepTagged reports whether every THEN of s has a tag. A step with no THEN
// item has none to tag.
func stepTagged(s requirement.Step, tagged map[string]bool) bool {
	for _, th := range s.Thens {
		if !tagged[th.ID] {
			return false
		}
	}
	return len(s.Thens) > 0
}

// untaggedErrors fails each THEN of docs that no tag names.
func untaggedErrors(docs []requirement.Document, tags []marks.Tag) []string {
	tagged := taggedThens(tags)
	var errs []string
	for _, th := range allThens(docs) {
		if !tagged[th.ID] {
			errs = append(errs, requirement.At(th.path, th.Line, fmt.Sprintf("THEN %s has no tagged test case; give its case the tag %s#%s", th.ID, th.ID, th.Hash())))
		}
	}
	return errs
}
