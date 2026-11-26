# swift-storybook-requirements

Fails each Swift Storybook page not named by a requirement acceptance, and each acceptance without exactly one page ([#79](https://github.com/may-journal/fitness-checks/issues/79), [#129](https://github.com/may-journal/fitness-checks/issues/129)). The Storybook package in common-swift reads each page's acceptance to badge new and changed pages from that acceptance's git history ([common-swift #53](https://github.com/may-journal/common-swift/issues/53)). A page named any other way, or an acceptance with no page, would drop out unnoticed.

## Behavior

[Requirement 0001](docs/requirements/0001-each-acceptance-has-one-storybook-page.md) states what this check does, and a test case named for each step proves it.

## Storybook apps

The modules are common-swift's `Storybook` and `StorybookSandbox`, swift-storybook's `StorybookKit`, and `SwiftUIStorybook`. A plain `#Preview` needs only SwiftUI, so it marks no Storybook app. A manifest may name Storybook only as a dependency of a dependency.

## Pages

A `#Preview` macro is a page named by its first argument. A `BookPage` call is a page named by its own `title:` argument. A name starts with one acceptance ID, alone or followed by a space and a label, such as `#Preview("0006.3 unlocked")`. A page's variants are labeled rows inside that page.

A failure reads:

```
App/BadgeView.swift:12 - page "day" does not start with an acceptance ID; name it such as "0001.1 day"
MosaicShareView+Previews.swift:77 - page "0010.4 month" repeats acceptance 0010.4 from line 73; draw each variant as a labeled row on one page
docs/requirements/0006-lock.md:17 - acceptance 0006.3 has no Storybook page; add one named "0006.3 <label>"
```

A repo can allow more pages per acceptance in `.fitnessrc.json`. Run `fitness explain swift-storybook-requirements` for the key and its default.

## Apart from other checks

This check owns only the page names and the pages each acceptance needs. It shares helpers, the Swift lexer and the requirement doc reader, and no judgment with any other check.
