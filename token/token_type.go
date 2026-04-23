package token

import "strconv"

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

	// operators
	ASSIGN // =
	DIVIDE // /

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

var tokenLiterals = map[TokenType]string{
	EOF:           "end of file",
	EQ:            "==",
	NOT_EQ:        "!=",
	ASSIGN:        "=",
	NOT:           "!",
	LT:            "<",
	LTE:           "<=",
	GT:            ">",
	GTE:           ">=",
	IDENT:         "identifier",
	NUMBER:        "number",
	STRING:        "string",
	UNKNOWN:       "unknown",
	ILLEGAL:       "ILLEGAL",
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
	DIVIDE:        "/",
}

func (tt TokenType) String() string {
	lit, ok := tokenLiterals[tt]
	if !ok {
		return "UNKNOWN(" + strconv.Itoa(int(tt)) + ")"
	}
	return lit
}

func Lookup(lit string) TokenType {
	switch lit {
	case "let":
		return LET
	case "function":
		return FUNCTION
	}
	return IDENT
}
