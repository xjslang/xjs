package ast

type Node interface {
	Type() string
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type ExprNode struct{}

//nolint:unused
func (ExprNode) exprNode() {}

type StmtNode struct{}

//nolint:unused
func (StmtNode) stmtNode() {}
