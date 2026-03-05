---
name: kata-contributor
description: Guides an agent on how to contribute to the Kata CLI codebase itself, outlining testing, linting, i18n, and commit conventions. Use this skill when asked to change or improve the Kata CLI code.
license: Apache-2.0
---

# Contributing to Kata CLI

The Kata CLI is a Go project designed to be automation-friendly, deterministic, and follow clean architecture principles. When modifying the Kata codebase, AI agents must adhere to strict quality and testing standards.

## Quality & Standards

When modifying the Kata codebase, agents **must ensure the following checks pass** before considering a task complete:

1. **Formatting**: Always run `go fmt ./...` (or equivalent) after code changes.
2. **Testing**: Always run `task test` to ensure no regressions were introduced.
3. **Linting**: Always run `task lint` to ensure code style and quality standards are met.
4. **i18n (Internationalization)**: The project uses a custom i18n system. If adding or removing translation keys (`translator.T("key")`), you **must** run the i18n tools:
   - Updates must be applied to all supported language files (e.g., `en.toml`, `pt-BR.toml`) in `internal/i18n/locales/`.
   - Run `go generate ./...` (which runs `tools/i18n-codegen`).
   - The validation tool `tools/i18n-lint` verifies keys against `internal/i18n/locales/en.toml`.

*Note: The `internal/assets/catalog` and `internal/assets/scaffold` directories contain intentionally broken code or templates and are explicitly skipped by linting and testing tools.*

## Key Architectural Principles

- **Clean Architecture**: Dependencies point inward. CLI Commands -> Application -> Domain -> Infrastructure.
- **SOLID and CUPID**: Follow these principles when designing new features or refactoring.
- **No Inference**: The binary must act as a coordinator. Do not add logic to guess failure reasons, parse test outputs, or parse language-specific semantics.
- **Explicit State**: State is explicitly managed in `.kata/current_state.json`.

## Coding Conventions

- **Godoc**: Use godoc style comments for public functions, variables, and constants.
- **Logging**: Debug messages must use the `debug` log level via `log/slog`. Direct use of `fmt.Print` is discouraged.
- **Errors**: Application-specific errors must be translated for the user at the UI/CLI boundary. Sentinel errors (e.g., `ErrTestsFailing`) remain in English for internal logic.
- **UI**: The `internal/interactive` package utilizes `github.com/charmbracelet/huh` for building interactive terminal user interfaces.

## Commit Workflow

When performing commits on the Kata codebase, follow this safety and quality protocol:

1. **Verification**: Run `task lint` and `task test`. Ensure all tests pass. Resolve any regressions, linting issues, or unused i18n keys.
2. **Staging**: Stage the relevant changes using `git add`.
3. **User Review**: Present the staged changes to the user (e.g., using `git diff --staged`) and request explicit approval.
4. **Commit**: Only execute `git commit` after the user has approved the changes. Follow conventional commit messages.
