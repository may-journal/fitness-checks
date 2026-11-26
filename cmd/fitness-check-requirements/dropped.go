package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// droppedErrors flags each acceptance ID history held that a numbered doc
// holds as neither an acceptance nor a Removed line. It names the Removed
// line that keeps the ID, and a bullet doc keeps its gaps.
func droppedErrors(root string, docs []requirement.Document, hist history) []string {
	var errs []string
	for _, d := range docs {
		for _, n := range dropped(root, d, hist) {
			errs = append(errs, requirement.At(d.Path, d.TitleLine, fmt.Sprintf("%s.%d was an acceptance and this doc no longer holds it; "+
				"a removed acceptance keeps its number, so add `%d. %s in YYYY.MM.DD.HHMM: <reason>` to Requirements, with the changelog version that removes it",
				d.ID, n, n, requirement.RemovedName)))
		}
	}
	return errs
}

// dropped is the numbers, in order, of the acceptance IDs history held that
// numbered doc d drops. Only a doc missing some acceptance is read for its
// Removed lines.
func dropped(root string, d requirement.Document, hist history) []int {
	gone := missing(d, hist)
	if !d.Gherkin || len(gone) == 0 {
		return nil
	}
	kept := removedLines(root, d.Path)
	return slices.DeleteFunc(gone, func(n int) bool { return kept[fmt.Sprintf("%s.%d", d.ID, n)] != "" })
}

// missing is the numbers, in order, of the acceptance IDs history held that
// d holds as no acceptance. An ID whose acceptance moved to another number
// is the move check's to name, so it is not missing.
func missing(d requirement.Document, hist history) []int {
	holds := map[string]bool{}
	for _, a := range d.Acceptances {
		holds[a.ID] = true
		holds[hist.homes[homeKey(d.ID, a.Name)]] = true
	}
	var gone []int
	for id := range hist.held {
		if n, ok := strings.CutPrefix(id, d.ID+"."); ok && !holds[id] {
			num, _ := strconv.Atoi(n)
			gone = append(gone, num)
		}
	}
	slices.Sort(gone)
	return gone
}

// removedLines files the IDs of the Removed lines the doc at path holds. It
// reads the doc as a diff that adds it whole, as history reads a commit.
func removedLines(root, path string) map[string]string {
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		return nil
	}
	_, names := commitNet(diffHeader + "a/" + path + " b/" + path + "\n+" + strings.ReplaceAll(string(b), "\n", "\n+"))
	return names[removedSide]
}
