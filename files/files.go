package files

import (
	"os"
)

func DirExists(location string) bool {
	if _, err := os.Open(location); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateDirIfNotExists(loc string) error {
	if DirExists(loc) {
		return nil
	}
	return os.Mkdir(loc, os.ModeDir)
}

func FileSep() string {
	return string(os.PathSeparator)
}
