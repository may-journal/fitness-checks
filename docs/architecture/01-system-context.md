# System Context

```mermaid
C4Context
    title Fitness checks

    Person(developer, "1 Developer")
    Person(agent, "2 AI agent")
    System(checks, "3 Fitness checks", "One program per check, each judging a repo on its own.")
    System_Ext(engine, "4 Fitness runner", "Engine that starts each check and renders the results")
    System_Ext(consumer, "5 Consumer repo", "Application under test")
    System_Ext(tools, "6 Lint format test tools", "eslint, prettier, swiftlint, git, etc.")

    Rel(developer, engine, "A")
    Rel(agent, engine, "B")
    Rel(engine, checks, "C")
    Rel(checks, consumer, "D")
    Rel(checks, tools, "E")
    Rel(checks, engine, "F")

    UpdateElementStyle(developer, $bgColor="#334155", $borderColor="#1e293b", $fontColor="#ffffff")
    UpdateElementStyle(agent, $bgColor="#334155", $borderColor="#1e293b", $fontColor="#ffffff")
    UpdateElementStyle(checks, $bgColor="#1d4ed8", $borderColor="#1e3a8a", $fontColor="#ffffff")
    UpdateElementStyle(engine, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
    UpdateElementStyle(consumer, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
    UpdateElementStyle(tools, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
```

```mermaid
C4Context
    title Legend

    Person(p, "Person", "Human or agent actor")
    System(s, "System", "Software we own")
    System_Ext(e, "External", "Outside our boundary")

    UpdateElementStyle(p, $bgColor="#334155", $borderColor="#1e293b", $fontColor="#ffffff")
    UpdateElementStyle(s, $bgColor="#1d4ed8", $borderColor="#1e3a8a", $fontColor="#ffffff")
    UpdateElementStyle(e, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
```

Numbers on nodes and letters on arrows match the callout table.

| #   | Description                                                                 | Why                                                     |
| --- | --------------------------------------------------------------------------- | ------------------------------------------------------- |
| 1   | Runs `fitness`, git hooks, or a single check locally.                       | Primary operator; same bar as CI.                       |
| 2   | Claude Code or another agent runs fitness before or after an edit.         | Agents share the same command and exit codes.           |
| 3   | One Go program per check, built from this repo, released on its own line.   | A check change ships without an engine release.         |
| 4   | The fitness-runner engine, which a repo installs beside the checks it uses. | The engine finds and starts each check program.         |
| 5   | Any repo; each check passes clean where it does not apply.                  | Checks run in the consumer's working directory.         |
| 6   | Commodity linters, formatters, test runners, and git.                       | Tool-wrapper checks run these and judge their output.   |
| A   | Local runs and git hooks.                                                   | Fast feedback before a commit.                          |
| B   | Agent hooks and skills run the same command.                                | No separate agent path.                                 |
| C   | The engine starts `fitness-check-<name>` with `--root` and `FITNESS_*` context. | The check protocol is the only contract.            |
| D   | Each check reads the staged files, the tree, or a body it is given.         | The working directory is the app root.                  |
| E   | Tool-wrapper checks exec the underlying tool; native engines need none.     | Orchestration, not reimplementation.                    |
| F   | Each check prints one JSON result and exits 0, 1 or 3.                      | Pass or fail is deterministic, not judged by a model.   |
