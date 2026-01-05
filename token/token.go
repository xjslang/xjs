// Package token defines the token types and structures used by the XJS lexer and parser.
package token

import "fmt"

type Type int

// Position represents a location in the source code (1-indexed).
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

const (
	// Special tokens
	ILLEGAL Type = iota
	EOF

	// Identifiers and literals
	IDENT      // variables, functions
	INT        // 123
	FLOAT      // 123.45
	STRING     // "hello"
	RAW_STRING //
	COMMENT    // // line comment
	BLANK_LINE // empty line separator

	// Operators
	ASSIGN       // =
	PLUS_ASSIGN  // +=
	MINUS_ASSIGN // -=
	PLUS         // +
	MINUS        // -
	MULTIPLY     // *
	DIVIDE       // /
	MODULO       // %

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

	// Custom tokens
	DYNAMIC_TOKENS_START = 1000
)

type Token struct {
	Type         Type
	Literal      string
	Start        Position // Starting position of the token (1-indexed)
	End          Position // Ending position of the token (1-indexed)
	AfterNewline bool     // true if this token follows a newline character
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %q, Start: %d:%d, End: %d:%d}",
		t.Type, t.Literal, t.Start.Line, t.Start.Column, t.End.Line, t.End.Column)
}

func (tt Type) String() string {
	switch tt {
	case ILLEGAL:
		return "illegal"
	case EOF:
		return "end of line"
	case IDENT:
		return "identifier"
	case INT:
		return "integer"
	case FLOAT:
		return "float number"
	case STRING:
		return "string"
	case RAW_STRING:
		return "raw string"
	case COMMENT:
		return "comment"
	case ASSIGN:
		return "="
	case PLUS_ASSIGN:
		return "+="
	case MINUS_ASSIGN:
		return "-="
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MULTIPLY:
		return "*"
	case DIVIDE:
		return "/"
	case MODULO:
		return "%"
	case EQ:
		return "=="
	case NOT_EQ:
		return "!="
	case LT:
		return "<"
	case GT:
		return ">"
	case LTE:
		return "<="
	case GTE:
		return ">="
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case INCREMENT:
		return "++"
	case DECREMENT:
		return "--"
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case COLON:
		return ":"
	case DOT:
		return "."
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case FUNCTION:
		return "function"
	case LET:
		return "let"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case WHILE:
		return "while"
	case FOR:
		return "for"
	case RETURN:
		return "return"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NULL:
		return "undefined"
	default:
		return fmt.Sprintf("unknown(%d)", tt)
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

func LookupIdent(ident string) Type {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
