package token

import "github.com/xjslang/xjs/source"

type Token struct {
	source.Position
	Type          TokenType
	Literal       string
	LeadingTrivia []string
	AfterNewline  bool
}
