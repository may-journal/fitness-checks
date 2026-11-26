### 2026.11.26.0100

> The requirement folders of changed checks run as up to six parallel CI jobs.

- Fix: `check-requirements` deals the folders into shares, heaviest first, so a batch of twelve checks stays under its two-minute step timeout.
- Test: root requirement 0001.2.2 proves each changed folder runs once in at most six shares.
- Docs: root requirement 0001.2.2 says each changed folder runs once in six shares.
