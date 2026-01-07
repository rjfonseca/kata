# Kata

> *Red tests light the path*
> *Green brings calm, refactor breathes*
> *Steps guide the journey*

Kata is a CLI tool for practicing **Test-Driven Development (TDD)** through structured, repeatable code katas.

It provides a clear workflow to initialize, start, execute, and advance katas while remaining **technology-agnostic**, **declarative**, and **automation-friendly**.

---

## What Kata Is

- A tool to **practice TDD**, not just run tests
- A framework for **incremental learning**
- A system designed to work equally well for:
  - humans
  - CI pipelines
  - AI agents

Kata does not know *how* to test.
It only knows *when tests pass*.

---

## Basic Workflow

```text
kata init
kata start <kata-name>
write code
kata run
kata next
(repeat)
```

## Key Concepts

- Scaffold: initial project structure applied once
- Steps: incremental changes, usually adding tests
- Runner: declarative definition of how tests are executed
- State: explicit, local execution state stored in .kata/

## Documentation

- 📖 Conceptual overview: docs/overview.md
- 🏗 Architecture and design decisions: docs/architecture.md
- 🧪 Writing runners: docs/runners.md
- 🧩 Writing katas: docs/katas.md
- 📐 Templates: docs/templates.md
- 🤖 AI agents: AGENTS.md
- 🤝 Contributing: CONTRIBUTING.md

## Design Principles

- TDD-first
- Declarative over imperative
- Minimal magic
- Clean Architecture
- SOLID and CUPID principles
