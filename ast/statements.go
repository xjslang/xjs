package ast

import "github.com/xjslang/xjs/scanner"

type EOF struct {
	EOFToken scanner.Token
}

func (node *EOF) Kind() string {
	return "EOF"
}

type Let struct {
	// keywords and delimiters
	LetToken    scanner.Token
	AssignToken scanner.Token
	SemiToken   scanner.Token

	Name  scanner.Token
	Value Node
}

func (node *Let) Kind() string {
	return "let statement"
}

type Block struct {
	// keywords and delimiters
	LbraceToken scanner.Token
	RbraceToken scanner.Token

	Statements []Node
}

func (node *Block) Kind() string {
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

func (node *Function) Kind() string {
	return "function declaration"
}
