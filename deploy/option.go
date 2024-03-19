package deploy

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

const keyLogger = "cert-deployer.logger"

type Options map[string]interface{}

type InvalidOptionError struct {
	Message string
}

func NewInvalidOptionError(message string) *InvalidOptionError {
	return &InvalidOptionError{Message: message}
}

func (e *InvalidOptionError) Error() string {
	return e.Message
}

func RecoverFromInvalidOptionError(handler func(err *InvalidOptionError)) {
	if v := recover(); v != nil {
		if e, ok := v.(*InvalidOptionError); ok {
			handler(e)
		} else {
			panic(v)
		}
	}
}

// MustReadString reads a string value from options. See MustReadOption.
func (o Options) MustReadString(key string) string {
	return MustReadOption[string](o, key)
}

// MustReadLogger gets the logger from options. See MustReadOption.
func (o Options) MustReadLogger() *zap.SugaredLogger {
	return MustReadOption[*zap.SugaredLogger](o, keyLogger)
}

// MustReadOption reads a value from options. Panic with InvalidOptionError if key does not exist or
// the type of value is not T.
func MustReadOption[T any | *zap.SugaredLogger](opt Options, key string) T {
	if v, ok := opt[key]; !ok {
		panic(NewInvalidOptionError(fmt.Sprintf("option '%s' does not exist", key)))
	} else if s, ok := v.(T); !ok {
		var tmp T
		panic(NewInvalidOptionError(fmt.Sprintf("type of option '%s' should be %T, actual is %s", key,
			reflect.TypeOf(tmp), reflect.TypeOf(v).String())))
	} else {
		return s
	}
}
