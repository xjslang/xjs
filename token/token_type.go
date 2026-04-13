package token

import "strconv"

type TokenType int

const (
	EOF TokenType = iota
	EQ
	NOT_EQ
	ASSIGN
	NOT
	LOWER
	LOWER_OR_EQ
	GREATER
	GREATER_OR_EQ
	IDENT
	NUMBER
	STRING
	UNKNOWN
	ILLEGAL
	LET
	SEMI
	FUNCTION
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	NEWLINE
	SINGLELINE_COMMENT
	MULTILINE_COMMENT
	DIVISION
)

var tokenLiterals = map[TokenType]string{
	EOF:                "EOF",
	EQ:                 "EQ",
	NOT_EQ:             "NOT_EQ",
	ASSIGN:             "ASSIGN",
	NOT:                "NOT",
	LOWER:              "LOWER",
	LOWER_OR_EQ:        "LOWER_OR_EQ",
	GREATER:            "GREATER",
	GREATER_OR_EQ:      "GREATER_OR_EQ",
	IDENT:              "IDENT",
	NUMBER:             "NUMBER",
	STRING:             "STRING",
	UNKNOWN:            "UNKNOWN",
	ILLEGAL:            "ILLEGAL",
	LET:                "LET",
	SEMI:               "SEMI",
	FUNCTION:           "FUNCTION",
	LPAREN:             "LPAREN",
	RPAREN:             "RPAREN",
	LBRACE:             "LBRACE",
	RBRACE:             "RBRACE",
	NEWLINE:            "NEWLINE",
	SINGLELINE_COMMENT: "SINGLELINE_COMMENT",
	MULTILINE_COMMENT:  "MULTILINE_COMMENT",
	DIVISION:           "DIVISION",
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
