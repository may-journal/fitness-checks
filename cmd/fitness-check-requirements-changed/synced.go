package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// recordPath is the file the org sync writes beside the files it copies,
// naming each with its blob hash.
const recordPath = ".github/synced-files.json"

// syncRecord is what the sync wrote: the repo its files come from, and each
// file's blob hash by path.
type syncRecord struct {
	Source string            `json:"source"`
	Files  map[string]string `json:"files"`
}

// syncedSet holds, for each changed file the sync record names, whether its
// staged copy is the one the sync wrote. Source is the repo to change instead.
type syncedSet struct {
	source string
	intact map[string]bool
}

// readSynced checks each change the staged sync record names against its
// staged blob. A repo without a record, or with one that does not parse,
// syncs nothing.
func readSynced(root string, changes []change) (syncedSet, error) {
	staged, err := stagedBlobs(root, recordAndEdits(changes))
	if err != nil || staged[recordPath] == "" {
		return syncedSet{}, err
	}
	text, err := gitOut(root, "cat-file", "blob", staged[recordPath])
	var record syncRecord
	if err != nil || json.Unmarshal([]byte(text), &record) != nil {
		return syncedSet{}, err
	}
	return record.against(staged), nil
}

// recordAndEdits is the sync record's path, then the path of each change
// that is not a deletion.
func recordAndEdits(changes []change) []string {
	paths := []string{recordPath}
	for _, c := range changes {
		if !c.deleted {
			paths = append(paths, c.path)
		}
	}
	return paths
}

// against marks each staged file the record names as intact when its blob
// matches the record.
func (r syncRecord) against(staged map[string]string) syncedSet {
	set := syncedSet{source: r.Source, intact: map[string]bool{}}
	for path, blob := range staged {
		if want, ok := r.Files[path]; ok && path != recordPath {
			set.intact[path] = blob == want
		}
	}
	return set
}

// stagedBlobs maps each of paths the index holds to its blob hash.
func stagedBlobs(root string, paths []string) (map[string]string, error) {
	out, err := gitOut(root, append([]string{"ls-files", "-s", "-z", "--"}, paths...)...)
	blobs := map[string]string{}
	for _, entry := range strings.Split(out, "\x00") {
		// Each entry is "<mode> <blob> <stage>\t<path>".
		meta, path, ok := strings.Cut(entry, "\t")
		if fields := strings.Fields(meta); ok && len(fields) == 3 {
			blobs[path] = fields[1]
		}
	}
	return blobs, err
}

// verdict reports whether the sync owns path, and the problem a hand edit to
// it makes, none when the sync wrote it unchanged.
func (s syncedSet) verdict(path string) (owned bool, problems []string) {
	intact, owned := s.intact[path]
	if !owned || intact {
		return owned, nil
	}
	return true, []string{fmt.Sprintf("%s is synced from %s and was edited by hand; change it there instead", path, s.source)}
}
