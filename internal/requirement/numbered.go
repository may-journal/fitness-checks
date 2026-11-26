package requirement

import "fmt"

// noThen is what a WHEN with no THEN under it fails with.
const noThen = "this WHEN needs a THEN nested under it"

// numbered tracks where an acceptance's items stand while checkNumbered
// reads them.
type numbered struct {
	id            string
	top           item
	givens, whens int
	thens         int
	indent        map[kind]int
	lastWhen      item
	maxSteps      int
}

// checkNumbered requires GIVENs under acceptance id's name, WHENs under the
// last GIVEN and THENs under each WHEN, each level numbered from 1. It
// reports whether the acceptance has that shape.
func checkNumbered(p *problems, id string, g []item, kinds map[int]kind) bool {
	st := numbered{id: id, top: g[0], indent: map[kind]int{}, maxSteps: p.maxSteps}
	for _, it := range g[1:] {
		n, msg, ok := st.next(it, kinds[it.n])
		if !ok {
			p.addAny(n, msg)
			return false
		}
	}
	n, msg := st.end()
	p.addAny(n, msg)
	return msg == ""
}

// addAny adds msg at line n unless msg is empty.
func (p *problems) addAny(n int, msg string) {
	if msg != "" {
		p.add(n, msg)
	}
}

// next reads item it of kind k, and reports whether it fits. When it
// doesn't, it also returns the line and message it fails with.
func (st *numbered) next(it item, k kind) (int, string, bool) {
	if msg, bad := unfit(it, k); bad {
		return it.n, msg, false
	}
	read := map[kind]func(item) (int, string, bool){given: st.readGiven, when: st.readWhen, then: st.readThen}
	return read[k](it)
}

// unfit is why an item of kind k is no step, and whether it is none. An AND
// or BUT has failed already.
func unfit(it item, k kind) (string, bool) {
	switch {
	case k == joined:
		return "", true
	case k == other:
		return "open each item with GIVEN, WHEN or THEN", true
	case !it.numbered:
		return "number each GIVEN, WHEN and THEN item; the numbered form holds no dashes", true
	}
	return "", false
}

// readGiven reads a GIVEN, which comes before any WHEN, nested under the name.
func (st *numbered) readGiven(it item) (int, string, bool) {
	if st.whens > 0 || !st.level(given, it, st.top, st.givens) {
		return it.n, "nest each GIVEN under the acceptance name, before its first WHEN", false
	}
	return st.count(it, &st.givens, givenWord)
}

// readWhen reads a WHEN nested under the last GIVEN, once the WHEN before it has
// a THEN.
func (st *numbered) readWhen(it item) (int, string, bool) {
	if st.whens > 0 && st.thens == 0 {
		return st.lastWhen.n, noThen, false
	}
	if st.givens == 0 || !st.level(when, it, item{indent: st.indent[given]}, st.whens) {
		return it.n, "nest each WHEN under the last GIVEN", false
	}
	st.lastWhen, st.thens = it, 0
	delete(st.indent, then)
	return st.count(it, &st.whens, whenWord)
}

// readThen reads a THEN nested under its WHEN.
func (st *numbered) readThen(it item) (int, string, bool) {
	if st.whens == 0 || !st.level(then, it, st.lastWhen, st.thens) {
		return it.n, "nest each THEN under its WHEN", false
	}
	return st.count(it, &st.thens, thenWord)
}

// level reports whether it sits deeper than parent, at the same indent as
// the seen items of kind k before it.
func (st *numbered) level(k kind, it, parent item, seen int) bool {
	if seen == 0 {
		st.indent[k] = it.indent
		return it.indent > parent.indent
	}
	return it.indent == st.indent[k]
}

// count requires it to be numbered one past *seen, then counts it.
func (st *numbered) count(it item, seen *int, word string) (int, string, bool) {
	if it.num != *seen+1 {
		return it.n, fmt.Sprintf("number this %s `%d.`; each level numbers from 1 in order", word, *seen+1), false
	}
	*seen++
	return it.n, "", true
}

// end fails an acceptance that stops before a WHEN and its THEN, or chains
// more steps than the limit.
func (st *numbered) end() (int, string) {
	switch {
	case st.whens == 0:
		return st.top.n, st.id + " needs a WHEN nested under its last GIVEN"
	case st.thens == 0:
		return st.lastWhen.n, noThen
	case st.whens > st.maxSteps:
		return st.top.n, fmt.Sprintf("%s chains %d steps, at most %d; move the rest to a new ID", st.id, st.whens, st.maxSteps)
	}
	return 0, ""
}
