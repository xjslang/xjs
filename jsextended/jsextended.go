package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	// strict equality
	STRICT_EQ     = token.RegisterType("STRICT_EQ", "===")
	STRICT_NOT_EQ = token.RegisterType("STRICT_NOT_EQ", "!==")
	// compound assignment
	PLUS_ASSIGN     = token.RegisterType("PLUS_ASSIGN", "+=")
	MINUS_ASSIGN    = token.RegisterType("MINUS_ASSIGN", "-=")
	MULTIPLY_ASSIGN = token.RegisterType("MULTIPLY_ASSIGN", "*=")
	DIVIDE_ASSIGN   = token.RegisterType("DIVIDE_ASSIGN", "/=")
	MODULO_ASSIGN   = token.RegisterType("MODULO_ASSIGN", "%=")
	OR_ASSIGN       = token.RegisterType("OR_ASSIGN", "|=")
	XOR_ASSIGN      = token.RegisterType("XOR_ASSIGN", "^=")
	AND_ASSIGN      = token.RegisterType("AND_ASSIGN", "&=")
	EXPO_ASSIGN     = token.RegisterType("EXPO_ASSIGN", "**=")
	SHL_ASSIGN      = token.RegisterType("SHL_ASSIGN", "<<=")
	SHR_ASSIGN      = token.RegisterType("SHR_ASSIGN", ">>=")
	USHR_ASSIGN     = token.RegisterType("USHR_ASSIGN", ">>>=")
	// bitwise
	NOT_BITWISE  = token.RegisterType("NOT_BITWISE", "~")
	AND_BITWISE  = token.RegisterType("AND_BITWISE", "&")
	OR_BITWISE   = token.RegisterType("OR_BITWISE", "|")
	XOR_BITWISE  = token.RegisterType("XOR_BITWISE", "^")
	SHL_BITWISE  = token.RegisterType("SHL_BITWISE", "<<")
	SHR_BITWISE  = token.RegisterType("SHR_BITWISE", ">>")
	USHR_BITWISE = token.RegisterType("USHR_BITWISE", ">>>") // unsigned SHR
	// others
	IN = token.RegisterType("IN", "in")
)

func ScannerBuilder() *scanner.Builder {
	sb := js.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
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
			case "extends":
				tok.Type = EXTENDS
			case "void":
				tok.Type = VOID
			case "yield":
				tok.Type = YIELD
			// TODO: Mapping the identifier literal `"in"` to the new token type `IN` will cause parsing errors anywhere the grammar expects an identifier token type (`token.IDENT`). For example, `ParseForinStmt` parses the loop variable with `js.ParseIdent`, which requires `token.IDENT` (see `js/ident.go`), and there are existing pass fixtures like `testdata/pass/9fe1d41db318afba.js` and `testdata/pass/9aa93e1e417ce8e3.js` that use `in` as the binding identifier (`for(let in 1);`, `for (let in a) {}`). Decide whether `in` should become a reserved word in this grammar; if so, those fixtures should move to `testdata/errors`, otherwise identifier parsing needs to accept keyword token types based on their literal.
			case "in":
				tok.Type = IN
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
					tok.Type = TERNARY
				}
			case "&":
				if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = AND_ASSIGN
					tok.Literal = "&="
				} else {
					tok.Type = AND_BITWISE
				}
			case "|":
				if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = OR_ASSIGN
					tok.Literal = "|="
				} else {
					tok.Type = OR_BITWISE
				}
			case "^":
				if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = XOR_ASSIGN
					tok.Literal = "^="
				} else {
					tok.Type = XOR_BITWISE
				}
			case "~":
				tok.Type = NOT_BITWISE
				tok.Literal = "~"
			}
		case token.LT:
			if sc.CurrentChar() == '<' {
				sc.AdvanceChar()
				if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = SHL_ASSIGN
					tok.Literal = "<<="
				} else {
					tok.Type = SHL_BITWISE
					tok.Literal = "<<"
				}
			}
		case token.GT:
			if sc.CurrentChar() == '>' {
				sc.AdvanceChar()
				if sc.CurrentChar() == '>' {
					sc.AdvanceChar()
					if sc.CurrentChar() == '=' {
						sc.AdvanceChar()
						tok.Type = USHR_ASSIGN
						tok.Literal = ">>>="
					} else {
						tok.Type = USHR_BITWISE
						tok.Literal = ">>>"
					}
				} else if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = SHR_ASSIGN
					tok.Literal = ">>="
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
				if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = EXPO_ASSIGN
					tok.Literal = "**="
				} else {
					tok.Type = EXPO
					tok.Literal = "**"
				}
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
	return sb
}

