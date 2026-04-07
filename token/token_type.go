package token

import "fmt"

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
)

func (tt TokenType) String() string {
	switch tt {
	case EOF:
		return "EOF"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case NOT:
		return "NOT"
	case LOWER:
		return "LOWER"
	case LOWER_OR_EQ:
		return "LOWER_OR_EQ"
	case GREATER:
		return "GREATER"
	case GREATER_OR_EQ:
		return "GREATER_OR_EQ"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case UNKNOWN:
		return "UNKNOWN"
	case ILLEGAL:
		return "ILLEGAL"
	}
	return fmt.Sprintf("UNKNOWN(%d)", tt)
}
