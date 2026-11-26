# node-version

Requires the machine's active Node version to satisfy the repo's `.nvmrc`. It only compares; run `nvm use` yourself if needed. It asks `node --version` for the current version. A missing or broken node binary reports info with an install hint, since no version was read to judge.

## Behavior

[Requirement 0001](docs/requirements/0001-node-matches-the-repo-pin.md) states what this check does, and a test case named for each step proves it.
