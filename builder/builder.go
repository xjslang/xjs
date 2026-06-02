package builder

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type Builder struct {
	scanner *scanner.Scanner
	parser  *parser.Parser
}

func New() *Builder {
	return &Builder{
		scanner: &scanner.Scanner{},
		parser:  &parser.Parser{},
	}
}

func (b *Builder) UseScanner(scanner func(sc *scanner.Scanner, next func() token.Token) token.Token) {
	b.scanner.UseScanner(scanner)
}

func (b *Builder) UseUnaryParser(parser func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	b.parser.UseUnaryParser(parser)
}

func (b *Builder) UseBinaryParser(parser func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error)) {
	b.parser.UseBinaryParser(parser)
}

func (b *Builder) UseStmtParser(parser func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error)) {
	b.parser.UseStmtParser(parser)
}

func (b *Builder) UseExprParser(parser func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	b.parser.UseExprParser(parser)
}

func (b *Builder) Install(plugin func(b *Builder)) *Builder {
	plugin(b)
	return b
}

func (b *Builder) Build(src []byte) *parser.Parser {
	b.scanner.Init(src)
	b.parser.Init(b.scanner)
	return b.parser
}
