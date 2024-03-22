package domain

import "os"

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FileReaderFunc func(name string) ([]byte, error)

func (f FileReaderFunc) ReadFile(name string) ([]byte, error) {
	return f(name)
}

// IsDir determine whether the name exists and must be a dir.
func IsDir(name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
