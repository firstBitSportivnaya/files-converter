package fileutil

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dest.
func CopyFile(src, dest string) error {
	sourceFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceFileInfo.Mode())
	if err != nil {
		return err
	}
	defer func() {
		_ = destFile.Close()
	}()

	_, err = io.Copy(destFile, sourceFile)

	return err
}

func CopyDir(src, dest string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return nil
		}
		desPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(desPath, info.Mode())
		}

		return CopyFile(path, desPath)
	})

	return err
}
