# gofmt

Lists Go files that gofmt would change. It is on by default.

## Behavior

- Skips `testdata` and `vendor` fixtures, which are not its own code, and ignored paths. Run `fitness explain gofmt` for the files it reads.
- Runs `gofmt -l` over them and reports each listed file as `path: not gofmt-formatted (run: gofmt -w path)`.
- A repo with no `.go` files passes clean (`filesChecked: 0`), so the check is safe in every repo.
- A file gofmt cannot parse fails as `path:line:column: message`.
- A missing Go toolchain reports info with a one-line install hint, never a failure.
- A gofmt that could not start, was killed, or exits non-zero with neither finding reports info quoting its output.
- `filesChecked` counts the Go files checked.
