# Writing Katas

This document explains how to create and structure **Katas** for the Kata CLI.

---

## What is a Kata?

A kata is a self-contained coding exercise designed to practice specific skills.
It consists of:

- **Metadata**: Name, description, and configuration.
- **Scaffold**: The initial code and file structure.
- **Steps**: A sequence of incremental changes (usually tests) to guide the user.

---

## Kata Structure

A kata is defined in a directory under `katas/catalog/`.

```text
katas/catalog/<kata-name>/
├── config.toml          # Configuration
├── README.md            # Description
├── scaffold/            # Initial project state
│   ├── Taskfile.yml     # Entry point for execution
│   └── ...
└── steps/
    ├── 01-step-name/    # First step
    ├── 02-step-name/    # Second step
    └── ...
```

### Hidden Katas (Go specific)

If you are developing Go katas inside the Kata repository itself, you might encounter build errors if the kata source code is incomplete (which is intentional for katas).
To prevent the Go build system from trying to compile your kata, prefix the kata directory name with an underscore `_`.

Example:
`katas/catalog/_my-go-kata/`

The Kata CLI automatically:
1. Strips the `_` prefix when listing katas.
2. Resolves the name correctly when starting (e.g., `kata start my-go-kata` works even if the directory is `_my-go-kata`).

---

## Configuration (`config.toml`)

The `config.toml` file defines the kata's metadata and dependencies.

```toml
[runner]
name = "go"  # The runner to use (must match a directory in katas/runners/)
```

---

## Scaffold

The `scaffold/` directory contains the files that will be copied to the user's workspace when they run `kata start`.

### The `Taskfile.yml`

The scaffold **must** include a `Taskfile.yml` that acts as the entry point for the Kata CLI.
Typically, it includes the runner's Taskfile and delegates the `test` task.

```yaml
version: '3'

includes:
  runner: ./Taskfile.testrunner.yml

tasks:
  test:
    cmds:
      - task: runner:test
```

---

## Lifecycle Hooks

Katas can define special tasks in their `Taskfile.yml` to execute custom logic during the kata lifecycle.
These hooks are optional.

| Hook | Trigger | Description |
| :--- | :--- | :--- |
| `on_kata_start` | `kata start` | Runs immediately after the scaffold is applied. Useful for setup (e.g., initializing a database, installing dependencies). |
| `before_run` | `kata run` | Runs before tests are executed. Useful for pre-test checks or compilation. |
| `on_run` | `kata run` | Runs after tests finish, regardless of success or failure. |
| `on_run_success` | `kata run` | Runs only if tests pass. |
| `on_run_fail` | `kata run` | Runs only if tests fail. |
| `on_kata_finish` | `kata next` | Runs when the kata is completed (no more steps). |

### Example with Hooks

```yaml
version: '3'

includes:
  runner: ./Taskfile.testrunner.yml

tasks:
  test:
    cmds:
      - task: runner:test

  on_kata_start:
    cmds:
      - echo "Welcome to the Kata! Installing dependencies..."
      - npm install

  on_run_fail:
    cmds:
      - echo "Tests failed. Don't give up!"

  on_kata_finish:
    cmds:
      - echo "Congratulations! You have completed the kata."
```

### Failure Handling

- If `on_kata_start` or `before_run` fails (non-zero exit code), the command aborts.
- If `on_run_*` or `on_kata_finish` fails, the error is logged, but the command generally proceeds (e.g., `kata run` will still report the test result).

---

## Steps

Steps are located in `steps/` and are applied sequentially by `kata next`.
Each step directory contains files that will be overlaid onto the user's workspace.

- **Files**: Overwrite existing files or create new ones.
- **Deletions**: (Advanced) Currently, steps only add/overwrite files.

Steps are usually ordered alphabetically. Using a numeric prefix (e.g., `01-`, `02-`) is recommended.

---

## Best Practices

1.  **Keep it Simple**: The scaffold should be minimal.
2.  **Focus on the Goal**: Ensure steps guide the user logically.
3.  **Test Your Kata**: Run through the kata yourself to ensure steps work as expected.
