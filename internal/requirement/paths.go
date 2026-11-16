package requirement

import (
	"path"
	"strings"

	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

// Dir is the folder that holds a repo's requirement docs.
const Dir = "docs/requirements/"

// DocPaths lists the requirement docs tracked under root, leaving out the
// template.
func DocPaths(root string) []string {
	var out []string
	for _, f := range walkfs.FilesByExt(root, ".md") {
		if strings.HasPrefix(f, Dir) && path.Base(f) != "template.md" {
			out = append(out, f)
		}
	}
	return out
}
