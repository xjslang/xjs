package ast

import "github.com/xjslang/xjs/token"

type Node interface {
	Type() string
}

type Ident struct {
	Value token.Token
}

func (node *Ident) Type() string {
	return "Ident"
}

type ExprStmt struct {
	SemiToken token.Token

	Expr Node
}

func (node *ExprStmt) Type() string {
	return "ExprStmt"
}

type UnaryExpr struct {
	Operator token.Token
	Value    Node
}

func (node *UnaryExpr) Type() string {
	return "UnaryExpr"
}

type BinaryExpr struct {
	LeftValue  Node
	Operator   token.Token
	RightValue Node
}

func (node *BinaryExpr) Type() string {
	return "BinaryExpr"
}

type BasicLit struct {
	Value token.Token
}

func (node *BasicLit) Type() string {
	return "BasicLit"
}
