# XJS (eXtensible JavaScript)

XJS is an extensible JavaScript parser. That is, we can design our own constructs so that the parser knows how to "interpret" them, even if they are not part of the ECMAScript standard.

**What it is not?**

XJS **is not and will not be** a complete JavaScript parser, as that is outside the scope of the project. For example, XJS does not support `arrow functions` or `try / catch / finally`. However, we can create plugins that support those and other constructs.

## How does it work?

The `xjs` package exposes the `Plugin` and `Printer` functions to create our own custom parsers and printers.

The first step is to create a custom parser that transforms the input code into an AST (Abstract Syntax Tree). That input code may contain our custom constructs:

```go
import (
  "github.com/xjslang/xjs"
  "github.com/xjslang/xjs/builder"
  "github.com/xjslang/xjs/js"
  "github.com/xjslang/xjs/printer"
)

// create a builder and install the plugins
// that will "enrich" our parser
b := builder.New().
  // "base" plugin
  Install(xjs.Plugin).
  // our custom plugins
  Install(arrowFuncPlugin).
  Install(tryCatchPlugin).
  Install(strictEqPlugin).
  Install(anotherPlugin)

// now the parser knows how to interpret
// our custom constructs
data, _ := os.ReadFile(file)
p := b.Build(data) // returns the "enriched" parser
result, err := js.ParseProgram(p) // returns the AST
```

The second step is to create a custom printer. The printer is responsible for transforming the AST back into code. Among other things, we can create compilers and formatters:

```go
// here we are creating a compiler
c := &printer.Printer{}
c.UsePrinter(xjs.Printer) // "base" printer
c.UsePrinter(compiler)    //  tells it how to compile custom nodes
c.Init(printer.Compact()) // don't forget to init!
c.Print(result)
jsCode, err := c.Output() // returns the compiled code

// here we are creating a formatter
fmt := &printer.Printer{}
fmt.UsePrinter(xjs.Printer) // "base" printer
fmt.UsePrinter(formatter)   // tells it how to format custom nodes
fmt.Init()                  // don't forget to init!
fmt.Print(result)
formattedCode, err := fmt.Output() // returns the formatted code
```

## Show me an example!

You can see a very simple example [HERE](https://github.com/xjslang/djs).

## I like the it! Can I use it in production?

The project is still in alpha and is not ready for production. Building a JavaScript parser requires significant effort and I do what I can in my spare time. However, you can help me by finding bugs or sharing your ideas.

**If you are also an experienced Go programmer, your suggestions are very welcome.**
