package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func Parse(input []byte) (*js.Program, error) {
	p := PluginBuilder().Build(input)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func PluginBuilder() *plugin.Builder {
	return plugin.New().Install(func(b *plugin.Builder) {
		token.RegisterUnaryType(js.FUNCTION)

		b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
			if tok, err = next(); err != nil {
				return
			}
			if tok.Type == token.IDENT {
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
				case "import":
					tok.Type = js.IMPORT
				case "export":
					tok.Type = js.EXPORT
				}
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
			case js.IMPORT:
				return js.ParseImportStmt(p)
			case js.EXPORT:
				return js.ParseExportStmt(p)
			case token.IDENT:
				switch p.PeekToken.Type {
				case token.COLON:
					return js.ParseLabelStmt(p)
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
			switch p.CurrentToken.Type {
			case js.FUNCTION:
				return js.ParseFunctionExpr(p)
			case token.LPAREN:
				return js.ParseGroupExpr(p)
			case token.LBRACE:
				return js.ParseObjExpr(p)
			case token.LBRACKET:
				return js.ParseArrayExpr(p)
			}
			return js.ParseUnaryExpr(p)
		})
		b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
			switch p.CurrentToken.Type {
			case token.ASSIGN:
				return js.ParseAssignExpr(p, left)
			case token.LPAREN:
				return js.ParseCallExpr(p, left)
			case token.LBRACKET:
				return js.ParseIndexExpr(p, left)
			case token.DOT:
				return js.ParseMemberExpr(p, left)
			case token.INCREMENT:
				return js.ParseIncExpr(p, left)
			case token.DECREMENT:
				return js.ParseDecExpr(p, left)
			}
			return js.ParseBinaryExpr(p, left)
		})
	})
}

func PrinterBuilder() *printer.Builder {
	return printer.NewBuilder().UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		switch v := node.(type) {
		case *js.Program:
			return js.PrintProgram(pr, v)
		case *js.BlockStmt:
			return js.PrintBlockStmt(pr, v)
		case *js.IfStmt:
			return js.PrintIfStmt(pr, v)
		case *js.WhileStmt:
			return js.PrintWhileStmt(pr, v)
		case *js.FunctionDecl:
			return js.PrintFunctionDecl(pr, v)
		case *js.LetStmt:
			return js.PrintLetStmt(pr, v)
		case *js.ForStmt:
			return js.PrintForStmt(pr, v)
		case *js.FunctionExpr:
			return js.PrintFunctionExpr(pr, v)
		case *js.CallExpr:
			return js.PrintCallExpr(pr, v)
		case *js.IndexExpr:
			return js.PrintIndexExpr(pr, v)
		case *js.GroupExpr:
			return js.PrintGroupExpr(pr, v)
		case *js.ObjExpr:
			return js.PrintObjExpr(pr, v)
		case *js.ArrayExpr:
			return js.PrintArrayExpr(pr, v)
		case *js.IncExpr:
			return js.PrintIncExpr(pr, v)
		case *js.DecExpr:
			return js.PrintDecExpr(pr, v)
		case *js.AssignExpr:
			return js.PrintAssignExpr(pr, v)
		case *js.UnaryExpr:
			return js.PrintUnaryExpr(pr, v)
		case *js.BinaryExpr:
			return js.PrintBinaryExpr(pr, v)
		case *js.Ident:
			return js.PrintIdent(pr, v)
		case *js.Variable:
			pr.Print(v.Token)
			return nil
		case *js.Literal:
			pr.Print(v.Value)
			return nil
		case *js.ExprStmt:
			return js.PrintExprStmt(pr, v)
		case *js.ReturnStmt:
			return js.PrintReturnStmt(pr, v)
		case *js.BreakStmt:
			return js.PrintBreakStmt(pr, v)
		case *js.ContinueStmt:
			return js.PrintContinueStmt(pr, v)
		case *js.LabelStmt:
			return js.PrintLabelStmt(pr, v)
		case *js.MemberExpr:
			return js.PrintMemberExpr(pr, v)
		case *js.SemiStmt:
			return js.PrintSemiStmt(pr, v)
		case *js.ImportStmt:
			return js.PrintImportStmt(pr, v)
		case *js.ExportStmt:
			return js.PrintExportStmt(pr, v)
		}
		return next(node)
	})
}
