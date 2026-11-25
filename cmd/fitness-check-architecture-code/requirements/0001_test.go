// Package requirements proves the architecture-code check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo whose code view draws Go, Swift, TypeScript or Java.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// The shared helpers and architecture fixtures, named as the root
// requirement tests name them.
var (
	acFixture          = acceptance.AcFixture
	acView             = acceptance.AcView
	architectureChecks = acceptance.ArchitectureChecks
	commit             = acceptance.Commit
	example            = acceptance.Example
	fitness            = acceptance.Fitness
	git                = acceptance.Git
	passedWith         = acceptance.PassedWith
	passedWithNoFiles  = acceptance.PassedWithNoFiles
	prEnv              = acceptance.PullRequestEnv
	sees               = acceptance.Sees
	with               = acceptance.With
)

const (
	acCode = acceptance.AcCode
	acRoot = acceptance.AcRoot
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 classes matching their declarations pass", func(t *testing.T) {
		t.Parallel()
		for _, files := range []map[string]string{cvFixture(), cvNative(), cvScript(), cvJava()} {
			out, code := cvRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, cvCheck, "1")
		}
		apart := with(cvFixture(), map[string]string{"src/store/load.go": "package store\n\n// Load reads the record back.\nfunc (r Record) Load() Record { return r }\n"})
		apart[cvView] = cvCode("\nclass load[\"4 store.Record.Load() Record\"]\nload --> Record : 5", "Reads it back, in [load.go](../src/store/load.go).", "Load returns the record.")
		out, code := cvRun(t, apart)
		sees(t, out, code, 0)
		files := cvEdit(cvFixture(), cvView, "+ID\n", "+ID: string\n+Load() (Record, error)\n")
		files = cvEdit(files, cvStore, cvSaves, cvSaves+"\n// Load reads the record back.\nfunc (r Record) Load() (Record, error) { return r, nil }\n")
		out, code = cvRun(t, files)
		sees(t, out, code, 0)
		embeds := cvEdit(cvFixture(), cvStore, "package store\n", "package store\n\nimport \"sync\"\n")
		embeds = cvEdit(embeds, cvStore, "\tID string\n", "\tID string\n\tsync.Mutex\n")
		out, code = cvRun(t, cvEdit(embeds, cvView, "+ID\n", "+ID\n+Mutex\n"))
		sees(t, out, code, 0)
		compact := cvEdit(cvJava(), cvJavaDir+"Store.java", "int count) {}", "int count) {\n    public Entry {}\n  }")
		out, code = cvRun(t, cvEdit(compact, cvView, "+count: int\n", "+count: int\n+Entry(String, int)\n"))
		sees(t, out, code, 0)
		out, code = cvRun(t, cvEdit(cvJava(), cvJavaDir+"DiskStore.java", "    Writer.write", "    int mask = 0xE+Writer.write"))
		sees(t, out, code, 0)
		out, code = cvRun(t, cvEdit(cvJava(), cvView, "+records: Map~String, Integer~", "+Map~String, Integer~ records"))
		sees(t, out, code, 0)
		spaced := cvEdit(cvFixture(), cvView, cvSave, "class save[\"2 store.Save( id string ) Record\"]")
		out, code = cvRun(t, cvEdit(spaced, cvView, "+ID\n", "+ID : string\n"))
		sees(t, out, code, 0)
		broken := with(cvFixture(), map[string]string{"src/store/testdata/broken.go": "package broken\n\nfunc {\n"})
		out, code = cvRun(t, broken)
		sees(t, out, code, 0)
		passedWith(t, out, cvCheck, "1")
		extended := cvEdit(cvNative(), cvModel, "{ true }\n}\n", "{ true }\n}\n\nextension Store {\n    func reset() {}\n}\n")
		out, code = cvRun(t, cvEdit(extended, cvView, "class Store[\"2 Store\"]", "class Store[\"2 Store\"] {\n+reset()\n}"))
		sees(t, out, code, 0)
	})
	t.Run("0001.1.2 a class whose linked lines lack its symbol fails", func(t *testing.T) {
		t.Parallel()
		moved := with(cvFixture(), map[string]string{cvStore: cvRecord, "src/store/save.go": "package store\n\n" + cvSaves})
		out, code := cvRun(t, moved)
		sees(t, out, code, 1, cvView+":8: save: store.Save(id string) Record is declared at src/store/save.go:4, not in the linked src/store/store.go")
		out, code = cvRun(t, cvEdit(moved, cvView, "Writes one record, in [store.go](../src/store/store.go)", "Writes one record, in [save.go](../src/store/save.go)"))
		sees(t, out, code, 0)
		renamed := cvEdit(cvFixture(), cvStore, "type Record struct", "type Entry struct")
		renamed = cvEdit(renamed, cvStore, "Record {\n\treturn Record{", "Entry {\n\treturn Entry{")
		out, code = cvRun(t, renamed)
		sees(t, out, code, 1, cvView+":5: Record: src/store/store.go declares no store.Record")
		out, code = cvRun(t, cvEdit(cvScript(), cvView, "disk, in [disk.ts](../src/store/disk.ts)", "disk, in [store.ts](../src/store/store.ts)"))
		sees(t, out, code, 1, cvView+":13: DiskStore: DiskStore is declared at "+cvTSDisk+":3, not in the linked "+cvTSStore)
	})
	t.Run("0001.1.3 a kind, signature, member or field type that differs fails", func(t *testing.T) {
		t.Parallel()
		cases := []struct{ file, old, new, want string }{
			{cvStore, "Save(id string) Record {", "Save(id string, force bool) Record {", ":8: save: store.Save at src/store/store.go:9 has signature (string, bool) Record, not store.Save(id string) Record"},
			{cvView, cvSave, "class save[\"2 store.Save\"]", ":8: save: store.Save at src/store/store.go:9 is a function, not a type"},
			{cvStore, "\tID string\n", "\tKey string\n", ":6: Record: store.Record at src/store/store.go:4 declares no field ID"},
			{cvView, "+ID\n", "+ID: int\n", ":6: Record: field store.Record.ID at src/store/store.go:5 is string, not int"},
			{cvView, "+ID\n", "+ID\n+Load() Record\n", ":7: Record: store.Record at src/store/store.go:4 declares no method Load"},
			{cvView, "+ID\n", "<<interface>>\n+ID\n", ":5: Record: store.Record at src/store/store.go:4 is a type, not an interface"},
		}
		for _, c := range cases {
			out, code := cvRun(t, cvEdit(cvFixture(), c.file, c.old, c.new))
			sees(t, out, code, 1, cvView+c.want)
		}
		embeds := cvEdit(cvFixture(), cvStore, "package store\n", "package store\n\nimport \"sync\"\n")
		embeds = cvEdit(embeds, cvStore, "\tID string\n", "\tID string\n\tsync.Mutex\n")
		out, code := cvRun(t, cvEdit(embeds, cvView, "+ID\n", "+ID\n+Locker\n"))
		sees(t, out, code, 1, cvView+":7: Record: store.Record at src/store/store.go:6 declares no field Locker")
		script := []struct{ file, old, new, want string }{
			{cvView, "+save(string) boolean", "+save(number) boolean", cvView + ":11: Store: method Store.save at " + cvTSStore + ":13 has signature (string) boolean, not save(number) boolean"},
			{cvView, "number~\n+save", "string~\n+save", cvView + ":10: Store: field Store.records at " + cvTSStore + ":10 is Map<string, number>, not Map<string, string>"},
		}
		for _, c := range script {
			out, code := cvRun(t, cvEdit(cvScript(), c.file, c.old, c.new))
			sees(t, out, code, 1, c.want)
		}
		java := []struct{ file, old, new, want string }{
			{cvView, "+save(String) boolean", "+save(int) boolean", cvView + ":11: Store: method shop.Store.save at " + cvJavaDir + "Store.java:10 has signature (String) boolean, not save(int) boolean"},
			{cvView, "+count: int", "+count: long", cvView + ":20: Entry: field shop.Store.Entry.count at " + cvJavaDir + "Store.java:15 is int, not long"},
			{cvView, "+count: int", "+static final int count", cvView + ":20: Entry: member \"+static final int count\" names no field or method"},
		}
		for _, c := range java {
			out, code := cvRun(t, cvEdit(cvJava(), c.file, c.old, c.new))
			sees(t, out, code, 1, c.want)
		}
		shelf := cvEdit(cvJava(), cvJavaDir+"Store.java", "{}\n}\n", "{}\n}\n\nclass Shelf {\n  record Entry(long size) {}\n}\n")
		out, code = cvRun(t, cvEdit(shelf, cvView, "+count: int\n", "+count: int\n+size: long\n"))
		sees(t, out, code, 1, cvView+":21: Entry: shop.Store.Entry at "+cvJavaDir+"Store.java:15 declares no field size")
	})
	t.Run("0001.1.4 a name declared twice lists each declaration", func(t *testing.T) {
		t.Parallel()
		opener := "class opener[\"4 store.open()\"]"
		both := "Opens storage, in [linux](../src/store/open_linux.go) or [darwin](../src/store/open_darwin.go)."
		files := with(cvFixture(), map[string]string{
			"src/store/open_linux.go":  "//go:build linux\n\npackage store\n\nfunc open() {}\n",
			"src/store/open_darwin.go": "//go:build darwin\n\npackage store\n\nfunc open() {}\n",
		})
		files[cvView] = cvCode("\n"+opener, both)
		out, code := cvRun(t, files)
		sees(t, out, code, 1, cvView+":10: opener: store.open is declared 2 times, at src/store/open_darwin.go:5 under go:build darwin, src/store/open_linux.go:5 under go:build linux; link one file or line range, or state its signature")
		files[cvView] = cvCode("\n"+opener, "Opens storage, in [linux](../src/store/open_linux.go).")
		out, code = cvRun(t, files)
		sees(t, out, code, 0)
		native := cvNative()
		native[cvModel] = strings.Replace(native[cvModel], "    func save(_ id: String) -> Bool { true }\n", "#if DEBUG\n    func save(_ id: String) -> Bool { true }\n#else\n    func save(_ id: String) -> Bool { false }\n#endif\n    @available(iOS 17, *) func save(_ ids: [String]) -> Bool { true }\n", 1)
		native[cvView] = strings.Replace(native[cvView], "HomeScreen --> Store : 5", "HomeScreen --> Store : 5\nclass saveOne[\"6 Store.save(String) Bool\"]\nclass saveAll[\"7 Store.save([String]) Bool\"]", 1)
		native[cvView] += "| 6 | Saves one, in [Store.swift](../App/Model/Store.swift#L7-L12). | Shared. |\n| 7 | Saves many, in [Store.swift](../App/Model/Store.swift). | Shared. |\n"
		out, code = cvRun(t, native)
		sees(t, out, code, 1, cvView+":12: saveOne: Store.save is declared 2 times, at "+cvModel+":9 under #if DEBUG, "+cvModel+":11 under #else of #if DEBUG")
		native[cvView] = strings.Replace(native[cvView], "Store.swift#L7-L12", "Store.swift#L8-L9", 1)
		out, code = cvRun(t, native)
		sees(t, out, code, 0)
	})
	t.Run("0001.1.5 a class linking a test, generated or other file fails", func(t *testing.T) {
		t.Parallel()
		tested := with(cvFixture(), map[string]string{"src/store/store_test.go": "package store\n\n// Fake stands in for a record.\ntype Fake struct{}\n"})
		tested[cvView] = cvCode("\nclass Fake[\"4 store.Fake\"]", "A fake, in [store_test.go](../src/store/store_test.go).")
		out, code := cvRun(t, tested)
		sees(t, out, code, 1, cvView+":10: Fake: src/store/store_test.go declares no store.Fake")
		generated := with(cvFixture(), map[string]string{".gitattributes": "src/store/gen.go linguist-generated\n", "src/store/gen.go": "package store\n\n// Gen is generated.\ntype Gen struct{}\n"})
		generated[cvView] = cvCode("\nclass Gen[\"4 store.Gen\"]", "Generated, in [gen.go](../src/store/gen.go).")
		out, code = cvRun(t, generated)
		sees(t, out, code, 1, cvView+":10: Gen: src/store/gen.go is a generated source; the code view shows product code")
		out, code = cvRun(t, cvEdit(cvFixture(), cvView, "A saved record, "+cvLink, "A saved record, in [check.sh](../tools/check.sh)."))
		sees(t, out, code, 1, cvView+":14: unsupported evidence tools/check.sh; Go, Swift, TypeScript, and Java sources only")
	})
	t.Run("0001.1.6 two classes naming one declaration or a label naming none fails", func(t *testing.T) {
		t.Parallel()
		files := cvFixture()
		files[cvView] = cvCode("\nclass again[\"4 store.Record\"]", "The record again, "+cvLink)
		out, code := cvRun(t, files)
		sees(t, out, code, 1, cvView+":10: again and Record at "+cvView+":5 both name store.Record")
		files = cvFixture()
		files[cvView] = cvCode("\nclass Contract[\"4 Save contract\"]\nsave --> Contract : 5", "What saving promises.", "Save keeps the promise.")
		out, code = cvRun(t, files)
		sees(t, out, code, 1, cvView+":10: Contract \"Save contract\" names no declaration; name a Go, Swift, TypeScript, or Java symbol, or mark it <<abstract>>")
		files[cvView] = cvCode("\nclass Contract[\"4 Save contract\"] {\n<<abstract>>\n+promise: one record\n}\nsave --> Contract : 5", "What saving promises.", "Save keeps the promise.")
		out, code = cvRun(t, files)
		sees(t, out, code, 0)
	})
	t.Run("0001.1.7 an exported type the view leaves out fails", func(t *testing.T) {
		t.Parallel()
		files := cvEdit(cvFixture(), cvStore, cvSaves, cvSaves+"\n// Batch is many records.\ntype Batch struct{}\n\ntype batch struct{}\n")
		out, code := cvRun(t, files)
		sees(t, out, code, 1, cvStore+":14: store.Batch is an exported type in src/store, which "+cvView+" shows, but no node names it")
		files[cvView] = cvCode("\nclass Batch[\"4 store.Batch\"]", "Many records, "+cvLink)
		out, code = cvRun(t, files)
		sees(t, out, code, 0)
		generator := with(cvFixture(), map[string]string{"src/store/gen.go": "//go:build ignore\n\npackage main\n\n// Tool writes records by hand.\ntype Tool struct{}\n\nfunc main() {}\n"})
		out, code = cvRun(t, generator)
		sees(t, out, code, 0)
		out, code = cvRun(t, cvEdit(generator, "src/store/gen.go", "//go:build ignore\n\npackage main", "package store"))
		sees(t, out, code, 1, "src/store/gen.go:4: store.Tool is an exported type in src/store, which "+cvView+" shows, but no node names it")
		out, code = cvRun(t, cvEdit(cvScript(), "src/store/cache.ts", "", "export default class Cache {}\n\nclass Local {}\n"))
		sees(t, out, code, 1, "src/store/cache.ts:1: Cache is an exported type in src/store, which "+cvView+" shows, but no node names it")
		out, code = cvRun(t, cvEdit(cvJava(), cvJavaDir+"Cache.java", "", "package shop;\n\n// public class Ghost {}\npublic class Cache {}\n\nclass Local {}\n"))
		sees(t, out, code, 1, cvJavaDir+"Cache.java:4: shop.Cache is an exported type in src/main/java/shop, which "+cvView+" shows, but no node names it")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a source relation no drawn edge allows fails", func(t *testing.T) {
		t.Parallel()
		files := cvEdit(cvFixture(), cvView, cvEdge, "class load[\"3 store.Load() string\"]")
		files = cvEdit(files, cvView, "Save returns the record.", "Loads a name, "+cvLink)
		files[cvStore] += "\n// Load names a Record in prose only.\nfunc Load() string { return \"Record\" } // Record\n"
		out, code := cvRun(t, files)
		sees(t, out, code, 1, cvStore+":9: save references Record")
		files[cvView] = strings.Replace(files[cvView], "class load", cvEdge+"\nclass load", 1)
		files[cvView] = strings.Replace(files[cvView], "| 3 | Loads a name, "+cvLink, "| 3 | Save returns the record. | Shared by later checks. |\n| 4 | Loads a name, "+cvLink, 1)
		files[cvView] = strings.Replace(files[cvView], "\"3 store.Load", "\"4 store.Load", 1)
		out, code = cvRun(t, files)
		sees(t, out, code, 0)
		native := cvEdit(cvNative(), cvView, "HomeScreen --> Store : 5", "")
		native[cvHome] = "import SwiftUI\n\nstruct HomeScreen {\n    let title = \"Store\" // Store\n    let store = Store()\n}\n"
		out, code = cvRun(t, native)
		sees(t, out, code, 1, cvHome+":5: HomeScreen calls Store, but "+cvView+" "+cvNoEdge+"HomeScreen -> Store"+cvAllowIt)
		// The deinit sits in Store but outside the saveOne method's lines.
		native = cvEdit(cvNative(), cvModel, "{ true }\n", "{ true }\n    deinit {\n        _ = HomeScreen()\n    }\n")
		native = cvEdit(native, cvView, "HomeScreen --> Store : 5", "HomeScreen --> Store : 5\nclass saveOne[\"6 Store.save(String) Bool\"]")
		native[cvView] += "| 6 | Saves one, in [Store.swift](../App/Model/Store.swift). | Shared. |\n"
		out, code = cvRun(t, native)
		sees(t, out, code, 1, cvModel+":10: Store calls HomeScreen, but "+cvView+" "+cvNoEdge+"Store -> HomeScreen"+cvAllowIt)
		native = cvEdit(native, cvView, "Store.save(String) Bool\"]", "Store.save(String) Bool\"]\nStore --> HomeScreen : 7")
		native[cvView] += "| 7 | The store makes a screen. | Shared. |\n"
		out, code = cvRun(t, native)
		sees(t, out, code, 0)
		for _, c := range []struct{ old, want string }{
			{"DiskStore --> write : 8", cvTSDisk + ":6: DiskStore calls write, but " + cvView + " " + cvNoEdge + "DiskStore -> write" + cvAllowIt},
			{"\nHandler --> Store : 11", cvTSView + ":3: Handler references Store, but " + cvView + " " + cvNoEdge + "Handler -> Store" + cvAllowIt},
		} {
			out, code := cvRun(t, cvEdit(cvScript(), cvView, c.old, ""))
			sees(t, out, code, 1, c.want)
		}
		out, code = cvRun(t, cvEdit(cvJava(), cvView, "DiskStore --> write : 9\n", ""))
		sees(t, out, code, 1, cvJavaDir+"DiskStore.java:6: DiskStore calls write, but "+cvView+" "+cvNoEdge+"DiskStore -> write"+cvAllowIt)
		helper := cvEdit(cvFixture(), cvView, cvEdge, cvEdge+"\nclass load[\"4 store.Load() string\"]")
		helper[cvView] += "| 4 | Loads a name, " + cvLink + " | Shared. |\n"
		helper[cvStore] += "\n// Load names the record through a helper.\nfunc Load() string { return name() }\n\nfunc name() string { _ = Record{}; return \"\" }\n"
		out, code = cvRun(t, helper)
		sees(t, out, code, 1, cvStore+":14: load references Record, but "+cvView+" "+cvNoEdge+"load -> Record"+cvAllowIt)
	})
	t.Run("0001.2.2 each arrow allows a use its direction permits", func(t *testing.T) {
		t.Parallel()
		for _, edge := range []string{"save ..> Record : 3", "save *-- Record : 3", "save o-- Record : 3", "Record -- save : 3", "Record .. save : 3"} {
			out, code := cvRun(t, cvEdit(cvFixture(), cvView, cvEdge, edge))
			sees(t, out, code, 0)
		}
		out, code := cvRun(t, cvEdit(cvFixture(), cvView, cvEdge, "Record --> save : 3"))
		sees(t, out, code, 1, cvView+":9: "+cvView+" draws association Record -> save, but no source shows Record calls or references or embeds save")
	})
	t.Run("0001.2.3 a drawn edge with no source relation of its kind fails", func(t *testing.T) {
		t.Parallel()
		out, code := cvRun(t, cvEdit(cvNative(), cvHome, "let store = Store()", "let count = 0"))
		sees(t, out, code, 1, cvView+":11: "+cvView+" draws association HomeScreen -> Store, but no source shows HomeScreen calls or references or embeds Store")
		out, code = cvRun(t, cvSink("Write"))
		sees(t, out, code, 1, cvView+":12: "+cvView+" draws realization Disk -> Sink, but no source shows Disk conforms to Sink")
		out, code = cvRun(t, cvEdit(cvNative(), cvModel, "final class Store: Saving {", "final class Store {"))
		sees(t, out, code, 1, cvView+":10: "+cvView+" draws realization Store -> Saving, but no source shows Store conforms to Saving")
		base := cvEdit(cvNative(), cvModel, "final class Store: Saving {", "class Base {}\n\nfinal class Store: Base, Saving {")
		base = cvEdit(base, cvView, "HomeScreen --> Store : 5", "HomeScreen --> Store : 5\nclass Base[\"6 Base\"]\nBase <|-- Store : 7")
		base = cvEdit(base, cvView, "| 5 | The screen makes a store. | Shared by later checks. |\n", "| 5 | The screen makes a store. | Shared by later checks. |\n| 6 | The base, in [Store.swift](../App/Model/Store.swift). | Shared. |\n| 7 | Store extends it. | Shared. |\n")
		out, code = cvRun(t, base)
		sees(t, out, code, 0)
		out, code = cvRun(t, cvEdit(base, cvModel, "Store: Base, Saving", "Store: Saving"))
		sees(t, out, code, 1, cvView+":13: "+cvView+" draws inheritance Store -> Base, but no source shows Store embeds or inherits from Base")
		out, code = cvRun(t, cvEdit(cvScript(), cvTSStore, "Store implements Saving", "Store"))
		sees(t, out, code, 1, cvView+":18: "+cvView+" draws realization Store -> Saving, but no source shows Store conforms to Saving")
		out, code = cvRun(t, cvEdit(cvJava(), cvJavaDir+"Store.java", "Store implements Saving", "Store"))
		sees(t, out, code, 1, cvView+":22: "+cvView+" draws realization Store -> Saving, but no source shows Store conforms to Saving")
	})
	t.Run("0001.2.4 an edge row's link must hold its relation", func(t *testing.T) {
		t.Parallel()
		out, code := cvRun(t, cvSink("Put"))
		sees(t, out, code, 0)
		sink := cvEdit(cvSink("Put"), cvView, "Disk is a sink.", "Disk is a sink, see [source](../src/store/store.go#L22).")
		out, code = cvRun(t, sink)
		sees(t, out, code, 0)
		out, code = cvRun(t, cvEdit(sink, cvView, "store.go#L22)", "store.go#L1)"))
		sees(t, out, code, 1, cvView+":22: evidence src/store/store.go#L1 shows no Disk -> Sink relation")
		linked := cvEdit(cvFixture(), cvView, "Save returns the record.", "Save returns the record, see [source](../src/store/store.go#L10).")
		out, code = cvRun(t, linked)
		sees(t, out, code, 0)
		out, code = cvRun(t, cvEdit(linked, cvView, "store.go#L10)", "store.go#L1)"))
		sees(t, out, code, 1, cvView+":16: evidence src/store/store.go#L1 shows no save -> Record relation")
	})
	t.Run("0001.2.5 an edge touching an abstract class needs no use", func(t *testing.T) {
		t.Parallel()
		files := cvFixture()
		files[cvView] = cvCode("\nclass Contract[\"4 Save contract\"] {\n<<abstract>>\n+promise: one record\n}\nsave --> Contract : 5", "What saving promises.", "Save keeps the promise.")
		out, code := cvRun(t, files)
		sees(t, out, code, 0)
		files[cvView] = strings.Replace(files[cvView], "Save keeps the promise.", "Save keeps the promise, see [source](../src/store/store.go#L10).", 1)
		out, code = cvRun(t, files)
		sees(t, out, code, 1, cvView+":23: unsupported evidence src/store/store.go#L10; save -> Contract touches an abstract node, which no source declares")
	})
	t.Run("0001.2.6 any other arrow is unsupported", func(t *testing.T) {
		t.Parallel()
		out, code := cvRun(t, cvEdit(cvFixture(), cvView, cvEdge, "save <--> Record : 3"))
		sees(t, out, code, 1, cvView+":9: unsupported relation <--> from save to Record; draw -->, ..>, <|--, or ..|>")
	})
	t.Run("0001.2.7 only changes that can move the views are judged", func(t *testing.T) {
		t.Parallel()
		out, code := cvRun(t, map[string]string{"src/main.go": "package main\n"})
		sees(t, out, code, 0)
		passedWithNoFiles(t, out, cvCheck)
		// The base already lacks the save edge; the pull request adds a test only.
		repo := example(t, "happyRepo", cvEdit(cvFixture(), cvView, cvEdge, ""))
		git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		commit(t, repo, map[string]string{"src/store/store_test.go": cvTestFile}, "test(store): add the store test")
		for _, check := range architectureChecks {
			out, _ := fitness(t, repo, prEnv, check)
			passedWithNoFiles(t, out, check)
		}
		out, code = fitness(t, repo, nil, "--all", cvCheck)
		sees(t, out, code, 1, cvStore+":9: save references Record, but "+cvView+" "+cvNoEdge+"save -> Record")
		for _, changed := range []string{"", "go.mod", cvStore, cvView} {
			out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, cvCheck)
			sees(t, out, code, 1, cvStore+":9: save references Record, but "+cvView+" "+cvNoEdge+"save -> Record"+cvAllowIt)
		}
	})
}

