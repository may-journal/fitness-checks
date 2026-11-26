# open-questions

Fails a list of open questions in a doc, page, PR or issue body. [AGENTS.md](https://github.com/may-journal/fitness-runner/blob/main/AGENTS.md#ask-in-chat) holds the rule.

## Behavior

[Requirement 0001](docs/requirements/0001-questions-are-asked-in-chat.md)

## Notes

Run `fitness explain open-questions` for the files it reads. The [plan-structure](../fitness-check-plan-structure/README.md) check shares this definition of the name. The Claude Code hook judges an Artifact page with it before the page publishes.

## Example

```
fitness open-questions
fitness open-questions --body-file body.md
```

It reads no `.fitnessrc.json` key.

