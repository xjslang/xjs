package ast

import "github.com/xjslang/xjs/scanner"

type Let struct {
	Name  scanner.Token
	Value Node
}

func (node *Let) Kind() string {
	return "let statement"
}

type Block struct {
	Statements []Node
}

func (node *Block) Kind() string {
	return "block statement"
}

type Function struct {
	Name scanner.Token
	Body *Block
}

func (node *Function) Kind() string {
	return "function declaration"
}