const (
	cvCheck   = "architecture-code"
	cvView    = acRoot + acCode
	cvStore   = "src/store/store.go"
	cvModel   = "App/Model/Store.swift"
	cvHome    = "App/Screens/HomeScreen.swift"
	cvLink    = "in [store.go](../src/store/store.go)."
	cvSave    = "class save[\"2 store.Save(id string) Record\"]"
	cvEdge    = "save --> Record : 3"
	cvRecord  = "package store\n\n// Record is one saved record.\ntype Record struct {\n\tID string\n}\n\n"
	cvSaves   = "// Save writes one record.\nfunc Save(id string) Record {\n\treturn Record{ID: id}\n}\n"
	cvNoEdge  = "draws no "
	cvAllowIt = " edge that allows it"
)

// cvTestFile is a Go test file in package store with one passing test.
const cvTestFile = "package store\n\nimport \"testing\"\n\nfunc TestStore(t *testing.T) {}\n"

// cvCode is the Go code view; the record sits on line 5, save on 8, the edge
// on 9.
func cvCode(extra string, rows ...string) string {
	diagram := "classDiagram\nclass Record[\"1 store.Record\"] {\n+ID\n}\n" + cvSave + "\n" + cvEdge + extra
	return acView(diagram, append([]string{"A saved record, " + cvLink, "Writes one record, " + cvLink, "Save returns the record."}, rows...)...)
}

