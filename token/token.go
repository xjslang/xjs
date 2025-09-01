// Package token defines the token types and structures used by the XJS lexer and parser.
package token

import "fmt"

// Type represents the different types of tokens
type Type int

const (
	// Special tokens
	ILLEGAL Type = iota
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
	COLON     // :
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

// Token represents a lexer token with its type, literal value, and position information
type Token struct {
	Type    Type
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
func (tt Type) String() string {
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
	case COLON:
		return "COLON"
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

// Keywords maps string literals to their token types
var Keywords = map[string]Type{
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

// LookupIdent checks if an identifier is a keyword and returns the appropriate token type
func LookupIdent(ident string) Type {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
