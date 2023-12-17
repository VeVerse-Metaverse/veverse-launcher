// Package utils provides utility functions for the launcher.
package utils

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractArchive extracts the given archive to the given destination path.
func ExtractArchive(ctx context.Context, archivePath string, destinationPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to open archive: %s", err)
		return fmt.Errorf("failed to open archive: %w", err)
	}

	defer func(r *zip.ReadCloser) {
		if err1 := r.Close(); err1 != nil {
			runtime.LogErrorf(ctx, "failed to close archive: %s", err1)
		}
	}(r)

	err = os.MkdirAll(destinationPath, 0755)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to create destination directory: %s", err)
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			runtime.LogErrorf(ctx, "failed to open file in archive: %s", err)
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		defer func(rc io.ReadCloser) {
			if err1 := rc.Close(); err1 != nil {
				runtime.LogErrorf(ctx, "failed to close archive file: %s", err1)
			}
		}(rc)

		path := filepath.Join(destinationPath, f.Name)

		if !strings.HasPrefix(path, filepath.Clean(destinationPath)+string(os.PathSeparator)) {
			runtime.LogErrorf(ctx, "illegal file path: %s", path)
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				runtime.LogErrorf(ctx, "failed to create directory: %s", err)
				return fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				runtime.LogErrorf(ctx, "failed to create directory: %s", err)
				return fmt.Errorf("failed to create directory: %w", err)
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				runtime.LogErrorf(ctx, "failed to open file: %s", err)
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer func(f *os.File) {
				if err1 := f.Close(); err1 != nil {
					runtime.LogErrorf(ctx, "failed to close file: %s", err1)
				}
			}(f)

			_, err = io.Copy(f, rc)
			if err != nil {
				runtime.LogErrorf(ctx, "failed to write file: %s", err)
				return fmt.Errorf("failed to write file: %w", err)
			}
		}

		return nil
	}

	for _, f := range r.File {
		err = extractAndWriteFile(f)
		if err != nil {
			runtime.LogErrorf(ctx, "failed to extract file: %s", err)
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	return nil
}
