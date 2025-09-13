package lexer

import "github.com/xjslang/xjs/token"

func baseNextToken(l *Lexer) token.Token {
	var tok token.Token
	l.skipWhitespace()
	line := l.line
	column := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "==", Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: line, Column: column}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.NOT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '*':
		tok = token.Token{Type: token.MULTIPLY, Literal: string(l.ch), Line: line, Column: column}
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Skip the comment and get the next token
		} else {
			tok = token.Token{Type: token.DIVIDE, Literal: string(l.ch), Line: line, Column: column}
		}
	case '%':
		tok = token.Token{Type: token.MODULO, Literal: string(l.ch), Line: line, Column: column}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: line, Column: column}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: line, Column: column}
	case ':':
		tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: line, Column: column}
	case '.':
		tok = token.Token{Type: token.DOT, Literal: string(l.ch), Line: line, Column: column}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: line, Column: column}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: line, Column: column}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
		tok.Line = line
		tok.Column = column
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString('\'')
		tok.Line = line
		tok.Column = column
	case '`':
		tok.Type = token.RAW_STRING
		tok.Literal = l.readRawString()
		tok.Line = line
		tok.Column = column
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = line
		tok.Column = column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = column
			// Don't call readChar() here because readIdentifier() already does it
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			tok.Line = line
			tok.Column = column
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	}

	l.readChar()
	return tok
}
