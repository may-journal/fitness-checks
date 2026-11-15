### 2026.11.15.0234

> The README check passes on main again, since the local-branch README comes from this repo's generator.

- Docs: regenerate `cmd/fitness-check-local-branch/README.md` with `scripts/generate-readmes.sh`, so its header names this repo's generator.
- Fix: `scripts/generate-readmes.sh -check` and requirement 0165.1 pass again.
- Chore: the stale header dates from before `local-branch` moved here from the engine.
