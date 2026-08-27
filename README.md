# XJS (eXtensible JavaScript parser)

A tool for creating JavaScript dialects.

## What is a Scanner, a Parser and a Printer?

The parsing process is divided into three responsibilities:

- `scanner.Scanner`: splits or "tokenizes" the input into tokens, which are later processed by the parser.
- `parser.Parser`: interprets the tokens and produces the AST (Abstract Syntax Tree).
- `printer.Printer`: turns AST nodes back into code, whether compiled, formatted, or otherwise. Each printer has a different purpose.

## How does it work?

It works by intercepting the parsing flow at strategic points and adding custom scanners, parsers, and printers. This is done through the following methods:

- `scanner.Builder.UseScanner`: add your own custom tokens.
- `parser.Builder.UseStmtParser`: add your own statement parsers.
- `parser.Builder.UsePrefixOpParser`: add your own prefix operators.
- `parser.Builder.UseInfixOpParser`: add your own infix operators.
- `parser.Builder.UseExprParser`: add your own expression parsers (*).
- `printer.Builder.UsePrinter`: add your own node printers.

(*): You usually do not need `UseExprParser`, since `UsePrefixOpParser` and `UseInfixOpParser` are enough.

## Examples

Real-world dialects built with XJS:

- [djs](https://github.com/xjslang/djs): JavaScript with `defer`.
- [hjs](https://github.com/hjslang/hjs): JavaScript with native HTML tags.

See also the [`examples/`](examples/) directory for basic samples.

## Benchmarks

```bash
go test -run="^$" -bench=BenchmarkCompare -benchmem -count=10 . > result.out
benchstat -col /parser result.out
```

Results (lower is faster):
```
Goja     ░░░░░░░░░░░░░░░░░░░░░░░░░░░ (9.65% faster)
tdewolff ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ (2.37% faster)
XJS      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
```

**How to interpret the results?**

Comparisons are tricky, especially when the tools being compared serve different purposes. Therefore, the results above simply tell us that XJS is fast enough and uses memory and CPU resources in a modest way. Nothing more :)

## Contributing

Help is welcome. The project is still under development, and maintaining a JavaScript parser is a lot of work.

See the open issues at https://github.com/xjslang/xjs/issues.
