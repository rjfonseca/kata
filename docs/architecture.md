# Architecture

This document describes the architectural principles, boundaries, and design decisions behind **Kata**.

It is intended for maintainers, contributors, and anyone interested in understanding *why* the system is designed the way it is — not just *how* it works.

---

## Architectural Goals

Kata is designed to:

- Support **Test-Driven Development (TDD)** explicitly
- Remain **technology-agnostic**
- Be **automation-friendly** (CI, scripts, AI agents)
- Favor **declarative configuration** over imperative logic
- Keep the core binary **small, predictable, and extensible**

These goals shape every major architectural decision.

---

## High-Level Design

At a high level, Kata follows **Clean Architecture** principles:

```text
┌───────────────────────────┐
│        CLI Commands       │  ← orchestration, flags, UX
├───────────────────────────┤
│        Application        │  ← use cases (start, run, next)
├───────────────────────────┤
│           Domain          │  ← kata state and rules
├───────────────────────────┤
│     Infrastructure        │  ← filesystem, task runner
└───────────────────────────┘

Dependencies always point inward.

## What the Binary Does (and Does Not Do)

### The Binary Does

- Manage kata lifecycle (init, start, run, next)
- Apply scaffolds and steps
- Track and persist kata state
- Delegate test execution to external tools
- Interpret test success or failure

### The Binary Does Not

- Know how tests are written
- Know how tests are executed
- Parse test output
- Contain any language-specific logic
- Contain Docker, HTTP, or framework semantics

This separation is intentional and critical.

## Declarative Execution Model

### Why Declarative?

Imperative test execution logic quickly becomes brittle:

- every new language requires code
- every new tool requires branching logic
- automation becomes harder to reason about

Kata avoids this by using a declarative execution model.

### Runners

A **runner** defines *how* tests are executed.

Runners are implemented using:

- `Taskfile.testrunner.yml`
- optional supporting files (e.g. `compose.override.yaml`)

The Kata binary simply asks:

> “Did the tests pass?”

It does not care how that answer was produced.

### Task as a Library

Kata uses **go-task** as a Go library, not as an external binary.

This provides:

- no external dependencies
- consistent behavior across environments
- reliable exit codes
- direct integration with the Go runtime

The Taskfile remains the source of truth for execution logic.

## Domain Model: State

The domain core of Kata is the kata state, stored in `.kata/current_state.json`

The state explicitly records:

- which kata is active
- which step is current
- whether tests are passing
- whether the kata is finished
- when tests were last executed

### Why Explicit State?

No hidden behavior

- No inference from filesystem state
- Easy debugging
- Safe automation
- Clear interaction for AI agents

## Advancing Steps

The `next` operation is governed by simple domain rules:

- steps can only advance if tests are passing
- steps advance sequentially
- once no steps remain, the kata is marked as finished
- finished katas can still run tests and refactor code

These rules live in the domain layer, not in CLI commands.

## Templates

Kata supports scaffold files with the `.tmpl` extension.

Current Behavior

- `.tmpl` files are materialized during kata start
- the `.tmpl` extension is removed
- file contents are copied as-is

### Future Behavior

Templates will be rendered using Go’s `text/template` with a dynamic context derived from:

- kata configuration
- user configuration
- repository metadata
- allowed environment variables

Templates exist to enable flexible scaffolds, not to introduce a DSL.

## Why Not Plugins?

Kata intentionally avoids a plugin system.

Plugins tend to:

- blur responsibility boundaries
- hide execution paths
- complicate security and debugging
- make automation harder

Declarative configuration provides most of the benefits of plugins with fewer downsides.

## Why Not Combine YAML Files?

Kata does not merge or combine YAML files programmatically.

Instead, it relies on:

- native include mechanisms (`Taskfile`)
- native override mechanisms (`docker compose`)

This avoids reimplementing semantics already handled by existing tools.

## Automation and AI Agents

Kata is designed to be safe and predictable for automation:

- all commands are non-interactive by default
- state is explicit and machine-readable
- logs are written to stable locations
- exit codes are meaningful

This makes Kata suitable for CI pipelines and AI-assisted workflows.

## Summary

Kata’s architecture is built around a few strong constraints:

- the binary is a coordinator, not an executor
- execution logic lives in declarative files
- state is explicit and local
- behavior is predictable and inspectable

These constraints enable flexibility without sacrificing clarity, and allow the system to grow without accumulating hidden complexity.
