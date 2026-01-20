package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Find all used keys in .go files
	usedKeys, err := findUsedKeys(".")
	if err != nil {
		return fmt.Errorf("failed to find used keys: %w", err)
	}

	// 2. Load defined keys from .toml files
	definedKeys, err := loadDefinedKeys("internal/i18n/locales")
	if err != nil {
		return fmt.Errorf("failed to load defined keys: %w", err)
	}

	var errors []string

	// 3. Validation: Missing Keys in Default Locale (en)
	defaultLocale := "en"
	if defaultKeys, ok := definedKeys[defaultLocale]; ok {
		for key := range usedKeys {
			if _, exists := defaultKeys[key]; !exists {
				errors = append(errors, fmt.Sprintf("Missing translation for key '%s' in %s.toml", key, defaultLocale))
			}
		}
	} else {
		errors = append(errors, fmt.Sprintf("Default locale file %s.toml not found", defaultLocale))
	}

	// 4. Validation: Unused Keys in ALL Locales
	for locale, keys := range definedKeys {
		for key := range keys {
			if _, used := usedKeys[key]; !used {
				errors = append(errors, fmt.Sprintf("Unused key '%s' in %s.toml", key, locale))
			}
		}
	}

	if len(errors) > 0 {
		for _, msg := range errors {
			fmt.Fprintln(os.Stderr, msg)
		}
		return fmt.Errorf("found %d i18n validation errors", len(errors))
	}

	fmt.Println("i18n validation passed!")
	return nil
}

// findUsedKeys scans all .go files in the root directory and subdirectories
// for calls to function named "T" with a string literal as the first argument.
func findUsedKeys(root string) (map[string]bool, error) {
	usedKeys := make(map[string]bool)
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir // Skip hidden dirs like .git
			}
			if path == "internal/assets/catalog" || path == "internal/assets/scaffold" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// Skip tests? Maybe. Tests might verify translations.
		// If a test uses a key, it counts as usage.
		// However, tests might use fake keys.
		// Let's include tests for now. If it causes noise, we can exclude.

		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing file %s: %w", path, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if it is a method call (SelectorExpr) or function call (Ident)
			// We look for method calls like translator.T("key") or t.T("key")
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				// Potentially handling `T("key")` if T is a local function?
				// Unlikely in this codebase structure.
				return true
			}

			if sel.Sel.Name == "T" {
				if len(call.Args) > 0 {
					if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						key, err := strconv.Unquote(lit.Value)
						if err == nil {
							usedKeys[key] = true
						}
					}
				}
			}
			return true
		})

		return nil
	})

	return usedKeys, err
}

// loadDefinedKeys reads all .toml files in the dir and returns map[locale]map[key]bool
func loadDefinedKeys(dir string) (map[string]map[string]bool, error) {
	definedKeys := make(map[string]map[string]bool)

	files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		locale := strings.TrimSuffix(filepath.Base(file), ".toml")

		var domainMessages map[string]map[string]string
		if _, err := toml.DecodeFile(file, &domainMessages); err != nil {
			return nil, fmt.Errorf("decoding %s: %w", file, err)
		}

		keys := make(map[string]bool)
		for domain, msgs := range domainMessages {
			for k := range msgs {
				keys[domain+"."+k] = true
			}
		}
		definedKeys[locale] = keys
	}

	return definedKeys, nil
}
