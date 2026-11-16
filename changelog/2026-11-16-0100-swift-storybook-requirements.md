### 2026.11.16.0100

> Every Storybook page in a Swift app is named by the acceptance it shows, and every acceptance has a page. So the catalog can badge new and changed pages.

- Feat: `swift-storybook-requirements` fails a `#Preview` or `BookPage` with no name, or one not named by a defined acceptance ID.
- Feat: it fails each acceptance no page names, at its doc line, and two pages sharing a full name.
- Fix: it fails a name held in a variable or interpolated, reads every Swift file, and only a `BookPage`'s own `title:`.
- Test: requirements 0168.1 to 0168.8 run it on Swift pages, comments, strings, missing pages, and repos without Swift or docs.
- Test: the catalog and the size baseline list it, and its README is the fixed stub.
