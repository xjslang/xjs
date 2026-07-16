package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func Plugin(b *plugin.Builder) {
	token.RegisterUnaryType(FUNCTION)

	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
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
			}
		}
		return
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
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
	b.UseExprParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		return ParseExpr(p)
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case FUNCTION:
			return ParseFunctionExpr(p)
		case token.LPAREN:
			return ParseGroupExpr(p)
		case token.LBRACE:
			return ParseObjExpr(p)
		case token.LBRACKET:
			return ParseArrayExpr(p)
		}
		return ParseUnaryExpr(p)
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case token.ASSIGN:
			return ParseAssignExpr(p, left)
		case token.LPAREN:
			return ParseCallExpr(p, left)
		case token.LBRACKET:
			return ParseIndexExpr(p, left)
		case token.DOT:
			return ParseMemberExpr(p, left)
		case token.INCREMENT:
			return ParseIncExpr(p, left)
		case token.DECREMENT:
			return ParseDecExpr(p, left)
		}
		return ParseBinaryExpr(p, left)
	})
}

func Printer(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
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
	case *FunctionExpr:
		return PrintFunctionExpr(pr, v)
	case *CallExpr:
		return PrintCallExpr(pr, v)
	case *IndexExpr:
		return PrintIndexExpr(pr, v)
	case *GroupExpr:
		return PrintGroupExpr(pr, v)
	case *ObjExpr:
		return PrintObjExpr(pr, v)
	case *ArrayExpr:
		return PrintArrayExpr(pr, v)
	case *IncExpr:
		return PrintIncExpr(pr, v)
	case *DecExpr:
		return PrintDecExpr(pr, v)
	case *AssignExpr:
		return PrintAssignExpr(pr, v)
	case *UnaryExpr:
		return PrintUnaryExpr(pr, v)
	case *BinaryExpr:
		return PrintBinaryExpr(pr, v)
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
	case *MemberExpr:
		return PrintMemberExpr(pr, v)
	case *SemiStmt:
		return PrintSemiStmt(pr, v)
	case *ImportStmt:
		return PrintImportStmt(pr, v)
	case *ExportStmt:
		return PrintExportStmt(pr, v)
	}
	return next(node)
}
