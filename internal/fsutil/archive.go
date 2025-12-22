package fsutil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CreateArchive creates a .tar.gz file from the source directory.
// It skips any paths provided in the ignore list.
func CreateArchive(srcDir string, destFile string, ignorePaths []string) (err error) {
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	gw := gzip.NewWriter(f)
	defer func() {
		if cerr := gw.Close(); err == nil {
			err = cerr
		}
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		if cerr := tw.Close(); err == nil {
			err = cerr
		}
	}()

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		// Check if path should be ignored
		for _, ignore := range ignorePaths {
			if rel == ignore || strings.HasPrefix(rel, ignore+string(filepath.Separator)) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		_, err = io.Copy(tw, file)
		return err
	})
}

// ExtractArchive unpacks a .tar.gz file into the destination directory.
func ExtractArchive(srcFile string, destDir string) (err error) {
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := gr.Close(); err == nil {
			err = cerr
		}
	}()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, filepath.FromSlash(header.Name))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(f, tr)
			closeErr := f.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		default:
			return fmt.Errorf("unsupported type: %c in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}