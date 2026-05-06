package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

var (
	// maths infix operators
	plusType     = token.RegisterType("+")
	minusType    = token.RegisterType("-")
	multiplyType = token.RegisterType("*")
	divideType   = token.RegisterType("/")
	moduloType   = token.RegisterType("%")

	// keywords
	letType  = token.RegisterType("let")
	funcType = token.RegisterType("function")
)

func opPrecedence(tt token.TokenType) int {
	switch tt {
	case plusType, minusType:
		return 1
	case multiplyType, divideType, moduloType:
		return 2
	}
	return 0
}

func MathPlugin(b *parser.Builder) {
	b.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		switch l.CurrentChar {
		case '+':
			l.AdvanceChar()
			return token.Token{Type: plusType, Literal: plusType.String()}
		case '-':
			l.AdvanceChar()
			return token.Token{Type: minusType, Literal: minusType.String()}
		case '*':
			l.AdvanceChar()
			return token.Token{Type: multiplyType, Literal: multiplyType.String()}
		case '/':
			if c := l.PeekChar(); c == '/' || c == '*' {
				break
			}
			l.AdvanceChar()
			return token.Token{Type: divideType, Literal: divideType.String()}
		case '%':
			l.AdvanceChar()
			return token.Token{Type: moduloType, Literal: moduloType.String()}
		}
		return next()
	})
	fn := func(op token.Token, left, right ast.Expression) ast.Expression {
		return &InfixOperator{
			LeftValue:  left,
			Operator:   op,
			RightValue: right,
		}
	}
	b.RegisterInfixOperator(plusType, opPrecedence(plusType), fn)
	b.RegisterInfixOperator(minusType, opPrecedence(minusType), fn)
	b.RegisterInfixOperator(multiplyType, opPrecedence(multiplyType), fn)
	b.RegisterInfixOperator(divideType, opPrecedence(divideType), fn)
	b.RegisterInfixOperator(moduloType, opPrecedence(moduloType), fn)
}

func LetPlugin(b *parser.Builder) {
	b.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		tok := next()
		if lit := letType.String(); tok.Type == token.IDENT && tok.Literal == lit {
			return token.Token{Type: letType, Literal: lit}
		}
		return tok
	})
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == letType {
			stmt, err := ParseLetStatement(p)
			if err != nil {
				return nil, err
			}
			return stmt, nil
		}
		return next()
	})
}

func FunctionPlugin(b *parser.Builder) {
	b.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		tok := next()
		if lit := funcType.String(); tok.Type == token.IDENT && tok.Literal == lit {
			return token.Token{Type: funcType, Literal: lit}
		}
		return tok
	})
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == funcType {
			stmt, err := ParseFunctionDeclaration(p)
			if err != nil {
				return nil, err
			}
			return stmt, nil
		}
		return next()
	})
}
