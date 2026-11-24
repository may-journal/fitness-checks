# markdown-no-bold-italic

Fails bold and italic emphasis in markdown files: `**bold**`, `__bold__`, `*italic*`, or `_italic_`.

## Behavior

- Pass: no Markdown files, or no Markdown file contains bold (double asterisk or double underscore) or italic (single asterisk or single underscore) emphasis.
- Fail: any Markdown file contains disallowed emphasis; one error per occurrence with file path and snippet.

`CHANGELOG.md` is judged too, and `fitness explain markdown-no-bold-italic` lists the files it reads. Use bold or italic only when explicitly required. Ignores emphasis inside fenced code, inline code, and Markdown link blocks, so underscores in URLs or link text are not flagged. The bold scope Release Please gives release bullets is allowed too.
