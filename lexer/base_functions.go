package lexer

import "github.com/xjslang/xjs/token"

func baseNextToken(l *Lexer) token.Token {
	var tok token.Token
	l.skipWhitespace()
	Line := l.Line
	Column := l.Column

	switch l.CurrentChar {
	case '=':
		if l.PeekChar() == '=' {
			l.ReadChar()
			tok = token.Token{Type: token.EQ, Literal: "==", Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.NOT, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '<':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '>':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '&':
		if l.PeekChar() == '&' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '|':
		if l.PeekChar() == '|' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '+':
		if l.PeekChar() == '+' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '-':
		if l.PeekChar() == '-' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: string(ch) + string(l.CurrentChar), Line: Line, Column: Column}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '*':
		tok = token.Token{Type: token.MULTIPLY, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '/':
		if l.PeekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Skip the comment and get the next token
		} else {
			tok = token.Token{Type: token.DIVIDE, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	case '%':
		tok = token.Token{Type: token.MODULO, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case ':':
		tok = token.Token{Type: token.COLON, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '.':
		tok = token.Token{Type: token.DOT, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.CurrentChar), Line: Line, Column: Column}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
		tok.Line = Line
		tok.Column = Column
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString('\'')
		tok.Line = Line
		tok.Column = Column
	case '`':
		tok.Type = token.RAW_STRING
		tok.Literal = l.readRawString()
		tok.Line = Line
		tok.Column = Column
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = Line
		tok.Column = Column
	default:
		if isLetter(l.CurrentChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = Line
			tok.Column = Column
			// Don't call ReadChar() here because readIdentifier() already does it
			return tok
		} else if isDigit(l.CurrentChar) {
			tok.Literal, tok.Type = l.readNumber()
			tok.Line = Line
			tok.Column = Column
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.CurrentChar), Line: Line, Column: Column}
		}
	}

	l.ReadChar()
	return tok
}
