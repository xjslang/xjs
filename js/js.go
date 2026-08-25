// The `js` package provides parsing and printing functions.
//
// # Key functions
//
//   - [ScannerBuilder]: extend the default JS scanner
//   - [ParserBuilder]: extend the default JS parser
//   - [PrinterBuilder]: extend the default JS printer
//
// # Other functions:
//   - Parse[*]: parsing functions
//   - Print[*]: printing functions
package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	DEFAULT = token.RegisterType("DEFAULT", "default")
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
	sb := &scanner.Builder{}
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) token.Token) token.Token {
		tok := next(sc)
		switch tok.Type {
		case token.IDENT:
			switch tok.Literal {
			case "function":
				tok.Type = FUNCTION
			case "let":
				tok.Type = LET
			case "if":
				tok.Type = IF
			case "else":
				tok.Type = ELSE
			case "while":
				tok.Type = WHILE
			case "for":
				tok.Type = FOR
			case "return":
				tok.Type = RETURN
			case "break":
				tok.Type = BREAK
			case "continue":
				tok.Type = CONTINUE
			case "import":
				tok.Type = IMPORT
			case "export":
				tok.Type = EXPORT
			case "delete":
				tok.Type = DELETE
			case "default":
				tok.Type = DEFAULT
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
			// TODO: Mapping the identifier literal `"in"` to the new token type `IN` will cause parsing errors anywhere the grammar expects an identifier token type (`token.IDENT`). For example, `ParseForinStmt` parses the loop variable with `ParseIdent`, which requires `token.IDENT` (see `js/ident.go`), and there are existing pass fixtures like `testdata/pass/9fe1d41db318afba.js` and `testdata/pass/9aa93e1e417ce8e3.js` that use `in` as the binding identifier (`for(let in 1);`, `for (let in a) {}`). Decide whether `in` should become a reserved word in this grammar; if so, those fixtures should move to `testdata/errors`, otherwise identifier parsing needs to accept keyword token types based on their literal.
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
		return tok
	})
	return sb
}

