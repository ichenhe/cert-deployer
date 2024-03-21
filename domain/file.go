package domain

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FileReaderFunc func(name string) ([]byte, error)

func (f FileReaderFunc) ReadFile(name string) ([]byte, error) {
	return f(name)
}
