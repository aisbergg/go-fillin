// Package fillin renders text templates by replacing {name} placeholders
// with values from a map.
//
// The default delimiter pair is '{' '}'; override with [WithDelimiters].
// Placeholder names are trimmed of surrounding whitespace before lookup,
// so "{ name }" resolves the same key as "{name}". Matching is
// case-sensitive unless [WithCaseInsensitive] is supplied.
//
// Two call shapes are supported:
//
//	t := fillin.New(src)                       // compile once, render many
//	out, err := t.Execute(ctx)
//
//	out, err := fillin.Render(src, ctx)         // one-shot, with options
//
// Behaviour for missing keys is controlled by [WithErrorMode].
package fillin

import "io"

// New compiles src and returns a reusable [*Template]. Options configure
// delimiters, case sensitivity, and missing-key behaviour; see the
// [With] helpers. [Strict] error mode is the default when no error-mode
// option is supplied.
func New(src string, opts ...Option) (*Template, error) {
	t := &Template{
		src:     src,
		pre:     "{",
		post:    "}",
		errMode: Strict,
	}
	for _, opt := range opts {
		opt(t)
	}
	t.parsed = t.parse(src, t.pre, t.post)
	return t, nil
}

// Render is a one-shot convenience wrapper that creates a [*Template]
// via [New] and immediately calls [Template.Render].
func Render(src string, ctx map[string]any, opts ...Option) (string, error) {
	t, err := New(src, opts...)
	if err != nil {
		return "", err
	}
	return t.Render(ctx)
}

// RenderTo is a one-shot convenience wrapper that creates a [*Template]
// via [New] and immediately calls [Template.RenderTo].
func RenderTo(w io.Writer, src string, ctx map[string]any, opts ...Option) (int, error) {
	t, err := New(src, opts...)
	if err != nil {
		return 0, err
	}
	return t.RenderTo(w, ctx)
}

// Append renders src against ctx and appends the result to dst, returning
// the extended slice. On error, dst is returned unchanged.
func Append(dst []byte, src string, ctx map[string]any, opts ...Option) ([]byte, error) {
	t, err := New(src, opts...)
	if err != nil {
		return dst, err
	}
	return t.Append(dst, ctx)
}
