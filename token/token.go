package token

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) String() string {
	return t.Literal
}
