### 2026.11.14.2359

> A missing major tag is created, since GitHub answers a move of a missing tag with a 422 and not a 404.

- Fix: `scripts/release-major.sh` creates the tag when the move says the reference does not exist.
- Test: the stand-in gh answers a missing tag as GitHub does, with HTTP 422.
- Chore: the first run with the App token failed with HTTP 422 on the first tag.
