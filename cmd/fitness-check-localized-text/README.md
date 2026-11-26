# localized-text

Fails SwiftUI text and requirement test literals that skip the app's [String Catalog](https://developer.apple.com/documentation/xcode/localizing-and-varying-text-with-a-string-catalog) ([#527](https://github.com/may-journal/fitness-runner/issues/527)). The Swift language module owns the reading and its judgment. A literal passed to `Text` compiles as a [LocalizedStringKey](https://developer.apple.com/documentation/swiftui/localizedstringkey), so a missing entry only shows in another language.

## Behavior

[Requirement 0001](docs/requirements/0001-labels-read-the-string-catalog.md) states what this check does, and a test case named for each step proves it.

An app opts in by tracking a `.xcstrings` file, and may name its reference type and allow literals. Run `fitness explain localized-text` for the keys:

```
{
  "localizedText": {
    "reference": "Words",
    "allow": [{ "pattern": "Beta [0-9]+", "reason": "build tags stay in English for testers" }]
  }
}
```

Org repos cannot turn the check off, as [ADR 0009](https://github.com/may-journal/fitness-runner/blob/main/fitness/adr-0009/README.md) decides.
