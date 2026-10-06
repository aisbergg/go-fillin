# ADR-004: Use functional options for Template configuration

- Status: Decided
- Date: 2026-10-05

## Context

Templates need to be configurable across three independent dimensions:
delimiter pair (e.g. `{` `}` vs `{{` `}}`), case sensitivity of name
matching, and error handling mode (plus optional default-value for the
missing case). The original ADR-003 described modes as flags attached to
the constructor but did not pin a Go-idiomatic configuration pattern.

## Evaluation Criteria

- Idiomatic Go (matches stdlib patterns like `http.Server` options)
- Extensible without breaking existing callers
- Plays well with the dual API (ADR-002): same options accepted by
  `New` and the `Render` convenience wrapper

## Options

1. **Config struct**: `fillin.New(src, Config{Pre: "{{", Mode: Strict})`.
   Forces callers to populate a struct; all fields show up on every call
   even when unused.
2. **Enum-style setters**: `t := fillin.New(src); t.OnError(); t.Delims("{{", "}}")`.
   Imperative, breaks the "configure-then-use" flow and cannot be
   forwarded to `Render` as a single value.
3. **Functional options**: `WithXxx` functions applied to the Template via
   a variadic `...Option` parameter on `New` (and `Render`).

## Decision

Use functional options. The package exposes an `Option` type and
`With*` constructors; `New` (and `Render`) takes `...Option`. Each option
mutates an unexported field on the constructed `*Template`.

Currently provided options:

- `WithDelimiters(pre, post string)`
- `WithErrorMode(mode ErrorMode, defaultValue ...string)`
- `WithCaseInsensitive()`

Unrecognised flags cannot exist: misspelt options become compile errors
because the `Option` type is concrete.

## Implications

- **Positive**: Idiomatic Go, easy to extend (new option does not break
  existing call sites).
- **Positive**: Keeps `Render` symmetrical with `New`; one set of options
  serves both call shapes.
- **Negative**: Slightly more allocation than a config struct on hot paths
  (options captured in closures), negligible for template startup cost.
- **Negative**: Reflection-free, so misuse is caught at compile time only
  when the spelling is wrong; semantic misuse (e.g. `WithErrorMode(Empty)`
  then expecting defaults) is not diagnosed.

## Related

- ADR-002 (dual API)
- ADR-003 (error modes, now expressed through `WithErrorMode`)
