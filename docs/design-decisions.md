# Design Decisions

This document records the most important architectural and design decisions made in **Kata**.

Its purpose is to preserve context, explain trade-offs, and prevent already-solved problems from being repeatedly revisited as the project evolves.

This is not a full ADR system, but a lightweight and practical alternative.

---

## 1. Declarative Execution Instead of Hardcoded Test Logic

**Decision:**
Kata does not implement any test execution logic in Go.

**Rationale:**
Embedding test logic in the binary would require:
- language-specific code
- framework-specific abstractions
- constant maintenance as ecosystems evolve

By delegating execution to declarative tools (Taskfile, Docker Compose), Kata remains:
- language-agnostic
- extensible without recompilation
- predictable in automation

**Consequence:**
The binary only interprets *success or failure*, never test semantics.

---

## 2. Runners as Declarative Files, Not Code

**Decision:**
Runners are defined using files (e.g. `Taskfile.testrunner.yml`), not Go plugins.

**Rationale:**
Code-based plugins:
- introduce hidden execution paths
- complicate security and debugging
- make automation harder
- tightly couple the binary to execution details

Declarative runners:
- are transparent
- are versionable
- can be modified without recompiling Kata
- are easier for humans and AI agents to reason about

**Consequence:**
Adding a new runner requires no Go code changes.

---

## 3. Task Is Used as a Go Library

**Decision:**
Kata uses `go-task` as a Go library instead of invoking the `task` binary.

**Rationale:**
Relying on an external binary would:
- introduce environment dependencies
- cause version drift
- complicate CI and automation

Using Task as a library ensures:
- consistent behavior across environments
- no external runtime dependency
- reliable exit code propagation

**Consequence:**
The Taskfile remains the source of truth, but execution is fully controlled by the binary.

---

## 4. No YAML Merging or Interpretation

**Decision:**
Kata does not programmatically merge or manipulate YAML files.

**Rationale:**
YAML merging semantics are:
- tool-specific
- complex
- error-prone to reimplement

Instead, Kata relies on:
- Taskfile’s native include mechanism
- Docker Compose’s native override mechanism

**Consequence:**
Kata avoids duplicating logic already handled by established tools.

---

## 5. Explicit Local State

**Decision:**
Kata stores its execution state explicitly in `.kata/current_state.json`.

**Rationale:**
Inferring state from filesystem contents leads to:
- hidden rules
- fragile heuristics
- unclear automation behavior

Explicit state provides:
- clarity
- debuggability
- safe automation
- predictable behavior for AI agents

**Consequence:**
All lifecycle decisions (`run`, `next`) depend only on explicit state.

---

## 6. No Rollback or Arbitrary Navigation

**Decision:**
Kata does not support rollback or arbitrary step navigation.

**Rationale:**
Rollback introduces:
- ambiguous state transitions
- complex filesystem diffs
- unclear test expectations

Kata models a forward-only TDD flow:
- red → green → refactor → advance

**Consequence:**
The mental model stays simple and aligned with TDD.

---

## 7. `.tmpl` Files Instead of Embedded Modules

**Decision:**
Scaffold files that would otherwise violate Go `embed` constraints (e.g. `go.mod`) are stored as `.tmpl`.

**Rationale:**
Go does not allow embedding files that belong to a different module.

Using `.tmpl`:
- allows full scaffolds to be embedded
- avoids multiple Go modules in the binary
- enables future dynamic rendering

**Consequence:**
Templates are materialized during `kata start`.

---

## 8. Templates Are Minimal by Design

**Decision:**
Templates are intentionally simple and use Go’s standard `text/template`.

**Rationale:**
Complex templating systems:
- become DSLs
- increase cognitive load
- obscure generated output

Kata templates are meant for:
- light parameterization
- not logic-heavy generation

**Consequence:**
Generated files remain readable and predictable.

---

## 9. No Mandatory Interactivity

**Decision:**
All commands are non-interactive by default.

**Rationale:**
Interactivity:
- breaks automation
- complicates CI
- introduces branching UX logic

Interactive behavior is optional and layered on top.

**Consequence:**
Kata works equally well for humans, scripts, CI, and AI agents.

---

## 10. The Binary Is a Coordinator

**Decision:**
The Kata binary coordinates actions but does not perform domain-specific work.

**Rationale:**
Centralizing logic in the binary would:
- increase coupling
- reduce flexibility
- make evolution harder

Delegation keeps the system:
- composable
- predictable
- easy to extend

**Consequence:**
Most behavior lives in files, not code.

---

## Summary

These decisions intentionally trade flexibility in *how* things are implemented for clarity in *how the system behaves*.

The result is a tool that:
- is simple to reason about
- scales across languages and ecosystems
- remains friendly to automation and AI agents
- avoids hidden complexity