// cvFixture is a Go repository whose code view matches its source.
func cvFixture() map[string]string {
	return with(acFixture(acRoot), map[string]string{
		"go.mod": "module example.com/app\n\ngo 1.26\n",
		cvStore:  cvRecord + cvSaves,
		cvView:   cvCode(""),
	})
}

// cvNative is a Swift app's code view: Saving on line 5, Store 8, HomeScreen
// 9, edges 10 and 11.
func cvNative() map[string]string {
	files := acFixture(acRoot)
	files[cvModel] = "import Foundation\n\nprotocol Saving {\n    func save(_ id: String) -> Bool\n}\n\nfinal class Store: Saving {\n    func save(_ id: String) -> Bool { true }\n}\n"
	files[cvHome] = "import SwiftUI\n\nstruct HomeScreen {\n    let store = Store()\n}\n"
	files[cvView] = acView("classDiagram\nclass Saving[\"1 Saving\"] {\n<<interface>>\n}\nclass Store[\"2 Store\"]\nclass HomeScreen[\"3 HomeScreen\"]\nStore ..|> Saving : 4\nHomeScreen --> Store : 5",
		"Saves records, in [Store.swift](../App/Model/Store.swift).", "Keeps records, in [Store.swift](../App/Model/Store.swift).",
		"Lists records, in [HomeScreen.swift](../App/Screens/HomeScreen.swift).", "The store saves.", "The screen makes a store.")
	return files
}

