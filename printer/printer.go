package printer

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer/internal/formatter"
	"github.com/xjslang/xjs/token"
)

type PrinterOption = formatter.FormatterOption

func WithIndent(value string) PrinterOption {
	return formatter.WithIndent(value)
}

type Printer struct {
	fmt         formatter.Formatter
	ensureLine  bool
	ensureSpace bool
	printer     func(*Printer, ast.Node)
}

func (p *Printer) Init(opts ...PrinterOption) {
	p.fmt.Init(opts...)
	if p.printer == nil {
		p.printer = defaultPrinter
	}
}

func (p *Printer) IncreaseIndent() {
	p.fmt.IncreaseIndent()
}

func (p *Printer) DecreaseIndent() {
	p.fmt.DecreaseIndent()
}

func (p *Printer) printNode(node ast.Node) {
	p.printer(p, node)
}

func (p *Printer) String() string {
	return p.fmt.String()
}

func (p *Printer) Bytes() []byte {
	return []byte(p.String())
}

func (p *Printer) printString(s string) {
	if len(s) == 0 {
		return
	}
	p.printSeparatorIfNeeded()
	p.fmt.PrintString(s)
}

func (p *Printer) printRune(r rune) {
	p.printSeparatorIfNeeded()
	p.fmt.PrintRune(r)
}

func (p *Printer) PrintTrivia(trivia []token.Token) {
	for _, tok := range trivia {
		switch tok.Type {
		case token.NEWLINE:
			p.fmt.PrintRune('\n')
		case token.LINE_COMMENT:
			p.printSpaceIfNeeded()
			p.printIndentIfNeeded()
			p.fmt.PrintString("//" + tok.Literal)
		case token.BLOCK_COMMENT:
			p.printIndentIfNeeded()
			p.fmt.PrintString("/*" + tok.Literal + "*/")
		}
	}
}

func (p *Printer) printToken(tok token.Token) {
	p.PrintTrivia(tok.LeadingTrivia)
	p.printString(tok.Literal)
}

func (p *Printer) EnsureLine() {
	p.ensureLine = true
}

func (p *Printer) EnsureSpace() {
	p.ensureSpace = true
}

func (p *Printer) printLineIfNeeded() {
	if !isNewLine(p.fmt.LastChar) {
		p.fmt.PrintRune('\n')
	}
}

func (p *Printer) printSpaceIfNeeded() {
	if !isWhitespace(p.fmt.LastChar) {
		p.fmt.PrintRune(' ')
	}
}

func (p *Printer) printIndentIfNeeded() {
	if isNewLine(p.fmt.LastChar) {
		p.fmt.PrintIndent()
	}
}

func (p *Printer) printSeparatorIfNeeded() {
	if p.ensureLine {
		p.printLineIfNeeded()
		p.ensureLine = false
	}
	if p.ensureSpace {
		p.printSpaceIfNeeded()
		p.ensureSpace = false
	}
	p.printIndentIfNeeded()
}

func (p *Printer) Print(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			p.printString(v)
		case rune:
			p.printRune(v)
		case ast.Node:
			p.printNode(v)
		case token.Token:
			p.printToken(v)
		default:
			panic("Unsupported type")
		}
	}
}

func (p *Printer) LnPrint(arg any) {
	p.EnsureLine()
	p.Print(arg)
}

func (p *Printer) SpPrint(args ...any) {
	for _, arg := range args {
		p.EnsureSpace()
		p.Print(arg)
	}
}
