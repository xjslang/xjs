package ast

import "github.com/xjslang/xjs/token"

type Node interface {
	Type() string
}

type Program struct {
	EOFToken   token.Token
	Statements []Node
}

func (node *Program) Type() string {
	return "Program"
}
