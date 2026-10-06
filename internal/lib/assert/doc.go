// Package assert provides runtime assertion functions for validating
// preconditions, postconditions, and invariants in production code.
//
// Assertions detect programmer errors — conditions that should never be false
// if the code is correct. When an assertion fails, the program panics with a
// message indicating the file, line, and condition. This is appropriate for
// internal consistency checks where recovery is not meaningful.
//
// Usage:
//
//	package main
//
//	import "github.com/aisbergg/go-fillin/internal/lib/assert"
//
//	func process(items []string) {
//	    assert.True(len(items) > 0, "items must not be empty")
//	    assert.Equal(status, expectedStatus, "status mismatch after processing")
//	}
//
// Build tag: Assertions are disabled by default (no-ops). Enable with:
//
//	go build -tags assert ./...
//
// When disabled, assertion calls are compiled out, and their argument
// expressions are not evaluated. This allows assertions to be used in
// performance-sensitive code without overhead when disabled.

package assert