// cvRun runs `fitness-install -- architecture-code` on happyRepo with files
// written over it.
func cvRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, cvCheck)
}

// cvEdit replaces old with new in one file of files.
func cvEdit(files map[string]string, file, old, new string) map[string]string {
	files[file] = strings.Replace(files[file], old, new, 1)
	return files
}

// cvSink adds a Go interface and a type that satisfies it, drawn as a
// realization on line 12.
func cvSink(put string) map[string]string {
	files := cvFixture()
	files[cvStore] += "\n// Sink takes records.\ntype Sink interface {\n\tPut()\n}\n\n// Disk writes records.\ntype Disk struct{}\n\n// " + put + " writes one.\nfunc (Disk) " + put + "() {}\n"
	files[cvView] = cvCode("\nclass Sink[\"4 store.Sink\"]\nclass Disk[\"5 store.Disk\"]\nDisk ..|> Sink : 6", "Takes records, "+cvLink, "Writes records, "+cvLink, "Disk is a sink.")
	return files
}

const (
	cvTSStore = "src/store/store.ts"
	cvTSDisk  = "src/store/disk.ts"
	cvTSView  = "src/store/panel.tsx"
)

// cvScript is a TypeScript repository matching its code view: store members
// on lines 10 and 11, subclass 13, edges 18 to 20.
func cvScript() map[string]string {
	files := acFixture(acRoot)
	files[cvTSStore] = "import { Sink } from './sink';\n\n// Saving writes records.\nexport interface Saving {\n  save(id: string): boolean;\n}\n\n" +
		"// Store keeps records.\nexport class Store implements Saving {\n  readonly records: Map<string, number> = new Map();\n  private label = 'Sink'; // Sink names nothing here\n\n" +
		"  save(id: string): boolean {\n    return this.records.has(id);\n  }\n}\n"
	files[cvTSDisk] = "import { Store } from './store';\n\nexport class DiskStore\n  extends Store {\n  flush(force?: boolean): Promise<void> {\n    return write(this.records);\n  }\n}\n\n" +
		"export function write(records: Map<string, number>): Promise<void> {\n  return Promise.resolve();\n}\n"
	files["src/store/sink.ts"] = "class Sink {}\n\nexport { Sink };\n"
	files[cvTSView] = cvPanel
	files[cvView] = acView("classDiagram\nclass Saving[\"1 Saving\"] {\n<<interface>>\n+save(id: string) boolean\n}\nclass Store[\"2 Store\"] {\n+records: Map~string, number~\n+save(string) boolean\n}\n"+
		"class DiskStore[\"3 DiskStore\"] {\n+flush(force: boolean) Promise~void~\n}\nclass write[\"4 write(Map~string, number~) Promise~void~\"]\nclass Sink[\"5 Sink\"]\n"+
		"Store ..|> Saving : 6\nStore <|-- DiskStore : 7\nDiskStore --> write : 8\n"+
		"class Handler[\"9 Handler\"]\nclass Panel[\"10 Panel\"] {\n+host\n+rows: Array~string~[]\n+open(number, string) void\n}\nHandler --> Store : 11",
		"Saves records, in [store.ts](../src/store/store.ts).", "Keeps records, in [store.ts](../src/store/store.ts).",
		"Writes to disk, in [disk.ts](../src/store/disk.ts).", "Writes records, in [disk.ts](../src/store/disk.ts).",
		"Takes records, in [sink.ts](../src/store/sink.ts).", "The store saves.", "A disk store is a store.", "Flushing writes.",
		"Takes a store, in [panel.tsx](../src/store/panel.tsx).", "Shows records, in [panel.tsx](../src/store/panel.tsx).", "A handler reads the store.")
	return files
}

