package ast

type Node interface {
	Type() string
}

type Expr interface {
	Node
	exprNode()
}

type ExprNode struct{}

//nolint:unused
func (ExprNode) exprNode() {}
