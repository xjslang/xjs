# XJS (eXtensible JavaScript parser)

A tool for creating JavaScript dialects.

## Why XJS?

It all started with an innocent question: "Would it be possible to use `defer` (from Go) in JS?" And the answer is "Yes, but with caveats."

We could write a Babel plugin or extend Acorn to transpile `defer()` to JS. But those and other approaches shared a common problem: "the extension had to respect JS grammar rules." And `defer` didn't.

That's why I decided to build a JS parser capable of transpiling exotic or experimental features into standard JS.

**What can you do with XJS?**

You can create your own dialect without having to build a language from scratch. You'd simply add your own operators and statements to JS. The only limit is your imagination.

## Examples

Real-world dialects built with XJS:

- [DJS](https://github.com/xjslang/djs) - JavaScript with `defer`.
- [HJS](https://github.com/hjslang/hjs) - JavaScript with native HTML tags.

See also the [`examples/`](examples/) directory for basic samples.

## How does it work?

It works by intercepting the parsing flow at strategic points and adding custom **scanners**, **parsers**, and **printers**. This is done through the following methods:

- `scanner.Builder.UseScanner`: add your own custom tokens.
- `parser.Builder.UseStmtParser`: add your own statement parsers.
- `parser.Builder.UsePrefixOpParser`: add your own prefix operators.
- `parser.Builder.UseInfixOpParser`: add your own infix operators.
- `parser.Builder.UseExprParser`: add your own expression parsers (*).
- `printer.Builder.UsePrinter`: add your own node printers.

(*): You usually do not need `UseExprParser`, since `UsePrefixOpParser` and `UseInfixOpParser` are enough.

## What is a Scanner, a Parser and a Printer?

The parsing process is divided into three responsibilities:

- `scanner.Scanner`: splits or "tokenizes" the input into tokens, which are later processed by the parser.
- `parser.Parser`: interprets the tokens and produces the AST (Abstract Syntax Tree).
- `printer.Printer`: turns AST nodes back into code, whether compiled, formatted, or otherwise. Each printer has a different purpose.

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
