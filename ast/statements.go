package ast

import "github.com/xjslang/xjs/scanner"

type Let struct {
	Name  scanner.Token
	Value Node
}

func (node *Let) Type() string {
	return "let statement"
}

type Block struct {
	Statements []Node
}

func (node *Block) Type() string {
	return "block statement"
}

type Function struct {
	Name scanner.Token
	Body *Block
}

func (node *Function) Type() string {
	return "function declaration"
}
