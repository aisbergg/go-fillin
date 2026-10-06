package benchmark

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	ftpl "github.com/valyala/fasttemplate"

	fillin "github.com/aisbergg/go-fillin/pkg/fillin"
)

var global string

func buildFillinTemplate(n int) string {
	parts := make([]string, 0, n)
	for i := range n {
		parts = append(parts, fmt.Sprintf("val%d={key%d}", i, i))
	}
	return strings.Join(parts, " ")
}

func buildGoTemplate(n int) string {
	parts := make([]string, 0, n)
	for i := range n {
		parts = append(parts, fmt.Sprintf("val%d={{index . \"key%d\"}}", i, i))
	}
	return strings.Join(parts, " ")
}

func makeScaledContext(n int) map[string]any {
	ctx := make(map[string]any, n)
	for i := range n {
		ctx[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}
	return ctx
}

var sizes = []int{1, 10, 50, 100}

func BenchmarkCompareParse(b *testing.B) {
	for _, n := range sizes {
		fillinSrc := buildFillinTemplate(n)
		goSrc := buildGoTemplate(n)
		ftSrc := buildFillinTemplate(n)

		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.Run("fillin", func(b *testing.B) {
				for b.Loop() {
					_, _ = fillin.New(fillinSrc)
				}
			})

			b.Run("fasttemplate", func(b *testing.B) {
				for b.Loop() {
					_ = ftpl.New(ftSrc, "{", "}")
				}
			})

			b.Run("go-template", func(b *testing.B) {
				for b.Loop() {
					_, _ = template.New("bench").Parse(goSrc)
				}
			})
		})
	}
}

func BenchmarkCompareRender(b *testing.B) {
	for _, n := range sizes {
		fillinSrc := buildFillinTemplate(n)
		goSrc := buildGoTemplate(n)
		ftSrc := buildFillinTemplate(n)
		ctx := makeScaledContext(n)

		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.Run("fillin", func(b *testing.B) {
				t, _ := fillin.New(fillinSrc)
				for b.Loop() {
					s, _ := t.Render(ctx)
					global = s
				}
			})

			b.Run("fasttemplate", func(b *testing.B) {
				t := ftpl.New(ftSrc, "{", "}")
				for b.Loop() {
					s := t.ExecuteString(ctx)
					global = s
				}
			})

			b.Run("go-template", func(b *testing.B) {
				t, _ := template.New("bench").Parse(goSrc)
				for b.Loop() {
					var buf bytes.Buffer
					_ = t.Execute(&buf, ctx)
					global = buf.String()
				}
			})
		})
	}
}

func BenchmarkCompareOneShot(b *testing.B) {
	fillinSrc := buildFillinTemplate(10)
	goSrc := buildGoTemplate(10)
	ctx := makeScaledContext(10)

	b.Run("fillin", func(b *testing.B) {
		for b.Loop() {
			s, _ := fillin.Render(fillinSrc, ctx)
			global = s
		}
	})

	b.Run("fasttemplate", func(b *testing.B) {
		for b.Loop() {
			s := ftpl.ExecuteString(fillinSrc, "{", "}", ctx)
			global = s
		}
	})

	b.Run("go-template", func(b *testing.B) {
		for b.Loop() {
			t, _ := template.New("bench").Parse(goSrc)
			var buf bytes.Buffer
			_ = t.Execute(&buf, ctx)
			global = buf.String()
		}
	})
}

func BenchmarkCaseInsensitive(b *testing.B) {
	src := "Hallo {NAME}, willkommen in {CITY}!"
	ctx := map[string]any{"name": "Max", "city": "Berlin"}

	b.Run("fillin", func(b *testing.B) {
		t, _ := fillin.New(src, fillin.WithCaseInsensitive())
		for b.Loop() {
			s, _ := t.Render(ctx)
			global = s
		}
	})
}

func BenchmarkAPIMethods(b *testing.B) {
	src := buildFillinTemplate(100)
	ctx := makeScaledContext(100)

	b.Run("Render", func(b *testing.B) {
		t, _ := fillin.New(src)
		for b.Loop() {
			s, _ := t.Render(ctx)
			global = s
		}
	})

	b.Run("Append", func(b *testing.B) {
		t, _ := fillin.New(src)
		buf := make([]byte, 0, len(src)+1000) // pre-allocate
		for b.Loop() {
			buf = buf[:0] // reset for next iteration
			buf, _ = t.Append(buf, ctx)
			global = string(buf)
		}
	})

	b.Run("RenderTo", func(b *testing.B) {
		t, _ := fillin.New(src)
		var buf bytes.Buffer
		buf.Grow(len(src) + 1000) // pre-allocate
		for b.Loop() {
			buf.Reset()
			t.RenderTo(&buf, ctx)
			global = buf.String()
		}
	})
}
