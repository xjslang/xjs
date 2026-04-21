package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
)

func usage() {
	fmt.Printf("xjscli is a tool for testing the different capabilities of the XJS parser.\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\txjscli example.js - parses \"example.js\" and displays the formatted output\n")
	fmt.Printf("\txjscli -h         - show help\n\n")
}

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()
	if help {
		usage()
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}

	inputFile := flag.Arg(0)
	data, err := os.ReadFile(inputFile)
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
