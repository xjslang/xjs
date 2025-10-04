package lexer

import "github.com/xjslang/xjs/token"

// NewToken creates a new token.
func (l *Lexer) NewToken(tokenType token.Type, literal string) token.Token {
	return token.Token{
		Type:         tokenType,
		Literal:      literal,
		Line:         l.Line,
		Column:       l.Column,
		AfterNewline: l.hadNewlineBefore,
	}
}

func baseNextToken(l *Lexer) token.Token {
	var tok token.Token

	switch l.CurrentChar {
	case '=':
		if l.PeekChar() == '=' {
			l.ReadChar()
			tok = l.NewToken(token.EQ, "==")
		} else {
			tok = l.NewToken(token.ASSIGN, string(l.CurrentChar))
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.NOT_EQ, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.NOT, string(l.CurrentChar))
		}
	case '<':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.LTE, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.LT, string(l.CurrentChar))
		}
	case '>':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.GTE, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.GT, string(l.CurrentChar))
		}
	case '&':
		if l.PeekChar() == '&' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.AND, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar))
		}
	case '|':
		if l.PeekChar() == '|' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.OR, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar))
		}
	case '+':
		if l.PeekChar() == '+' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.INCREMENT, string(ch)+string(l.CurrentChar))
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.PLUS_ASSIGN, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.PLUS, string(l.CurrentChar))
		}
	case '-':
		if l.PeekChar() == '-' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.DECREMENT, string(ch)+string(l.CurrentChar))
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.MINUS_ASSIGN, string(ch)+string(l.CurrentChar))
		} else {
			tok = l.NewToken(token.MINUS, string(l.CurrentChar))
		}
	case '*':
		tok = l.NewToken(token.MULTIPLY, string(l.CurrentChar))
	case '/':
		if l.PeekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Skip the comment and get the next token
		} else {
			tok = l.NewToken(token.DIVIDE, string(l.CurrentChar))
		}
	case '%':
		tok = l.NewToken(token.MODULO, string(l.CurrentChar))
	case ',':
		tok = l.NewToken(token.COMMA, string(l.CurrentChar))
	case ';':
		tok = l.NewToken(token.SEMICOLON, string(l.CurrentChar))
	case ':':
		tok = l.NewToken(token.COLON, string(l.CurrentChar))
	case '.':
		tok = l.NewToken(token.DOT, string(l.CurrentChar))
	case '(':
		tok = l.NewToken(token.LPAREN, string(l.CurrentChar))
	case ')':
		tok = l.NewToken(token.RPAREN, string(l.CurrentChar))
	case '{':
		tok = l.NewToken(token.LBRACE, string(l.CurrentChar))
	case '}':
		tok = l.NewToken(token.RBRACE, string(l.CurrentChar))
	case '[':
		tok = l.NewToken(token.LBRACKET, string(l.CurrentChar))
	case ']':
		tok = l.NewToken(token.RBRACKET, string(l.CurrentChar))
	case '"':
		tok = l.NewToken(token.STRING, l.readString('"'))
	case '\'':
		tok = l.NewToken(token.STRING, l.readString('\''))
	case '`':
		tok = l.NewToken(token.RAW_STRING, l.readRawString())
	case 0:
		tok = l.NewToken(token.EOF, "")
	default:
		if isLetter(l.CurrentChar) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			tok = l.NewToken(tokType, literal)
			// Don't call ReadChar() here because readIdentifier() already does it
			return tok
		} else if isDigit(l.CurrentChar) {
			literal, tokType := l.readNumber()
			tok = l.NewToken(tokType, literal)
			return tok
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar))
		}
	}

	l.ReadChar()
	return tok
}
