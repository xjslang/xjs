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
	UNKNOWN
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
	case UNKNOWN:
		return "UNKNOWN"
	}
	return fmt.Sprintf("UNKNOWN(%d)", tt)
}
