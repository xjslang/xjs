// Package js implements a simplified version of JavaScript, without arrow functions, try/catch, etc.
//
// # Key functions
//
//   - [ParseProgram]: parse a JS program
//   - [ScannerBuilder]: extend the default JS scanner
//   - [ParserBuilder]: extend the default JS parser
//   - [PrinterBuilder]: extend the default JS printer
//
// # Utilities
//
//   - [ExpectSemi]: expect next token to be a semicolon
//   - [IsSemi]: check if a given token is a semicolon
//
// # Other functions
//
//   - Parse*: parsing functions
//   - Print*: printing functions
package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var DEFAULT = token.RegisterType("DEFAULT", "default")

// ScannerBuilder extends the simplified JavaScript scanner, allowing the input text to be split into custom tokens.
func ScannerBuilder() *scanner.Builder {
	sb := &scanner.Builder{}
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(sc); err != nil {
			return
		}
		if tok.Type == token.IDENT {
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
			}
		}
		return
	})
	return sb
}

// ParserBuilder extends the simplified JavaScript parser, allowing custom language features to be parsed.
func ParserBuilder() *parser.Builder {
	token.RegisterPrefixOp(FUNCTION)
	token.RegisterPrefixOp(DELETE)

	pb := &parser.Builder{}
	pb.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
		switch p.CurrentToken.Type {
		case FUNCTION:
			return ParseFunctionDecl(p)
		case LET:
			return ParseLetStmt(p)
		case IF:
			return ParseIfStmt(p)
		case WHILE:
			return ParseWhileStmt(p)
		case FOR:
			return ParseForStmt(p)
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
			switch p.PeekToken.Type {
			case token.COLON:
				return ParseLabelStmt(p)
			}
		case token.SEMICOLON:
			return ParseSemiStmt(p)
		}
		return ParseStmt(p)
	})
	pb.UseExprParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		return ParseExpr(p)
	})
	pb.UsePrefixOpParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case FUNCTION:
			return ParsePrefixFunctionOp(p)
		case DELETE:
			return ParsePrefixDeleteOp(p)
		case token.LPAREN:
			return ParsePrefixParenOp(p)
		case token.LBRACE:
			return ParsePrefixBraceOp(p)
		case token.LBRACKET:
			return ParsePrefixBracketOp(p)
		}
		return ParsePrefixOp(p)
	})
	pb.UseInfixOpParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
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
		}
		return ParseInfixOp(p, left)
	})
	return pb
}

// PrinterBuilder extends the simplified JavaScript printer, allowing custom language features to be printed.
func PrinterBuilder() *printer.Builder {
	pb := &printer.Builder{}
	pb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		switch v := node.(type) {
		case *Program:
			return PrintProgram(pr, v)
		case *BlockStmt:
			return PrintBlockStmt(pr, v)
		case *IfStmt:
			return PrintIfStmt(pr, v)
		case *WhileStmt:
			return PrintWhileStmt(pr, v)
		case *FunctionDecl:
			return PrintFunctionDecl(pr, v)
		case *LetStmt:
			return PrintLetStmt(pr, v)
		case *ForStmt:
			return PrintForStmt(pr, v)
		case *PrefixFunctionOp:
			return PrintPrefixFunctionOp(pr, v)
		case *InfixParenOp:
			return PrintInfixParenOp(pr, v)
		case *InfixBracketOp:
			return PrintInfixBracketOp(pr, v)
		case *PrefixParenOp:
			return PrintPrefixParenOp(pr, v)
		case *PrefixBraceOp:
			return PrintPrefixBraceOp(pr, v)
		case *PrefixBracketOp:
			return PrintPrefixBracketOp(pr, v)
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
		}
		return next(node)
	})
	return pb
}
