### 2026.11.19.2300

> Only a Swift app that uses Storybook needs a page for each acceptance, so a library with no Storybook passes.

- Fix: `swift-storybook-requirements` judges a repo only when a Swift file imports a Storybook module or declares a `BookPage`.
- Fix: any other repo passes with a "no Storybook in this repo" note; an import in a comment or string is not counted.
- Feat: `Storybook`, `StorybookSandbox`, `StorybookKit` and `SwiftUIStorybook` count, after any attribute or a semicolon.
- Test: requirement 0006.8 passes a repo without Storybook and fails a Storybook app with no pages.
