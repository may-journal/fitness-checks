package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// configKey is one .fitnessrc.json key a check reads, as its Go struct
// field declares it with a `check:"<name>"` tag.
type configKey struct {
	// Check is the check that reads the key.
	Check string
	// Key is the field's JSON key.
	Key string
	// Doc is the field's doc comment.
	Doc string
	// Fields are the nested keys, when the field is a struct.
	Fields []configKey
}

// configKeys reads struct fields tagged `check:"<name>"` by check name,
// erroring on any without a doc comment. It reads this module's internal
// folder and the go/pkg folder of the fitness-runner version go.mod pins.
func configKeys(root string) (map[string][]configKey, error) {
	engine, err := enginePkg(root)
	if err != nil {
		return nil, err
	}
	keys := map[string][]configKey{}
	for _, dir := range []string{filepath.Join(root, "internal"), engine} {
		err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil || !entry.IsDir() {
				return err
			}
			if path != dir && skipped(entry.Name()) {
				return filepath.SkipDir
			}
			return readTagged(path, keys)
		})
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}

// skipped reports whether a folder holds fixtures or vendored code, never a
// check's config.
func skipped(name string) bool {
	return name == "testdata" || name == "vendor"
}

// enginePkg is the go/pkg folder of the fitness-runner module go.mod pins,
// downloading the module when the cache lacks it.
func enginePkg(root string) (string, error) {
	cmd := exec.Command("go", "mod", "download", "-json", "github.com/may-journal/fitness-runner/go")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("download the fitness-runner module: %w", err)
	}
	var module struct{ Dir string }
	if err := json.Unmarshal(out, &module); err != nil || module.Dir == "" {
		return "", fmt.Errorf("find the fitness-runner module in %q: %v", out, err)
	}
	return filepath.Join(module.Dir, "pkg"), nil
}

// readTagged adds the check-tagged fields of the package in dir to keys. It
// parses the package only when one of its files holds such a tag.
func readTagged(dir string, keys map[string][]configKey) error {
	fset := token.NewFileSet()
	files, err := taggedPackage(fset, dir)
	if err != nil || len(files) == 0 {
		return err
	}
	types := structTypes(files)
	for _, file := range files {
		if err := taggedFields(fset, file, types, keys); err != nil {
			return err
		}
	}
	return nil
}

// taggedPackage parses every package source in dir with comments, or none
// when no source holds a check tag.
func taggedPackage(fset *token.FileSet, dir string) ([]*ast.File, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if !slices.ContainsFunc(entries, func(entry os.DirEntry) bool { return holdsTag(dir, entry) }) {
		return nil, nil
	}
	return parseSources(fset, dir, parser.ParseComments)
}

// holdsTag reports whether entry is a package source in dir holding a check
// tag.
func holdsTag(dir string, entry os.DirEntry) bool {
	if !packageSource(entry) {
		return false
	}
	src, err := os.ReadFile(filepath.Join(dir, entry.Name()))
	return err == nil && bytes.Contains(src, []byte(`check:"`))
}

// structTypes indexes the struct types the files declare by name, so a field
// of a named struct type shows its keys.
func structTypes(files []*ast.File) map[string]*ast.StructType {
	types := map[string]*ast.StructType{}
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if spec, ok := n.(*ast.TypeSpec); ok {
				if st, ok := spec.Type.(*ast.StructType); ok {
					types[spec.Name.Name] = st
				}
			}
			return true
		})
	}
	return types
}

// structOf is the struct type t names, so a list of objects shows its keys.
// It looks through slices, pointers and maps, and names types declares.
func structOf(t ast.Expr, types map[string]*ast.StructType) (*ast.StructType, bool) {
	switch t := t.(type) {
	case *ast.StructType:
		return t, true
	case *ast.Ident:
		st, ok := types[t.Name]
		return st, ok
	}
	if elem := element(t); elem != nil {
		return structOf(elem, types)
	}
	return nil, false
}

// element is the type a slice, pointer or map holds, or nil.
func element(t ast.Expr) ast.Expr {
	switch t := t.(type) {
	case *ast.ArrayType:
		return t.Elt
	case *ast.StarExpr:
		return t.X
	case *ast.MapType:
		return t.Value
	}
	return nil
}

