package ast

import "github.com/xjslang/xjs/scanner"

type Node interface {
	Kind() string
}

type Program struct {
	EOFToken   scanner.Token
	Statements []Node
}

func (node *Program) Kind() string {
	return "Program"
}
