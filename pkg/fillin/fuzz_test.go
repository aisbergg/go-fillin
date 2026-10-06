package fillin_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aisbergg/go-fillin/pkg/fillin"
)

func FuzzParse(f *testing.F) {
	f.Add("hello {world}")
	f.Add("{unclosed")
	f.Add("no placeholders")
	f.Add("")
	f.Add("{{nested}}")

	f.Fuzz(func(t *testing.T, src string) {
		_, err := fillin.New(src)
		// Should never panic; parse errors are acceptable
		_ = err
	})
}

func FuzzRender_Missing(f *testing.F) {
	f.Add("hello {name}", []byte{0x01})
	f.Add("{a} and {b}", []byte{0x02, 0x03})

	f.Fuzz(func(t *testing.T, src string, seed []byte) {
		tpl, _ := fillin.New(src, fillin.WithErrorMode(fillin.Keep))
		if tpl == nil {
			return
		}

		// Context has unrelated keys - no placeholders will be filled
		ctx := make(map[string]any)
		for i := range seed {
			ctx[fmt.Sprintf("unrelated%d", i)] = seed[i]
		}

		_, err := tpl.Render(ctx)
		// Should never panic; may return error depending on mode
		_ = err
	})
}

func FuzzRender_Matching(f *testing.F) {
	f.Add("hello {name}", []byte{0x01})
	f.Add("{a} and {b}", []byte{0x02, 0x03})

	f.Fuzz(func(t *testing.T, src string, seed []byte) {
		tpl, _ := fillin.New(src, fillin.WithErrorMode(fillin.Empty))
		if tpl == nil {
			return
		}

		// Extract placeholder names from template and create matching context
		ctx := make(map[string]any)
		re := regexp.MustCompile(`\{([^}]*)\}`)
		matches := re.FindAllStringSubmatch(src, -1)
		for i, match := range matches {
			key := strings.TrimSpace(match[1])
			if key != "" {
				ctx[key] = fmt.Sprintf("value%d", i)
			}
		}

		out, err := tpl.Render(ctx)
		// Should never panic; with Empty mode, no error expected
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_ = out
	})
}
