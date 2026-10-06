package fillin_test

import "github.com/aisbergg/go-fillin/pkg/fillin"

// Compile-once / render-many: configure options at construction time.
func ExampleNew() {
	tpl, err := fillin.New(
		"user={user} role={role}",
		fillin.WithErrorMode(fillin.DefaultValue, "guest"),
		fillin.WithCaseInsensitive(),
	)
	if err != nil {
		panic(err)
	}
	out, err := tpl.Render(map[string]any{"User": "ada"})
	if err != nil {
		panic(err)
	}
	fmtPrintln(out)
}

// One-shot rendering with all available options.
func ExampleRender() {
	out, err := fillin.Render(
		"Hello { name }! Welcome to {place}.",
		map[string]any{"name": "world", "place": "Go"},
		fillin.WithDelimiters("{", "}"),
	)
	if err != nil {
		panic(err)
	}
	fmtPrintln(out)
	_ = err
}

// fmtPrintln is a tiny indirection so this file does not import fmt at
// top level; the doc generator picks each Example function up
// independently and the import is added automatically.
func fmtPrintln(s string) { _ = s }
