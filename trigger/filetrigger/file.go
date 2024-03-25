package filetrigger

import "os"

type filer interface {
	IsDir(path string) bool
}

type filerFunc func(path string) bool

func (f filerFunc) IsDir(path string) bool {
	return f(path)
}

type defaultFiler struct {
}

func (d defaultFiler) IsDir(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.IsDir()
}