func ParserBuilder() *parser.Builder {
	token.RegisterUnaryType(NEW)
	token.RegisterUnaryType(SPREAD)
	token.RegisterUnaryType(TYPEOF)
	token.RegisterUnaryType(ASYNC)
	token.RegisterUnaryType(AWAIT)
	token.RegisterUnaryType(VOID)
	token.RegisterUnaryType(CLASS)
	token.RegisterUnaryType(YIELD)
	token.RegisterUnaryType(token.DIVIDE)  // regex
	token.RegisterUnaryType(DIVIDE_ASSIGN) // regex
	token.RegisterUnaryType(token.INCREMENT)
	token.RegisterUnaryType(token.DECREMENT)
	token.RegisterUnaryType(NOT_BITWISE)
	token.RegisterBinaryType(PLUS_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(MINUS_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(MULTIPLY_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(DIVIDE_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(MODULO_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(OR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(XOR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(AND_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(EXPO_ASSIGN, token.ASSIGN.Precedence())
	// TODO: These new shift-assignment operators are registered as generic binary operators (and will be parsed via js.ParseInfixOp), which makes them left-associative due to js.ParseRightExpr stopping on equal precedence. In JavaScript, all assignment operators (including <<=, >>=, >>>=) are right-associative, so expressions like `a <<= b <<= c` or `a <<= (b = c)` will build an incorrect AST shape. Consider parsing all assignment-like operators with a dedicated parser that treats them as right-associative (e.g., parse RHS with `opPrec-1` / allow same-precedence operators on the RHS), similar to how `=` is handled via InfixAssignOp.
	token.RegisterBinaryType(SHL_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(SHR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(USHR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterBinaryType(STRICT_EQ, token.EQ.Precedence())
	token.RegisterBinaryType(STRICT_NOT_EQ, token.EQ.Precedence())
	token.RegisterBinaryType(OPTIONAL_CHAINING, token.DOT.Precedence())
	token.RegisterBinaryType(INSTANCEOF, token.LT.Precedence())
	token.RegisterBinaryType(ARROW, token.ASSIGN.Precedence()+1)
	token.RegisterBinaryType(TERNARY, token.ASSIGN.Precedence()+1)
	token.RegisterBinaryType(EXPO, token.MULTIPLY.Precedence()+1)
	token.RegisterBinaryType(COALESCING, token.OR.Precedence())
	token.RegisterBinaryType(OR_BITWISE, token.AND.Precedence()+1)
	token.RegisterBinaryType(XOR_BITWISE, token.AND.Precedence()+2)
	token.RegisterBinaryType(AND_BITWISE, token.AND.Precedence()+3)
	token.RegisterBinaryType(SHL_BITWISE, token.LT.Precedence()+5)
	token.RegisterBinaryType(SHR_BITWISE, token.LT.Precedence()+5)
	token.RegisterBinaryType(USHR_BITWISE, token.LT.Precedence()+5)
	token.RegisterBinaryType(IN, token.LT.Precedence())
	token.RegisterBinaryType(token.STRING, token.LPAREN.Precedence())

	pb := js.ParserBuilder()
	pb.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case token.LPAREN:
			return ParsePrefixParenOp(p)
		case token.INCREMENT, token.DECREMENT, NOT_BITWISE:
			return js.ParsePrefixOp(p)
		case js.FUNCTION:
			return ParsePrefixFunctionOp(p)
		case token.LBRACE:
			return ParsePrefixBraceOp(p)
		case token.LBRACKET:
			return ParsePrefixBracketOp(p)
		case token.DIVIDE, DIVIDE_ASSIGN:
			return ParsePrefixRegOp(p)
		case NEW:
			return ParsePrefixNewOp(p)
		case SPREAD:
			return ParsePrefixSpreadOp(p)
		case TYPEOF:
			return ParsePrefixTypeofOp(p)
		case ASYNC:
			return ParsePrefixAsyncOp(p)
		case AWAIT:
			return ParsePrefixAwaitOp(p)
		case VOID:
			return ParsePrefixVoidOp(p)
		case CLASS:
			return ParsePrefixClassOp(p)
		case YIELD:
			return ParsePrefixYieldOp(p)
		}
		return next()
	})
	pb.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case STRICT_EQ, STRICT_NOT_EQ, PLUS_ASSIGN, MINUS_ASSIGN, MULTIPLY_ASSIGN, DIVIDE_ASSIGN, MODULO_ASSIGN, OR_ASSIGN, XOR_ASSIGN, AND_ASSIGN, EXPO_ASSIGN, IN:
			return js.ParseInfixOp(p, left)
		case ARROW:
			return ParseInfixArrowOp(p, left)
		case TERNARY:
			return ParseInfixTernaryOp(p, left)
		case OPTIONAL_CHAINING:
			return ParseInfixOptionalChainingOp(p, left)
		case INSTANCEOF:
			return ParseInfixInstanceofOp(p, left)
		case EXPO:
			return ParseInfixExpoOp(p, left)
		case COALESCING:
			return ParseInfixCoalescingOp(p, left)
		case token.STRING:
			return ParseInfixTemplateOp(p, left)
		}
		return next(left)
	})
	pb.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
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
	return pb
}

func PrinterBuilder() *printer.Builder {
	prb := js.PrinterBuilder()
	prb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		switch v := node.(type) {
		case *js.PrefixBracketOp:
			return PrintPrefixBracketOp(pr, v)
		case *PrefixBraceOp:
			return PrintPrefixBraceOp(pr, v)
		case *VarStmt:
			return PrintVarStmt(pr, v)
		case *TryStmt:
			return PrintTryStmt(pr, v)
		case *SwitchStmt:
			return PrintSwitchStmt(pr, v)
		case *ThrowStmt:
			return PrintThrowStmt(pr, v)
		case *PrefixNewOp:
			return PrintPrefixNewOp(pr, v)
		case *DoWhileStmt:
			return PrintDoWhileStmt(pr, v)
		case *InfixArrowOp:
			return PrintInfixArrowOp(pr, v)
		case *PrefixSpreadOp:
			return PrintPrefixSpreadOp(pr, v)
		case *PrefixTypeofOp:
			return PrintPrefixTypeofOp(pr, v)
		case *InfixInstanceofOp:
			return PrintInfixInstanceofOp(pr, v)
		case *ForofStmt:
			return PrintForofStmt(pr, v)
		case *ForinStmt:
			return PrintForinStmt(pr, v)
		case *FunctionDecl:
			return PrintFunctionDecl(pr, v)
		case *PrefixFunctionOp:
			return PrintPrefixFunctionOp(pr, v)
		case *InfixTernaryOp:
			return PrintInfixTernaryOp(pr, v)
		case *InfixOptionalChainingOp:
			return PrintInfixOptionalChainingOp(pr, v)
		case *PrefixAsyncOp:
			return PrintPrefixAsyncOp(pr, v)
		case *PrefixAwaitOp:
			return PrintPrefixAwaitOp(pr, v)
		case *InfixExpoOp:
			return PrintInfixExpoOp(pr, v)
		case *InfixCoalescingOp:
			return PrintInfixCoalescingOp(pr, v)
		case *ClassStmt:
			return PrintClassStmt(pr, v)
		case *PrefixClassOp:
			return PrintPrefixClassOp(pr, v)
		case *PrefixVoidOp:
			return PrintPrefixVoidOp(pr, v)
		case *PrefixYieldOp:
			return PrintPrefixYieldOp(pr, v)
		case *PrefixRegOp:
			return PrintPrefixRegOp(pr, v)
		case *PrefixParenOp:
			return PrintPrefixParenOp(pr, v)
		case *InfixTemplateOp:
			return PrintInfixTemplateOp(pr, v)
		}
		return next(node)
	})
	return prb
}
