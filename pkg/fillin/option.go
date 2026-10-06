package fillin

import (
	"fmt"
)

// ErrorMode controls what happens when a placeholder has no entry in
// the ctx passed to [Template.Render].
type ErrorMode int

const (
	// Strict (the default) returns a [*Error] whose Undefined slice
	// lists every placeholder missing from ctx.
	Strict ErrorMode = iota
	// Empty replaces every missing placeholder with "".
	Empty
	// Keep leaves the original placeholder literal in the output.
	Keep
	// DefaultValue replaces every missing placeholder with the value
	// passed as the variadic second argument to [WithErrorMode].
	DefaultValue
)

// Option configures a [*Template]. Construct via the [With] helpers.
type Option func(*Template)

// WithDelimiters overrides the default '{' '}' placeholder markers. Use
// this when templates contain literal braces (e.g. JSON, Go code).
func WithDelimiters(pre, post string) Option {
	return func(t *Template) { t.pre, t.post = pre, post }
}

// WithErrorMode selects how missing keys are handled. For
// [DefaultValue], pass the replacement as the second argument:
//
//	fillin.New(src, fillin.WithErrorMode(fillin.DefaultValue, "n/a"))
//
// Other modes ignore the variadic argument. Calling WithErrorMode more
// than once keeps the last value.
func WithErrorMode(mode ErrorMode, defaultValue ...string) Option {
	return func(t *Template) {
		t.errMode = mode
		if mode == DefaultValue && len(defaultValue) > 0 {
			t.defaultVal = defaultValue[0]
		}
	}
}

// WithCaseInsensitive makes placeholder lookup case-insensitive. With
// this option, "{Name}" in the template resolves the key "name" in ctx.
func WithCaseInsensitive() Option {
	return func(t *Template) { t.caseInsensitive = true }
}

// appendStringified appends a ctx value to buf, handling []byte and
// string directly without a string round-trip. fmt.Stringer and error
// come first so user-defined types win over the default %v fallback.
func appendStringified(buf []byte, v any) []byte {
	switch s := v.(type) {
	case nil:
		return buf
	case string:
		return append(buf, s...)
	case []byte:
		return append(buf, s...)
	case error:
		return append(buf, s.Error()...)
	case fmt.Stringer:
		return append(buf, s.String()...)
	default:
		return fmt.Appendf(buf, "%v", v)
	}
}
