# AGENTS

This document describes how **AI agents** are expected to interact with the Kata project.

Kata is intentionally designed to be **automation-friendly**.  
AI agents are treated as first-class users of the system, alongside humans and CI pipelines.

This document defines the operational contract between Kata and an agent.

---

## Core Principle

> An agent must never infer state.  
> All relevant state is explicit and machine-readable.

Agents should rely on files, exit codes, and commands — never heuristics.

---

## Supported Interaction Model

Agents are expected to interact with Kata exclusively through:

- CLI commands
- exit codes
- files under `.kata/`

Agents **must not**:
- modify Kata internals
- bypass commands
- interpret test output semantically

---

## Primary Commands

Agents should use the following commands only.

### `kata init`

Initializes a kata repository by copying the embedded catalog and runners.

Expected behavior:
- creates `katas/`
- does not modify application code
- safe to run once per repository

---

### `kata start <kata-name>`

Starts a kata.

Expected behavior:
- applies the kata scaffold
- applies the first step
- initializes `.kata/current_state.json`

Agents should:
- verify success via exit code
- inspect `.kata/current_state.json` if needed

---

### `kata run`

Runs the kata tests.

Expected behavior:
- executes `task test`
- streams logs to stdout/stderr
- updates `.kata/current_state.json`
- writes `.kata/last_test.log`

Agents should:
- rely on the exit code to determine success or failure
- consult `.kata/last_test.log` for diagnostics if tests fail

Agents **must not** parse or interpret test output beyond failure context.

---

### `kata next`

Advances the kata to the next step.

Rules:
- only allowed when tests are passing
- applies the next step or completes the kata

Agents should:
- call `kata next` only after a successful `kata run`
- handle failure via exit code

---

## State Inspection

The authoritative state is stored in:

```text
.kata/current_state.json
```

This file contains:

- kata name
- current step
- whether tests are passing
- whether the kata is finished
- last test execution timestamp

Agents may read this file at any time.

Agents must not modify this file directly.

---

## File System Conventions

`.kata/` is owned by Kata
- application source code is owned by the user or agent
- `katas/` contains the catalog and runners

Agents must:

- treat `.kata/` as read-only except via Kata commands
- avoid deleting or rewriting files under `.kata/`

---

## Error Handling

Agents should rely on:

- command exit codes
- explicit error messages
- .kata/last_test.log

Agents must not:

- guess failure reasons
- assume partial success
- attempt recovery by manipulating internal files

---

## Non-Interactive Mode

Kata commands are non-interactive by default.

Agents should:

- avoid relying on prompts
- prefer deterministic execution paths
- use flags explicitly when available

This ensures reproducible behavior across environments.

---

## Refactoring and Iteration

After a kata is completed:

- tests may still be executed using `kata run`
- refactoring is allowed and encouraged
- `kata next` will no longer advance steps

Agents may continue improving the solution without advancing the kata.

---

## What Agents Should Not Do

Agents must not:

- edit files under `.kata/`
- introduce new test frameworks
- modify runner definitions
- bypass Taskfiles
- implement custom execution logic

All execution behavior must remain declarative and explicit.

---

## Design Rationale

Kata is designed so that:

- humans and agents use the same commands
- behavior is observable and predictable
- automation does not require special cases
- failures are explicit and inspectable

This allows agents to operate safely without privileged access or hidden mechanisms.

---

## Summary

For AI agents, Kata provides:

- explicit state
- deterministic commands
- reliable exit codes
- stable log locations

In return, agents are expected to:

- follow the defined workflow
- respect ownership boundaries
- avoid inference and heuristics
- operate only through supported interfaces
