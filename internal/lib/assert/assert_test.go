//go:build assert

package assert_test

import (
	"testing"

	"github.com/aisbergg/go-fillin/internal/lib/assert"
)

func TestTrue(t *testing.T) {
	assert.True(true, "condition is true")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	assert.True(false, "condition should be true")
}

func TestFalse(t *testing.T) {
	assert.False(false, "condition is false")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	assert.False(true, "condition should be false")
}

func TestEqual(t *testing.T) {
	assert.Equal(42, 42, "values are equal")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	assert.Equal(42, 43, "values should match")
}

func TestNotEqual(t *testing.T) {
	assert.NotEqual(42, 43, "values are different")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	assert.NotEqual(42, 42, "values should differ")
}

func TestNotNil(t *testing.T) {
	v := &struct{}{}
	assert.NotNil(v, "pointer is not nil")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	var nilPtr *int
	assert.NotNil(nilPtr, "pointer should not be nil")
}

func TestNil(t *testing.T) {
	var nilPtr *int
	assert.Nil(nilPtr, "pointer is nil")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	v := &struct{}{}
	assert.Nil(v, "should be nil")
}

func TestError(t *testing.T) {
	err := assertError("test error")
	assert.Error(err, "error is not nil")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	assert.Error(nil, "should have error")
}

func TestNoError(t *testing.T) {
	assert.NoError(nil, "no error")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	err := assertError("test error")
	assert.NoError(err, "should not have error")
}

func TestFormatArgs(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		expected := "assertion failed at "
		got := r.(string)
		if !contains(got, expected) || !contains(got, "user 42 not found") {
			t.Fatalf("unexpected panic message: %s", got)
		}
	}()
	assert.Equal(42, 43, "user %d not found", 42)
}

func assertError(msg string) error {
	return &testError{msg: msg}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
