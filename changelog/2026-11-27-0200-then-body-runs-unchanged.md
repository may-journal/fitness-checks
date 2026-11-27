### 2026.11.27.0200

> A body inside `req.Then` runs as it was written, so tagging a test case changes none of its lines.

- Feat: `req.Then` runs its body in a nested case and hands it that case's `*testing.T`, 0011.3.
- Feat: `t.Errorf`, `t.Fatalf` and the acceptance helpers work in the body with no edit, 0011.3.
- Docs: requirement 0011 states that a body runs as it was written.
