// Package requirements proves the localized-text check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example app with a String Catalog.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

// catalogFix is how a failure tells a literal with no catalog entry to be
// added and read through a Words reference.
const catalogFix = ` skips the String Catalog; add it to the String Catalog and use a Words reference`

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a literal in a text view fails naming the call and the entry", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{
			"App/ShareView.swift": localizedView,
			"App/StateView.swift": localizedStateView,
		})
		acceptance.Sees(t, out, code, 1,
			`App/ShareView.swift:7 - "Share" passed to Button skips the String Catalog; use a Words reference for catalog entry "verb.share"`,
			`App/ShareView.swift:9 - "\(count) shared" passed to Text`+catalogFix,
			`App/StateView.swift:11 - "Photos" passed to Label skips the String Catalog; use a Words reference for catalog entry "noun.photo"`)
	})
	t.Run("0001.1.2 a literal in a modifier fails naming it", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/ShareView.swift": localizedView})
		acceptance.Sees(t, out, code, 1, `App/ShareView.swift:8 - "Library" passed to navigationTitle`+catalogFix)
	})
	t.Run("0001.1.3 a literal in a ternary or ?? branch fails", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/StateView.swift": localizedStateView})
		acceptance.Sees(t, out, code, 1,
			`App/StateView.swift:8 - "Shared" passed to Text`+catalogFix,
			`App/StateView.swift:8 - "Private" passed to Text`+catalogFix,
			`App/StateView.swift:9 - "Nobody" passed to Text`+catalogFix)
	})
	t.Run("0001.1.4 a literal in an interpolated ternary fails", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/CountView.swift": localizedCountView})
		acceptance.Sees(t, out, code, 1,
			`App/CountView.swift:7 - "photo" passed to Text`+catalogFix,
			`App/CountView.swift:7 - "photos" passed to Text`+catalogFix)
	})
	t.Run("0001.1.5 a multi-line literal is named whole", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/HelpView.swift": "import SwiftUI\n\nstruct HelpView: View {\n    var body: some View {\n        Text(\"\"\"\n            Share a \"photo\"\n            with a friend\n            \"\"\")\n    }\n}\n"})
		acceptance.Sees(t, out, code, 1,
			"App/HelpView.swift:5 - \"\"\"\n            Share a \"photo\"\n            with a friend\n            \"\"\" passed to Text"+catalogFix)
	})
	t.Run("0001.1.6 references and values pass", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/PhotoView.swift": `import SwiftUI

struct PhotoView: View {
    let caption: String
    let count: Int
    let titles: [String: String]

    var body: some View {
        VStack {
            Label(Words.Noun.photo, systemImage: "photo")
            Text(caption)
            Text("·")
            Text("\(count)")
            Text("\(titles["first"] ?? Words.Noun.photo)")
            Button(Words.Verb.share) {}
                .accessibilityLabel(Words.Verb.share)
        }
    }
}
`})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a test typing a catalog label fails naming the reference", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"Tests/Requirements/ShareTests.swift": `import Testing

struct ShareTests {
    func tapsShare(_ user: User) async throws {
        try await user.taps(.button("Share"))
    }
}
`})
		acceptance.Sees(t, out, code, 1,
			`Tests/Requirements/ShareTests.swift:5 - "Share" types catalog entry "verb.share"; read it through a Words reference`)
	})
	t.Run("0001.2.2 captions and prices and plain data pass", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"Tests/Requirements/Support/Guest.swift": `struct Guest {
    let user: User

    func sharesCaptionedPhoto() async throws {
        try await user.writes("Swim at the lake")
        try await user.sees(.text("$9.99"))
        try await user.sees(.image("Sparkle"))
        try await user.doesNotSee(.button("Archive"))
        try await user.taps(.button(Words.Verb.share))
        let label = "Share"
        try await user.writes(label)
        try await user.adds(Photo(title: "Photos", caption: "Share"))
    }
}
`})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 a literal an allow pattern matches passes", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{
			".fitnessrc.json": `{"localizedText": {"reference": "Words", "allow": [
  {"pattern": "Beta [0-9]+", "reason": "build tags stay in English for testers"}
]}}` + "\n",
			"App/BuildTag.swift": "import SwiftUI\n\nstruct BuildTag: View {\n    var body: some View { Text(\"Beta 3\") }\n}\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.3.2 an allow pattern without a reason fails", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{
			".fitnessrc.json": `{"localizedText": {"reference": "Words", "allow": [
  {"pattern": "Beta [0-9]+"},
  {"pattern": "(Beta", "reason": "build tags stay in English for testers"}
]}}` + "\n",
		})
		acceptance.Sees(t, out, code, 1,
			`localizedText.allow pattern "Beta [0-9]+" needs a reason`,
			`localizedText.allow pattern "(Beta" is not a regular expression`)
	})
	t.Run("0001.3.3 an app without a String Catalog passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{
			"App/Localizable.xcstrings": "",
			"App/ShareView.swift":       localizedView,
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "localized-text")
	})
}
