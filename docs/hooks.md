# Hooks

Kata supports a set of optional hooks that can be defined in your kata's `Taskfile.yml`. These hooks allow you to execute custom logic at specific points in the kata lifecycle.

## Available Hooks

### `on_start_hook`
Executed after `kata start` has successfully completed scaffolding the kata.
This is useful for initializing environments, such as starting Docker containers.

### `pre_run_hook`
Executed before the test task in `kata run`.
Useful for setting up prerequisites before tests run.

### `on_success_hook`
Executed after the test task succeeds in `kata run`.

### `on_failure_hook`
Executed after the test task fails in `kata run`.

### `on_kata_finish_hook`
Executed when the kata transitions to the "finished" state (i.e., when `kata next` is called and there are no more steps).

## Example `Taskfile.yml`

```yaml
version: '3'

includes:
  runner:
    taskfile: ./Taskfile.testrunner.yml
    internal: true

tasks:
  on_start_hook:
    cmds:
      - echo "Kata started! initializing environment..."
      - docker compose up -d

  pre_run_hook:
    cmds:
      - echo "About to run tests..."

  on_success_hook:
    cmds:
      - echo "Tests passed! Great job!"

  on_failure_hook:
    cmds:
      - echo "Tests failed. Keep trying!"

  on_kata_finish_hook:
    cmds:
      - echo "Congratulations! You have completed the kata."
```
