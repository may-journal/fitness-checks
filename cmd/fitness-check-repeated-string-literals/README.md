# repeated-string-literals

Fails a string literal that appears three or more times across the repo's JS and TS sources. It should be a shared constant. When the repeats are a small closed set like `'active'` and `'pending'`, a union type or enum fits better. It is on by default.

## Why not just ESLint or duplicate-code

- `duplicate-code` finds copied blocks of 5+ lines; a single repeated literal is one token on one line, so it never sees it.
- `sonarjs/no-duplicate-string` (ESLint) only flags strings of 10+ chars containing a separator, which excludes the short, identifier-like tokens (`'active'`, `'GET'`) this check targets.

## What passes and fails

A value repeated only twice passes; a value used 3+ times anywhere across the scanned files fails, one error per value, most-repeated first:

```
// src/a.ts
if (order.status === 'active') setStatus('active');
// src/b.ts
const rows = all.filter((r) => r.status === 'active');
```

It fails with:

```
"active" appears 3 times (src/a.ts:1, src/a.ts:1, src/b.ts:2) — extract a shared constant
```

## Behavior

[Requirement 0001](docs/requirements/0001-a-repeated-string-becomes-a-shared-constant.md) states what this check does, and a test case named for each step proves it.

- Run `fitness explain repeated-string-literals` for the files it reads.
- A hand lexer reads each file's code view for single- and double-quoted literals.

## Allow list

Allow a value for the rare repeat no constant can fix. That means a value that must live once per runtime island — TS sources, published bin scripts, or automation scripts. Or one spelling shared by genuinely different things:

```
{
  "repeatedStringLiterals": { "allow": ["my-check-name", "kindValue"] }
}
```

Keep the list a shrinking baseline; the fix is to centralize the constant. A `.fitnessrc.json` that fails to load leaves the allow list empty. Run `fitness explain repeated-string-literals` for the key.

## Limitations

- Template literals are out of scope — only plain `'…'` and `"…"` literals count.
- Thresholds and the idiomatic-token set are fixed in the binary; there is no inline marker and no enum-candidate detection yet (tracked in #44).
