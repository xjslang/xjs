package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	STRICT_EQ     = token.RegisterType("===")
	STRICT_NOT_EQ = token.RegisterType("!==")
)

func Plugin(b *plugin.Builder) {
	token.RegisterUnaryType(NEW)
	token.RegisterUnaryType(SPREAD)
	token.RegisterUnaryType(TYPEOF)
	token.RegisterBinaryType(STRICT_EQ, token.EQ.Precedence())
	token.RegisterBinaryType(STRICT_NOT_EQ, token.EQ.Precedence())
	token.RegisterBinaryType(ARROW, -1) // lowest precedence possible
	token.RegisterBinaryType(QUESTION_MARK, -1)

	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.IDENT:
			switch tok.Literal {
			case "const":
				tok.Type = CONST
			case "var":
				tok.Type = VAR
			case "try":
				tok.Type = TRY
			case "catch":
				tok.Type = CATCH
			case "finally":
				tok.Type = FINALLY
			case "switch":
				tok.Type = SWITCH
			case "case":
				tok.Type = CASE
			case "default":
				tok.Type = DEFAULT
			case "throw":
				tok.Type = THROW
			case "new":
				tok.Type = NEW
			case "do":
				tok.Type = DO
			case "typeof":
				tok.Type = TYPEOF
			}
		case token.UNKNOWN:
			switch tok.Literal {
			case "?":
				tok.Type = QUESTION_MARK
			}
		case token.EQ:
			if sc.CurrentChar() == '=' {
				sc.AdvanceChar()
				tok.Type = STRICT_EQ
				tok.Literal = "==="
			}
		case token.NOT_EQ:
			if sc.CurrentChar() == '=' {
				sc.AdvanceChar()
				tok.Type = STRICT_NOT_EQ
				tok.Literal = "!=="
			}
		case token.DOT:
			if sc.CurrentChar() == '.' && sc.PeekChar() == '.' {
				sc.AdvanceChar()
				sc.AdvanceChar()
				tok.Type = SPREAD
				tok.Literal = "..."
			}
		case token.ASSIGN:
			if sc.CurrentChar() == '>' {
				sc.AdvanceChar()
				tok.Type = ARROW
				tok.Literal = "=>"
			}
		}
		return
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case token.LBRACE:
			return ParseObjExpr(p)
		case token.LBRACKET:
			return ParseArrayExpr(p)
		case token.LPAREN:
			return parser.Switch(p, func(p *parser.Parser) (ast.Expr, error) {
				return js.ParseGroupExpr(p)
			}, func(p *parser.Parser) (ast.Expr, error) {
				return ParseSequenceExpr(p)
			})
		case NEW:
			return ParseNewExpr(p)
		case SPREAD:
			return ParseSpreadExpr(p)
		case TYPEOF:
			return ParseTypeofExpr(p)
		}
		return next()
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case STRICT_EQ, STRICT_NOT_EQ:
			return js.ParseBinaryExpr(p, left)
		case ARROW:
			return ParseArrowFunc(p, left)
		case QUESTION_MARK:
			return ParseTernaryExpr(p, left)
		}
		return next(left)
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
		switch p.CurrentToken.Type {
		case js.LET, CONST, VAR:
			return ParseVarStmt(p)
		case js.FOR:
			return parser.Switch(p, func(p *parser.Parser) (ast.Stmt, error) {
				return ParseForofStmt(p)
			}, func(p *parser.Parser) (ast.Stmt, error) {
				return js.ParseForStmt(p)
			})
		case TRY:
			return ParseTryStmt(p)
		case SWITCH:
			return ParseSwitchStmt(p)
		case THROW:
			return ParseThrowStmt(p)
		case DO:
			return ParseDoWhileStmt(p)
		}
		return next()
	})
}

func Printer(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *js.ArrayExpr:
		return PrintArrayExpr(pr, v)
	case *ObjExpr:
		return PrintObjExpr(pr, v)
	case *VarStmt:
		return PrintVarStmt(pr, v)
	case *TryStmt:
		return PrintTryStmt(pr, v)
	case *SwitchStmt:
		return PrintSwitchStmt(pr, v)
	case *ThrowStmt:
		return PrintThrowStmt(pr, v)
	case *NewExpr:
		return PrintNewExpr(pr, v)
	case *DoWhileStmt:
		return PrintDoWhileStmt(pr, v)
	case *ArrowFuncExpr:
		return PrintArrowFunc(pr, v)
	case *SpreadExpr:
		return PrintSpreadExpr(pr, v)
	case *TypeofExpr:
		return PrintTypeofExpr(pr, v)
	case *ForofStmt:
		return PrintForofStmt(pr, v)
	case *TernaryExpr:
		return PrintTernaryExpr(pr, v)
	case *SequenceExpr:
		return PrintSequenceExpr(pr, v)
	}
	return next(node)
}
