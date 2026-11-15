# Containers

```mermaid
C4Container
    title Fitness checks

    Person(developer, "1 Developer")

    System_Boundary(ship, "Shipped programs") {
        Container(check, "2 Check programs", "Go", "cmd/fitness-check-*")
    }

    System_Ext(engine, "3 Fitness runner", "engine and shared go/pkg packages")
    System_Ext(consumer, "4 Consumer repo", ".fitnessrc.json and the sources under test")
    System_Ext(peers, "5 Peer tooling", "each language module's tools, and git")
    System_Ext(registry, "6 npm registry", "dependency-currency queries")

    Rel(developer, engine, "A")
    Rel(engine, check, "B")
    Rel(check, consumer, "C")
    Rel(check, peers, "D")
    Rel(check, registry, "E")
    Rel(check, engine, "F")

    UpdateElementStyle(developer, $bgColor="#334155", $borderColor="#1e293b", $fontColor="#ffffff")
    UpdateElementStyle(check, $bgColor="#0f766e", $borderColor="#115e59", $fontColor="#ffffff")
    UpdateElementStyle(engine, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
    UpdateElementStyle(consumer, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
    UpdateElementStyle(peers, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
    UpdateElementStyle(registry, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
```

```mermaid
C4Container
    title Legend

    Person(p, "Person", "Human or agent actor")
    System_Boundary(b, "System boundary") {
        Container(c, "Container", "Tech", "Runnable or deployable unit")
    }
    System_Ext(e, "External", "Outside our boundary")

    UpdateElementStyle(p, $bgColor="#334155", $borderColor="#1e293b", $fontColor="#ffffff")
    UpdateElementStyle(c, $bgColor="#0f766e", $borderColor="#115e59", $fontColor="#ffffff")
    UpdateElementStyle(e, $bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")
```

Numbers on nodes and letters on arrows match the callout table. Each container names the commands it is built from; a trailing `*` names every sibling folder.

| #   | Description                                                                          | Why                                                   |
| --- | ------------------------------------------------------------------------------------ | ----------------------------------------------------- |
| 1   | Same actor as the system context.                                                    | Every run starts at the engine.                       |
| 2   | One static program per check, each built from its own `cmd/fitness-check-*` folder.  | A repo installs only the checks it uses.              |
| 3   | The engine, and the shared packages each check imports from its Go module.           | One copy of each helper, pinned in `go.mod`.          |
| 4   | Optional `.fitnessrc.json` with per-check options, and the files checks scan.        | Config and the code under test live with the app.     |
| 5   | Tool-wrapper checks exec these when installed.                                       | Peer tools are found at runtime, never bundled.       |
| 6   | dependency-currency queries the registry over HTTPS.                                 | No npm binary needed; offline degrades to pass.       |
| A   | The developer runs `fitness` or a hook.                                              | One front door.                                       |
| B   | The engine starts each check program as its own process.                             | Checks cannot share state.                            |
| C   | Each check reads the config and the tree, staged or full.                            | The working directory is the app under test.          |
| D   | Checks find peer tools the way their language module does, then on PATH.             | A missing tool reports info with an install hint.     |
| E   | Abbreviated registry metadata, honoring the registry `.npmrc` sets.                  | Currency judged without shelling to npm.              |
| F   | Each check links the engine's public packages and answers in its protocol.           | Both sides speak one versioned contract.              |
