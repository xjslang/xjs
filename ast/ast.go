package ast

type Node interface {
	node()
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

func (ExprNode) node()     {}
func (ExprNode) exprNode() {}

type StmtNode struct{}

func (StmtNode) node()     {}
func (StmtNode) stmtNode() {}
