### 2026.11.27.0200

> A body inside `req.Then` runs as it was written, so tagging a test case changes none of its lines.

- Feat: `req.Then` runs its body in a nested case and hands it that case's `*testing.T`, so `t.Errorf` and the acceptance helpers work unchanged, 0011.3.
