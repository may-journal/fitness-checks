# local-branch

Runs PR and Plan judgments locally. It runs the existing PR and Plan checks against this branch's live GitHub metadata, without posting comments or changing issues.

## Usage

```
fitness local-branch
fitness branch-check
fitness --policy=external --checks=local-branch
```

Full local Fitness checks the current branch. Pre-push runs one meta-check over each live pushed destination branch and local tip, including branches outside the checkout. Deletions and tags skip; pre-commit excludes it. In Actions or a repo without a GitHub origin, it passes without network reads.

## Behavior

The check matches the current branch's open PR, including drafts and forks, then validates its closing issues and branch Plan trailers. History starts at the fetched base branch when it contains the API base commit, else at that commit. Upstream Plans and entries merged since never count. Ordinary issues receive PR checklist rules.

No PR or detached HEAD skips. Ambiguous PRs, unfinished checklists, inaccessible metadata, and missing base history fail with retry guidance.

GitHub reads share a 30-second deadline, and the runner allows 35 seconds for cancellation and reporting. Inside pre-push, the gate opens the push's deadline and hands over the PR and Plans it read. This check reads them no more and spends what is left on the rest. When the runner's probe finds GitHub's API host unreachable, the check reads nothing and reports info naming the host and proxy.

## Implementation

The judgments live in the engine's internal `metadata` package, which no check imports. This program declares the check, with its describe metadata, and fitness answers for it with those judgments, even when this program is installed. Run alone, the program could not judge, so it reports that and exits 3.

Only Actions applies GitHub effects. The meta-check calls judgments directly, so it cannot recursively run the suite.

The check's place in the system is in the [architecture](../../docs/architecture/03-components.md).
