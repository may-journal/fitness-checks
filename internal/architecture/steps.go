package architecture

import "regexp"

// stepID is a cited step such as 0003.6.2, which a journey label reads as
// doc 0003.
var stepID = regexp.MustCompile(`^\d{4}\.\d+\.\d+$`)

// docOf is the doc a cited step reads as, else id.
func docOf(id string) string {
	if stepID.MatchString(id) {
		return id[:4]
	}
	return id
}

// docsOf reads each cited step of ids as its doc.
func docsOf(ids []string) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = docOf(id)
	}
	return out
}

// mark records the doc a valid label ID names, reading a cited step as its
// doc.
func (j *journey) mark(id string, at Location) {
	if j.valid(id, at) {
		j.named[docOf(id)] = true
	}
}

// step requires a cited step to exist in some requirement doc.
func (j *journey) step(id string, at Location) bool {
	if !j.m.hasID(id) {
		j.fail(at, "no step %s in docs/requirements", id)
		return false
	}
	return true
}
