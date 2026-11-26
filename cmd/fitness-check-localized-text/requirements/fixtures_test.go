package requirements

import (
	"maps"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// localizedCatalog is the fixture app's String Catalog: a live verb, a plural
// noun, a stale word, and Xcode's empty key.
const localizedCatalog = `{
  "sourceLanguage" : "en",
  "strings" : {
    "" : { "shouldTranslate" : false },
    "noun.photo" : {
      "localizations" : { "en" : { "variations" : { "plural" : {
        "one" : { "stringUnit" : { "state" : "translated", "value" : "Photo" } },
        "other" : { "stringUnit" : { "state" : "translated", "value" : "Photos" } }
      } } } }
    },
    "verb.archive" : {
      "extractionState" : "stale",
      "localizations" : { "en" : { "stringUnit" : { "state" : "translated", "value" : "Archive" } } }
    },
    "verb.share" : {
      "extractionState" : "extracted_with_value",
      "localizations" : { "en" : { "stringUnit" : { "state" : "new", "value" : "Share" } } }
    }
  },
  "version" : "1.0"
}
`

// localizedWords is the fixture app's reference type, which reads each
// label from the catalog by key.
const localizedWords = `import Foundation

enum Words {
    enum Verb {
        static var share: String {
            String(localized: "verb.share", defaultValue: "Share", comment: "Verb: share a photo")
        }
    }

    enum Noun {
        static var photo: String {
            String(localized: "noun.photo", defaultValue: "Photo", comment: "Noun: one photo")
        }
    }
}
`

// localizedApp is the fixture app with a String Catalog, its Words
// reference, and a config naming that reference, with files written over
// it. An empty body deletes a fixture file.
func localizedApp(t *testing.T, files map[string]string) string {
	t.Helper()
	app := map[string]string{
		".fitnessrc.json":           `{"localizedText": {"reference": "Words"}}` + "\n",
		"App/Localizable.xcstrings": localizedCatalog,
		"App/Words.swift":           localizedWords,
	}
	maps.Copy(app, files)
	return acceptance.Example(t, "happyRepo", app)
}

// localizedText runs `fitness-install -- localized-text` on the fixture
// app with files written over it.
func localizedText(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, localizedApp(t, files), nil, "localized-text")
}

// localizedView is a view that passes literals to a button and a title.
const localizedView = `import SwiftUI

struct ShareView: View {
    let count: Int

    var body: some View {
        Button("Share") {}
            .navigationTitle("Library")
        Text("\(count) shared")
    }
}
`

// localizedStateView passes literals as ternary and ?? branches.
const localizedStateView = `import SwiftUI

struct StateView: View {
    let isShared: Bool
    let name: String?

    var body: some View {
        Text(isShared ? "Shared" : "Private")
        Text(name ?? "Nobody")
        Text(isShared ? Words.Verb.share : Words.Noun.photo)
        Label("Photos", systemImage: isShared ? "person" : "lock")
    }
}
`

// localizedCountView passes literals as branches inside an interpolation.
const localizedCountView = `import SwiftUI

struct CountView: View {
    let count: Int

    var body: some View {
        Text("\(count == 1 ? "photo" : "photos")")
        Text("\(count)\(Words.Noun.photo)")
    }
}
`
