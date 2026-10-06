package fillin

import (
	"io"
	"slices"
	"strings"
)

// Template holds a compiled template ready for repeated [Template.Render]
// calls. Construct via [New]; the zero value is not usable.
type Template struct {
	src string

	pre             string
	post            string
	caseInsensitive bool
	errMode         ErrorMode
	defaultVal      string

	parsed []part
}

// Source returns the original template string passed to [New].
func (t *Template) Source() string { return t.src }

// Render renders the template against ctx according to the options
// supplied to [New]. In [Strict] mode, which is the default, Render
// returns a non-nil [*Error] when one or more placeholders have no
// entry in ctx.
func (t *Template) Render(ctx map[string]any) (string, error) {
	buf, missing := t.render(nil, ctx)
	if err := strictErr(t.errMode, missing); err != nil {
		return "", err
	}
	return string(buf), nil
}

// Append appends the rendered template to dst, returning the extended
// slice. On error, dst is returned unchanged.
func (t *Template) Append(dst []byte, ctx map[string]any) ([]byte, error) {
	buf, missing := t.render(dst, ctx)
	if err := strictErr(t.errMode, missing); err != nil {
		return dst, err
	}
	return buf, nil
}

// RenderTo renders the template against ctx and writes the result to w.
// Returns the number of bytes written. In [Strict] mode, returns a non-nil
// [*Error] without writing any bytes when one or more placeholders miss.
func (t *Template) RenderTo(w io.Writer, ctx map[string]any) (int, error) {
	buf, missing := t.render(nil, ctx)
	if err := strictErr(t.errMode, missing); err != nil {
		return 0, err
	}
	return w.Write(buf)
}

// render appends the rendered template to buf, returning the extended
// slice and the list of placeholder names missing from ctx. Strict-mode
// reporting is deferred to the caller so partial output never leaks.
func (t *Template) render(buf []byte, ctx map[string]any) ([]byte, []string) { //nolint: revive
	if buf == nil {
		// estimate output size: template length + extra for value substitution
		est := len(t.src) + 10*len(t.parsed)
		buf = make([]byte, 0, est)
	}
	var missing []string
	for _, p := range t.parsed {
		if p.kind == partLiteral {
			buf = append(buf, p.text(t.src)...)
			continue
		}
		name := p.text(t.src) // already trimmed at parse time
		val, ok := t.lookup(ctx, name)
		if !ok {
			missing = append(missing, name)
			switch t.errMode { //nolint: revive
			case Empty, Strict:
				// strict defers error to caller
			case Keep:
				buf = append(buf, t.pre...)
				buf = append(buf, name...)
				buf = append(buf, t.post...)
			case DefaultValue:
				buf = append(buf, t.defaultVal...)
			}
			continue
		}
		buf = appendStringified(buf, val)
	}
	return buf, missing
}

// strictErr returns a sorted [*Error] when mode is Strict and missing is
// non-empty, otherwise nil.
func strictErr(mode ErrorMode, missing []string) error {
	if mode != Strict || len(missing) == 0 {
		return nil
	}
	slices.Sort(missing)
	return &Error{Undefined: missing}
}

// part is one segment of a parsed template, stored as an offset into the
// original source string to avoid per-part string allocations. For
// placeholders, off/len point to the trimmed name (whitespace removed at
// parse time).
type part struct {
	kind partKind
	off  uint32 // byte offset into Template.src
	len  uint16 // length in bytes (max 64KB per segment)
}

// text returns the string content of this part from the template source.
func (p *part) text(src string) string {
	end := int(p.off) + int(p.len)
	return src[p.off:end]
}

type partKind uint8

const (
	partLiteral partKind = iota
	partPlaceholder
)

// parse splits src into literal segments and placeholder bodies using
// pre/post as the delimiter pair. An opening pre with no matching post
// is emitted as a trailing literal so templates with stray braces still
// render predictably. Parts store byte offsets into the original source
// string to avoid per-part string allocations.
func (t *Template) parse(src, pre, post string) []part {
	preLen, postLen := len(pre), len(post)
	count := strings.Count(src, pre)
	if count == 0 {
		return []part{{kind: partLiteral, off: 0, len: uint16(len(src))}} //nolint: gosec
	}
	parts := make([]part, 0, count*2+1)
	pos := 0
	for pos < len(src) {
		next := strings.Index(src[pos:], pre)
		if next < 0 {
			parts = append(
				parts,
				part{kind: partLiteral, off: uint32(pos), len: uint16(len(src) - pos)}, //nolint: gosec
			)
			break
		}
		next += pos // absolute offset in src
		// try to find a matching post after the pre
		bodyStart := next + preLen
		end := strings.Index(src[bodyStart:], post)
		if end < 0 {
			parts = append(
				parts,
				part{kind: partLiteral, off: uint32(pos), len: uint16(len(src) - pos)}, //nolint: gosec
			)
			break
		}
		end += bodyStart // absolute offset in src
		bodyEnd := end
		parts = append(parts, part{kind: partLiteral, off: uint32(pos), len: uint16(next - pos)}) //nolint: gosec
		// trim placeholder name at parse time, store offset to trimmed content
		raw := src[bodyStart:bodyEnd]
		startTrim := 0
		for startTrim < len(raw) && (raw[startTrim] == ' ' || raw[startTrim] == '\t' || raw[startTrim] == '\n' || raw[startTrim] == '\r') {
			startTrim++
		}
		endTrim := len(raw)
		for endTrim > startTrim && (raw[endTrim-1] == ' ' || raw[endTrim-1] == '\t' || raw[endTrim-1] == '\n' || raw[endTrim-1] == '\r') {
			endTrim--
		}
		parts = append(
			parts,
			part{
				kind: partPlaceholder,
				off:  uint32(bodyStart + startTrim), //nolint: gosec
				len:  uint16(endTrim - startTrim),   //nolint: gosec
			},
		)
		pos = end + postLen
	}
	return parts
}

// lookup returns the raw ctx value for name and whether it was found.
func (t *Template) lookup(ctx map[string]any, name string) (any, bool) {
	if v, ok := ctx[name]; ok {
		return v, true
	}
	if !t.caseInsensitive {
		return nil, false
	}
	for k, v := range ctx {
		if strings.EqualFold(k, name) {
			return v, true
		}
	}
	return nil, false
}
