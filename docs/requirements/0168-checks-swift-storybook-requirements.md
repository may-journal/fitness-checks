# 0168 Checks Swift Storybook Requirements

## Why

Every Storybook page in my Swift app is named by the requirement acceptance it shows. So the catalog can badge new and changed pages from that acceptance's history.

## Measurement

Storybook pages named by a defined acceptance
-
Storybook pages in a Swift app with requirement docs

## Requirements

- 0168.1
    - Given Developer has previews and book pages each named by a defined acceptance, two sharing one ID under different labels
        - When swift-storybook-requirements runs
            - Then it passes
- 0168.2
    - Given Developer has a preview with no name, and a book page whose only title belongs to a call inside it
        - When swift-storybook-requirements runs
            - Then it fails naming each page's file and line, and a name to use
- 0168.3
    - Given Developer has a preview and a book page whose names do not start with an acceptance ID
        - When swift-storybook-requirements runs
            - Then it fails naming each page's file, line, and name
- 0168.4
    - Given Developer has a page named by an ID no requirement doc defines
        - When swift-storybook-requirements runs
            - Then it fails naming that ID and the docs folder
- 0168.5
    - Given Developer has two pages with the same full name
        - When swift-storybook-requirements runs
            - Then it fails the second, naming where the first is
- 0168.6
    - Given Developer has a book page whose title is a variable, and a preview whose name interpolates
        - When swift-storybook-requirements runs
            - Then it fails naming each name as written
- 0168.7
    - Given Developer has unnamed previews only in comments and strings, or a Swift app with no requirement docs
        - When swift-storybook-requirements runs
            - Then it passes
