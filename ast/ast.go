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

type Decl interface {
	Stmt
	declNode()
}

// default implementations
type (
	BaseNode struct{}
	BaseExpr struct{ BaseNode }
	BaseStmt struct{ BaseNode }
	BaseDecl struct{ BaseStmt }
)

func (BaseNode) node()     {}
func (BaseExpr) exprNode() {}
func (BaseStmt) stmtNode() {}
func (BaseDecl) declNode() {}
