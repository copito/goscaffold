package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

// exists returns whether the given file or directory exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// PathCopy copies file/folder to another location
func PathCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// Get delta relative path
func DeltaRelativePath(basePath string, currentPath string) string {
	// Find common prefix
	baseDir := filepath.ToSlash(filepath.Clean(basePath))
	currentDir := filepath.ToSlash(filepath.Clean(currentPath))

	prefixLen := 0
	for i := 0; i < len(baseDir) && i < len(currentDir); i++ {
		if baseDir[i] != currentDir[i] {
			break
		}
		prefixLen++
	}

	// Construct relative path
	relPath := currentDir[prefixLen:]
	if relPath != "" && relPath[0] == '/' {
		relPath = relPath[1:]
	}
	return relPath
}
