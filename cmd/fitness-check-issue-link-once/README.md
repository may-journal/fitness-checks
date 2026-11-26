# issue-link-once

Keeps each Issue's timeline readable ([#153](https://github.com/may-journal/fitness-runner/issues/153)). Every commit that links an Issue adds an entry there, so a branch links an Issue from one commit at most.

## Behavior

[Requirement 0001](docs/requirements/0001-each-issue-is-linked-once.md) states what this check does, and a test case named for each step proves it.

A link is `#NN`, `owner/repo#NN`, or an Issue or PR URL on github.com. A bare `#NN` belongs to origin's repo. An earlier commit counts only when no other pushed branch holds it. A branch stacked on a pushed epic branch starts where that branch ends.

## Where it runs

The `commit-msg` hook judges the proposed message, which it passes as the context-inline `--message` value. An amend leaves out HEAD, the commit it replaces. With no open PR, or no answer from gh, it warns and skips the PR rule. When GitHub's API host misses the runner's probe, it skips only the open PR read, with a note naming the host and proxy.

## On a pull request

The org `fitness-suite` workflow judges every branch commit, oldest first, in every repo ([#157](https://github.com/may-journal/fitness-runner/issues/157)). The one link written before the PR opened passes, judged by author date, which a rebase keeps. A failure stays red until you rewrite that commit without the link and force-push with `--force-with-lease`.

Anywhere else, such as a local suite run or a push to `main`, it passes inert.
