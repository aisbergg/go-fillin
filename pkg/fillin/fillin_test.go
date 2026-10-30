package fillin_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/aisbergg/go-fillin/internal/lib/test/assert"
	"github.com/aisbergg/go-fillin/internal/lib/test/require"
	"github.com/aisbergg/go-fillin/pkg/fillin"
)

type stringerT struct{ s string }

func (s stringerT) String() string { return s.s }

// renderWith normalizes all three entry points to (string, error) for uniform testing.
func renderWith(method, src string, ctx map[string]any, opts ...fillin.Option) (string, error) {
	switch method {
	case "Render":
		return fillin.Render(src, ctx, opts...)
	case "Append":
		out, err := fillin.Append(nil, src, ctx, opts...)
		return string(out), err
	case "RenderTo":
		var buf strings.Builder
		_, err := fillin.RenderTo(&buf, src, ctx, opts...)
		return buf.String(), err
	}
	panic("unknown method")
}

func TestBasicRendering(t *testing.T) {
	src := "Hello {name}, you are #{n}."
	ctx := map[string]any{"name": "world", "n": 42}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "Hello world, you are #42.")
		})
	}
}

func TestTrimsPlaceholderName(t *testing.T) {
	src := "{  greeting  } {greeting}"
	ctx := map[string]any{"greeting": "hi"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "hi hi")
		})
	}
}

func TestEmptyPlaceholder(t *testing.T) {
	src := "a{}b"
	ctx := map[string]any{}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.Keep))
			require.NoError(t, err)
			assert.Equal(t, got, "a{}b")
		})
	}
}

func TestUnterminatedPlaceholder(t *testing.T) {
	src := "literal {oops rest of text"
	ctx := map[string]any{}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.Keep))
			require.NoError(t, err)
			assert.Equal(t, got, "literal {oops rest of text")
		})
	}
}

func TestCustomDelimiters(t *testing.T) {
	src := "<<name>> has JSON like {x}"
	ctx := map[string]any{"name": "ada", "x": 1}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithDelimiters("<<", ">>"))
			require.NoError(t, err)
			assert.Equal(t, got, "ada has JSON like {x}")
		})
	}
}

func TestCaseInsensitive(t *testing.T) {
	src := "{Name} {ROLE}"
	ctx := map[string]any{"name": "ada", "role": "admin"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithCaseInsensitive())
			require.NoError(t, err)
			assert.Equal(t, got, "ada admin")
		})
	}
}

func TestStringifies(t *testing.T) {
	src := "{b}/{i}/{f}/{s}/{sr}"
	ctx := map[string]any{
		"b":  true,
		"i":  -7,
		"f":  1.5,
		"s":  stringerT{"X"},
		"sr": []byte("raw"),
	}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "true/-7/1.5/X/raw")
		})
	}
}

func TestNilValue(t *testing.T) {
	src := "[{x}]"
	ctx := map[string]any{"x": nil}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "[]")
		})
	}
}

func TestStrictMissing(t *testing.T) {
	src := "{a} {b} {a}"
	ctx := map[string]any{"a": "A"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			_, err := renderWith(method, src, ctx)
			require.Error(t, err)

			var fe *fillin.Error
			assert.True(t, errors.As(err, &fe), "expected *fillin.Error")
			assert.True(t, errors.Is(err, fillin.ErrUndefined), "expected ErrUndefined")
			assert.EqualUnordered(t, fe.Undefined, []string{"b"})
		})
	}
}

func TestEmptyMode(t *testing.T) {
	src := "[{a}][{b}]"
	ctx := map[string]any{"a": "A"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.Empty))
			require.NoError(t, err)
			assert.Equal(t, got, "[A][]")
		})
	}
}

func TestKeepMode(t *testing.T) {
	src := "[{a}][{b}]"
	ctx := map[string]any{"a": "A"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.Keep))
			require.NoError(t, err)
			assert.Equal(t, got, "[A][{b}]")
		})
	}
}

func TestDefaultMode(t *testing.T) {
	src := "[{a}][{b}]"
	ctx := map[string]any{"a": "A"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.DefaultValue, "n/a"))
			require.NoError(t, err)
			assert.Equal(t, got, "[A][n/a]")
		})
	}
}

func TestDefaultModeExtraArgIgnored(t *testing.T) {
	src := "[{a}]"
	ctx := map[string]any{"a": "A"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.Empty, "ignored"))
			require.NoError(t, err)
			assert.Equal(t, got, "[A]")
		})
	}
}

func TestMultipleOccurrences(t *testing.T) {
	src := "{x}-{x}-{x}"
	ctx := map[string]any{"x": "Y"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "Y-Y-Y")
		})
	}
}

func TestWhitespaceOnlyBody(t *testing.T) {
	src := "a{ \t }b"
	ctx := map[string]any{}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx, fillin.WithErrorMode(fillin.Keep))
			require.NoError(t, err)
			assert.Equal(t, got, "a{}b")
		})
	}
}

func TestInternalWhitespacePreservedAfterTrim(t *testing.T) {
	src := "[{a b}]"
	ctx := map[string]any{"a b": "ok"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "[ok]")
		})
	}
}

func TestCaseInsensitiveMissing(t *testing.T) {
	src := "{X}"
	ctx := map[string]any{"y": "Y"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			_, err := renderWith(method, src, ctx, fillin.WithCaseInsensitive())
			require.Error(t, err)
		})
	}
}

func TestNoPlaceholders(t *testing.T) {
	src := "just text"
	ctx := map[string]any{}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "just text")
		})
	}
}

func TestConsecutiveRepeatedPlaceholders(t *testing.T) {
	src := "{a}{b}"
	ctx := map[string]any{"a": "1", "b": "2"}

	for _, method := range []string{"Render", "Append", "RenderTo"} {
		t.Run(method, func(t *testing.T) {
			got, err := renderWith(method, src, ctx)
			require.NoError(t, err)
			assert.Equal(t, got, "12")
		})
	}
}

func TestTemplateSource(t *testing.T) {
	tpl, err := fillin.New("Hello {name}")
	require.NoError(t, err)
	assert.Equal(t, tpl.Source(), "Hello {name}")
}

// Covers errors.Error() method + joinQuoted() for multiple missing keys.
func TestErrorFormatting(t *testing.T) {
	// Single missing key
	_, err := fillin.Render("{a}", map[string]any{})
	require.Error(t, err)
	assert.Equal(t, `fillin: undefined placeholder "a"`, err.Error())

	// Multiple missing keys (triggers joinQuoted)
	_, err = fillin.Render("{a} {b}", map[string]any{})
	require.Error(t, err)
	assert.Equal(t, `fillin: undefined placeholders "a", "b"`, err.Error())
}
