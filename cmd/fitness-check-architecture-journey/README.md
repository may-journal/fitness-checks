# architecture-journey

Keeps the user journey diagram and the requirement docs covering each other. Every requirement doc ID appears in a journey diagram label, and every journey node that promises behavior names one. It reads the shared model from `internal/architecture.Load`.

## Scope

A test, changelog or doc change cannot move the views. A scoped run with no changed view, requirement doc, product or tool source, manifest or `.gitattributes` passes with zero files. Every other run reads the whole journey view and every requirement doc.

## Requirement IDs in the diagram

The journey view (`05-user-journey.md` under `architecture/` or `docs/architecture/`) names each requirement doc in its Mermaid diagram, on the step that fulfils it:

```
flowchart TB
    user((1 User)):::persona
    open(2 Open the app<br/>Req 0001, 0002):::action
    subgraph saving["Saving<br/>Req 0003"]
        done[3 Record saved]:::screen
    end
    user --> open
    open -->|saves 0004| done
```

Labels follow these rules:

- Any four-digit number in a node, subgraph, or edge label is a requirement ID. Write it after the callout number, never first.
- A Mermaid character code such as `#9829;` is not an ID.
- An acceptance such as `0001.1` is not a doc ID. Acceptances keep their one owning test under the `requirements` check.
- One label names at most eight IDs. A step that fulfils more splits into steps, such as one node per check family.
- One ID may sit on several steps. Tables and prose outside the diagram cover nothing.

Personas and layout subgraphs promise no behavior and need no ID. Legends stay out of the model. Every other node, such as a screen, action, or system step, names an ID in its own label or an enclosing subgraph label. An edge label covers its IDs but no node.
## Generated region

A journey view may hold one node per requirement doc, written from the doc's ID and `# NNNN Title` heading. Two Mermaid comments mark the region:

```
%% journey-nodes:start
%% journey-nodes:end
```

Between them, each doc in ID order becomes `reqNNNN[(N Title<br/>Req NNNN)]:::system`, numbered from 2 after the persona. Subgraphs around the region, edges and the callout table stay hand-written.

- A region whose lines differ from the docs fails at the start marker, naming `fitness architecture-journey --fix`.
- `--fix` rewrites each region from the docs and leaves every other line alone.
- A view with no markers stays hand-written and passes, so a repo opts in by adding them.
- A repeated marker, or a start with no end after it, fails at its line.

## Steps that read as another doc

Each doc's file name holds its title, such as `0002-save.md`. A step reads as a doc when its label holds every title word of that doc. The step must also share no title word with any doc it names. Short words, `and`, `the` and a plural `s` are left out of the comparison.

- Two steps that read as each other's docs have crossed labels, and both fail.
- One such step alone may only be worded loosely, so it adds a note and the check still passes.
- A step that shares one title word with a doc it names is never judged.

The check compares words only. It cannot tell whether a step is the behavior its doc describes when the two share no words.

## Failures

Every error names the view line it concerns and the requirement ID.

- A requirement doc no label names, at the diagram's fence line.
- A label naming an acceptance ID or a doc that does not exist.
- A label naming more than eight IDs.
- A behavior node with no ID, at its diagram line.
- A generated region that differs from the docs, or broken markers, at the marker's line.
- A generated region that differs from the docs, or broken markers, at the marker's line.
- Two steps whose labels read as each other's docs, at each step's line.

A repo without architecture views passes with zero files. A missing journey view belongs to `architecture-contract`.

## Requirements

[Requirement 0005](https://github.com/may-journal/fitness-checks/blob/main/docs/requirements/0005-diagrams-match-the-code.md) proves this check with installed-command cases.
