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

Each node is one whole requirement doc, grouped by phase. The `architecture-journey` check reads the IDs in the labels.

```mermaid
flowchart TB
    dev((1 Developer or agent)):::persona

    subgraph get["Get the checks"]
        install[(2 Checks release and install on their own<br/>Req 0001)]:::system
        apply[(3 Checks run only where they apply<br/>Req 0002)]:::system
    end

    subgraph judge["Judge a change"]
        history[(4 History stays clean<br/>Req 0003)]:::system
        writing[(5 Writing stays clear<br/>Req 0004)]:::system
        diagrams[(6 Diagrams match the code<br/>Req 0005)]:::system
        tested[(7 Code is tested by its own tools<br/>Req 0006)]:::system
        source[(8 The repo holds only working source<br/>Req 0007)]:::system
        docs[(9 Check docs are Markdown<br/>Req 0008)]:::system
    end

    red[10 Finding names the file, the line and the fix<br/>Req 0003, 0004, 0005, 0006, 0007]:::blocked

    dev -->|A| install
    install --> apply
    apply -->|B| judge
    judge --> red
    red -->|C| apply
    judge -->|D| install

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
