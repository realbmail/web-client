package utils

import (
	"os"
	"path/filepath"
)

func WorkHome(dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	baseDir := filepath.Join(currentDir, string(filepath.Separator), dir)
	return baseDir
}
