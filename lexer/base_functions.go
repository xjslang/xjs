package lexer

import "github.com/xjslang/xjs/token"

// NewToken creates a new token using the current position as both start and end.
// For single-character tokens, start and end positions are the same.
func (l *Lexer) NewToken(tokenType token.Type, literal string) token.Token {
	return token.Token{
		Type:         tokenType,
		Literal:      literal,
		Line:         l.Line,
		Column:       l.Column,
		StartLine:    l.Line,
		StartColumn:  l.Column,
		AfterNewline: l.hadNewlineBefore,
	}
}

// NewTokenAt creates a new token with explicit start position.
// This should be used when the start position was captured before reading a multi-character literal.
func (l *Lexer) NewTokenAt(tokenType token.Type, literal string, startLine, startColumn int) token.Token {
	return token.Token{
		Type:         tokenType,
		Literal:      literal,
		Line:         l.Line,
		Column:       l.Column,
		StartLine:    startLine,
		StartColumn:  startColumn,
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
		// Capture position BEFORE reading the string
		startLine, startColumn := l.Line, l.Column
		tok = l.NewTokenAt(token.STRING, l.readString('"'), startLine, startColumn)
	case '\'':
		// Capture position BEFORE reading the string
		startLine, startColumn := l.Line, l.Column
		tok = l.NewTokenAt(token.STRING, l.readString('\''), startLine, startColumn)
	case '`':
		// Capture position BEFORE reading the raw string
		startLine, startColumn := l.Line, l.Column
		tok = l.NewTokenAt(token.RAW_STRING, l.readRawString(), startLine, startColumn)
	case 0:
		tok = l.NewToken(token.EOF, "")
	default:
		if isLetter(l.CurrentChar) {
			// Capture position BEFORE reading the identifier
			startLine, startColumn := l.Line, l.Column
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			tok = l.NewTokenAt(tokType, literal, startLine, startColumn)
			// Don't call ReadChar() here because readIdentifier() already does it
			return tok
		} else if isDigit(l.CurrentChar) {
			// Capture position BEFORE reading the number
			startLine, startColumn := l.Line, l.Column
			literal, tokType := l.readNumber()
			tok = l.NewTokenAt(tokType, literal, startLine, startColumn)
			return tok
		} else {
			tok = l.NewToken(token.ILLEGAL, string(l.CurrentChar))
		}
	}

	l.ReadChar()
	return tok
}
