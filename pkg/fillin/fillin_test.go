package fillin_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/aisbergg/go-fillin/pkg/fillin"
)

type stringerT struct{ s string }

func (s stringerT) String() string { return s.s }

func TestRender_basic(t *testing.T) {
	got, err := fillin.Render(
		"Hello {name}, you are #{n}.",
		map[string]any{"name": "world", "n": 42},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "Hello world, you are #42."; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_trimsPlaceholderName(t *testing.T) {
	got, err := fillin.Render(
		"{  greeting  } {greeting}",
		map[string]any{"greeting": "hi"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "hi hi"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_emptyPlaceholder(t *testing.T) {
	got, err := fillin.Render(
		"a{}b",
		map[string]any{},
		fillin.WithErrorMode(fillin.Keep),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "a{}b"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_unterminatedPlaceholder(t *testing.T) {
	got, err := fillin.Render(
		"literal {oops rest of text",
		map[string]any{},
		fillin.WithErrorMode(fillin.Keep),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "literal {oops rest of text"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_customDelimiters(t *testing.T) {
	got, err := fillin.Render(
		"<<name>> has JSON like {x}",
		map[string]any{"name": "ada", "x": 1},
		fillin.WithDelimiters("<<", ">>"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "ada has JSON like {x}"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_caseInsensitive(t *testing.T) {
	tpl, err := fillin.New("{Name} {ROLE}", fillin.WithCaseInsensitive())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := tpl.Render(map[string]any{"name": "ada", "role": "admin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "ada admin"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_stringifies(t *testing.T) {
	got, err := fillin.Render(
		"{b}/{i}/{f}/{s}/{sr}",
		map[string]any{
			"b":  true,
			"i":  -7,
			"f":  1.5,
			"s":  stringerT{"X"},
			"sr": []byte("raw"),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "true/-7/1.5/X/raw"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_nilValue(t *testing.T) {
	got, err := fillin.Render("[{x}]", map[string]any{"x": nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "[]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_strict_missing(t *testing.T) {
	_, err := fillin.Render(
		"{a} {b} {a}",
		map[string]any{"a": "A"},
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var fe *fillin.Error
	if !errors.As(err, &fe) {
		t.Fatalf("expected *fillin.Error, got %T", err)
	}
	if !errors.Is(err, fillin.ErrUndefined) {
		t.Fatalf("expected ErrUndefined, got %v", err)
	}
	if want := []string{"b"}; !equal(fe.Undefined, want) {
		t.Fatalf("Undefined: got %v, want %v", fe.Undefined, want)
	}
}

func TestRender_emptyMode(t *testing.T) {
	got, err := fillin.Render("[{a}][{b}]", map[string]any{"a": "A"},
		fillin.WithErrorMode(fillin.Empty))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "[A][]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_keepMode(t *testing.T) {
	got, err := fillin.Render(
		"[{a}][{b}]",
		map[string]any{"a": "A"},
		fillin.WithErrorMode(fillin.Keep),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "[A][{b}]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_defaultMode(t *testing.T) {
	got, err := fillin.Render(
		"[{a}][{b}]",
		map[string]any{"a": "A"},
		fillin.WithErrorMode(fillin.DefaultValue, "n/a"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "[A][n/a]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_defaultMode_extraArgIgnored(t *testing.T) {
	got, err := fillin.Render(
		"[{a}]",
		map[string]any{"a": "A"},
		fillin.WithErrorMode(fillin.Empty, "ignored"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "[A]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_multipleOccurrences(t *testing.T) {
	got, err := fillin.Render("{x}-{x}-{x}", map[string]any{"x": "Y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "Y-Y-Y"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_whitespaceOnlyBody(t *testing.T) {
	got, err := fillin.Render(
		"a{ \t }b",
		map[string]any{},
		fillin.WithErrorMode(fillin.Keep),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "a{}b"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_internalWhitespacePreservedAfterTrim(t *testing.T) {
	got, err := fillin.Render("[{a b}]", map[string]any{"a b": "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "[ok]"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_caseInsensitive_missing(t *testing.T) {
	_, err := fillin.Render(
		"{X}",
		map[string]any{"y": "Y"},
		fillin.WithCaseInsensitive(),
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRender_noPlaceholders(t *testing.T) {
	got, err := fillin.Render("just text", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "just text"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRender_consecutiveRepeatedPlaceholders(t *testing.T) {
	got, err := fillin.Render("{a}{b}", map[string]any{"a": "1", "b": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "12"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderTo_basic(t *testing.T) {
	var buf strings.Builder
	tpl, _ := fillin.New("Hello {name}!")
	n, err := tpl.RenderTo(&buf, map[string]any{"name": "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "Hello Alice!"
	if buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
	if n != len(want) {
		t.Fatalf("returned %d bytes, expected %d", n, len(want))
	}
}

func TestRenderTo_strict_missing(t *testing.T) {
	var buf strings.Builder
	tpl, _ := fillin.New("Hello {name}!")
	_, err := tpl.RenderTo(&buf, map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing placeholder")
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty output on error, got %q", buf.String())
	}
}

func TestRenderTo_empty_mode(t *testing.T) {
	var buf strings.Builder
	tpl, _ := fillin.New("Hello {name}!", fillin.WithErrorMode(fillin.Empty))
	_, err := tpl.RenderTo(&buf, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "Hello !"; buf.String() != want {
		t.Fatalf("got %q, want %q", buf.String(), want)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
