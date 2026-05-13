package ast

type Node interface {
	Kind() string
}

type Program struct {
	Statements []Node
}

func (node *Program) Kind() string {
	return "Program"
}