// taggedFields adds each check-tagged field in file to keys.
func taggedFields(fset *token.FileSet, file *ast.File, types map[string]*ast.StructType, keys map[string][]configKey) error {
	var err error
	ast.Inspect(file, func(n ast.Node) bool {
		field, ok := n.(*ast.Field)
		if !ok || err != nil {
			return err == nil
		}
		check := tag(field, "check")
		if check == "" {
			return true
		}
		var key configKey
		key, err = readKey(fset, field, types)
		key.Check = check
		keys[check] = append(keys[check], key)
		return false
	})
	return err
}

// readKey reads a field's JSON key, doc comment, and nested keys.
func readKey(fset *token.FileSet, field *ast.Field, types map[string]*ast.StructType) (configKey, error) {
	key := configKey{Key: jsonKey(field), Doc: fieldDoc(field)}
	if key.Doc == "" {
		return key, fmt.Errorf("%s: config key %q has no doc comment", fset.Position(field.Pos()), key.Key)
	}
	nested, ok := structOf(field.Type, types)
	if !ok {
		return key, nil
	}
	var err error
	key.Fields, err = nestedKeys(fset, nested, types)
	return key, err
}

// nestedKeys reads the config keys among a struct's fields.
func nestedKeys(fset *token.FileSet, nested *ast.StructType, types map[string]*ast.StructType) ([]configKey, error) {
	var keys []configKey
	for _, inner := range slices.DeleteFunc(slices.Clone(nested.Fields.List), func(f *ast.Field) bool { return !configField(f) }) {
		sub, err := readKey(fset, inner, types)
		if err != nil {
			return nil, err
		}
		keys = append(keys, sub)
	}
	return keys, nil
}

// configField reports whether field is a key of the config file: one named,
// exported field that its json tag does not hide.
func configField(field *ast.Field) bool {
	return len(field.Names) == 1 && field.Names[0].IsExported() && tag(field, "json") != "-"
}

// jsonKey is a field's JSON key: its json tag name, else its Go name, as
// encoding/json reads it.
func jsonKey(field *ast.Field) string {
	if name := tag(field, "json"); name != "" || len(field.Names) != 1 {
		return name
	}
	return field.Names[0].Name
}

// fieldDoc is a field's doc comment, else its line comment, on one line,
// its Go name replaced by its JSON key.
func fieldDoc(field *ast.Field) string {
	text := field.Doc.Text()
	if text == "" {
		text = field.Comment.Text()
	}
	text = strings.Join(strings.Fields(text), " ")
	if len(field.Names) == 1 {
		if rest, ok := strings.CutPrefix(text, field.Names[0].Name+" "); ok {
			text = "`" + jsonKey(field) + "` " + rest
		}
	}
	return text
}

// tag is the name part of a field's struct tag key, or "".
func tag(field *ast.Field, key string) string {
	if field.Tag == nil {
		return ""
	}
	raw, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return ""
	}
	name, _, _ := strings.Cut(reflect.StructTag(raw).Get(key), ",")
	return name
}

// configMarkdown renders a check's config keys as a section, each key a
// paragraph with nested keys in a flat list, or "" for none.
func configMarkdown(keys []configKey) string {
	if len(keys) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("## Configuration\n\nThese `.fitnessrc.json` keys tune the check.\n")
	for _, key := range keys {
		b.WriteString("\n" + keyLine(key) + "\n")
		if len(key.Fields) > 0 {
			b.WriteString("\n")
		}
		writeNested(&b, key.Fields, "")
	}
	return b.String()
}

// writeNested writes nested keys as one flat list, naming deeper keys by path
// such as allow.pattern, so no list nests past prose-budget's depth.
func writeNested(b *strings.Builder, keys []configKey, prefix string) {
	for _, key := range keys {
		path := prefix + key.Key
		b.WriteString("- " + keyLine(renamed(key, path)) + "\n")
		writeNested(b, key.Fields, path+".")
	}
}

// renamed is key under the name path, with its doc comment's opening name
// changed to match.
func renamed(key configKey, path string) configKey {
	if rest, ok := strings.CutPrefix(key.Doc, "`"+key.Key+"`"); ok {
		key.Doc = "`" + path + "`" + rest
	}
	key.Key = path
	return key
}

// keyLine is a key's doc comment, led by the key itself when the comment
// does not already open with it.
func keyLine(key configKey) string {
	name := "`" + key.Key + "`"
	if strings.HasPrefix(key.Doc, name) {
		return key.Doc
	}
	return name + ": " + key.Doc
}
