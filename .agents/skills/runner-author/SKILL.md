---
name: runner-author
description: Guides an agent on how to create a new Runner (Taskfile.testrunner.yml, docker-compose, etc.) for a specific language or framework. Use this skill when asked to create a new test runner.
license: Apache-2.0
---

# Creating a Kata Runner

A Kata Runner defines the test execution logic (the "how" of testing) independently of the Kata's domain logic (the "what" of testing). As an AI agent, you can define new runners to support various languages, frameworks, or testing tools within the Kata CLI.

## What is a Runner?

A Runner is a declarative definition of how tests are executed. The Kata binary does not embed any testing logic; it asks the runner "Did the tests pass?" based on the runner's exit code. Runners live in the `katas/runners/` directory.

### Why declarative?
A declarative model avoids imperative code that becomes brittle. The Taskfile acts as the source of truth for execution logic. The CLI uses `go-task` as a library to interpret the Taskfile, so there are no external binaries required to run tests.

## Runner Structure

A typical runner directory looks like this:

```text
katas/runners/<runner-name>/
├── Taskfile.testrunner.yml
└── ... (optional supporting files like compose.override.yaml)
```

The core component is the `Taskfile.testrunner.yml`.

### Example `Taskfile.testrunner.yml`

This file is included by the Kata's scaffold `Taskfile.yml`. It defines the specific testing tool and execution environment.

```yaml
version: '3'

tasks:
  test:
    desc: Run tests using Python pytest
    cmds:
      - pytest --verbose
```

## Integrating a Runner into a Kata

To use a newly created runner, a Kata must specify the runner's name in its `config.toml` file:

```toml
[runner]
name = "<runner-name>"
```

The Kata's `scaffold/Taskfile.yml` includes the runner's execution logic by referencing its `Taskfile.testrunner.yml`.

```yaml
version: '3'

includes:
  runner: ./Taskfile.testrunner.yml

tasks:
  test:
    cmds:
      - task: runner:test
```

## Best Practices

- **Exit Codes**: Ensure your test execution tool (e.g., `pytest`, `go test`, `npm test`) returns `0` on success and non-zero on failure. The Kata CLI relies heavily on these exit codes to interpret success or failure.
- **Dependencies**: Provide the environment for the tests to run. This could mean using a Docker container (via Docker Compose) or relying on local system dependencies if necessary.
- **No Test Framework Logic**: The runner should simply act as an orchestrator. Do not try to parse test output or handle framework-specific semantics inside the Kata core or the Taskfile itself.
- **Isolate Environment**: For consistency, avoid hidden behavior and inferring state from the filesystem.

## Summary

When defining a runner, you are abstracting test execution. A well-designed runner is declarative, handles its own environment dependencies, and clearly communicates test success/failure via standard exit codes.
