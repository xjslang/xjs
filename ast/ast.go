package ast

import "github.com/xjslang/xjs/scanner"

type Node interface {
	Type() string
}

type Program struct {
	EOFToken   scanner.Token
	Statements []Node
}

func (node *Program) Type() string {
	return "Program"
}
