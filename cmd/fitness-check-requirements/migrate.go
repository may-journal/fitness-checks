package main

import (
	"errors"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// fix deletes the leftover Source line from each doc at paths when args ask
// for `--fix`, before the check judges.
func fix(root string, paths, args []string) error {
	if flag := unknownFlag(args); flag != "" {
		return flagError{"unknown flag " + flag + "; the only flag is " + fixFlag + ", which deletes Source lines"}
	}
	if !slices.Contains(args, fixFlag) {
		return nil
	}
	var err error
	for _, path := range paths {
		err = errors.Join(err, dropSource(root, path))
	}
	return err
}

// flagError is a flag mistake, which fails the check instead of stopping it.
type flagError struct{ msg string }

func (e flagError) Error() string { return e.msg }

// fixResult fails the check on a flag mistake and stops it on any other error.
func fixResult(err error) (checkkit.Result, error) {
	var bad flagError
	if errors.As(err, &bad) {
		return checkkit.Fail(0, bad.Error()), nil
	}
	return checkkit.Result{}, err
}

// dropSource deletes the leftover Source line from the doc at path.
func dropSource(root, path string) error {
	docs, _, _ := requirement.Load(root, []string{path})
	if len(docs) == 0 {
		return nil
	}
	return requirement.DropSource(root, docs[0])
}

// unknownFlag is the first dash argument other than the fix flag, such as the
// removed --up and --down, or "" when there is none.
func unknownFlag(args []string) string {
	for _, a := range args {
		if strings.HasPrefix(a, "-") && a != fixFlag {
			return a
		}
	}
	return ""
}
