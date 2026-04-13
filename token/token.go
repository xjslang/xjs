package token

type Token struct {
	Type          TokenType
	Literal       string
	LeadingTrivia []string
}
