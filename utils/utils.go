package utils

import (
	"os"
	"path"
)

func Dir() string {
	return path.Join(os.Getenv("HOME"), ".config/rofi-gitlab")
}

func Path(fileName string) string {
	return path.Join(Dir(), fileName)
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true

	} else {
		return false
	}
}

func CreateDir() error {
	return os.MkdirAll(Dir(), 0755)
}
