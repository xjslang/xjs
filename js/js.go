package js

import (
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// ExpectSemi expects a semicolon or any other symbol that acts as a
// "statement terminator", such as '}' or ')'. If the statement terminator is a
// semicolon, then it consumes it and advances to the next token.
func ExpectSemi(p *parser.Parser) (tok token.Token, err error) {
	tok = p.CurrentToken
	switch tok.Type {
	case token.SEMICOLON:
		p.AdvanceToken()
		// ignore consecutive semicolons without leading trivia
		for {
			if t := p.CurrentToken; t.Type != token.SEMICOLON || len(t.LeadingTrivia) > 0 {
				break
			}
			p.AdvanceToken()
		}
		return
	case token.RBRACE, token.RPAREN, token.EOF:
		tok = token.Token{
			Type:     token.SEMICOLON,
			Literal:  token.SEMICOLON.String(),
			Position: tok.Position,
		}
		return
	default:
		if tok.AfterNewline {
			tok = token.Token{
				Type:     token.SEMICOLON,
				Literal:  token.SEMICOLON.String(),
				Position: tok.Position,
			}
			return
		}
	}
	err = p.Error(token.SEMICOLON.String() + " expected")
	return
}
