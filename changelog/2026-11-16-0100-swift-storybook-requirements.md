### 2026.11.16.0100

> Every Storybook page in a Swift app is named by the acceptance it shows. So the catalog can badge new and changed pages.

- Feat: `swift-storybook-requirements` fails a `#Preview` or `BookPage` with no name, or one not named by a defined acceptance ID.
- Feat: it fails two pages with the same full name, and a name held in a variable or built by interpolation.
- Fix: it reads every Swift file in a scoped run, and only a `BookPage` call's own `title:` argument.
- Test: requirements 0168.1 to 0168.7 run it on Swift pages, comments, strings and a repo without requirement docs.
- Test: the catalog and the size baseline list it, and its README is the fixed stub.
