# Code

```mermaid
classDiagram
    class CheckProgram["1 fitness-check-name"] {
        <<abstract>>
        +argv: root, passthrough
        +stdout: JSON result
        +exit: verdict code
    }
    class Main["2 program.Main(string, checkkit.Check)"]
    class Summary["3 program.Summary(string) string"]
    class described["4 program.described"] {
        +Describe
        +Protocol
    }
    class checkkitMain["5 checkkit.Main"] {
        <<abstract>>
    }

    CheckProgram --> Main : A
    Main --> Summary : B
    Main --> described : C
    Main --> checkkitMain : D

    style CheckProgram fill:#6d28d9,stroke:#4c1d95,color:#ffffff
    style Main fill:#6d28d9,stroke:#4c1d95,color:#ffffff
    style Summary fill:#6d28d9,stroke:#4c1d95,color:#ffffff
    style described fill:#6d28d9,stroke:#4c1d95,color:#ffffff
    style checkkitMain fill:#6d28d9,stroke:#4c1d95,color:#ffffff
```

```mermaid
classDiagram
    class Example {
        +property
        +method()
    }
    class Parent
    class Child
    Parent <|-- Child : inherits
    Parent --> Child : uses

    style Example fill:#6d28d9,stroke:#4c1d95,color:#ffffff
    style Parent fill:#6d28d9,stroke:#4c1d95,color:#ffffff
    style Child fill:#6d28d9,stroke:#4c1d95,color:#ffffff
```

Numbers on classes and letters on relationships match the callout table.

| #   | Description | Why |
| --- | ----------- | --- |
| 1   | One program per check, started with `--root` and the engine's context; its `main` lives in its `cmd/fitness-check-*/main.go`. | The protocol is the whole interface. |
| 2   | Answers `--describe` itself, else runs the check through `checkkit.Main`, in [program.go](../../internal/program/program.go). | The handshake adds the protocol number. |
| 3   | The first sentence of the check's README in [program.go](../../internal/program/program.go). | The summary is written once, in the README. |
| 4   | The check's metadata plus the protocol number, `program.Protocol`, in [program.go](../../internal/program/program.go). | The engine reads both from one JSON line and fails a run on a mismatch. |
| 5   | The engine library's scaffolding that runs a check and emits its result. | Exit codes and the JSON shape live once. |
| A   | Each check's `main` passes its embedded README and its check. | The check carries its own metadata. |
| B   | `--describe` fills a missing summary from the README. | No generated summary table to keep current. |
| C   | `--describe` prints the metadata with the protocol. | One JSON line per handshake. |
| D   | Every other run goes to the scaffolding. | Runs behave as the bundled check did. |