func ParserBuilder() *parser.Builder {
	token.RegisterPrefixOp(FUNCTION)
	token.RegisterPrefixOp(DELETE)
	token.RegisterPrefixOp(NEW)
	token.RegisterPrefixOp(SPREAD)
	token.RegisterPrefixOp(TYPEOF)
	token.RegisterPrefixOp(ASYNC)
	token.RegisterPrefixOp(AWAIT)
	token.RegisterPrefixOp(VOID)
	token.RegisterPrefixOp(CLASS)
	token.RegisterPrefixOp(YIELD)
	token.RegisterPrefixOp(token.DIVIDE)  // regex
	token.RegisterPrefixOp(DIVIDE_ASSIGN) // regex
	token.RegisterPrefixOp(token.INCREMENT)
	token.RegisterPrefixOp(token.DECREMENT)
	token.RegisterPrefixOp(NOT_BITWISE)
	token.RegisterInfixOp(PLUS_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(MINUS_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(MULTIPLY_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(DIVIDE_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(MODULO_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(OR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(XOR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(AND_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(EXPO_ASSIGN, token.ASSIGN.Precedence())
	// TODO: These new shift-assignment operators are registered as generic infix operators (and will be parsed via js.ParseInfixOp), which makes them left-associative due to ParseRightExpr stopping on equal precedence. In JavaScript, all assignment operators (including <<=, >>=, >>>=) are right-associative, so expressions like `a <<= b <<= c` or `a <<= (b = c)` will build an incorrect AST shape. Consider parsing all assignment-like operators with a dedicated parser that treats them as right-associative (e.g., parse RHS with `opPrec-1` / allow same-precedence operators on the RHS), similar to how `=` is handled via InfixAssignOp.
	token.RegisterInfixOp(SHL_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(SHR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(USHR_ASSIGN, token.ASSIGN.Precedence())
	token.RegisterInfixOp(STRICT_EQ, token.EQ.Precedence())
	token.RegisterInfixOp(STRICT_NOT_EQ, token.EQ.Precedence())
	token.RegisterInfixOp(OPTIONAL_CHAINING, token.DOT.Precedence())
	token.RegisterInfixOp(INSTANCEOF, token.LT.Precedence())
	token.RegisterInfixOp(ARROW, token.ASSIGN.Precedence()+1)
	token.RegisterInfixOp(TERNARY, token.ASSIGN.Precedence()+1)
	token.RegisterInfixOp(EXPO, token.MULTIPLY.Precedence()+1)
	token.RegisterInfixOp(COALESCING, token.OR.Precedence())
	token.RegisterInfixOp(OR_BITWISE, token.AND.Precedence()+1)
	token.RegisterInfixOp(XOR_BITWISE, token.AND.Precedence()+2)
	token.RegisterInfixOp(AND_BITWISE, token.AND.Precedence()+3)
	token.RegisterInfixOp(SHL_BITWISE, token.LT.Precedence()+5)
	token.RegisterInfixOp(SHR_BITWISE, token.LT.Precedence()+5)
	token.RegisterInfixOp(USHR_BITWISE, token.LT.Precedence()+5)
	token.RegisterInfixOp(IN, token.LT.Precedence())
	token.RegisterInfixOp(token.STRING, token.LPAREN.Precedence())

	pb := &parser.Builder{}
	pb.UsePrefixOpParser(func(p *parser.Parser, next func(*parser.Parser) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case DELETE:
			return ParsePrefixDeleteOp(p)
		case token.LPAREN:
			return ParsePrefixParenOp(p)
		case FUNCTION:
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
		return ParsePrefixOp(p)
	})
	pb.UseInfixOpParser(func(p *parser.Parser, left ast.Expr, next func(*parser.Parser, ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case token.ASSIGN:
			return ParseInfixAssignOp(p, left)
		case token.LPAREN:
			return ParseInfixParenOp(p, left)
		case token.LBRACKET:
			return ParseInfixBracketOp(p, left)
		case token.DOT:
			return ParseInfixDotOp(p, left)
		// TODO: Adding the comma operator here makes precedence interactions with assignment important: js.ParseInfixAssignOp currently parses its RHS with p.ParseExpr(), which will now eagerly consume commas. That yields an AST equivalent to a = (b, c) for a = b, c, but JavaScript semantics require (a = b), c because comma has lower precedence than assignment.
		case token.COMMA:
			return ParseInfixCommaOp(p, left)
		case token.INCREMENT:
			return ParsePostfixIncOp(p, left)
		case token.DECREMENT:
			return ParsePostfixDecOp(p, left)
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
		return ParseInfixOp(p, left)
	})
	pb.UseStmtParser(func(p *parser.Parser, next func(*parser.Parser) (ast.Stmt, error)) (ast.Stmt, error) {
		switch p.CurrentToken.Type {
		case IF:
			return ParseIfStmt(p)
		case WHILE:
			return ParseWhileStmt(p)
		case RETURN:
			return ParseReturnStmt(p)
		case BREAK:
			return ParseBreakStmt(p)
		case CONTINUE:
			return ParseContinueStmt(p)
		case IMPORT:
			return ParseImportStmt(p)
		case EXPORT:
			return ParseExportStmt(p)
		case token.IDENT:
			if p.PeekToken.Type == token.COLON {
				return ParseLabelStmt(p)
			}
		case token.SEMICOLON:
			return ParseSemiStmt(p)
		case LET, CONST, VAR:
			return ParseVarStmt(p)
		case FOR:
			return parser.Switch(p, func(p *parser.Parser) (ast.Stmt, error) {
				return ParseForinStmt(p)
			}, func(p *parser.Parser) (ast.Stmt, error) {
				return ParseForofStmt(p)
			}, func(p *parser.Parser) (ast.Stmt, error) {
				return ParseForStmt(p)
			})
		case FUNCTION:
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
		return ParseStmt(p)
	})
	pb.UseExprParser(func(p *parser.Parser, next func(*parser.Parser) (ast.Expr, error)) (ast.Expr, error) {
		return ParseExpr(p)
	})
	return pb
}

func PrinterBuilder() *printer.Builder {
	prb := &printer.Builder{}
	prb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
		switch v := node.(type) {
		case *Program:
			return PrintProgram(pr, v)
		case *BlockStmt:
			return PrintBlockStmt(pr, v)
		case *IfStmt:
			return PrintIfStmt(pr, v)
		case *WhileStmt:
			return PrintWhileStmt(pr, v)
		case *ForStmt:
			return PrintForStmt(pr, v)
		case *InfixParenOp:
			return PrintInfixParenOp(pr, v)
		case *InfixBracketOp:
			return PrintInfixBracketOp(pr, v)
		case *PostfixIncOp:
			return PrintPostfixIncOp(pr, v)
		case *PostfixDecOp:
			return PrintPostfixDecOp(pr, v)
		case *InfixAssignOp:
			return PrintInfixAssignOp(pr, v)
		case *PrefixOp:
			return PrintPrefixOp(pr, v)
		case *InfixOp:
			return PrintInfixOp(pr, v)
		case *Ident:
			return PrintIdent(pr, v)
		case *Variable:
			pr.Print(v.Token)
			return nil
		case *Literal:
			pr.Print(v.Value)
			return nil
		case *ExprStmt:
			return PrintExprStmt(pr, v)
		case *ReturnStmt:
			return PrintReturnStmt(pr, v)
		case *BreakStmt:
			return PrintBreakStmt(pr, v)
		case *ContinueStmt:
			return PrintContinueStmt(pr, v)
		case *LabelStmt:
			return PrintLabelStmt(pr, v)
		case *InfixDotOp:
			return PrintInfixDotOp(pr, v)
		case *InfixCommaOp:
			return PrintInfixCommaOp(pr, v)
		case *SemiStmt:
			return PrintSemiStmt(pr, v)
		case *ImportStmt:
			return PrintImportStmt(pr, v)
		case *ExportStmt:
			return PrintExportStmt(pr, v)
		case *PrefixDeleteOp:
			return PrintPrefixDeleteOp(pr, v)
		case *PrefixBracketOp:
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
		return next(pr, node)
	})
	return prb
}
