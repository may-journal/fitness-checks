package main

import (
	"os/exec"
	"regexp"
	"slices"
	"strings"

	gherkin "github.com/may-journal/fitness-checks/cmd/fitness-check-requirements/migrations/0003.2-gherkin"
	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-runner/go/pkg/comments"
	"github.com/may-journal/fitness-runner/go/pkg/gitx"
	"github.com/may-journal/fitness-runner/go/pkg/mdx"
)

// listMarker is a bullet or ordered-list marker at a line start that a formatter may rewrite, as `*` to `-`.
var listMarker = regexp.MustCompile(`(?m)^(\s*)(?:[-*+]|\d+[.)])(\s)`)

// sourceLine is a `Source:` line older docs end their Measurement with.
// Nothing reads it, so deleting it changes no wording.
var sourceLine = regexp.MustCompile(`(?m)^Source:.*$`)

// reformats returns the changed docs whose wording did not change since
// base, so a formatter's rewrite, such as `biome format`, needs no test
// change. It is empty when there is no base, and one git call reads every
// doc's two versions.
func reformats(root string, roots []string, base string, changed []string) (map[string]bool, error) {
	found := map[string]bool{}
	docs := slices.DeleteFunc(slices.Clone(changed), func(f string) bool { return docNumber(roots, f) == "" })
	if base == "" || len(docs) == 0 {
		return found, nil
	}
	from := docsAtBase(root, roots, base, docs)
	var objects []string
	for _, f := range docs {
		objects = append(objects, base+":"+from[f], gitx.JudgedObject(f))
	}
	text, err := gitx.ShowAll(root, objects)
	for _, f := range docs {
		found[f] = sameWording(text, f, base+":"+from[f], gitx.JudgedObject(f))
	}
	return found, err
}

// sameWording reports whether text holds both versions of doc f, with the
// same wording in each. A doc of numbered GIVEN, WHEN and THEN items reads
// as its bullets.
func sameWording(text map[string]string, f, before, after string) bool {
	was, ok := text[before]
	if !ok {
		return false
	}
	is, ok := text[after]
	if !ok {
		return false
	}
	was, wasGherkin := gherkin.Bullets(f, was)
	is, isGherkin := gherkin.Bullets(f, is)
	moved := wasGherkin != isGherkin
	return wording(was, true, moved) == wording(is, false, moved)
}

// docsAtBase maps each of docs to its path at base, read from one listing
// of base's docs folders. A doc renamed since base maps to the base doc of
// its number in its package, since a number is never reused there.
func docsAtBase(root string, roots []string, base string, docs []string) map[string]string {
	var dirs []string
	for _, f := range docs {
		if dir := packages.Join(packages.Of(roots, f), docsDir); !slices.Contains(dirs, dir) {
			dirs = append(dirs, dir)
		}
	}
	out, _ := exec.Command("git", append([]string{"-C", root, "ls-tree", "--name-only", base}, dirs...)...).Output()
	listed := strings.Fields(string(out))
	from := map[string]string{}
	for _, f := range docs {
		from[f] = renamedFrom(roots, f, listed)
	}
	return from
}

// renamedFrom is f when base lists it, else the one listed doc with f's
// number. With none, or more than one, it is f.
func renamedFrom(roots []string, f string, listed []string) string {
	if slices.Contains(listed, f) {
		return f
	}
	var same []string
	for _, doc := range listed {
		if docNumber(roots, doc) == docNumber(roots, f) {
			same = append(same, doc)
		}
	}
	if len(same) == 1 {
		return same[0]
	}
	return f
}

// wording is a doc's words in order without front matter, HTML comments, a base Source line, or formatter-changeable spacing and markers.
// A doc moved between bullets and numbered items reads each comma as the
// "and" that joins an extra GIVEN or THEN.
func wording(doc string, dropSource, moved bool) string {
	body := comments.WithoutHTMLComments(mdx.StripFrontMatter(doc))
	if dropSource {
		body = sourceLine.ReplaceAllString(body, "")
	}
	if moved {
		body = commaJoins.Replace(body)
	}
	body = listMarker.ReplaceAllString(body, "$1-$2")
	return strings.Join(strings.Fields(body), " ")
}

// commaJoins reads a comma as the "and" that joins an extra GIVEN or THEN,
// so only that split compares equal. A "but" stays a "but".
var commaJoins = strings.NewReplacer(", but ", " but ", ", and ", " and ", ",", " and ")
