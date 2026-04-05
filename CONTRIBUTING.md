# Contributing to XJS

Thank you for your interest in contributing!

## Commit convention

Commits follow the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
type(scope): short title
```

The scope is optional. The title must be written in imperative English (`add`, `fix`, `update`), completing the sentence _"If applied, this commit will..."_

### Types

| Type       | Description                                  |
| ---------- | -------------------------------------------- |
| `feat`     | New functionality                            |
| `fix`      | Bug fix                                      |
| `chore`    | Maintenance (configs, dependencies, tooling) |
| `docs`     | Documentation                                |
| `test`     | Tests                                        |
| `refactor` | Refactoring without behavior change          |
| `perf`     | Performance improvements                     |

### Examples

```
feat(lexer): add token consumer for strings
fix(token): correct position tracking
chore: update linter config
docs: add contributing guide
```
