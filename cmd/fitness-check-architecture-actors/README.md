# architecture-actors

Requires each requirement's Given to start with an actor the system context declares. No role slips into a requirement that the architecture does not know, and a dropped role fails every Given that still names it.

The check reads the model that [architecture-contract](https://pkg.go.dev/github.com/may-journal/fitness-checks/cmd/fitness-check-architecture-contract) reads and the Given text that [requirements](https://pkg.go.dev/github.com/may-journal/fitness-checks/cmd/fitness-check-requirements) parses. It adds no fields, columns, or rows.

## Actors

An actor is a person in `01-system-context.md`: a C4 `Person` or `Person_Ext`. Its name is the label without the callout number, and its Mermaid ID is the one alias the diagram declares.

- `Person(developer, "1 Developer")` declares the name `Developer` and the alias `developer`.
- `Person_Ext(host, "2 Host")` declares the name `Host` and the alias `host`.

The system context is a `C4Context` diagram, as its [template](https://github.com/may-journal/fitness-checks/blob/main/internal/architecture/templates/01-system-context.md) draws it. A flowchart node classed `person` is no actor, so a Given that names one fails until the view becomes C4. Systems, legends, and people in other views are never actors.

## Grammar

A Given starts with one actor, or several joined by `and`, then the state. A spelling matches only as a whole word, with the diagram's case.

- `Given Host has no network` passes.
- `Given host has no network` passes, by ID.
- `Given Host and Participant have joined` passes.
- `Given Hosts have joined` fails: no whole-word actor.
- `Given Admin has joined` fails: undeclared.
- `Given I have joined` fails unless `I` is an actor's name or ID.
- `Given no network` fails: no actor.
- `Given Host and` fails: no actor after `and`.

Not every actor needs a Given; the check constrains the roles requirements use.

## Failures

- No actor at the start, or after `and`: the Given line, the context diagram line, and the declared names.
- A spelling that is one actor's ID and another actor's name: the Given line and both context lines.

## Scope

A test, changelog or doc change cannot move the views. A scoped run with no changed view, requirement doc, product or tool source, manifest or `.gitattributes` passes with zero files. Every other run reads every requirement doc and the whole model, so a removed actor fails even when no requirement changed.

## Behavior

- No architecture views: pass with `filesChecked: 0`.
- Otherwise `filesChecked` counts the requirement docs read.
- Doc shape and missing views belong to the requirements and `architecture-contract` checks.

## Requirements

[Requirement 0005](https://github.com/may-journal/fitness-checks/blob/main/docs/requirements/0005-diagrams-match-the-code.md) proves this check with installed-command cases.
