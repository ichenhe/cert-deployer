package utils

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"reflect"
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

// MustReadStringOption read a string value from options map. Panic if key does not exist or the
// value is not a string type.
func MustReadStringOption(options map[string]interface{}, key string) string {
	if v, ok := options[key]; !ok {
		panic(fmt.Errorf("option '%s' does not exist", key))
	} else if s, ok := v.(string); !ok {
		panic(fmt.Errorf("option '%s' should be string, actual is %s", key,
			reflect.TypeOf(v).String()))
	} else {
		return s
	}
}

// MustReadOption read a value from options map. Panic if key does not exist or the type of
// value is not T.
func MustReadOption[T string | int | *zap.SugaredLogger](options map[string]interface{}, key string) T {
	if v, ok := options[key]; !ok {
		panic(fmt.Errorf("option '%s' does not exist", key))
	} else if s, ok := v.(T); !ok {
		var tmp T
		panic(fmt.Errorf("type of option '%s' should be %T, actual is %s", key,
			reflect.TypeOf(tmp), reflect.TypeOf(v).String()))
	} else {
		return s
	}
}

func MapSlice[S any, T any](src []S, mapper func(S) T) []T {
	n := make([]T, len(src))
	for i, e := range src {
		n[i] = mapper(e)
	}
	return n
}
