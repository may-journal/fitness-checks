# localized-text

Fails SwiftUI text and requirement test literals that skip the app's [String Catalog](https://developer.apple.com/documentation/xcode/localizing-and-varying-text-with-a-string-catalog) ([#527](https://github.com/may-journal/fitness-runner/issues/527)). The Swift language module owns the reading and its judgment. A literal passed to `Text` compiles as a [LocalizedStringKey](https://developer.apple.com/documentation/swiftui/localizedstringkey), so a missing entry only shows in another language.

## Behavior

[Requirement 0001](docs/requirements/0001-labels-read-the-string-catalog.md)

## Notes

An app opts in by tracking a `.xcstrings` file, and may name its reference type and allow literals.

Org repos cannot turn the check off, as [ADR 0009](https://github.com/may-journal/fitness-runner/blob/main/fitness/adr-0009/README.md) decides.

## Example

```
fitness localized-text
```

Run `fitness explain localized-text` for what each key does. This sets each key to its default:

```json
{
  "localizedText": {
    "reference": "",
    "allow": []
  }
}
```

