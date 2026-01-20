# Kata – Overview

Kata is a CLI tool designed to support **Test-Driven Development (TDD)** through structured, repeatable coding exercises known as *code katas*.

The tool provides a clear and explicit workflow to create, start, execute, and progress through katas, while remaining **technology-agnostic**, **declarative**, and **automation-friendly**.

This document explains the **mental model** behind Kata, how its main concepts fit together, and how the overall workflow is intended to be used.

---

## Mental Model

At its core, Kata enforces a simple idea:

> The tool does not know *how* to test.  
> It only knows *when tests pass*.

Everything else — languages, frameworks, test tools, containers — is delegated to declarative configuration provided by each kata.

---

## Core Concepts

### Kata

A kata is a structured coding exercise composed of:

- an initial **scaffold**
- a sequence of incremental **steps**
- a **runner** that defines how tests are executed

Each kata is self-contained and versioned alongside its tests.

---

### Scaffold

The scaffold represents the initial project structure applied when a kata is started.

It may include:

- source code
- configuration files
- Docker Compose setup
- Taskfile definitions

The scaffold is applied **once**, during `kata start`.

---

### Steps

Steps define the incremental evolution of the kata.

Each step typically:

- adds or modifies acceptance tests
- introduces a new requirement

Steps are applied sequentially using `kata next`.

A kata can only advance to the next step when tests are passing.

---

### Runner

A runner defines **how tests are executed**, not what is being tested.

Runners are implemented declaratively using:

- `Taskfile.testrunner.yml`
- optional supporting files (for example, `compose.override.yaml`)

The Kata binary does not embed any testing logic.  
Different katas may use different runners.

---

## Directory Structure

A typical kata repository looks like this:

```text
katas/
├── catalog/
│   └── <kata-name>/
│       ├── README.md
│       ├── config.toml
│       ├── scaffold/
│       └── steps/
│           ├── 01-...
│           ├── 02-...
│           └── ...
└── runners/
    └── <runner-name>/
        ├── Taskfile.testrunner.yml
        └── ...

```
During execution, Kata creates a local working directory:

```text
.kata/
├── current_state.json
├── manifest.json
└── last_run.log
```

The .kata/ directory contains execution state and logs and must not be committed.

## Workflow
The intended workflow follows classic TDD:
```text
kata init (for the first time to initialize the kata repo)

kata start <kata-name>

(write code)
kata run
kata next
(repeat)
```

## Commands

### `kata init`

Initializes a kata repository by copying:
- the embedded kata catalog
- the embedded runners

This allows users to extend, customize, or version their own katas and runners locally.

### `kata start <kata-name>`

Starts a kata by:

1. Resolving the kata from the local catalog (found in **Workspace Root**) or the embedded catalog.
2. Copying the kata scaffold into the **Project Root** (default: current directory).
3. Copying the configured runner into the Project Root.
4. Initializing the kata state.

**Flags:**
- `--dir <path>`: Create and start the kata in a specific subdirectory (Project Root).

The kata state is stored in `.kata/current_state.json` inside the Project Root.

### `kata run`

Runs the kata tests by executing the test task defined in the kata’s Taskfile.

Internally, Kata uses the Task Go library, not the task binary, so no external dependency is required.

This command:

- streams test output directly to the terminal
- updates the kata state based on the test result
- stores the last execution log in `.kata/last_run.log`

Exit codes are propagated correctly, making the command suitable for CI pipelines and automation.

### `kata next`

Advances the kata to the next step.

Rules:
- advancement is only allowed when tests are passing
- if no steps remain, the kata is marked as completed
- after completion, users may continue refactoring and running tests

### `kata task [name] [args...]`

Executes an arbitrary task defined in the kata's `Taskfile.yml`.

- If `name` is provided, runs that specific task.
- If no name is provided, lists all available tasks.
- Arguments can be passed to the task.

Example:
```bash
kata task lint
kata task benchmark --duration=5s
```

## Templates

Some files in a kata scaffold may be defined as templates using the .tmpl extension.

### How templates work

- Files ending with `.tmpl` are materialized during `kata start`
- The `.tmpl` extension is removed in the destination file
- Contents are rendered using Go's `text/template`

Example:

```
go.mod.tmpl → go.mod
```

Templates have access to dynamic variables provided by:

- kata configuration
- user configuration
- repository metadata


## Design Principles

Kata is built around the following principles:

- TDD-first: the workflow enforces red → green → refactor
- Declarative over imperative: behavior is defined in files, not code
- Technology-agnostic: no test framework or language is hardcoded
- Automation-friendly: works in CI and with AI agents
- Minimal magic: explicit files and predictable behavior

Architecturally, the project follows:

- Clean Architecture
- SOLID principles
- CUPID principles

## Intended Audience

Kata is designed for:

- engineers practicing TDD
- teams onboarding new developers
- technical leaders running workshops
- CI pipelines validating kata progress
- AI agents assisting with coding exercises
