package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
)

func usage() {
	fmt.Printf("xjscli is a tool for testing the different capabilities of the XJS parser.\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Println("\txjscli example.js")
	fmt.Println("\techo \"code\" | xjscli -stdin ")
	fmt.Printf("\nOptions:\n\n")
	fmt.Println("\t-help  show this help")
	fmt.Println("\t-stdin read input from stdin (pipe or redirect)")
	fmt.Printf("\nExamples:\n\n")
	fmt.Println("\txjscli examples.js")
	fmt.Println("\techo 'function foo(){}' | xjscli -stdin")
	fmt.Println()
}

func main() {
	var helpFlag bool
	var stdinFlag bool

	flag.BoolVar(&helpFlag, "help", false, "show help")
	flag.BoolVar(&stdinFlag, "stdin", false, "read from stdin")
	flag.Parse()

	if helpFlag {
		usage()
		os.Exit(0)
	}
	if (stdinFlag && flag.NArg() != 0) || (!stdinFlag && flag.NArg() != 1) {
		usage()
		os.Exit(2)
	}

	var data []byte
	var err error
	if stdinFlag {
		stat, statErr := os.Stdin.Stat()
		if statErr != nil {
			fmt.Fprintln(os.Stderr, statErr)
			os.Exit(1)
		}
		if stat.Mode()&os.ModeCharDevice != 0 {
			fmt.Fprintln(os.Stderr, "Error: -stdin requires piped input")
			fmt.Fprintln(os.Stderr, "Example: echo \"code\" | xjscli -stdin")
			os.Exit(1)
		}

		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(flag.Arg(0))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	program, err := parser.Parse(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pr := printer.New()
	program.PrintTo(pr)
	fmt.Print(pr.String())
}
