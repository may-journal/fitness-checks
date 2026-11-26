# readme-behavior

Keeps a README's Behavior section to one link to its own requirements. A check README must also show its full API in an Example section.

## What it reads

A README with a `## Behavior` section holds exactly one Markdown link there and nothing else. The link points at that README's own `docs/requirements`, as a folder or a doc in it, and the target must exist. A failure names the file and the line.

A README holds no GIVEN, WHEN or THEN list item in any case outside a code fence; a step belongs in `docs/requirements`.

The README of a check program, at `cmd/fitness-check-<name>/README.md`, needs a `## Example` section too. The Example shows the command `fitness <name>` and each flag the program declares. It also sets each config key its `--describe` declares to its default in a `.fitnessrc.json` snippet.

The check reads flags and keys from the `fitness-check-<name> --describe` output of the program that sits beside it. It sees `--body-file` when `--describe` says `blob`, and the inline flag that `--describe` names, such as `--message`. It cannot see a flag the program parses and never declares, such as `--fix`. When it finds no program to describe, it judges only that the Example has the command, and says so as info.

## Behavior

[Requirement 0001](docs/requirements/0001-readmes-show-their-requirements-and-api.md)

## Example

It takes no flag and no config key.

```
fitness readme-behavior
```
