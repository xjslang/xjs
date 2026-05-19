package ast

import "github.com/xjslang/xjs/token"

type Let struct {
	// keywords and delimiters
	LetToken    token.Token
	AssignToken token.Token
	SemiToken   token.Token

	Name  token.Token
	Value Node
}

func (node *Let) Type() string {
	return "Let"
}

type Block struct {
	// keywords and delimiters
	LbraceToken token.Token
	RbraceToken token.Token

	Statements []Node
}

func (node *Block) Type() string {
	return "block statement"
}

type Function struct {
	// keywords and delimiters
	FunctionToken token.Token
	LparenToken   token.Token
	RparenToken   token.Token

	Name token.Token
	Body *Block
}

func (node *Function) Type() string {
	return "function declaration"
}
