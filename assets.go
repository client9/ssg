package ssg

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CopyAssets copies all files from srcDir to destDir, preserving the
// directory structure. Directories and files whose names begin with '.'
// are skipped.
func CopyAssets(srcDir, destDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("CopyAssets walkdir error at %q: %w", path, err)
		}
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(destDir, rel)

		if err := os.MkdirAll(filepath.Dir(dest), 0750); err != nil {
			return fmt.Errorf("CopyAssets: mkdir %s: %w", filepath.Dir(dest), err)
		}
		return copyFile(path, dest)
	})
}

func copyFile(src, dest string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, in)
	return err
}
