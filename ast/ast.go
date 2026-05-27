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

type Program struct {
	EOFToken token.Token

	Stmts []Node
}

func (node *Program) Type() string {
	return "Program"
}

type Block struct {
	LbraceToken token.Token
	RbraceToken token.Token

	Stmts []Node
}

func (node *Block) Type() string {
	return "Block"
}

type ExprStmt struct {
	SemiToken token.Token

	Expr Node
}

func (node *ExprStmt) Type() string {
	return "ExprStmt"
}

type LetStmt struct {
	LetToken    token.Token
	AssignToken token.Token
	SemiToken   token.Token

	Name  token.Token
	Value Node
}

func (node *LetStmt) Type() string {
	return "LetStmt"
}

type FuncDecl struct {
	FunctionToken token.Token
	LparenToken   token.Token
	RparenToken   token.Token

	Name token.Token
	Body *Block
}

func (node *FuncDecl) Type() string {
	return "FuncDecl"
}

type PrefixExpr struct {
	Operator token.Token
	Value    Node
}

func (node *PrefixExpr) Type() string {
	return "PrefixExpr"
}

type InfixExpr struct {
	LeftValue  Node
	Operator   token.Token
	RightValue Node
}

func (node *InfixExpr) Type() string {
	return "BinaryExpr"
}

type CallExpr struct {
	LparenToken token.Token
	RparenToken token.Token

	Function  Node
	Arguments []Node
}

func (node *CallExpr) Type() string {
	return "CallExpr"
}

type BasicLit struct {
	Value token.Token
}

func (node *BasicLit) Type() string {
	return "BasicLit"
}

type ParenExpr struct {
	LparenToken token.Token
	RparenToken token.Token

	Value Node
}

func (node *ParenExpr) Type() string {
	return "ParenExpr"
}
