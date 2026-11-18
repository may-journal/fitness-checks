# 0020 Checks Swift Storybook Requirements

## Why

Every Storybook page in my Swift app is named by the requirement acceptance it shows, and every acceptance has a page. So the catalog can badge new and changed pages from that acceptance's history.

## Measurement

Storybook pages named by a defined acceptance
-
Storybook pages in a Swift app with requirement docs

## Requirements

- 0020.1
    - Given Developer has previews and book pages each named by a defined acceptance, two sharing one ID under different labels
        - When swift-storybook-requirements runs
            - Then it passes
- 0020.2
    - Given Developer has a preview with no name, and a book page whose only title belongs to a call inside it
        - When swift-storybook-requirements runs
            - Then it fails naming each page's file and line, and a name to use
- 0020.3
    - Given Developer has a preview and a book page whose names do not start with an acceptance ID
        - When swift-storybook-requirements runs
            - Then it fails naming each page's file, line, and name
- 0020.4
    - Given Developer has a page named by an ID no requirement doc defines
        - When swift-storybook-requirements runs
            - Then it fails naming that ID and the docs folder
- 0020.5
    - Given Developer has two pages with the same full name
        - When swift-storybook-requirements runs
            - Then it fails the second, naming where the first is
- 0020.6
    - Given Developer has a book page whose title is a variable, and a preview whose name interpolates
        - When swift-storybook-requirements runs
            - Then it fails naming each name as written
- 0020.7
    - Given Developer has unnamed previews only in comments and strings, no Swift files, or a Swift app with no requirement docs
        - When swift-storybook-requirements runs
            - Then it passes
- 0020.8
    - Given Developer has a Swift app with an acceptance that no page names, even with no pages at all
        - When swift-storybook-requirements runs
            - Then it fails naming the doc line that declares it, and a page name to add
