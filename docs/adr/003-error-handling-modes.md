# ADR-003: Multiple error handling modes for undefined placeholders

- Status: Decided
- Date: 2026-10-05
- Updated: 2026-10-05 — renamed modes + switched to functional option (see ADR-004)

## Context

When a template contains a placeholder that has no corresponding value in the context map, different use cases require different behaviors. Email templates might want defaults, validation scenarios might want errors, and logging templates might prefer silent empty strings.

## Evaluation Criteria

- Flexibility across use cases
- Predictability of behavior
- Simplicity of configuration
- Clear error messages when needed

## Options

1. **Single behavior (error)**: Always fail on undefined placeholder — strict but inflexible.
2. **Single behavior (empty string)**: Always output empty — lenient but hides mistakes.
3. **Configurable modes**: Allow users to choose behavior per template instance.

## Decision

Support four configurable error handling modes, selected via the
`WithErrorMode` functional option (see ADR-004):

1. `Strict` (default): return an error listing all undefined placeholders
2. `Empty`: replace undefined placeholders with the empty string
3. `Keep`: leave the original placeholder literal in the output
4. `DefaultValue`: replace with the value passed as the second argument
   to `WithErrorMode`

    fillin.New(src, fillin.WithErrorMode(fillin.DefaultValue, "n/a"))

## Implications

- **Positive**: Covers all common use cases without forcing one size fits all
- **Positive**: Default is strict (error), which helps catch bugs during development
- **Negative**: More configuration surface; users must understand the options
- **Negative**: Requires storing mode state in template struct

## Related

- ADR-004 (functional options)
