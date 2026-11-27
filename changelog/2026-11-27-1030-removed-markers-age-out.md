### 2026.11.27.1030

> A Removed marker stops holding its number once two releases pass, so docs shorten without surprising anyone.

- Feat: `requirements` lets a `Removed` line go once two released versions are newer, and later acceptances move up one number.
- Feat: `requirements` fails that renumber unless each moved acceptance's test and step cases are renamed with it.
- Feat: `requirements` fails deleting a younger `Removed` line, saying how to drop old markers and rename their tests.
- Docs: `requirements` gains requirement 0010 and a README section, with a case for each pass and fail.
