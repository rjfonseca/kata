package scaffold

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Copier applies scaffolds and steps to the project root.
type Copier struct {
	Root     string
	Manifest *Manifest
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
		// If the source file is a template, remove the .tmpl extension
		if strings.HasSuffix(target, ".tmpl") {
			// TODO(template): render .tmpl files using text/template
			target = strings.TrimSuffix(target, ".tmpl")
		}

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		return c.copyFile(src, path, target, sourceLabel)
	})
}

func (c *Copier) copyFile(src fs.FS, srcPath, dstPath, sourceLabel string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(dstPath); err == nil {
		return c.handleOverwrite(src, srcPath, dstPath, sourceLabel)
	}

	return c.writeNewFile(src, srcPath, dstPath, sourceLabel)
}

func (c *Copier) handleOverwrite(src fs.FS, srcPath, dstPath, sourceLabel string) error {
	entry, known := c.Manifest.Files[dstPath]
	if !known {
		return errors.New("refusing to overwrite user-managed file: " + dstPath)
	}

	current, err := ChecksumFile(dstPath)
	if err != nil {
		return err
	}

	if current != entry.Checksum {
		return errors.New("file modified by user: " + dstPath)
	}

	return c.writeNewFile(src, srcPath, dstPath, sourceLabel)
}

func (c *Copier) writeNewFile(src fs.FS, srcPath, dstPath, sourceLabel string) error {
	in, err := src.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

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
