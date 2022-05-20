package utils

import "fmt"

// ErrorCollection is an error that packed a series of errors. Unlike fmt.Errorf, all the errors
// here are at the same level. The length of errors must > 0.
//
// A typical use case is a function that performs operations on a lot of data, which are
// independent of each other and do not want to be interrupted by just an error.
type ErrorCollection struct {
	errors []error
}

// NewErrorCollection create an ErrorCollection based on given errors. If the given errors is nil
// or empty, then return nil.
func NewErrorCollection(errs []error) *ErrorCollection {
	if errs != nil && len(errs) > 0 {
		return &ErrorCollection{errors: errs}
	} else {
		return nil
	}
}

func (e *ErrorCollection) Error() string {
	s := fmt.Sprintf("%d errors:\n", len(e.errors))
	for i, err := range e.errors {
		s += fmt.Sprintf("No. %d: %s", i, err.Error())
		if i != len(e.errors)-1 {
			s += "\n"
		}
	}
	return s
}
