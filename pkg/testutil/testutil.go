package testutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break // reached the root directory
		}
		dir = parentDir
	}
	panic("project root with go.mod file not found")
}

func SetupTest(t *testing.T, pathsToCopy map[string]string) (string, func()) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		_ = os.Chdir(wd)
	}
	tempDir := t.TempDir()
	projectRoot := findProjectRoot()

	for source, destination := range pathsToCopy {
		sourcePath := strings.ReplaceAll(source, "@", projectRoot)
		if destination == "" {
			destination = filepath.Join(tempDir, filepath.Base(sourcePath))
		} else {
			destination = filepath.Join(tempDir, destination)
		}

		if err = copyPath(sourcePath, destination); err != nil {
			panic(err)
		}
	}
	if err := os.Chdir(tempDir); err != nil {
		panic(err)
	}
	return tempDir, cleanup
}

func copyPath(sourcePath, destinationPath string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDirectory(sourcePath, destinationPath)
	}
	return copyFile(sourcePath, destinationPath)
}

func copyDirectory(sourceDir, destinationDir string) error {
	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return err
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceEntryPath := filepath.Join(sourceDir, entry.Name())
		destinationEntryPath := filepath.Join(destinationDir, entry.Name())

		if entry.IsDir() {
			if err := copyDirectory(sourceEntryPath, destinationEntryPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourceEntryPath, destinationEntryPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(sourceFile, destinationFile string) error {
	if err := os.MkdirAll(filepath.Dir(destinationFile), os.ModePerm); err != nil {
		return err
	}

	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
