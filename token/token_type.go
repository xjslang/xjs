package token

import (
	"strconv"
	"sync"
)

type TokenType int

const (
	// special keywords
	EOF TokenType = iota
	IDENT
	ILLEGAL
	UNKNOWN

	// literals
	STRING
	NUMBER
	BOOLEAN

	// operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %

	// comparison operators
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	LTE    // <=
	GT     // >
	GTE    // >=

	// logical operators
	NOT

	// delimiters
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	NEWLINE   // \n

	// line and block comments
	LINE_COMMENT  // //
	BLOCK_COMMENT // /* ... */

	// keywords
	LET
	FUNCTION
)

// initial value of "custom types" created by RegisterType()
const initCustomType TokenType = 1000

var tokenLiterals = map[TokenType]string{
	EOF:           "end of file",
	EQ:            "==",
	NOT_EQ:        "!=",
	ASSIGN:        "=",
	PLUS:          "+",
	MINUS:         "-",
	MULTIPLY:      "*",
	DIVIDE:        "/",
	MODULO:        "%",
	NOT:           "!",
	LT:            "<",
	LTE:           "<=",
	GT:            ">",
	GTE:           ">=",
	IDENT:         "identifier",
	NUMBER:        "number",
	STRING:        "string",
	UNKNOWN:       "unknown",
	ILLEGAL:       "illegal",
	LET:           "let",
	SEMICOLON:     ";",
	FUNCTION:      "function",
	LPAREN:        "(",
	RPAREN:        ")",
	LBRACE:        "{",
	RBRACE:        "}",
	NEWLINE:       "new line",
	LINE_COMMENT:  "line comment",
	BLOCK_COMMENT: "block comment",
}

var (
	nextType   TokenType = initCustomType
	registerMu sync.RWMutex
)

func (tt TokenType) String() string {
	registerMu.RLock()
	defer registerMu.RUnlock()
	lit, ok := tokenLiterals[tt]
	if !ok {
		return "unknown(" + strconv.Itoa(int(tt)) + ")"
	}
	return lit
}

func Lookup(lit string) TokenType {
	switch lit {
	case "let":
		return LET
	case "function":
		return FUNCTION
	case "true", "false":
		return BOOLEAN
	}
	return IDENT
}

func RegisterType(lit string) TokenType {
	registerMu.Lock()
	defer registerMu.Unlock()
	typ := nextType
	tokenLiterals[typ] = lit
	nextType++
	return typ
}
