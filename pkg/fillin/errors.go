package fillin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrUndefined is the sentinel returned (or wrapped) when Strict mode
// rejects the render. Use [errors.Is] to detect it.
var ErrUndefined = errors.New("fillin: undefined placeholder(s)")

// Error is the concrete error type returned by [Template.Render] when
// one or more placeholders cannot be resolved. Fields are exported so
// callers may inspect the failure with [errors.As].
type Error struct {
	// Undefined lists every placeholder that had no matching ctx entry.
	// Sorted lexicographically so output is stable across runs.
	Undefined []string
}

// Error formats the failure so the message stays useful when multiple
// placeholders are missing.
func (e *Error) Error() string {
	switch len(e.Undefined) {
	case 0:
		return "fillin: undefined placeholder(s)"
	case 1:
		return fmt.Sprintf("fillin: undefined placeholder %q", e.Undefined[0])
	default:
		return "fillin: undefined placeholders " + joinQuoted(e.Undefined, ", ")
	}
}

// Is lets [errors.Is](err, [ErrUndefined]) succeed for any *Error.
func (e *Error) Is(target error) bool {
	return target == ErrUndefined
}

func joinQuoted(names []string, sep string) string {
	quoted := make([]string, len(names))
	for i, n := range names {
		quoted[i] = strconv.Quote(n)
	}
	return strings.Join(quoted, sep)
}
