// Command fillin renders a text template by replacing {key} placeholders
// with values drawn from KEY=VALUE arguments or JSON files.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aisbergg/go-fillin/pkg/fillin"
)

const stdinMarker = "-"

const usage = `Usage: fillin [OPTIONS] TEMPLATE [KEY=VALUE ...]
       echo "TEMPLATE" | fillin [OPTIONS] -
       fillin [OPTIONS] TEMPLATE [KEY=VALUE ...] -f -

Options:
  -h, -help   show this help and exit
  -d OPEN CLOSE   placeholder delimiters (default "{ }")
  -m MODE         missing-key mode: strict|empty|keep|default (default strict)
  -D VALUE        default value for -m default
  -i              case-insensitive key match
  -f FILE         JSON file with string:any context; repeatable; -f - reads stdin

Examples:
  fillin "Hi {name}" name=Alice
  echo "Order #{id} for {customer}" | fillin - id=1001 customer=Bob
  fillin -d '[[' ']]' "[[greeting]] world" greeting=hello
  fillin -m default -D "anon" "Hello {user}"
  fillin -m empty -f ctx.json "report/{date}"
`

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	if before, after, ok := strings.Cut(v, "="); ok {
		*s = append(*s, before, after)
		return nil
	}
	*s = append(*s, v)
	return nil
}

// mergeDelimsArgs rewrites os.Args so that `-d OPEN CLOSE` becomes a
// single `-d OPEN=CLOSE` token, which stdlib flag can pass to Set. Value
// joins on '=' (delimiters that contain '=' are unsupported).
func mungedArgs() []string {
	src := os.Args[1:]
	out := make([]string, 0, len(src))
	for i := 0; i < len(src); i++ {
		a := src[i]
		out = append(out, a)
		if a == "-d" && i+2 < len(src) {
			out = append(out, src[i+1]+"="+src[i+2])
			i += 2
		}
	}
	return out
}

func main() {
	var (
		delims       stringSlice
		files        stringSlice
		errMode      = flag.String("m", "strict", "missing-key mode: strict|empty|keep|default")
		defaultValue = flag.String("D", "", "default value for -m default")
		insensitive  = flag.Bool("i", false, "case-insensitive key match")
	)
	flag.Var(&delims, "d", "placeholder delimiters OPEN CLOSE")
	flag.Var(&files, "f", "JSON file with string:any context; repeatable; -f - reads stdin")
	flag.Usage = func() { printToStderr(usage) }
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	_ = flag.CommandLine.Parse(mungedArgs())

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	mode, err := parseMode(*errMode)
	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}
	if len(delims) != 0 && len(delims) != 2 {
		printToStderr("-d requires exactly two arguments: OPEN CLOSE")
		os.Exit(1)
	}

	stdinForTemplate := args[0] == stdinMarker
	stdinForJSON := false
	for _, f := range files {
		if f == stdinMarker {
			stdinForJSON = true
		}
	}
	if stdinForTemplate && stdinForJSON {
		printToStderr("stdin can be either template (-) or JSON (-f -), not both")
		os.Exit(1)
	}

	tplSrc, err := readTemplate(args[0], stdinForTemplate)
	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	ctx, err := buildContext(files, args[1:])
	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	pre, post := "{", "}"
	if len(delims) == 2 {
		pre, post = delims[0], delims[1]
	}
	opts := []fillin.Option{
		fillin.WithDelimiters(pre, post),
		fillin.WithErrorMode(mode, *defaultValue),
	}
	if *insensitive {
		opts = append(opts, fillin.WithCaseInsensitive())
	}
	tpl, err := fillin.New(tplSrc, opts...)
	if err != nil {
		printToStderr(err)
		os.Exit(1)
	}

	if _, err := tpl.RenderTo(os.Stdout, ctx); err != nil {
		printToStderr(err)
		os.Exit(1)
	}
}

func printToStderr(s any) {
	fmt.Fprintln(os.Stderr, s) //nolint:revive
}

func parseMode(s string) (fillin.ErrorMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "strict":
		return fillin.Strict, nil
	case "empty":
		return fillin.Empty, nil
	case "keep":
		return fillin.Keep, nil
	case "default":
		return fillin.DefaultValue, nil
	default:
		return 0, fmt.Errorf("invalid value %q for -m option", s)
	}
}

func readTemplate(arg string, fromStdin bool) (string, error) {
	if fromStdin {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("read template from stdin: %w", err)
		}
		return string(b), nil
	}
	return arg, nil
}

func buildContext(files stringSlice, kvs []string) (map[string]any, error) {
	ctx := map[string]any{}
	for _, f := range files {
		var data []byte
		var err error
		if f == stdinMarker {
			data, err = io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("read stdin: %w", err)
			}
		} else {
			data, err = os.ReadFile(f)
			if err != nil {
				return nil, fmt.Errorf("read file %s: %w", f, err)
			}
		}
		if err := json.Unmarshal(data, &ctx); err != nil {
			return nil, fmt.Errorf("parse JSON from %s: %w", f, err)
		}
	}
	for _, kv := range kvs {
		before, after, ok := strings.Cut(kv, "=")
		if !ok {
			return nil, fmt.Errorf("invalid key=value (missing '='): %s", kv)
		}
		ctx[before] = after
	}
	return ctx, nil
}
