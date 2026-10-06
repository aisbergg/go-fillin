//go:build !assert

package assert

// True panics if condition is false.
func True(condition bool, msg string, args ...any) {}

// False panics if condition is true.
func False(condition bool, msg string, args ...any) {}

// Equal panics if a and b are not equal.
func Equal[T comparable](a, b T, msg string, args ...any) {}

// NotEqual panics if a and b are equal.
func NotEqual[T comparable](a, b T, msg string, args ...any) {}

// NotNil panics if v is nil. Handles typed nils (e.g., *int(nil)).
func NotNil(v any, msg string, args ...any) {}

// Nil panics if v is not nil. Handles typed nils (e.g., *int(nil)).
func Nil(v any, msg string, args ...any) {}

// Error panics if err is nil.
func Error(err error, msg string, args ...any) {}

// NoError panics if err is not nil.
func NoError(err error, msg string, args ...any) {}
