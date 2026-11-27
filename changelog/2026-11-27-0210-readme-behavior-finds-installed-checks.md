### 2026.11.27.0210

> The readme-behavior check reads the flags and keys of other checks wherever the engine installs them.

- Fix: `readme-behavior` finds the program of another check as the engine does, so an installed set no longer hides a missing flag.
- Docs: the installed-layout acceptance moves from requirement 0001 to the new requirement 0002 of the check.
- Test: the acceptance builds each check in its own folder and runs with one check or both.
