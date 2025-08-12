package main

import (
	"fmt"
)

// TokenType represents the different types of tokens
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT  // variables, functions
	INT    // 123
	FLOAT  // 123.45
	STRING // "hello"

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %

	// Comparison operators
	EQ     // ==
	NOT_EQ // !=
	EQ_STRICT     // ===
	NOT_EQ_STRICT // !==
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=

	// Logical operators
	AND // &&
	OR  // ||
	NOT // !

	// Increment/Decrement
	INCREMENT // ++
	DECREMENT // --

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	DOT       // .

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// Keywords
	FUNCTION
	LET
	IF
	ELSE
	WHILE
	FOR
	RETURN
	TRUE
	FALSE
	NULL
)

// Token represents a lexer token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %q, Line: %d, Col: %d}", 
		t.Type, t.Literal, t.Line, t.Column)
}

// String returns the token type name
func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULO:
		return "MODULO"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case EQ_STRICT:
		return "EQ_STRICT"
	case NOT_EQ_STRICT:
		return "NOT_EQ_STRICT"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case DOT:
		return "DOT"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case FUNCTION:
		return "FUNCTION"
	case LET:
		return "LET"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case RETURN:
		return "RETURN"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	default:
		return "UNKNOWN"
	}
}

// Language keywords
var keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
}

// Lexer main lexer structure
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line
	column       int  // current column
}

// New creates a new lexer
func NewLexer(input string) *Lexer {
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

// peekChar returns the next character without advancing position
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
func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	tokenType := INT
	
	for isDigit(l.ch) {
		l.readChar()
	}
	
	// Check if it's a decimal number
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
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

// NextToken scans the input and returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token
	
	l.skipWhitespace()
	
	// Capture position before processing the token
	line := l.line
	column := l.column
	
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = Token{Type: EQ_STRICT, Literal: "===", Line: line, Column: column}
			} else {
				tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
			}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.ch), Line: line, Column: column}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = Token{Type: NOT_EQ_STRICT, Literal: "!==", Line: line, Column: column}
			} else {
				tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
			}
		} else {
			tok = Token{Type: NOT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = Token{Type: LT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = Token{Type: GT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: INCREMENT, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = Token{Type: PLUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DECREMENT, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = Token{Type: MINUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '*':
		tok = Token{Type: MULTIPLY, Literal: string(l.ch), Line: line, Column: column}
	case '/':
		tok = Token{Type: DIVIDE, Literal: string(l.ch), Line: line, Column: column}
	case '%':
		tok = Token{Type: MODULO, Literal: string(l.ch), Line: line, Column: column}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch), Line: line, Column: column}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch), Line: line, Column: column}
	case '.':
		tok = Token{Type: DOT, Literal: string(l.ch), Line: line, Column: column}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Line: line, Column: column}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Line: line, Column: column}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '[':
		tok = Token{Type: LBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
		tok.Line = line
		tok.Column = column
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString('\'')
		tok.Line = line
		tok.Column = column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = line
		tok.Column = column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
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
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	}
	
	l.readChar()
	return tok
}

// newToken creates a new token
func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// lookupIdent checks if an identifier is a keyword
func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Example function to test the lexer
func main() {
	input := `
		let x = 5
		let y = 10.5
		let name = "Hello World"
		
		function add(a, b) {
			return a + b
		}
		
		if (x < y) {
			console.log("x is less than y")
		}
	`
	
	lexer := NewLexer(input)
	
	for {
		tok := lexer.NextToken()
		fmt.Println(tok)
		if tok.Type == EOF {
			break
		}
	}
}