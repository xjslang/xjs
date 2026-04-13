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

	// comments
	SINGLELINE_COMMENT // //
	MULTILINE_COMMENT  // /*

	// keywords
	LET
	FUNCTION
)

var tokenLiterals = map[TokenType]string{
	EOF:                "EOF",
	EQ:                 "EQ",
	NOT_EQ:             "NOT_EQ",
	ASSIGN:             "ASSIGN",
	NOT:                "NOT",
	LT:                 "LT",
	LTE:                "LTE",
	GT:                 "GT",
	GTE:                "GTE",
	IDENT:              "IDENT",
	NUMBER:             "NUMBER",
	STRING:             "STRING",
	UNKNOWN:            "UNKNOWN",
	ILLEGAL:            "ILLEGAL",
	LET:                "LET",
	SEMICOLON:          "SEMICOLON",
	FUNCTION:           "FUNCTION",
	LPAREN:             "LPAREN",
	RPAREN:             "RPAREN",
	LBRACE:             "LBRACE",
	RBRACE:             "RBRACE",
	NEWLINE:            "NEWLINE",
	SINGLELINE_COMMENT: "SINGLELINE_COMMENT",
	MULTILINE_COMMENT:  "MULTILINE_COMMENT",
	DIVIDE:             "DIVIDE",
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
