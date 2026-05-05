package js_test

import (
	"testing"

	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/printer"
)

func TestXxx(t *testing.T) {
	l := &lexer.Lexer{}
	l.Init([]byte("let x = 100"))

	b := js.Builder{}
	p := b.Build(l)
	pr, err := p.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	if n := len(pr.Statements); n != 1 {
		t.Fatalf("Expected 1 statement, got %d", n)
	}
	stmt0, ok := pr.Statements[0].(*js.LetStatement)
	if !ok {
		t.Fatalf("Expected *js.LetStatement, got %T", pr.Statements[0])
	}
	prt := printer.New()
	stmt0.PrintTo(prt)
	expected := "let x = 100;"
	if result := prt.String(); result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}
