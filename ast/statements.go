package ast

import "github.com/xjslang/xjs/scanner"

type Let struct {
	// keywords and delimiters
	LetToken    scanner.Token
	AssignToken scanner.Token
	SemiToken   scanner.Token

	Name  scanner.Token
	Value Node
}

func (node *Let) Type() string {
	return "Let"
}

type Block struct {
	// keywords and delimiters
	LbraceToken scanner.Token
	RbraceToken scanner.Token

	Statements []Node
}

func (node *Block) Type() string {
	return "block statement"
}

type Function struct {
	// keywords and delimiters
	FunctionToken scanner.Token
	LparenToken   scanner.Token
	RparenToken   scanner.Token

	Name scanner.Token
	Body *Block
}

func (node *Function) Type() string {
	return "function declaration"
}
