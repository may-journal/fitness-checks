// Command readmes writes each check's README from the doc comments of its
// cmd/fitness-check-<name> package, with the .fitnessrc.json keys it reads.
// With -check it writes nothing and fails naming each README that differs.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	check := flag.Bool("check", false, "fail on a stale README instead of writing it")
	flag.Parse()
	if err := run(".", *check); err != nil {
		fmt.Fprintln(os.Stderr, "readmes:", err)
		os.Exit(1)
	}
}

// run renders every check's README under root, then writes each stale one,
// or with check set fails naming them.
func run(root string, check bool) error {
	files, err := generated(root)
	if err != nil {
		return err
	}
	stale, err := staleFiles(root, files)
	if err != nil {
		return err
	}
	if check {
		return staleError(stale)
	}
	return write(root, files, stale)
}

// staleError names the stale READMEs, or is nil when none is stale.
func staleError(stale []string) error {
	if len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("stale, run scripts/generate-readmes.sh: %s", strings.Join(stale, ", "))
}

// write writes each path in paths under root from files.
func write(root string, files map[string][]byte, paths []string) error {
	for _, path := range paths {
		if err := os.WriteFile(filepath.Join(root, path), files[path], 0o644); err != nil {
			return err
		}
	}
	return nil
}

// staleFiles lists, in path order, each file under root that differs from
// what files holds for it.
func staleFiles(root string, files map[string][]byte) ([]string, error) {
	var stale []string
	for _, path := range slices.Sorted(maps.Keys(files)) {
		differs, err := differs(filepath.Join(root, path), files[path])
		if err != nil {
			return nil, err
		}
		if differs {
			stale = append(stale, path)
		}
	}
	return stale, nil
}

// generated renders each check's README by its path under root.
func generated(root string) (map[string][]byte, error) {
	dirs, err := filepath.Glob(filepath.Join(root, "cmd", "fitness-check-*"))
	if err != nil {
		return nil, err
	}
	keys, err := configKeys(root)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, dir := range dirs {
		d, err := readCheckDoc(dir)
		if err != nil {
			return nil, err
		}
		files[filepath.ToSlash(filepath.Join("cmd", d.Folder, "README.md"))] = checkReadme(withConfig(d, keys[d.Name]))
	}
	return files, nil
}

// withConfig adds to d's body the section on the config keys it reads.
func withConfig(d checkDoc, keys []configKey) checkDoc {
	if config := configMarkdown(keys); config != "" {
		d.Body += "\n" + config
	}
	return d
}

// differs reports whether the file at path holds other than want. A missing
// file differs; any other read error is returned.
func differs(path string, want []byte) (bool, error) {
	got, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return false, err
	}
	return !bytes.Equal(got, want), nil
}
