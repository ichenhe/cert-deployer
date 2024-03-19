package utils

import (
	"os"
)

// IsFile determine whether the name exists and must be a file.
func IsFile(name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// IsDir determine whether the name exists and must be a dir.
func IsDir(name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func MapSlice[S any, T any](src []S, mapper func(S) T) []T {
	n := make([]T, len(src))
	for i, e := range src {
		n[i] = mapper(e)
	}
	return n
}
