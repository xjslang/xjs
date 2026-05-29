package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/parser"
)

func usage() {
	fmt.Printf("xjscli is a tool for testing the different capabilities of the XJS parser.\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Println("\txjscli example.js")
	fmt.Println("\techo \"code\" | xjscli -stdin")
	fmt.Printf("\nOptions:\n\n")
	fmt.Println("\t-help  show this help")
	fmt.Println("\t-stdin read input from stdin (pipe or redirect)")
	fmt.Println("\t-check display only errors")
	fmt.Printf("\nExamples:\n\n")
	fmt.Println("\txjscli example.js")
	fmt.Println("\txjscli -check example.js")
	fmt.Println("\techo \"code\" | xjscli -stdin")
	fmt.Println("\techo \"code\" | xjscli -stdin -check")
	fmt.Println()
}

func main() {
	var helpFlag, stdinFlag, checkFlag bool
	flag.BoolVar(&helpFlag, "help", false, "show help")
	flag.BoolVar(&stdinFlag, "stdin", false, "read from stdin")
	flag.BoolVar(&checkFlag, "check", false, "display only errors")
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

	sc := xjs.NewScanner()
	sc.Init(data)
	p := xjs.NewParser()
	p.Init(sc)
	program, err := p.Parse()

	// prints errors
	if checkFlag {
		list := parser.ErrorList{}
		if l, ok := err.(parser.ErrorList); ok {
			list = l
		}
		result, err := json.Marshal(list)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(string(result))
		return
	}

	// prints the formatted output
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pr := xjs.NewPrinter()
	pr.Init()
	pr.Print(program)
	fmt.Print(pr.String())
}
