# ADR-001: Use {name} as default placeholder syntax with configurable delimiters

- Status: Decided
- Date: 2026-10-05

## Context

The template engine needs a clear, recognizable syntax for placeholders that users replace with values. Python's `string.Template` uses `$name`, Go's `text/template` uses `{{.Name}}`, and many other systems use `{name}`. We need to choose a default while allowing customization.

## Evaluation Criteria

- Familiarity to developers from Python, JavaScript, and other languages
- Low collision risk with literal text in templates
- Configurability for advanced use cases
- Simplicity of parsing

## Options

1. **`$name` (Python style)**: Very familiar but requires escaping in shell contexts; conflicts with shell variable expansion if templates are embedded in scripts.
2. **`{{name}}` (Mustache/Go template style)**: Widely recognized, low collision risk, but verbose for simple use cases.
3. **`{name}` (Simple braces)**: Clean, familiar from Python f-strings and JavaScript template literals, minimal parsing complexity.

## Decision

Use `{name}` as the default delimiter syntax while providing configurable pre/post delimiter strings. Users can customize to `{{`, `}}`, or any other string via constructor options. This balances familiarity with flexibility.

## Implications

- **Positive**: Intuitive default for most users; easy migration from Python f-strings or JS template literals
- **Positive**: Customizable for edge cases where `{` appears frequently in literal text (option mechanism: ADR-004)
- **Negative**: Slightly higher collision risk than `{{name}}` if templates contain JSON or similar structures (mitigated by configurable delimiters)

## Related

- ADR-004 (functional options)
