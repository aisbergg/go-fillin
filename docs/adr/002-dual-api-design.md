# ADR-002: Provide both function-based and struct-based APIs

- Status: Decided
- Date: 2026-10-05

## Context

Users have different usage patterns: simple one-shot rendering vs. compiling a template once and rendering it many times with different data. We need to support both efficiently without forcing users into an API that doesn't fit their use case.

## Evaluation Criteria

- Ergonomics for simple cases
- Performance for repeated renders with same template
- API surface area (fewer concepts is better)
- Consistency with Go idioms

## Options

1. **Function-only**: `fillin.Render(template, context)` — simple but re-parses template on every call, inefficient for repeated use.
2. **Struct-only**: `t := fillin.New(template); t.Execute(context)` — efficient but more boilerplate for simple cases.
3. **Both**: Provide a convenience function that wraps the struct approach internally.

## Decision

Provide both APIs:

- `fillin.Render(template string, ctx map[string]any) (string, error)` — convenience function for one-shot rendering
- `t := fillin.New(template)` followed by `t.Render(ctx)` — for compile-once, render-many scenarios

The function internally creates a temporary template struct and renders it immediately.

## Implications

- **Positive**: Best of both worlds; users choose based on their needs
- **Positive**: Matches Go's standard library patterns (e.g., `encoding/json` has both `Marshal` function and `Encoder` struct)
- **Negative**: Slightly larger API surface to document and maintain
- **Negative**: Users may be confused about which to use initially

## Related

(future) ADR-004 (functional options)
