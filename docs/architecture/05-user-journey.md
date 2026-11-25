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

Each node is one whole requirement doc, in ID order. The nodes between the `journey-nodes` markers are generated from the docs' IDs and titles; run `fitness architecture-journey --fix` to rewrite them. The edges and callout table are hand-written.

```mermaid
flowchart TB
    dev((1 Developer or agent)):::persona

    red[10 Finding names the file, the line and the fix<br/>Req 0003, 0004, 0005, 0006, 0007, 0008]:::blocked

    subgraph docs["Requirement docs"]
        %% journey-nodes:start
        req0001[(2 Checks Release And Install On Their Own<br/>Req 0001)]:::system
        req0002[(3 Checks Run Only Where They Apply<br/>Req 0002)]:::system
        req0003[(4 History Stays Clean<br/>Req 0003)]:::system
        req0004[(5 Writing Stays Clear<br/>Req 0004)]:::system
        req0005[(6 Diagrams Match The Code<br/>Req 0005)]:::system
        req0006[(7 Code Is Tested By Its Own Tools<br/>Req 0006)]:::system
        req0007[(8 The Repo Holds Only Working Source<br/>Req 0007)]:::system
        req0008[(9 Check Docs Are Markdown<br/>Req 0008)]:::system
        %% journey-nodes:end
    end

    dev -->|A| req0001
    req0001 --> req0002
    req0002 -->|B| req0003
    req0003 --> req0004
    req0004 --> req0005
    req0005 --> req0006
    req0006 --> req0007
    req0007 --> req0008
    req0008 --> red
    red -->|C| req0002
    req0008 -->|D| req0001

    classDef persona fill:#eef2ff,stroke:#6366f1,color:#312e81,stroke-width:2px
    classDef screen fill:#f8fafc,stroke:#64748b,color:#1e293b,stroke-width:2px
    classDef action fill:#fff7ed,stroke:#f97316,color:#9a3412,stroke-width:2px
    classDef system fill:#ecfdf5,stroke:#10b981,color:#064e3b,stroke-width:2px
    classDef blocked fill:#fef2f2,stroke:#ef4444,color:#991b1b,stroke-width:2px

    style docs fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
```

Numbers match the callout table. System context: [01-system-context.md](01-system-context.md). Code: [04-code.md](04-code.md).

| #   | Description                                        | Type      | Why                                                                    |
| --- | -------------------------------------------------- | --------- | ---------------------------------------------------------------------- |
| 1   | Developer or AI agent.                             | Persona   | Agents run the same command as a person.                               |
| 2   | Each check builds, tests, releases and installs.   | System    | A person takes one check's fix without the others.                     |
| 3   | Each check judges only the files it fits.          | System    | Any repo runs every check with nothing to configure, fast at any size. |
| 4   | Commit and changelog checks.                       | System    | Each change shows what it was, who helped and which issue it served.   |
| 5   | Prose, spelling, comment and body checks.          | System    | Any doc, issue or PR reads in a minute.                                |
| 6   | Mermaid and architecture view checks.              | System    | Every diagram renders and matches the code it draws.                   |
| 7   | Paired change, tooling and coverage checks.        | System    | Each language's own tools judge the code with the repo's config.       |
| 8   | Build output, duplicate, pin and workflow checks.  | System    | Reviewers see real source only.                                        |
| 9   | Check docs and their move to Markdown.             | System    | People edit each rule as Markdown that GitHub renders.                 |
| 10  | A red row in the results table.                    | Finding   | Each finding says where it is and what to change.                      |
| A   | Run `fitness`, a hook or one check in any repo.    | Connector | The engine installs the latest release first.                          |
| B   | Every check that fits the change runs.             | Connector | A check that does not fit passes with zero files.                      |
| C   | Fix the code, not the check, and rerun.            | Connector | The check measures the real code, so the code changes.                 |
| D   | A green change merges and CI releases it.          | Connector | The loop closes: a release reaches every repo.                         |
