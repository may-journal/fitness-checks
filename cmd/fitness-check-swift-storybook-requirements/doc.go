// Command fitness-check-swift-storybook-requirements fails each Swift
// Storybook page not named by a requirement acceptance, and each acceptance
// with no page ([#79]). The Storybook package in common-swift reads each
// page's acceptance to badge new and changed pages from that acceptance's
// git history ([common-swift #53]). A page named any other way, or an
// acceptance with no page, would drop out unnoticed.
//
// # Pages
//
//   - A `#Preview` macro is a page, and its first argument is its name.
//   - A `BookPage` call is a page, and its own `title:` argument is its name.
//   - A page in a comment or a string is not read.
//
// # Names
//
// A name starts with one acceptance ID, alone or followed by a space and a
// label, such as `#Preview("0006.3 unlocked")`. The check fails:
//
//   - a page with no name, such as `#Preview { … }`;
//   - a name that does not start with an acceptance ID;
//   - an ID that no doc in `docs/requirements` defines;
//   - a name held in a variable or built by interpolation;
//   - a page whose full name repeats an earlier page's, so a `--storybook`
//     launch opens exactly one page.
//
// Two pages may share an ID under different labels. A failure reads:
//
//	App/BadgeView.swift:12 - page "day" does not start with an acceptance ID; name it such as "0001.1 day"
//
// # Acceptances
//
// Every acceptance a doc in `docs/requirements` declares needs a page whose
// plain literal name starts with its ID. This holds even in an app with no
// pages yet. The failure names the doc line that declares it:
//
//	docs/requirements/0006-lock.md:17 - acceptance 0006.3 has no Storybook page; add one named "0006.3 <label>"
//
// A repo with no Swift files, or no requirement docs, passes with zero files.
// Every Swift file is read in every run, since names repeat and pages go
// missing across files.
//
// # Apart from other checks
//
// This check owns only the page names and the pages each acceptance needs.
// It shares helpers, the Swift lexer and the requirement doc reader, and no
// judgment with any other check.
//
// # Requirements
//
// [Requirement 0020] has one installed-command test per acceptance.
//
// [#79]: https://github.com/may-journal/fitness-checks/issues/79
// [common-swift #53]: https://github.com/may-journal/common-swift/issues/53
// [Requirement 0020]: https://github.com/may-journal/fitness-checks/blob/main/docs/requirements/0020-checks-swift-storybook-requirements.md
package main

import _ "embed"

// docSource is this file, whose comment --describe summarizes.
//
//go:embed doc.go
var docSource string
