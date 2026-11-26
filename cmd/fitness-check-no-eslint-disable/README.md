# no-eslint-disable

Fails any JS or TS source file that holds an ESLint disable directive, in file, block, `-line`, or `-next-line` form. This keeps lint hygiene on the fix-the-rule path instead of accumulating suppressions. It is on by default.

## Behavior

[Requirement 0001](docs/requirements/0001-source-holds-no-eslint-disable.md) states what this check does, and a test case named for each step proves it.
