package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func NewBuilder() *builder.Builder {
	return builder.New().Install(jsPlugin)
}

func NewPrinter(opts ...printer.Option) *printer.Printer {
	p := &printer.Printer{}
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		switch v := node.(type) {
		case *js.Program:
			js.PrintProgram(p, v)
		case *js.BlockStmt:
			js.PrintBlockStmt(p, v)
		case *js.IfStmt:
			js.PrintIfStmt(p, v)
		case *js.WhileStmt:
			js.PrintWhileStmt(p, v)
		case *js.FunctionDecl:
			js.PrintFunctionDecl(p, v)
		case *js.LetStmt:
			js.PrintLetStmt(p, v)
		case *js.AssignStmt:
			js.PrintAssignStmt(p, v)
		case *js.ForStmt:
			js.PrintForStmt(p, v)
		case *js.CallExpr:
			js.PrintCallExpr(p, v)
		case *js.IndexExpr:
			js.PrintIndexExpr(p, v)
		case *js.ParenExpr:
			js.PrintParenExpr(p, v)
		case *js.UnaryExpr:
			js.PrintUnaryExpr(p, v)
		case *js.BinaryExpr:
			js.PrintBinaryExpr(p, v)
		case *js.Ident:
			js.PrintIdent(p, v)
		case *js.Variable:
			p.Print(v.Name)
		case *js.Literal:
			p.Print(v.Value)
		case *js.ExprStmt:
			js.PrintExprStmt(p, v)
		case *js.ReturnStmt:
			js.PrintReturnStmt(p, v)
		case *js.SemiStmt:
			js.PrintSemiStmt(p, v)
		case *js.BreakStmt:
			js.PrintBreakStmt(p, v)
		case *js.ContinueStmt:
			js.PrintContinueStmt(p, v)
		case *js.LabelStmt:
			js.PrintLabelStmt(p, v)
		case *js.IncStmt:
			js.PrintIncStmt(p, v)
		case *js.DecStmt:
			js.PrintDecStmt(p, v)
		default:
			next(node)
		}
	})
	p.Init(opts...)
	return p
}

func Parse(input []byte) (*js.Program, error) {
	p := NewBuilder().Build(input)
	return js.ParseProgram(p)
}

func jsPlugin(b *builder.Builder) {
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) (tok token.Token) {
		tok = next()
		if tok.Type != token.IDENT {
			return
		}
		switch tok.Literal {
		case "function":
			tok.Type = js.FUNCTION
		case "let":
			tok.Type = js.LET
		case "if":
			tok.Type = js.IF
		case "else":
			tok.Type = js.ELSE
		case "while":
			tok.Type = js.WHILE
		case "for":
			tok.Type = js.FOR
		case "return":
			tok.Type = js.RETURN
		case "break":
			tok.Type = js.BREAK
		case "continue":
			tok.Type = js.CONTINUE
		}
		return
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
		switch p.CurrentToken.Type {
		case js.FUNCTION:
			return js.ParseFunctionDecl(p)
		case js.LET:
			return js.ParseLetStmt(p)
		case js.IF:
			return js.ParseIfStmt(p)
		case js.WHILE:
			return js.ParseWhileStmt(p)
		case js.FOR:
			return js.ParseForStmt(p)
		case js.RETURN:
			return js.ParseReturnStmt(p)
		case js.BREAK:
			return js.ParseBreakStmt(p)
		case js.CONTINUE:
			return js.ParseContinueStmt(p)
		case token.IDENT:
			switch p.PeekToken.Type {
			case token.ASSIGN:
				if !p.PeekToken.AfterNewline {
					return js.ParseAssignStmt(p)
				}
			case token.COLON:
				return js.ParseLabelStmt(p)
			case token.INCREMENT:
				if !p.PeekToken.AfterNewline {
					return js.ParseIncStmt(p)
				}
			case token.DECREMENT:
				if !p.PeekToken.AfterNewline {
					return js.ParseDecStmt(p)
				}
			}
		case token.SEMICOLON:
			return js.ParseSemiStmt(p)
		}
		return js.ParseStmt(p)
	})
	b.UseExprParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		return js.ParseExpr(p)
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		if p.CurrentToken.Type == token.LPAREN {
			return js.ParseParenExpr(p)
		}
		return js.ParseUnaryExpr(p)
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case token.LPAREN:
			return js.ParseCallExpr(p, left)
		case token.LBRACKET:
			return js.ParseIndexExpr(p, left)
		}
		return js.ParseBinaryExpr(p, left)
	})
}
