### 2026.11.15.0010

> The finish job moves the major tags with the automation App token, because the workflow token gets HTTP 403.

- Fix: the `finish` job mints the automation App token and uses it to move each major tag.
- Fix: the same token publishes the release index.
- Docs: the first run with the workflow token failed on the refs API with HTTP 403.
