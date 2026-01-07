# Runners

This document explains how **runners** work in Kata and how to create new ones.

A runner defines **how tests are executed**, not what is being tested.
It is entirely declarative and does not require writing Go code.

---

## What Is a Runner?

A runner is a reusable definition of test execution logic.

In Kata, runners:

- are implemented as files, not code
- rely on existing tools (Docker, test frameworks, CLIs)
- are copied into the project during `kata start`
- are referenced by katas via configuration

The Kata binary treats runners as opaque executors and only observes whether tests pass or fail.

---

## Runner Responsibilities

A runner is responsible for:

- executing tests
- producing a correct exit code
- streaming logs to the terminal
- optionally persisting logs for later inspection

A runner is **not** responsible for:

- defining application behavior
- managing kata steps
- interpreting test output
- deciding kata progression

---

## Runner Structure

A typical runner lives under:

```text
katas/runners/<runner-name>/
```

Minimum structure:

```text
katas/runners/<runner-name>/
├── Taskfile.testrunner.yml
└── ...
```

Optional files may include:

- compose.override.yaml
- scripts
- configuration files required by the runner

## The Taskfile

### Taskfile Name

Runners must define their execution logic in:

```text
Taskfile.testrunner.yml
```

This file is copied into the project root during `kata start`.

### Required Task

Every runner must define a task named `test`.

This is the task executed by `kata run`.

```yaml
version: '3'

tasks:
  test:
    cmds:
      - echo "running tests"
```

If the task exits with a non-zero status, the tests are considered failed.

## Exit Codes Matter

Kata determines test success solely based on the task exit code.

- `exit code == 0` → tests passing
- `exit code != 0` → tests failing

Runners must ensure that failures propagate correctly.

## Pipelines and `pipefail`

When using shell pipelines (for example, with `tee`), runners must enable `pipefail` to avoid masking errors.

Recommended configuration:

```yaml
version: '3'

set:
  - pipefail

tasks:
  test:
    cmds:
      - run-tests | tee .kata/last_test.log
```

Without `pipefail`, failing tests may incorrectly return exit code `0`.

## Logging

Runners should:

- stream logs directly to stdout/stderr
- optionally store logs in `.kata/last_test.log`

This enables:

- immediate feedback for users
- post-mortem inspection
- automation and AI agent consumption

Log parsing is intentionally avoided.

## Using Docker

Runners may use Docker or Docker Compose, but they should:

- avoid hardcoding application service names
- rely on the kata scaffold to define the system under test
- remain reusable across multiple katas

Common patterns include:

- a long-running test container kept idle
- test execution via docker compose exec
- Docker Compose override files

## Runner Configuration

Runners are selected by katas using `config.toml`:

```toml
runner = "hurl-docker"
```

Each kata may reference at most one runner.

The runner name must match a directory under `katas/runners/`.

## Best Practices

Keep runners minimal and focused

- Prefer declarative tools over custom scripts
- Avoid embedding environment-specific assumptions
- Use explicit file names and conventions
- Ensure failure conditions propagate reliably

## Anti-Patterns

Avoid the following:

- parsing test output in the runner
- hardcoding kata-specific behavior
- relying on interactive prompts
- swallowing exit codes
- implementing logic that belongs in the kata scaffold

## Adding a New Runner

To add a new runner:

- Create a directory under `katas/runners/<runner-name>`
- Add a `Taskfile.testrunner.yml`
- Include any required supporting files
- Reference the runner in a kata’s `config.toml`
- No Go code changes are required.

No Go code changes are required.

## Summary

Runners are the bridge between Kata and the testing ecosystem.

By keeping runners declarative, explicit, and minimal, Kata remains flexible, extensible, and safe for automation and AI-driven workflows.
