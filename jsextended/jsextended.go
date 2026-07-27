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
	// strict equality
	STRICT_EQ     = token.RegisterType("===")
	STRICT_NOT_EQ = token.RegisterType("!==")
	// compound assignment
	PLUS_ASSIGN     = token.RegisterType("+=")
	MINUS_ASSIGN    = token.RegisterType("-=")
	MULTIPLY_ASSIGN = token.RegisterType("*=")
	DIVIDE_ASSIGN   = token.RegisterType("/=")
	MODULO_ASSIGN   = token.RegisterType("%=")
	// bitwise
	AND_BITWISE  = token.RegisterType("&")
	OR_BITWISE   = token.RegisterType("|")
	XOR_BITWISE  = token.RegisterType("^")
	SHL_BITWISE  = token.RegisterType("<<")
	SHR_BITWISE  = token.RegisterType(">>")
	USHR_BITWISE = token.RegisterType(">>>") // unsigned SHR
)

func Plugin(b *plugin.Builder) {
	token.RegisterUnaryType(NEW)
	token.RegisterUnaryType(SPREAD)
	token.RegisterUnaryType(TYPEOF)
	token.RegisterUnaryType(ASYNC)
	token.RegisterUnaryType(AWAIT)
	token.RegisterBinaryType(PLUS_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(MINUS_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(MULTIPLY_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(DIVIDE_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(MODULO_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(STRICT_EQ, token.EQ.Precedence())
	token.RegisterBinaryType(STRICT_NOT_EQ, token.EQ.Precedence())
	token.RegisterBinaryType(OPTIONAL_CHAINING, token.DOT.Precedence())
	token.RegisterBinaryType(INSTANCEOF, token.LT.Precedence())
	token.RegisterBinaryType(ARROW, token.ASSIGN.Precedence()+1)
	token.RegisterBinaryType(QUESTION_MARK, -1)
	token.RegisterBinaryType(EXPO, token.MULTIPLY.Precedence()+1)
	token.RegisterBinaryType(COALESCING, token.OR.Precedence())
	token.RegisterBinaryType(OR_BITWISE, token.AND.Precedence()+1)
	token.RegisterBinaryType(XOR_BITWISE, token.AND.Precedence()+2)
	token.RegisterBinaryType(AND_BITWISE, token.AND.Precedence()+3)
	token.RegisterBinaryType(SHL_BITWISE, token.LT.Precedence()+5)
	token.RegisterBinaryType(SHR_BITWISE, token.LT.Precedence()+5)
	token.RegisterBinaryType(USHR_BITWISE, token.LT.Precedence()+5)

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
			case "instanceof":
				tok.Type = INSTANCEOF
			case "async":
				tok.Type = ASYNC
			case "await":
				tok.Type = AWAIT
			case "class":
				tok.Type = CLASS
			}
		case token.UNKNOWN:
			switch tok.Literal {
			case "?":
				switch sc.CurrentChar() {
				case '.':
					sc.AdvanceChar()
					tok.Type = OPTIONAL_CHAINING
					tok.Literal = "?."
				case '?':
					sc.AdvanceChar()
					tok.Type = COALESCING
					tok.Literal = "??"
				default:
					tok.Type = QUESTION_MARK
				}
			case "&":
				tok.Type = AND_BITWISE
			case "|":
				tok.Type = OR_BITWISE
			case "^":
				tok.Type = XOR_BITWISE
			}
		case token.LT:
			if sc.CurrentChar() == '<' {
				sc.AdvanceChar()
				tok.Type = SHL_BITWISE
				tok.Literal = "<<"
			}
		case token.GT:
			if sc.CurrentChar() == '>' {
				sc.AdvanceChar()
				if sc.CurrentChar() == '>' {
					sc.AdvanceChar()
					tok.Type = USHR_BITWISE
					tok.Literal = ">>>"
				} else {
					tok.Type = SHR_BITWISE
					tok.Literal = ">>"
				}
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
		case token.PLUS:
			if sc.CurrentChar() == '=' {
				sc.AdvanceChar()
				tok.Type = PLUS_ASSIGN
				tok.Literal = "+="
			}
		case token.MINUS:
			if sc.CurrentChar() == '=' {
				sc.AdvanceChar()
				tok.Type = MINUS_ASSIGN
				tok.Literal = "-="
			}
		case token.MULTIPLY:
			switch sc.CurrentChar() {
			case '*':
				sc.AdvanceChar()
				tok.Type = EXPO
				tok.Literal = "**"
			case '=':
				sc.AdvanceChar()
				tok.Type = MULTIPLY_ASSIGN
				tok.Literal = "*="
			}
		case token.DIVIDE:
			if sc.CurrentChar() == '=' {
				sc.AdvanceChar()
				tok.Type = DIVIDE_ASSIGN
				tok.Literal = "/="
			}
		case token.MODULO:
			if sc.CurrentChar() == '=' {
				sc.AdvanceChar()
				tok.Type = MODULO_ASSIGN
				tok.Literal = "%="
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
		case ASYNC:
			return ParseAsyncExpr(p)
		case AWAIT:
			return ParseAwaitExpr(p)
		}
		return next()
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case STRICT_EQ, STRICT_NOT_EQ, PLUS_ASSIGN, MINUS_ASSIGN, MULTIPLY_ASSIGN, DIVIDE_ASSIGN, MODULO_ASSIGN:
			return js.ParseBinaryExpr(p, left)
		case ARROW:
			return ParseArrowFunc(p, left)
		case QUESTION_MARK:
			return ParseTernaryExpr(p, left)
		case OPTIONAL_CHAINING:
			return ParseOptionalChainingExpr(p, left)
		case INSTANCEOF:
			return ParseInstanceofExpr(p, left)
		case EXPO:
			return ParseExpoExpr(p, left)
		case COALESCING:
			return ParseCoalescingExpr(p, left)
		}
		return next(left)
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
		switch p.CurrentToken.Type {
		case js.LET, CONST, VAR:
			return ParseVarStmt(p)
		case js.FOR:
			return parser.Switch(p, func(p *parser.Parser) (ast.Stmt, error) {
				return ParseForinStmt(p)
			}, func(p *parser.Parser) (ast.Stmt, error) {
				return ParseForofStmt(p)
			}, func(p *parser.Parser) (ast.Stmt, error) {
				return js.ParseForStmt(p)
			})
		case js.FUNCTION:
			return ParseFunctionDecl(p)
		case TRY:
			return ParseTryStmt(p)
		case SWITCH:
			return ParseSwitchStmt(p)
		case THROW:
			return ParseThrowStmt(p)
		case DO:
			return ParseDoWhileStmt(p)
		case CLASS:
			return ParseClassStmt(p)
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
	case *InstanceofExpr:
		return PrintInstanceofExpr(pr, v)
	case *ForofStmt:
		return PrintForofStmt(pr, v)
	case *ForinStmt:
		return PrintForinStmt(pr, v)
	case *FunctionDecl:
		return PrintFunctionDecl(pr, v)
	case *TernaryExpr:
		return PrintTernaryExpr(pr, v)
	case *SequenceExpr:
		return PrintSequenceExpr(pr, v)
	case *OptionalChainingExpr:
		return PrintOptionalChainingExpr(pr, v)
	case *AsyncExpr:
		return PrintAsyncExpr(pr, v)
	case *AwaitExpr:
		return PrintAwaitExpr(pr, v)
	case *ExpoExpr:
		return PrintExpoExpr(pr, v)
	case *CoalescingExpr:
		return PrintCoalescingExpr(pr, v)
	case *ClassStmt:
		return PrintClassStmt(pr, v)
	}
	return next(node)
}
