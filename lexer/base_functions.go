package lexer

import "github.com/xjslang/xjs/token"

// newToken creates a new token with ASI metadata
func (l *Lexer) newToken(tokenType token.Type, literal string, line, column int) token.Token {
	return token.Token{
		Type:         tokenType,
		Literal:      literal,
		Line:         line,
		Column:       column,
		AfterNewline: l.hadNewlineBefore,
	}
}

func baseNextToken(l *Lexer) token.Token {
	var tok token.Token
	l.skipWhitespace()
	line := l.Line
	column := l.Column

	switch l.CurrentChar {
	case '=':
		if l.PeekChar() == '=' {
			l.ReadChar()
			tok = l.newToken(token.EQ, "==", line, column)
		} else {
			tok = l.newToken(token.ASSIGN, string(l.CurrentChar), line, column)
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.NOT_EQ, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.NOT, string(l.CurrentChar), line, column)
		}
	case '<':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.LTE, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.LT, string(l.CurrentChar), line, column)
		}
	case '>':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.GTE, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.GT, string(l.CurrentChar), line, column)
		}
	case '&':
		if l.PeekChar() == '&' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.AND, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.CurrentChar), line, column)
		}
	case '|':
		if l.PeekChar() == '|' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.OR, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.CurrentChar), line, column)
		}
	case '+':
		if l.PeekChar() == '+' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.INCREMENT, string(ch)+string(l.CurrentChar), line, column)
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.PLUS_ASSIGN, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.PLUS, string(l.CurrentChar), line, column)
		}
	case '-':
		if l.PeekChar() == '-' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.DECREMENT, string(ch)+string(l.CurrentChar), line, column)
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.newToken(token.MINUS_ASSIGN, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.newToken(token.MINUS, string(l.CurrentChar), line, column)
		}
	case '*':
		tok = l.newToken(token.MULTIPLY, string(l.CurrentChar), line, column)
	case '/':
		if l.PeekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Skip the comment and get the next token
		} else {
			tok = l.newToken(token.DIVIDE, string(l.CurrentChar), line, column)
		}
	case '%':
		tok = l.newToken(token.MODULO, string(l.CurrentChar), line, column)
	case ',':
		tok = l.newToken(token.COMMA, string(l.CurrentChar), line, column)
	case ';':
		tok = l.newToken(token.SEMICOLON, string(l.CurrentChar), line, column)
	case ':':
		tok = l.newToken(token.COLON, string(l.CurrentChar), line, column)
	case '.':
		tok = l.newToken(token.DOT, string(l.CurrentChar), line, column)
	case '(':
		tok = l.newToken(token.LPAREN, string(l.CurrentChar), line, column)
	case ')':
		tok = l.newToken(token.RPAREN, string(l.CurrentChar), line, column)
	case '{':
		tok = l.newToken(token.LBRACE, string(l.CurrentChar), line, column)
	case '}':
		tok = l.newToken(token.RBRACE, string(l.CurrentChar), line, column)
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.CurrentChar), line, column)
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.CurrentChar), line, column)
	case '"':
		tok = l.newToken(token.STRING, l.readString('"'), line, column)
	case '\'':
		tok = l.newToken(token.STRING, l.readString('\''), line, column)
	case '`':
		tok = l.newToken(token.RAW_STRING, l.readRawString(), line, column)
	case 0:
		tok = l.newToken(token.EOF, "", line, column)
	default:
		if isLetter(l.CurrentChar) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			tok = l.newToken(tokType, literal, line, column)
			// Don't call ReadChar() here because readIdentifier() already does it
			return tok
		} else if isDigit(l.CurrentChar) {
			literal, tokType := l.readNumber()
			tok = l.newToken(tokType, literal, line, column)
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.CurrentChar), line, column)
		}
	}

	l.ReadChar()
	return tok
}
