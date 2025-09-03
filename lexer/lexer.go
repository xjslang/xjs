// Package lexer provides lexical analysis functionality for the XJS language.
// It tokenizes source code into a sequence of tokens that can be consumed by the parser.
package lexer

import "github.com/xjslang/xjs/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line
	column       int  // current column
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position in the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents "EOF"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or decimal)
func (l *Lexer) readNumber() (string, token.Type) {
	position := l.position
	tokenType := token.INT

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check if it's a decimal number
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = token.FLOAT
		l.readChar() // consume the '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], tokenType
}

// readString reads a string literal
func (l *Lexer) readString(delimiter byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == delimiter || l.ch == 0 {
			break
		}
		// TODO: Handle escape sequences like \"
	}
	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
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

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skipLineComment skips characters until the end of line for line comments (//)
func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}
