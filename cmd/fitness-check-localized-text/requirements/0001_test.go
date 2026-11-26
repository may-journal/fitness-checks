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
	views := func(t *testing.T) (string, int) {
		return localizedText(t, map[string]string{
			"App/ShareView.swift": localizedView,
			"App/StateView.swift": localizedStateView,
		})
	}
	t.Run("0001.1.2.1.1 it fails naming `file:line`", func(t *testing.T) {
		t.Parallel()
		out, code := views(t)
		acceptance.Sees(t, out, code, 1, "App/ShareView.swift:7 - ", "App/ShareView.swift:9 - ", "App/StateView.swift:11 - ")
	})
	t.Run("0001.1.2.1.2 it names the literal plus the call", func(t *testing.T) {
		t.Parallel()
		out, code := views(t)
		acceptance.Sees(t, out, code, 1, `"Share" passed to Button`, `"\(count) shared" passed to Text`, `"Photos" passed to Label`)
	})
	t.Run("0001.1.2.1.3 it names the entry or reference to use", func(t *testing.T) {
		t.Parallel()
		out, code := views(t)
		acceptance.Sees(t, out, code, 1,
			`use a Words reference for catalog entry "verb.share"`,
			`use a Words reference for catalog entry "noun.photo"`,
			catalogFix)
	})
	t.Run("0001.1.2.2.1 it fails naming that modifier", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/ShareView.swift": localizedView})
		acceptance.Sees(t, out, code, 1, `App/ShareView.swift:8 - "Library" passed to navigationTitle`+catalogFix)
	})
	t.Run("0001.1.2.3.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/StateView.swift": localizedStateView})
		acceptance.Sees(t, out, code, 1,
			`App/StateView.swift:8 - "Shared" passed to Text`+catalogFix,
			`App/StateView.swift:8 - "Private" passed to Text`+catalogFix,
			`App/StateView.swift:9 - "Nobody" passed to Text`+catalogFix)
	})
	t.Run("0001.1.2.4.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/CountView.swift": localizedCountView})
		acceptance.Sees(t, out, code, 1,
			`App/CountView.swift:7 - "photo" passed to Text`+catalogFix,
			`App/CountView.swift:7 - "photos" passed to Text`+catalogFix)
	})
	t.Run("0001.1.2.4.2 a reference in the interpolation passes", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/CountView.swift": "import SwiftUI\n\nstruct CountView: View {\n    let count: Int\n\n    var body: some View {\n        Text(\"\\(count)\\(Words.Noun.photo)\")\n    }\n}\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.5.1 it is named whole", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/HelpView.swift": "import SwiftUI\n\nstruct HelpView: View {\n    var body: some View {\n        Text(\"\"\"\n            Share a \"photo\"\n            with a friend\n            \"\"\")\n    }\n}\n"})
		acceptance.Sees(t, out, code, 1,
			"App/HelpView.swift:5 - \"\"\"\n            Share a \"photo\"\n            with a friend\n            \"\"\" passed to Text"+catalogFix)
	})
	t.Run("0001.1.2.6.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/PhotoView.swift": `import SwiftUI

struct PhotoView: View {
    let caption: String
    let titles: [String: String]

    var body: some View {
        VStack {
            Label(Words.Noun.photo, systemImage: "photo")
            Text(caption)
            Text("\(titles["first"] ?? Words.Noun.photo)")
            Button(Words.Verb.share) {}
                .accessibilityLabel(Words.Verb.share)
        }
    }
}
`})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.6.2 a literal without letters passes", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{"App/DotView.swift": `import SwiftUI

struct DotView: View {
    let count: Int

    var body: some View {
        VStack {
            Text("·")
            Text("\(count)")
        }
    }
}
`})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	typed := func(t *testing.T) (string, int) {
		return localizedText(t, map[string]string{"Tests/Requirements/ShareTests.swift": `import Testing

struct ShareTests {
    func tapsShare(_ user: User) async throws {
        try await user.taps(.button("Share"))
    }
}
`})
	}
	t.Run("0001.2.2.1.1 it fails naming `file:line`", func(t *testing.T) {
		t.Parallel()
		out, code := typed(t)
		acceptance.Sees(t, out, code, 1, "Tests/Requirements/ShareTests.swift:5 - ")
	})
	t.Run("0001.2.2.1.2 it names the literal plus the entry", func(t *testing.T) {
		t.Parallel()
		out, code := typed(t)
		acceptance.Sees(t, out, code, 1, `"Share" types catalog entry "verb.share"`)
	})
	t.Run("0001.2.2.1.3 it names the reference to read it through", func(t *testing.T) {
		t.Parallel()
		out, code := typed(t)
		acceptance.Sees(t, out, code, 1, `; read it through a Words reference`)
	})
	guest := func(t *testing.T, lines string) (string, int) {
		return localizedText(t, map[string]string{"Tests/Requirements/Support/Guest.swift": "struct Guest {\n    let user: User\n\n    func run() async throws {\n" + lines + "    }\n}\n"})
	}
	t.Run("0001.2.2.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := guest(t, `        try await user.writes("Swim at the lake")
        try await user.sees(.image("Sparkle"))
        try await user.taps(.button(Words.Verb.share))
        try await user.adds(Photo(title: "Photos", caption: "Share"))
`)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2.2.2 a price passes alike", func(t *testing.T) {
		t.Parallel()
		out, code := guest(t, "        try await user.sees(.text(\"$9.99\"))\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2.2.3 a literal assigned with `=` passes alike", func(t *testing.T) {
		t.Parallel()
		out, code := guest(t, "        let label = \"Share\"\n        try await user.writes(label)\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2.2.4 a stale entry passes alike", func(t *testing.T) {
		t.Parallel()
		out, code := guest(t, "        try await user.doesNotSee(.button(\"Archive\"))\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.2.1.1 that literal passes", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{
			".fitnessrc.json": `{"localizedText": {"reference": "Words", "allow": [
  {"pattern": "Beta [0-9]+", "reason": "build tags stay in English for testers"}
]}}` + "\n",
			"App/BuildTag.swift": "import SwiftUI\n\nstruct BuildTag: View {\n    var body: some View { Text(\"Beta 3\") }\n}\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	badAllow := func(t *testing.T) (string, int) {
		return localizedText(t, map[string]string{
			".fitnessrc.json": `{"localizedText": {"reference": "Words", "allow": [
  {"pattern": "Beta [0-9]+"},
  {"pattern": "(Beta", "reason": "build tags stay in English for testers"}
]}}` + "\n",
		})
	}
	t.Run("0001.3.2.2.1 the check fails naming the pattern", func(t *testing.T) {
		t.Parallel()
		out, code := badAllow(t)
		acceptance.Sees(t, out, code, 1, `localizedText.allow pattern "Beta [0-9]+" needs a reason`)
	})
	t.Run("0001.3.2.2.2 a pattern that does not compile fails alike", func(t *testing.T) {
		t.Parallel()
		out, code := badAllow(t)
		acceptance.Sees(t, out, code, 1, `localizedText.allow pattern "(Beta" is not a regular expression`)
	})
	t.Run("0001.3.2.3.1 the check passes with zero files", func(t *testing.T) {
		t.Parallel()
		out, code := localizedText(t, map[string]string{
			"App/Localizable.xcstrings": "",
			"App/ShareView.swift":       localizedView,
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "localized-text")
	})
}
