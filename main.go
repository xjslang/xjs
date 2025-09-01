// Package main implements the xjs compiler
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// Version of the xjs compiler
const Version = "0.1.0"

// Parse parses the input string and returns the AST and any errors
func Parse(input string) (program any, errors []string) {
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	errs := p.Errors()
	return prog, errs
}

var (
	version = flag.Bool("version", false, "Show version")
	output  = flag.String("o", "", "Output file (default: stdout)")
	verbose = flag.Bool("v", false, "Verbose output")
	demo    = flag.Bool("demo", false, "Run demo mode")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("xjs compiler version %s\n", Version)
		return
	}

	if *demo {
		runDemo()
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file.xjs>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s -demo (run demo mode)\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := args[0]

	if filepath.Ext(filename) != ".xjs" {
		fmt.Fprintf(os.Stderr, "Error: File must have .xjs extension\n")
		os.Exit(1)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Compiling %s...\n", filename)
	}

	program, errors := Parse(string(content))

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Parser errors in %s:\n", filename)
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "\t%s\n", err)
		}
		os.Exit(1)
	}

	result := program.(fmt.Stringer).String()

	if *output != "" {
		err := os.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Output written to %s\n", *output)
		}
	} else {
		fmt.Print(result)
	}
}

func runDemo() {
	input := `
		let x = 5
		let y = 10.5
		let name = "Hello World"
		
		function add(a, b) {
			return a + b
		}
		
		if (x < y) {
			console.log("x is less than y")
		}
		
		let numbers = [1, 2, 3]
		let person = {name: "John", age: 30}
	`

	fmt.Println("=== LEXER OUTPUT ===")
	l := lexer.New(input)

	for {
		tok := l.NextToken()
		fmt.Println(tok)
		if tok.Type == token.EOF {
			break
		}
	}

	fmt.Println("\n=== PARSER OUTPUT (using convenience function) ===")
	program, errors := Parse(input)

	if len(errors) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range errors {
			fmt.Println("\t" + err)
		}
		return
	}

	fmt.Println("AST:")
	fmt.Println(program.(fmt.Stringer).String())

	fmt.Printf("\nxjs version: %s\n", Version)
}
