# Translations

The `kata` CLI is fully internationalized (i18n). This document explains how the translation system works and how you can add new languages or customize existing translations.

## How it Works

The system uses the standard `golang.org/x/text/language` package to provide robust language matching and fallback.

1.  **Supported Languages:** The application automatically discovers all built-in translations (e.g., `en`, `pt-BR`).
2.  **Language Detection:** When you run a command with the `--lang` flag (e.g., `--lang pt`), you specify your desired language.
3.  **Best Match:** The system compares your desired language against the list of supported languages and picks the best possible match.
    *   If you ask for `pt`, and the application supports `pt-BR`, it will intelligently select `pt-BR`.
    *   If you ask for a language that is not supported, it will fall back to the default language, English (`en`).

## File Structure

Translation keys are organized into domain-based TOML files.

-   **Domains:** A domain typically corresponds to a command (e.g., `init`, `start`, `next`). There is also a `cli` domain for global application text, such as top-level flags and usage descriptions.
-   **Keys:** Within a domain, keys are named using `snake_case` and are prefixed to indicate their purpose.

### Key Prefixes

-   `usage`: The primary usage description for a command.
-   `args_usage`: A description of the command's arguments.
-   `flag_`: The help text for a command-line flag.
-   `error_`: An error message.
-   `log_`: An informational message printed to the console.

**Example from `en.toml`:**
```toml
# In this example, 'start' is the domain.
[start]
usage = "Start a kata"
flag_force_usage = "Overwrite existing kata state"
error_kata_name_required = "A kata name is required"
log_kata_started = "Kata started successfully"
```

In the code, a key is referenced as `domain.key`, for example: `translator.T("start.error_kata_name_required")`.

## Overriding Translations

You can easily override any translation without modifying the application's source code.

When you run `kata init`, a `katas/i18n` directory is created in your project. To override a translation, you simply create a TOML file in this directory that matches the language you want to change.

For example, to override the `cli.usage` message for Brazilian Portuguese, you would:

1.  Create or edit the file `katas/i18n/pt-BR.toml`.
2.  Add the specific key you want to override.

**Example `katas/i18n/pt-BR.toml`:**
```toml
# My custom overrides for Brazilian Portuguese
[cli]
usage = "Meu jeito customizado de usar o kata"

[start]
error_kata_name_required = "O nome do kata é obrigatório, por favor!"
```

When you run `kata --lang pt-BR`, the application will:
1.  Load the default, built-in `pt-BR` translations.
2.  Load your custom `katas/i18n/pt-BR.toml` file.
3.  Your custom messages will overwrite the default ones.

Any keys you *don't* specify in your override file will continue to use the default built-in translations.

## Adding a New Language

To add a completely new language (e.g., French), you would:
1.  Copy `katas/i18n/en.toml` to a new file, `katas/i18n/fr.toml`.
2.  Translate all the values in the file to French.
3.  Run the application with `--lang fr`.

The application will automatically detect and use your new translation file.

## Adding New Embedded Languages

If you are contributing to the `kata` CLI and want to add a new language directly into the application's built-in set of translations, follow these steps:

1.  **Create the Translation File:**
    *   Create a new TOML file named after the language tag (e.g., `fr.toml` for French) in the `internal/assets/i18n/locales/` directory.
    *   Populate this file with all the necessary translations, following the domain and key prefixing conventions (e.g., `[cli]`, `[init]`, `error_`, `log_`). You can use `internal/assets/i18n/locales/en.toml` as a template.

2.  **Regenerate Embedded Locales:**
    *   Navigate to the `internal/i18n` directory in your terminal:
        ```bash
        cd internal/i18n
        ```
    *   Run the `go generate` command. This command executes a tool that reads all the `.toml` files in the `locales/` directory and embeds them into `locales.go`:
        ```bash
        go generate
        ```
    *   This step is crucial for the application to recognize and include your new language in its built-in translations.

3.  **Test Your New Language:**
    *   Once `locales.go` is regenerated, you can test your new language by running the application with the `--lang` flag:
        ```bash
        go run ./cmd/kata/ --lang fr
        ```

This process ensures that your new language is compiled directly into the `kata` CLI, making it available out-of-the-box for all users.
