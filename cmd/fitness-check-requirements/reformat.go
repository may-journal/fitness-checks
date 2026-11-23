package main

import (
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"github.com/may-journal/fitness-runner/go/pkg/comments"
	"github.com/may-journal/fitness-runner/go/pkg/gitx"
	"github.com/may-journal/fitness-runner/go/pkg/mdx"
)

// listMarker is a bullet or ordered-list marker at a line start that a formatter may rewrite, as `*` to `-`.
var listMarker = regexp.MustCompile(`(?m)^(\s*)(?:[-*+]|\d+[.)])(\s)`)

// sourceLine is a `Source:` line older docs end their Measurement with.
// Nothing reads it, so deleting it changes no wording.
var sourceLine = regexp.MustCompile(`(?m)^Source:.*$`)

// unchanged lists the changed docs whose wording is the same as at base,
// and those whose Requirements section is.
type unchanged struct {
	reworded, acceptances map[string]bool
}

// reformats returns the changed docs whose wording, or whose acceptances, did not change since base.
// It is empty when there is no base, and one git call reads every doc's two versions.
func reformats(root, base string, changed []string) (unchanged, error) {
	found := unchanged{map[string]bool{}, map[string]bool{}}
	docs := slices.DeleteFunc(slices.Clone(changed), func(f string) bool { return docNumber(f) == "" })
	if base == "" || len(docs) == 0 {
		return found, nil
	}
	from := docsAtBase(root, base, docs)
	var objects []string
	for _, f := range docs {
		objects = append(objects, base+":"+from[f], gitx.JudgedObject(f))
	}
	text, err := gitx.ShowAll(root, objects)
	for _, f := range docs {
		found.compare(f, text, base+":"+from[f])
	}
	return found, err
}

// compare records whether doc f says the same in text as its version before.
func (u unchanged) compare(f string, text map[string]string, before string) {
	was, held := text[before]
	is, judged := text[gitx.JudgedObject(f)]
	if !held || !judged {
		return
	}
	u.reworded[f] = wording(was, true) == wording(is, false)
	u.acceptances[f] = wording(acceptanceText(was), false) == wording(acceptanceText(is), false)
}

// requirementsHeading opens the section that holds a doc's acceptances.
var requirementsHeading = regexp.MustCompile(`(?m)^## Requirements\s*$`)

// acceptanceText is a doc from its Requirements heading on, or "" without one.
func acceptanceText(doc string) string {
	body := comments.WithoutHTMLComments(mdx.StripFrontMatter(doc))
	at := requirementsHeading.FindStringIndex(body)
	if at == nil {
		return ""
	}
	return body[at[0]:]
}

// docsAtBase maps each of docs to its path at base, read from one listing
// of base's docs. A doc renamed since base maps to the base doc of its
// number, since a number is never reused.
func docsAtBase(root, base string, docs []string) map[string]string {
	out, _ := exec.Command("git", "-C", root, "ls-tree", "--name-only", base, docsDir).Output()
	listed := strings.Fields(string(out))
	from := map[string]string{}
	for _, f := range docs {
		from[f] = renamedFrom(f, listed)
	}
	return from
}

// renamedFrom is f when base lists it, else the one listed doc with f's
// number. With none, or more than one, it is f.
func renamedFrom(f string, listed []string) string {
	if slices.Contains(listed, f) {
		return f
	}
	var same []string
	for _, doc := range listed {
		if docNumber(doc) == docNumber(f) {
			same = append(same, doc)
		}
	}
	if len(same) == 1 {
		return same[0]
	}
	return f
}

// wording is a doc's words in order without front matter, HTML comments, a base Source line, or formatter-changeable spacing and markers.
func wording(doc string, dropSource bool) string {
	body := comments.WithoutHTMLComments(mdx.StripFrontMatter(doc))
	if dropSource {
		body = sourceLine.ReplaceAllString(body, "")
	}
	body = listMarker.ReplaceAllString(body, "$1-$2")
	return strings.Join(strings.Fields(body), " ")
}
