package domain

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMapSlice(t *testing.T) {
	mapper := func(x int, count *int) int {
		*count++
		return x + 1
	}
	type testCase struct {
		name  string
		count int
	}
	tests := []testCase{
		{"empty slice", 0},
		{"non-empty slice", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := make([]int, tt.count)
			want := make([]int, tt.count)
			for i := 0; i < tt.count; i++ {
				src[i] = i
				want[i] = i + 1
			}
			count := 0

			if got := MapSlice(src, func(s int) int {
				return mapper(s, &count)
			}); !reflect.DeepEqual(got, want) {
				t.Errorf("MapSlice() = %v, want %v", got, want)
			}
			assert.Equalf(t, tt.count, count, "mapper should be called %d time, but %d", tt.count, count)
		})
	}
}
