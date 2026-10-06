//go:build assert

package assert

import (
	"fmt"
	"reflect"
	"runtime"
)

func fail(msg string, args ...any) {
	formatted := msg
	if len(args) > 0 {
		formatted = fmt.Sprintf(msg, args...)
	}
	_, file, line, _ := runtime.Caller(2)
	panic(fmt.Sprintf("assertion failed at %s:%d: %s", file, line, formatted))
}

// True panics if condition is false.
func True(condition bool, msg string, args ...any) {
	if condition {
		return
	}
	fail(msg, args...)
}

// False panics if condition is true.
func False(condition bool, msg string, args ...any) {
	if !condition {
		return
	}
	fail(msg, args...)
}

// Equal panics if a and b are not equal.
func Equal[T comparable](a, b T, msg string, args ...any) {
	if a == b {
		return
	}
	fail("expected %v == %v; "+msg, append([]any{a, b}, args...)...)
}

// NotEqual panics if a and b are equal.
func NotEqual[T comparable](a, b T, msg string, args ...any) {
	if a != b {
		return
	}
	fail("expected %v != %v; "+msg, append([]any{a, b}, args...)...)
}

// NotNil panics if v is nil. Handles typed nils (e.g., *int(nil)).
func NotNil(v any, msg string, args ...any) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if !rv.IsNil() {
			return
		}
	default:
		if v != nil {
			return
		}
	}
	fail("expected non-nil; "+msg, args...)
}

// Nil panics if v is not nil. Handles typed nils (e.g., *int(nil)).
func Nil(v any, msg string, args ...any) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			return
		}
	default:
		if v == nil {
			return
		}
	}
	fail("expected nil, got %v; "+msg, append([]any{v}, args...)...)
}

// Error panics if err is nil.
func Error(err error, msg string, args ...any) {
	if err != nil {
		return
	}
	fail("expected error; "+msg, args...)
}

// NoError panics if err is not nil.
func NoError(err error, msg string, args ...any) {
	if err == nil {
		return
	}
	fail("unexpected error: %v; "+msg, append([]any{err}, args...)...)
}
