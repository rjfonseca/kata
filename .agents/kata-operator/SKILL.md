---
name: kata-operator
description: Guides agents on how to initialize, start, run, and advance Katas using the Kata CLI, and how to interpret the .kata/ state correctly. Use this skill when the user asks you to practice a kata, run tests, or use the kata CLI.
license: Apache-2.0
---

# Operating the Kata CLI

The Kata CLI is a tool for practicing Test-Driven Development (TDD). As an AI agent, you can operate Katas on behalf of the user by using the CLI and inspecting the explicit state.

## Core Principle

**Never infer state.** Rely exclusively on files, exit codes, and explicit commands.
Kata relies on declarative execution. Do not edit `.kata/` files, introduce new test frameworks, or bypass the Taskfile.

## Basic Workflow

The expected workflow for a Kata is:

1. `kata init` (once per repository to copy embedded katas/runners)
2. `kata start <kata-name>`
3. Write/fix code
4. `kata run`
5. `kata next`
6. Repeat steps 3-5 until the Kata is finished.

## CLI Commands

### 1. `kata init`
Initializes the workspace by copying the catalog of Katas and runners into `katas/`.
- Expected behavior: Creates `katas/` directory. Safe to run multiple times.

### 2. `kata start <kata-name>`
Starts a specific Kata.
- Expected behavior: Applies the Kata scaffold, applies the first step, and creates `.kata/current_state.json`.
- Agent action: Run this command, check the exit code, and verify `.kata/current_state.json` is created.

### 3. `kata run`
Executes the Kata's test suite.
- Expected behavior: Runs tests (via `task test`), outputs logs to stdout/stderr, updates `.kata/current_state.json`, and writes `.kata/last_run.log`.
- Agent action: **Rely on the exit code** to determine success (0) or failure (non-zero). If it fails, read `.kata/last_run.log` to understand why. Do not try to parse test output semantically beyond understanding the failure context.

### 4. `kata next`
Advances the Kata to the next step.
- Expected behavior: Only allowed when tests are passing. Applies the next step's files or completes the Kata.
- Agent action: Call this **only** after a successful `kata run` (exit code 0). Check the exit code of `kata next`.

## Inspecting State

The authoritative state of the active Kata is stored in `.kata/current_state.json`.
This file contains the kata name, current step, passing status, and whether the Kata is finished.

**Rules for State:**
- You may read `.kata/current_state.json` at any time.
- You **must not** modify `.kata/current_state.json` or any file inside `.kata/`. The `.kata/` directory is owned by the Kata CLI.

## Error Handling

If a command fails:
1. Check the exit code.
2. Read `.kata/last_run.log` or the command's stderr.
3. Fix the code in the user's project (not the `.kata/` folder).
4. Do not guess failure reasons or try to manipulate internal Kata state.

## Finishing a Kata

When `kata next` indicates no more steps, the Kata is finished. You can continue to run tests using `kata run` and refactor the user's code, but `kata next` will no longer advance steps.