// cvPanel is TSX without semicolons: a Store alias, braces in closing tags,
// and a class after ! or ++ with a < default.
const cvPanel = "import { Store } from './store'\n\nexport type Handler = (s: Store) => void\n\nlet renders = 0\n\n" +
	"export function Title(props: { show: boolean }) {\n  return <h1>{props.show && <b>Records</b>}</h1>\n}\n\n" +
	"const root = document.getElementById('root')!\nrenders++\nexport class Panel {\n  host = root!\n  rows: Array<string>[] = []\n\n" +
	"  open(limit: number = renders < 10 ? 1 : 2, label: string): void {\n    renders += limit\n  }\n}\n"

// cvJavaDir holds the Java sources of cvJava.
const cvJavaDir = "src/main/java/shop/"

// cvJava is a Java repository matching its code view: store members 10 and
// 11, record members 19 and 20, edges 22 to 25.
func cvJava() map[string]string {
	files := acFixture(acRoot)
	files[cvJavaDir+"Saving.java"] = "package shop;\n\n// Saving writes records.\npublic interface Saving {\n  boolean save(String id);\n}\n"
	files[cvJavaDir+"Store.java"] = "package shop;\n\nimport java.util.Map;\n\n/** Store keeps records; Sink names nothing here. */\npublic class Store implements Saving {\n  final Map<String, Integer> records = Map.of(), spare = null;\n  private String label = \"Sink\";\n\n  @Override\n  public boolean save(final String id) {\n    return records.containsKey(id);\n  }\n\n  public record Entry(String key, int count) {}\n}\n"
	files[cvJavaDir+"DiskStore.java"] = "package shop;\n\npublic class DiskStore\n    extends Store {\n  public void flush(boolean force) {\n    Writer.write(records);\n  }\n}\n"
	files[cvJavaDir+"Writer.java"] = "package shop;\n\npublic final class Writer {\n  static void write(java.util.Map<String, Integer> records) {}\n}\n"
	files[cvJavaDir+"Sink.java"] = "package shop;\n\nclass Sink {}\n"
	files["src/test/java/shop/StoreTest.java"] = "package shop;\n\npublic class StoreTest {}\n"
	files[cvView] = acView("classDiagram\nclass Saving[\"1 Saving\"] {\n<<interface>>\n+save(String id) boolean\n}\nclass Store[\"2 shop.Store\"] {\n+records: Map~String, Integer~\n+save(String) boolean\n}\n"+
		"class DiskStore[\"3 DiskStore\"] {\n+flush(boolean) void\n}\nclass Writer[\"4 Writer\"]\nclass write[\"5 Writer.write(java.util.Map~String, Integer~) void\"]\nclass Entry[\"6 Store.Entry\"] {\n+key: String\n+count: int\n}\n"+
		"Store ..|> Saving : 7\nStore <|-- DiskStore : 8\nDiskStore --> write : 9\nDiskStore --> Writer : 10",
		"Saves records, in [Saving.java](../"+cvJavaDir+"Saving.java).", "Keeps records, in [Store.java](../"+cvJavaDir+"Store.java).",
		"Writes to disk, in [DiskStore.java](../"+cvJavaDir+"DiskStore.java).", "Writes, in [Writer.java](../"+cvJavaDir+"Writer.java).",
		"Writes records, in [Writer.java](../"+cvJavaDir+"Writer.java).", "One entry, in [Store.java](../"+cvJavaDir+"Store.java).",
		"The store saves.", "A disk store is a store.", "Flushing writes.", "Flushing names the writer.")
	return files
}
