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

// default implementations
type (
	BaseNode struct{}
	BaseExpr struct{ BaseNode }
	BaseStmt struct{ BaseNode }
)

func (BaseNode) node()     {}
func (BaseExpr) exprNode() {}
func (BaseStmt) stmtNode() {}
