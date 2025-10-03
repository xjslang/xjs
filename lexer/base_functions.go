package lexer

import "github.com/xjslang/xjs/token"

// NewToken creates a new token.
func (l *Lexer) NewToken(tokenType token.Type, literal string, line, column int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		// TODO: we don't need `line` and `column` to be passed as parameters, since they can be obtained from the lexer itself
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
			tok = l.NewToken(token.EQ, "==", line, column)
		} else {
			tok = l.NewToken(token.ASSIGN, string(l.CurrentChar), line, column)
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.NOT_EQ, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.NOT, string(l.CurrentChar), line, column)
		}
	case '<':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.LTE, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.LT, string(l.CurrentChar), line, column)
		}
	case '>':
		if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.GTE, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.GT, string(l.CurrentChar), line, column)
		}
	case '&':
		if l.PeekChar() == '&' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.AND, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar), line, column)
		}
	case '|':
		if l.PeekChar() == '|' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.OR, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar), line, column)
		}
	case '+':
		if l.PeekChar() == '+' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.INCREMENT, string(ch)+string(l.CurrentChar), line, column)
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.PLUS_ASSIGN, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.PLUS, string(l.CurrentChar), line, column)
		}
	case '-':
		if l.PeekChar() == '-' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.DECREMENT, string(ch)+string(l.CurrentChar), line, column)
		} else if l.PeekChar() == '=' {
			ch := l.CurrentChar
			l.ReadChar()
			tok = l.NewToken(token.MINUS_ASSIGN, string(ch)+string(l.CurrentChar), line, column)
		} else {
			tok = l.NewToken(token.MINUS, string(l.CurrentChar), line, column)
		}
	case '*':
		tok = l.NewToken(token.MULTIPLY, string(l.CurrentChar), line, column)
	case '/':
		if l.PeekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Skip the comment and get the next token
		} else {
			tok = l.NewToken(token.DIVIDE, string(l.CurrentChar), line, column)
		}
	case '%':
		tok = l.NewToken(token.MODULO, string(l.CurrentChar), line, column)
	case ',':
		tok = l.NewToken(token.COMMA, string(l.CurrentChar), line, column)
	case ';':
		tok = l.NewToken(token.SEMICOLON, string(l.CurrentChar), line, column)
	case ':':
		tok = l.NewToken(token.COLON, string(l.CurrentChar), line, column)
	case '.':
		tok = l.NewToken(token.DOT, string(l.CurrentChar), line, column)
	case '(':
		tok = l.NewToken(token.LPAREN, string(l.CurrentChar), line, column)
	case ')':
		tok = l.NewToken(token.RPAREN, string(l.CurrentChar), line, column)
	case '{':
		tok = l.NewToken(token.LBRACE, string(l.CurrentChar), line, column)
	case '}':
		tok = l.NewToken(token.RBRACE, string(l.CurrentChar), line, column)
	case '[':
		tok = l.NewToken(token.LBRACKET, string(l.CurrentChar), line, column)
	case ']':
		tok = l.NewToken(token.RBRACKET, string(l.CurrentChar), line, column)
	case '"':
		tok = l.NewToken(token.STRING, l.readString('"'), line, column)
	case '\'':
		tok = l.NewToken(token.STRING, l.readString('\''), line, column)
	case '`':
		tok = l.NewToken(token.RAW_STRING, l.readRawString(), line, column)
	case 0:
		tok = l.NewToken(token.EOF, "", line, column)
	default:
		if isLetter(l.CurrentChar) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			tok = l.NewToken(tokType, literal, line, column)
			// Don't call ReadChar() here because readIdentifier() already does it
			return tok
		} else if isDigit(l.CurrentChar) {
			literal, tokType := l.readNumber()
			tok = l.NewToken(tokType, literal, line, column)
			return tok
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar), line, column)
		}
	}

	l.ReadChar()
	return tok
}
