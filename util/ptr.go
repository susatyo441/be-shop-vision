package util

import "github.com/susatyo441/go-ta-utils/functions"

func Ptr[T any](v T) *T {
	return functions.MakePointer(v)
}
