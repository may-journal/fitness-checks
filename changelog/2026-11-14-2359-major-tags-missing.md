### 2026.11.14.2359

> A missing major tag is created, since GitHub answers a move of a missing tag with a 422 and not a 404.

- Fix: `scripts/release-major.sh` creates the tag when a lookup finds none, instead of matching error text.
- Test: requirement 0163.6 proves a refused move of an existing tag stops the run without a create.
- Chore: the first run with the App token failed on the first tag with HTTP 422.
