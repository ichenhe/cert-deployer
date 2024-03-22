package domain

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestOptions_MustReadString(t *testing.T) {
	options := Options(map[string]any{
		"str": "hello",
		"int": 12,
	})
	tests := []struct {
		name  string
		key   string
		panic bool
		want  string
	}{
		{name: "ok", key: "str", panic: false, want: "hello"},
		{name: "type mismatch", key: "int", panic: true},
		{name: "key not exist", key: "x", panic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panic {
				assert.Equalf(t, tt.want, options.MustReadString(tt.key), "MustReadString(%v)", tt.key)
			} else {
				assert.Panics(t, func() { options.MustReadString(tt.key) })
			}
		})
	}
}

func TestMustReadOption(t *testing.T) {
	options := Options(map[string]any{
		"str": "hello",
		"int": 12,
	})
	type testCase[T interface{ any | *zap.SugaredLogger }] struct {
		name  string
		key   string
		panic bool
		want  T
	}
	tests := []testCase[int]{
		{name: "ok", key: "int", panic: false, want: 12},
		{name: "type mismatch", key: "str", panic: true},
		{name: "key not exist", key: "x", panic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panic {
				assert.Equalf(t, tt.want, MustReadOption[int](options, tt.key), "MustReadOption[%v]", tt.key)
			} else {
				assert.Panics(t, func() { MustReadOption[int](options, tt.key) })
			}
		})
	}
}

func TestRecoverFromInvalidOptionError(t *testing.T) {
	f := func() (err error) {
		options := Options(map[string]any{})
		defer RecoverFromInvalidOptionError(func(e *InvalidOptionError) {
			err = e
		})
		options.MustReadString("x")
		return
	}
	assert.NotNil(t, f(), "expected a cached error, but nil")
}
