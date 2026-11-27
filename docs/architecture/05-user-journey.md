# User Journey

## Legend

```mermaid
flowchart TB
    subgraph types["Node types"]
        persona((1 Persona)):::persona
        screen[2 Screen]:::screen
        action(3 Action):::action
        system[(4 System)]:::system
        finding[5 Finding]:::blocked
    end

    classDef persona fill:#eef2ff,stroke:#6366f1,color:#312e81,stroke-width:2px
    classDef screen fill:#f8fafc,stroke:#64748b,color:#1e293b,stroke-width:2px
    classDef action fill:#fff7ed,stroke:#f97316,color:#9a3412,stroke-width:2px
    classDef system fill:#ecfdf5,stroke:#10b981,color:#064e3b,stroke-width:2px
    classDef blocked fill:#fef2f2,stroke:#ef4444,color:#991b1b,stroke-width:2px

    style types fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
```

Numbers match the callout table. "Screen" is terminal output, since each check is a command-line program.

| #   | Description | Why                                                        |
| --- | ----------- | ---------------------------------------------------------- |
| 1   | Persona     | Developer or AI agent — same command, same exit codes.     |
| 2   | Screen      | What the invoker sees printed to the terminal.             |
| 3   | Action      | A command or file edit the invoker performs.               |
| 4   | System      | A check program's step: judge, report, exit.               |
| 5   | Finding     | A red result that sends the invoker back to fix.           |

## Journey

Each node is one whole requirement doc, grouped by phase. Nodes between `journey-nodes` markers are generated from the docs the start marker lists. Run `fitness architecture-journey --fix` to rewrite them. The phases, edges and callout table are hand-written.

```mermaid
flowchart TB
    dev((1 Developer or agent)):::persona

    subgraph get["Get the checks"]
        %% journey-nodes:start 0001 0002 0009
        req0001[(2 Checks Release And Install On Their Own<br/>Req 0001)]:::system
        req0002[(3 Checks Run Only Where They Apply<br/>Req 0002)]:::system
        req0009[(5 Releases Wake The Org Rerun<br/>Req 0009)]:::system
        %% journey-nodes:end
    end

    subgraph judge["Keep the docs"]
        %% journey-nodes:start 0008 0010
        req0008[(4 Check Docs Are Markdown<br/>Req 0008)]:::system
        req0010[(6 Requirement Docs Say Each Thing Once<br/>Req 0010)]:::system
        %% journey-nodes:end
    end

    red[7 Finding names the file, the line and the fix<br/>Req 0002, 0008]:::blocked

    dev -->|A| req0001
    req0001 --> req0002
    req0002 --> req0009
    req0009 -->|B| judge
    judge --> red
    red -->|C| req0002
    judge -->|D| req0001

    classDef persona fill:#eef2ff,stroke:#6366f1,color:#312e81,stroke-width:2px
    classDef screen fill:#f8fafc,stroke:#64748b,color:#1e293b,stroke-width:2px
    classDef action fill:#fff7ed,stroke:#f97316,color:#9a3412,stroke-width:2px
    classDef system fill:#ecfdf5,stroke:#10b981,color:#064e3b,stroke-width:2px
    classDef blocked fill:#fef2f2,stroke:#ef4444,color:#991b1b,stroke-width:2px

    style get fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style judge fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
```

Numbers match the callout table. System context: [01-system-context.md](01-system-context.md). Code: [04-code.md](04-code.md).

| #   | Description                                        | Type      | Why                                                                    |
| --- | -------------------------------------------------- | --------- | ---------------------------------------------------------------------- |
| 1   | Developer or AI agent.                             | Persona   | Agents run the same command as a person.                               |
| 2   | Each check builds, tests, releases and installs.   | System    | A person takes one check's fix without the others.                     |
| 3   | Each check judges only the files it fits.          | System    | Any repo runs every check with nothing to configure, fast at any size. |
| 4   | Check docs and their move to Markdown.             | System    | People edit each rule as Markdown that GitHub renders.                 |
| 5   | A published check wakes the org's fitness rerun.   | System    | No open pull request keeps a green from before the release.            |
| 6   | Requirement docs say each thing once.              | System    | A repeat tells the author to group steps, so each doc reads once.      |
| 7   | A red row in the results table.                    | Finding   | Each finding says where it is and what to change.                      |
| A   | Run `fitness`, a hook or one check in any repo.    | Connector | The engine installs the latest release first.                          |
| B   | Every check that fits the change runs.             | Connector | A check that does not fit passes with zero files.                      |
| C   | Fix the code, not the check, and rerun.            | Connector | The check measures the real code, so the code changes.                 |
| D   | A green change merges and CI releases it.          | Connector | The loop closes: a release reaches every repo.                         |
