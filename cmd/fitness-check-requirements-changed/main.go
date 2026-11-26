package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/comments"
	"github.com/may-journal/fitness-runner/go/pkg/gitx"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: checkkit.Describe{
			Name:             "requirements-changed",
			Scope:            checkkit.ScopeGit,
			ContextInlineArg: "--message",
			Tests:            true,
			Requires:         []checkkit.Requirement{checkkit.RequiresRepo, checkkit.RequiresDiff},
		},
		Run: run,
	})
}

func run(root string, _ []string) (checkkit.Result, error) {
	roots := packages.Roots()
	docs := requirementDocs(root, roots)
	if skips(root, docs) {
		return checkkit.Pass(0), nil
	}
	scope, err := diffScope(root)
	if err != nil {
		return checkkit.Result{}, err
	}
	changes, err := readChanges(root, scope.args)
	if err != nil {
		return checkkit.Result{}, err
	}
	tests, err := scanTests(root, roots, docs)
	if err != nil {
		return checkkit.Result{}, err
	}
	return judgeCode(root, scope.base, tests, changes)
}

// judgeCode judges changes once comment-only edits are set aside. With no
// product code changed, it passes before comparing any file.
func judgeCode(root, base string, tests testSet, changes []change) (checkkit.Result, error) {
	if len(productChanges(root, tests, changes)) == 0 {
		return checkkit.Pass(0), nil
	}
	judged, err := withoutCommentOnly(root, base, changes)
	if err != nil {
		return checkkit.Result{}, err
	}
	synced, err := readSynced(root, judged)
	if err != nil {
		return checkkit.Result{}, err
	}
	return judge(root, tests, synced, judged), nil
}

// skips reports a run with nothing to judge: a full scan, a replay, or a
// repo that has no requirement docs yet.
func skips(root string, docs []string) bool {
	return gitx.CommitChanges(root) == nil || gitx.Replaying(root) || len(docs) == 0
}

// judge fails each changed product file whose test or doc did not change, and each deleted one whose doc and tests stay.
// A file the sync wrote pairs with its source repo instead.
func judge(root string, tests testSet, synced syncedSet, changes []change) checkkit.Result {
	product := productChanges(root, tests, changes)
	changed := changedSet(changes)
	var errs []string
	for _, c := range product {
		errs = append(errs, c.problem(root, tests, synced, changed)...)
	}
	if len(errs) > 0 {
		return checkkit.Fail(len(product), errs...)
	}
	return checkkit.Pass(len(product))
}

// withoutCommentOnly drops each edit whose code, comments aside, matches base, since it changes no behavior and needs no test or doc.
func withoutCommentOnly(root, base string, changes []change) ([]change, error) {
	var edits []comments.Edit
	for _, c := range changes {
		if !c.deleted {
			edits = append(edits, c.edit())
		}
	}
	same, err := comments.CommentOnly(root, base, edits)
	return slices.DeleteFunc(slices.Clone(changes), func(c change) bool {
		return same[c.edit()]
	}), err
}

// edit is the change as an edit from its path at base to its path now.
func (c change) edit() comments.Edit {
	return comments.Edit{From: c.from, File: c.path}
}

// problem names what a product change leaves out: a doc naming a deleted file, or a test or doc for an edit.
// Both must be in its own package, which needs requirement docs, and a synced file fails only when edited by hand.
func (c change) problem(root string, tests testSet, synced syncedSet, changed map[string]bool) []string {
	if c.deleted {
		return leftBehind(root, c.path, tests, changed)
	}
	if owned, problems := synced.verdict(c.path); owned {
		return problems
	}
	pkg := packages.Of(tests.roots, c.path)
	if !tests.documented(pkg) {
		return nil
	}
	if what, which := missing(tests.anyTest(pkg, changed), tests.anyDoc(pkg, changed)); what != "" {
		return []string{fmt.Sprintf("%s changed without %s; change %s with it", c.path, what, which)}
	}
	return nil
}

// missing names the half a change lacks, and the pronoun for it, or "" when
// it has both.
func missing(tested, documented bool) (what, which string) {
	switch {
	case !tested && !documented:
		return "a test and a requirement doc", "both"
	case !tested:
		return "a test", "one"
	case !documented:
		return "a requirement doc", "one"
	}
	return "", ""
}

// changedSet is the set of paths that changed, a rename by its new name.
func changedSet(changes []change) map[string]bool {
	set := map[string]bool{}
	for _, c := range changes {
		set[c.path] = true
	}
	return set
}

// joinList names files as a reader lists them.
func joinList(files []string) string {
	return strings.Join(files, ", ")
}
