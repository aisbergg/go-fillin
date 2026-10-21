<a name="readme-top"></a>

[![GoDoc](https://pkg.go.dev/badge/github.com/aisbergg/go-fillin)](https://pkg.go.dev/github.com/aisbergg/go-fillin/pkg/fillin)
[![Coverage Status](https://codecov.io/gh/aisbergg/go-fillin/branch/main/graph/badge.svg)](https://codecov.io/gh/aisbergg/go-fillin)
[![License](https://img.shields.io/github/license/aisbergg/go-fillin)](https://pkg.go.dev/github.com/aisbergg/go-fillin)

<br />

<div align="center">
  <h2 align="center"><b>fillin</b></h2>

  <p align="center">
    Simple and fast string template engine for placeholder replacement
    <br />
    <br />
    <a href="https://pkg.go.dev/github.com/aisbergg/go-fillin/pkg/fillin">Explore the docs</a>
    ·
    <a href="https://github.com/aisbergg/go-fillin/issues/new?labels=bug&template=bug_report.md">Report Bug</a>
    ·
    <a href="https://github.com/aisbergg/go-fillin/issues/new?labels=enhancement&template=feature_request.md">Request Feature</a>
  </p>
</div>

I needed a simple, _fast_ template engine for string substitution in Go. Something like Python's lesser known String Template (I am not talking about f-strings or the new t-strings). And of course, what better way is there than just writing my own, thus _fillin_ was born!

Before I wrote _fillin_ I used Go's `text/template`, which is powerful, overkill for the use-case and I personally find it horrendous to use. Furthermore, I just found out that by just using `text/template` in your Go program you disable _Dead Code Elimination_ (DCE). Without DCE every little unused function is kept in the compiled binary, effectively bloating it.

<details open="open">
  <summary>Table of Contents</summary>

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Benchmark](#benchmark)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

</details>

## Features

- **Three entry points**, one for each common shape:
    - `Render(src, ctx) ⇒ (string, error)`: returns a freshly allocated string.
    - `Append(buf, src, ctx) ⇒ ([]byte, error)`: reuses caller-provided capacity.
    - `RenderTo(w, src, ctx)` or `Template.RenderTo(w, ctx)`: streams rendered text into any `io.Writer`.
- **Configurable delimiters**: default is `{key}`, but can be chosen freely with `WithDelimiters`.
- **Case-insensitive mode**: enable with `WithCaseInsensitive()` for relaxed key matching.
- **Four error modes** for missing placeholders: `Strict` (default, return error), `Empty` (leave blank), `Keep` (preserve original), `DefaultValue` (substitute provided value).
- **Compile once, render many**: parse templates into reusable `*Template` objects via `New()`.
- **Automatic type conversion**: context values are stringified using Go's default formatting.

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>

## Installation

```
go get -u github.com/aisbergg/go-fillin
```

Optional CLI (for playing around — not needed when using the library):

```
$ go install github.com/aisbergg/go-fillin/cmd/fillin

$ echo "Hi {name}" | fillin - name=Alice
Hi Alice

$ fillin -d '[[' ']]' "[[greeting]] world" greeting=hello
hello world

$ fillin -f ctx.json -m default -D anon "Hello {user}"
```

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>

## Usage

```go
package main

import (
	"fmt"
	"strings"

	"github.com/aisbergg/go-fillin/pkg/fillin"
)

func main() {
	//
	// One-shot rendering
	//

	src := "Hello {name}, welcome to {city}!"
	ctx := map[string]any{
		"name": "Alice",
		"city": "Berlin",
	}
	out, err := fillin.Render(src, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: Hello Alice, welcome to Berlin!

	//
	// Compile once, render many times
	//

	tpl, err := fillin.New("Order #{id} for {customer}")
	if err != nil {
		panic(err)
	}
	for _, order := range []map[string]any{
		{"id": 1001, "customer": "Bob"},
		{"id": 1002, "customer": "Carol"},
	} {
		out, _ := tpl.Render(order)
		fmt.Println(out)
	}
	// Output: Order #1001 for Bob
	//         Order #1002 for Carol

	//
	// Custom delimiters
	//

	src = "Hi [[name]], you have [[count]] messages"
	tpl, _ = fillin.New(src, fillin.WithDelimiters("[[", "]]"))
	out, _ = tpl.Render(map[string]any{"name": "Dave", "count": 5})
	fmt.Println(out)
	// Output: Hi Dave, you have 5 messages

	//
	// Case-insensitive matching
	//

	tpl, _ = fillin.New("Hello {Name}!", fillin.WithCaseInsensitive())
	out, _ = tpl.Render(map[string]any{"name": "Eve"})
	fmt.Println(out)
	// Output: Hello Eve!

	//
	// Missing key handling
	//

	tpl, _ = fillin.New("Hi {user}, your code is {code}",
		fillin.WithErrorMode(fillin.DefaultValue, "N/A"))
	out, _ = tpl.Render(map[string]any{"user": "Frank"})
	fmt.Println(out)
	// Output: Hi Frank, your code is N/A

	//
	// Append to existing buffer
	//

	buf := make([]byte, 0, 64)
	buf, err = fillin.Append(buf, "Item {id}: {name}",
		map[string]any{"id": 1, "name": "Widget"})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
	// Output: Item 1: Widget

	//
	// Render to an io.Writer
	//

	var bld strings.Builder
	tpl, _ = fillin.New("Template: {value}")
	tpl.RenderTo(&bld, map[string]any{"value": "test"})
	fmt.Println(bld.String())
	// Output: Template: test
}
```

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>

## Benchmark

Inside the [benchmarks](./benchmarks) directory are benchmarks comparing fillin against Go's standard `text/template` and [fasttemplate](https://github.com/valyala/fasttemplate). You can run them with:

```shell
cd benchmarks && go test -bench=BenchmarkCompare -benchtime=1s -benchmem
```

Results on Intel i7-8550U @ 1.80GHz:

| Benchmark                            | time/op                | B/op  | allocs/op |
| ------------------------------------ | ---------------------- | ----- | --------- |
| BenchmarkCompareOneShot/fillin       | 746 ns/op ± 23.8 ns/op | 816   | 4         |
| BenchmarkCompareOneShot/fasttemplate | 918 ns/op ± 3.31 ns/op | 421   | 13        |
| BenchmarkCompareOneShot/go-template  | 24.5 µs/op ± 487 ns/op | 12218 | 259       |

![](assets/bench_compare_oneshot.svg)

| Benchmark                                   | time/op                 | B/op  | allocs/op |
| ------------------------------------------- | ----------------------- | ----- | --------- |
| BenchmarkCompareParse/size=1/fillin         | 95.8 ns/op ± 2.78 ns/op | 160   | 2         |
| BenchmarkCompareParse/size=1/fasttemplate   | 171 ns/op ± 0.402 ns/op | 384   | 3         |
| BenchmarkCompareParse/size=1/go-template    | 2.79 µs/op ± 13.5 ns/op | 3104  | 31        |
| BenchmarkCompareParse/size=10/fillin        | 319 ns/op ± 1.39 ns/op  | 368   | 2         |
| BenchmarkCompareParse/size=10/fasttemplate  | 453 ns/op ± 2.31 ns/op  | 768   | 3         |
| BenchmarkCompareParse/size=10/go-template   | 13.7 µs/op ± 23.8 ns/op | 8600  | 143       |
| BenchmarkCompareParse/size=50/fillin        | 1.27 µs/op ± 1.52 ns/op | 1392  | 2         |
| BenchmarkCompareParse/size=50/fasttemplate  | 1.47 µs/op ± 8.41 ns/op | 2496  | 3         |
| BenchmarkCompareParse/size=50/go-template   | 61.2 µs/op ± 128 ns/op  | 32216 | 625       |
| BenchmarkCompareParse/size=100/fillin       | 2.69 µs/op ± 95.2 ns/op | 2800  | 2         |
| BenchmarkCompareParse/size=100/fasttemplate | 2.75 µs/op ± 25.6 ns/op | 4800  | 3         |
| BenchmarkCompareParse/size=100/go-template  | 124 µs/op ± 5.47 µs/op  | 62280 | 1226      |

![](assets/bench_compare_parse.svg)

| Benchmark                                    | time/op                 | B/op  | allocs/op |
| -------------------------------------------- | ----------------------- | ----- | --------- |
| BenchmarkCompareRender/size=1/fillin         | 76.3 ns/op ± 1.59 ns/op | 48    | 2         |
| BenchmarkCompareRender/size=1/fasttemplate   | 107 ns/op ± 0.286 ns/op | 24    | 2         |
| BenchmarkCompareRender/size=1/go-template    | 1.13 µs/op ± 4.34 ns/op | 568   | 16        |
| BenchmarkCompareRender/size=10/fillin        | 420 ns/op ± 1.11 ns/op  | 448   | 2         |
| BenchmarkCompareRender/size=10/fasttemplate  | 759 ns/op ± 0.886 ns/op | 422   | 13        |
| BenchmarkCompareRender/size=10/go-template   | 9.94 µs/op ± 31.5 ns/op | 3616  | 116       |
| BenchmarkCompareRender/size=50/fillin        | 1.95 µs/op ± 5.72 ns/op | 2496  | 2         |
| BenchmarkCompareRender/size=50/fasttemplate  | 3.44 µs/op ± 10.9 ns/op | 2905  | 56        |
| BenchmarkCompareRender/size=50/go-template   | 50.1 µs/op ± 521 ns/op  | 18467 | 559       |
| BenchmarkCompareRender/size=100/fillin       | 3.8 µs/op ± 3.71 ns/op  | 4864  | 2         |
| BenchmarkCompareRender/size=100/fasttemplate | 6.55 µs/op ± 18.1 ns/op | 5280  | 107       |
| BenchmarkCompareRender/size=100/go-template  | 99.2 µs/op ± 174 ns/op  | 36823 | 1110      |

![](assets/bench_compare_render.svg)

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>

## Contributing

If you have any suggestions, want to file a bug report or want to contribute to this project in some other way, please read the [contribution guideline](CONTRIBUTING.md).

And don't forget to give this project a star 🌟! Thanks again!

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>

## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>

## Contact

André Lehmann

- Email: aisberg@posteo.de
- [GitHub](https://github.com/aisbergg)
- [LinkedIn](https://www.linkedin.com/in/andre-lehmann-97408221a/)

<p align="right"><a href="#readme-top" alt="abc"><b>back to top ⇧</b></a></p>
