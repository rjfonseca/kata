package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rjfonseca/kata/internal/i18n"
)

// TemplateData holds the data that can be passed to templates for rendering.
type TemplateData struct {
	KataName string
	// Add other fields as needed, e.g., ModuleName string
}

// Copier applies scaffolds and steps to the project root.
type Copier struct {
	Root         string
	Manifest     *Manifest
	TemplateData TemplateData
	Translator   i18n.Translator
}

// Apply copies all files from the root of src into the project root.
// The src FS must represent the root directory to be applied.
func (c *Copier) Apply(src fs.FS, sourceLabel string) error {
	return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		target := filepath.Join(c.Root, path)
		isTemplate := strings.HasSuffix(target, ".tmpl")
		if isTemplate {
			target = strings.TrimSuffix(target, ".tmpl")
		}

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		if isTemplate {
			return c.renderTemplate(src, path, target, sourceLabel)
		}

		return c.copyFile(src, path, target, sourceLabel)
	})
}

func (c *Copier) renderTemplate(src fs.FS, srcPath, dstPath, sourceLabel string) (err error) {
	templateContent, err := fs.ReadFile(src, srcPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(srcPath)).Parse(string(templateContent))
	if err != nil {
		return err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, c.TemplateData); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(dstPath); err == nil {
		return c.handleOverwrite(src, srcPath, dstPath, sourceLabel)
	}

	return c.writeNewFileFromBytes(rendered.Bytes(), dstPath, sourceLabel)
}

func (c *Copier) copyFile(src fs.FS, srcPath, dstPath, sourceLabel string) error {
	in, err := src.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := in.Close(); err == nil {
			err = cerr
		}
	}()

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	checksum, err := ChecksumFile(dstPath)
	if err != nil {
		return err
	}

	c.Manifest.Files[dstPath] = FileEntry{
		Checksum: checksum,
		Source:   sourceLabel,
	}

	return nil
}

func (c *Copier) handleOverwrite(src fs.FS, srcPath, dstPath, sourceLabel string) error {
	entry, known := c.Manifest.Files[dstPath]
	if !known {
		return fmt.Errorf("%s: %s", c.Translator.T("copier.error_overwrite"), dstPath)
	}

	current, err := ChecksumFile(dstPath)
	if err != nil {
		return err
	}

	if current != entry.Checksum {
		return fmt.Errorf("%s: %s", c.Translator.T("copier.error_modified"), dstPath)
	}

	// If it was a template, render and overwrite
	if strings.HasSuffix(srcPath, ".tmpl") {
		templateContent, err := fs.ReadFile(src, srcPath)
		if err != nil {
			return err
		}
		tmpl, err := template.New(filepath.Base(srcPath)).Parse(string(templateContent))
		if err != nil {
			return err
		}
		var rendered bytes.Buffer
		if err := tmpl.Execute(&rendered, c.TemplateData); err != nil {
			return err
		}
		return c.writeNewFileFromBytes(rendered.Bytes(), dstPath, sourceLabel)
	}

	return c.writeNewFile(src, srcPath, dstPath, sourceLabel)
}

func (c *Copier) writeNewFile(src fs.FS, srcPath, dstPath, sourceLabel string) (err error) {
	in, err := src.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close() //nolint: errcheck

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close() //nolint: errcheck

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	checksum, err := ChecksumFile(dstPath)
	if err != nil {
		return err
	}

	c.Manifest.Files[dstPath] = FileEntry{
		Checksum: checksum,
		Source:   sourceLabel,
	}

	return nil
}

func (c *Copier) writeNewFileFromBytes(content []byte, dstPath, sourceLabel string) (err error) {
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close() //nolint: errcheck

	if _, err := out.Write(content); err != nil {
		return err
	}

	checksum, err := ChecksumFile(dstPath)
	if err != nil {
		return err
	}

	c.Manifest.Files[dstPath] = FileEntry{
		Checksum: checksum,
		Source:   sourceLabel,
	}

	return nil
}
