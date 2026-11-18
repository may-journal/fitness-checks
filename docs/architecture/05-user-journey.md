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

Each subgroup is one requirement doc, titled with its ID, in order from top to bottom. The `architecture-journey` check reads those titles, so every step inside names its requirement.

```mermaid
flowchart TB
    dev((1 Developer or agent)):::persona

    subgraph r1["Req 0001 Checks release and install on their own"]
        change(2 Change one check and open a PR):::action
        ci[(3 CI builds and tests only the checks it alters)]:::system
        release[(4 Main releases each changed check alone)]:::system
        index[5 Archives, checksums and a release index entry]:::screen
    end

    subgraph r2["Req 0002 Checks run only where they apply"]
        run(6 Run fitness, a hook or one check):::action
        reminder[7 Reminder and table of enabled checks]:::screen
        match[(8 Each check matches only its own files)]:::system
        skipped[9 Zero files, or unreadable input named as info]:::screen
    end

    subgraph r3["Req 0003 History stays clean"]
        commitChecks[(10 Commit subject, trailers and issue links)]:::system
        changelogChecks[(11 One dated changelog entry per change)]:::system
    end

    subgraph r4["Req 0004 Writing stays clear"]
        proseChecks[(12 Prose limits, spelling and wording)]:::system
        templateChecks[(13 Comments, templates and bodies)]:::system
    end

    subgraph r5["Req 0005 Diagrams match the code"]
        mermaidChecks[(14 Diagrams parse and pair callouts with rows)]:::system
        letters(15 Run with --fix to reletter connectors):::action
        viewChecks[(16 Architecture views match code and requirements)]:::system
    end

    subgraph r6["Req 0006 Code is tested by its own tools"]
        pairChecks[(17 Doc, test and code change together)]:::system
        toolChecks[(18 Each language's own tools and config)]:::system
        coverage[(19 Tests and coverage on the pinned runtime)]:::system
    end

    subgraph r7["Req 0007 The repo holds only working source"]
        sourceChecks[(20 No build output, plan files, copies or repeats)]:::system
        pinChecks[(21 Current pins, explained ignores, sound workflows)]:::system
    end

    red[22 Finding names the file, the line and the fix<br/>Req 0003, 0004, 0005, 0006, 0007]:::blocked
    fix(23 Fix the code, not the check<br/>Req 0001, 0002):::action

    dev --> change
    dev -->|A| run
    change --> run
    run --> reminder
    reminder --> match
    match -->|B| skipped
    match --> r3
    match --> r4
    match --> r5
    match --> r6
    match --> r7
    mermaidChecks -->|C| letters
    letters --> mermaidChecks
    pinChecks -->|D| skipped
    r3 --> red
    r4 --> red
    r5 --> red
    r6 --> red
    r7 --> red
    red --> fix
    fix --> run
    fix -->|E| ci
    ci -->|F| red
    ci -->|G| release
    release --> index
    index -->|H| run

    classDef persona fill:#eef2ff,stroke:#6366f1,color:#312e81,stroke-width:2px
    classDef screen fill:#f8fafc,stroke:#64748b,color:#1e293b,stroke-width:2px
    classDef action fill:#fff7ed,stroke:#f97316,color:#9a3412,stroke-width:2px
    classDef system fill:#ecfdf5,stroke:#10b981,color:#064e3b,stroke-width:2px
    classDef blocked fill:#fef2f2,stroke:#ef4444,color:#991b1b,stroke-width:2px

    style r1 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style r2 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style r3 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style r4 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style r5 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style r6 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
    style r7 fill:#ffffff,stroke:#e2e8f0,color:#475569,stroke-width:1px
```

Numbers match the callout table. System context: [01-system-context.md](01-system-context.md). Code: [04-code.md](04-code.md).

| #   | Description                                         | Type      | Why                                                                      |
| --- | --------------------------------------------------- | --------- | ------------------------------------------------------------------------ |
| 1   | Developer or AI agent.                              | Persona   | Agents run the same command as a person.                                 |
| 2   | A fix to one check, a shared helper or `go.mod`.    | Action    | A person takes one check's fix without the others.                       |
| 3   | CI on the pull request.                             | System    | It builds, vets, tests and sizes only the checks the change can alter.   |
| 4   | Release on a push to main.                          | System    | Each changed check gets the bump its commits call for, and no other.     |
| 5   | Per-platform archives, checksums and moved majors.  | Screen    | The index lists each published check so the engine installs it alone.   |
| 6   | `fitness`, a git hook, or one check program.        | Action    | Any repo runs every check with nothing to configure.                     |
| 7   | `read-repo-first` output.                           | Screen    | The invoker sees which checks this run enables before reading results.   |
| 8   | Each check's declared file patterns.                | System    | A check judges only what it fits, with a few git calls at any size.      |
| 9   | A pass with zero files, or skips listed as info.    | Screen    | A check that does not apply stays green, and still judges the rest.      |
| 10  | Commit message checks.                              | System    | Typed subjects, Plan and AI trailers, and one link per issue.            |
| 11  | Changelog entry checks.                             | System    | One dated entry of three to five typed bullets, and versions that agree. |
| 12  | Prose and spelling checks.                          | System    | Any doc reads in a minute: short, plain and spelled right.               |
| 13  | Comment, template and body checks.                  | System    | Code comments, docs, Plans and PR bodies follow their rules.             |
| 14  | Mermaid checks.                                     | System    | Every diagram renders and explains each callout in its table.            |
| 15  | `--fix` on connector labels.                        | Action    | Bad letters fix themselves instead of by hand.                           |
| 16  | Architecture view checks.                           | System    | Views match runtimes, folders, symbols and requirement docs.             |
| 17  | Paired change checks.                               | System    | A requirement doc and its one test change in the same commit.            |
| 18  | Go, Node, Swift and Java tooling checks.            | System    | Each language's formatter and linter judge code with the repo's config.  |
| 19  | Test and coverage checks.                           | System    | Every line is covered, on the runtime version the repo pins.             |
| 20  | Build output, design work and duplicate checks.     | System    | Reviewers see real source, with no copies or repeated values.            |
| 21  | Dependency, ignore file and workflow checks.        | System    | Stale pins, bare ignore lines and broken workflow refs never land.       |
| 22  | A red row in the results table.                     | Finding   | Each finding says where it is and what to change.                        |
| 23  | An edit to the code, doc or test.                   | Action    | The check measures the real code, so the code changes.                   |
| A   | Run in any repo.                                    | Connector | A person can start here without changing a check.                        |
| B   | Nothing the check can judge or read.                | Connector | A check that does not fit passes instead of failing.                     |
| C   | Connector letters out of order.                     | Connector | The fix gives each connector and row the next free letter.               |
| D   | The package registry is unreachable.                | Connector | A network outage adds info, not a failure.                               |
| E   | Every check is green on a check change.             | Connector | The fix goes up as a pull request.                                       |
| F   | CI or the size check fails.                         | Connector | The same findings show in CI as on the desk.                             |
| G   | The pull request merges.                            | Connector | Only the merged change releases.                                         |
| H   | The next run installs the new release.              | Connector | The loop closes: a release reaches every repo.                           |
