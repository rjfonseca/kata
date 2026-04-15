---
name: kata-author
description: Guides an agent on how to create a new Kata (scaffold, steps, configuration, Taskfile.yml) following the project's structure. Use this skill when asked to author or create a new Kata.
license: Apache-2.0
---

# Creating a Kata

A Kata is a self-contained coding exercise designed to practice specific skills. It consists of a scaffold, a sequence of steps, and configuration metadata. As an AI agent, you can help the user build new Katas by creating the required directory structure and files.

## What is a Kata?

A Kata provides:
- **Scaffold**: The initial code and file structure applied when a user runs `kata start`.
- **Steps**: Incremental changes (usually tests) applied via `kata next` to guide the user.
- **Runner**: A declarative execution engine that knows how to test the code (defined in `config.toml`).

## Directory Structure

To create a new Kata named `<kata-name>`, create the following structure under `katas/catalog/`:

```text
katas/catalog/<kata-name>/
├── config.toml          # Configuration and metadata
├── README.md            # Description and instructions
├── scaffold/            # Initial project state
│   ├── Taskfile.yml     # Entry point for execution
│   └── ...              # Other scaffold files
└── steps/
    ├── 01-step-name/    # First step
    ├── 02-step-name/    # Second step
    └── ...              # More steps
```

*(Note for Go Katas inside the Kata repository: prefix the directory name with an underscore, e.g., `_my-go-kata`, to prevent Go build errors for incomplete code. The Kata CLI will strip the underscore when listing and starting Katas.)*

## Kata Components

### 1. `config.toml`

The `config.toml` file defines the Kata's metadata and dependencies, including the runner.

```toml
[runner]
name = "go"  # Must match a directory in katas/runners/
```

### 2. Scaffold (`scaffold/`)

The scaffold contains the files copied to the user's workspace upon `kata start`.

**The `Taskfile.yml` is mandatory.** It acts as the entry point for the Kata CLI. A typical `Taskfile.yml` includes the runner's Taskfile and delegates the `test` task:

```yaml
version: '3'

includes:
  runner: ./Taskfile.testrunner.yml

tasks:
  test:
    cmds:
      - task: runner:test
```

### 3. Steps (`steps/`)

Steps are sequential directories that define incremental changes to the workspace. They are usually ordered with numeric prefixes (e.g., `01-`, `02-`).

- **Files added to a step** will overwrite existing files or create new ones in the user's workspace when `kata next` is called.
- The user is expected to write code to pass the tests introduced in each step.
- Ensure each step provides a clear, single concept to practice (e.g., a new failing test).

## Templates

Files in a Kata scaffold can be templates.

- Add the `.tmpl` extension to any file (e.g., `go.mod.tmpl`).
- When `kata start` is executed, the `.tmpl` extension is removed.
- Templates are rendered using Go's `text/template`, with access to dynamic variables from the Kata configuration, user configuration, and repository metadata.

## Lifecycle Hooks

You can define optional hooks in `scaffold/Taskfile.yml` to run custom logic during the Kata lifecycle:

- `on_kata_start`: Runs immediately after `kata start` (e.g., `npm install`).
- `before_run`: Runs before tests execute during `kata run`.
- `on_run`: Runs after tests finish (success or failure).
- `on_run_success`: Runs only if tests pass.
- `on_run_fail`: Runs only if tests fail.
- `on_kata_finish`: Runs when `kata next` completes the final step.

Example:
```yaml
tasks:
  test:
    cmds:
      - task: runner:test
  on_kata_start:
    cmds:
      - echo "Installing dependencies..."
      - npm install
```

## Best Practices

- **Keep it Simple**: The initial scaffold should be minimal.
- **Focus on the Goal**: Steps should guide the user logically toward solving the problem.
- **Test Your Kata**: Ensure that steps can be correctly applied and passed incrementally.
