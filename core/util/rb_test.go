package util

import (
	"slices"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func TestRB(t *testing.T) {
	is := is.New(t)

	rb := NewRB[string](10)
	for i := range 11 {
		rb.Append("abc" + strconv.Itoa(i))
	}

	is.Equal(slices.Collect(rb.Slice(12, 11, "")), []string(nil))
	is.Equal(slices.Collect(rb.Slice(0, 11, "pfd")), slices.Collect(func(yield func(string) bool) {
		if !yield("pfd") {
			return
		}
		for i := range 10 {
			if !yield("abc" + strconv.Itoa(i+1)) {
				return
			}
		}
	}))
	is.Equal(slices.Collect(rb.Slice(0, 10, "")), slices.Collect(func(yield func(string) bool) {
		for i := range 10 {
			if !yield("abc" + strconv.Itoa(i+1)) {
				return
			}
		}
	}))
	is.Equal(slices.Collect(rb.Slice(1, 10, "")), slices.Collect(func(yield func(string) bool) {
		for i := range 9 {
			if !yield("abc" + strconv.Itoa(i+1)) {
				return
			}
		}
	}))
}
