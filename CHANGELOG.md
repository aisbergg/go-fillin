# Changelog

All notable changes to this project will be documented in this file.

---

<a name="1.0.1"></a>

## [1.0.1](https://github.com/aisbergg/go-fillin/compare/v1.0.0...v1.0.1) (2026-10-21)

### Features

- Added a new `fillin` command-line tool that renders templates directly from the shell, reading the template from an argument or stdin and supplying values via `KEY=VALUE` arguments, `-f` JSON files, or stdin

### Performance Improvements

- Eliminated the allocation of the missing-names list in non-strict error modes, reducing overhead for templates that only use `Empty`, `Keep`, or `DefaultValue`

---

<a name="1.0.0"></a>

## [1.0.0](https://github.com/{username}/{repo}/compare/v0.0.0...v1.0.0) (2026-01-01)

Initial Release
