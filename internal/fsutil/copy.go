package fsutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyDir copies a directory tree from an fs.FS to the destination path.
// It fails if a destination file already exists.
func CopyDir(src fs.FS, srcDir, dstDir string) error {
	return fs.WalkDir(src, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		target := filepath.Join(dstDir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		if _, err := os.Stat(target); err == nil {
			return os.ErrExist
		}

		return copyFile(src, path, target)
	})
}

func copyFile(src fs.FS, srcPath, dstPath string) (err error) {
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

	return nil
}
