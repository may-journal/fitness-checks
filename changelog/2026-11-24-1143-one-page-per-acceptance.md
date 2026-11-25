### 2026.11.24.1143

> Each requirement acceptance has one Storybook page by default, so its variants sit as labeled rows on that page.

- Feat: `swift-storybook-requirements` fails a page whose acceptance ID an earlier page already names, whatever its label.
- Feat: `swiftStorybookRequirements.maxPagesPerAcceptance` in `.fitnessrc.json` raises the limit, by default 1.
- Docs: the README states the one-to-one rule and drops "Two pages may share an ID under different labels".
- Test: requirement 0006.8 proves one page, two pages, a preview beside a `BookPage`, and a limit of 2.
