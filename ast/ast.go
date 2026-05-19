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
	EOFToken   token.Token
	Statements []Node
}

func (node *Program) Type() string {
	return "Program"
}

type Block struct {
	// keywords and delimiters
	LbraceToken token.Token
	RbraceToken token.Token

	Statements []Node
}

func (node *Block) Type() string {
	return "Block"
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

type BinaryExpr struct {
	LeftValue  Node
	Operator   token.Token
	RightValue Node
}

func (node *BinaryExpr) Type() string {
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

type Integer struct {
	Value token.Token
}

func (node *Integer) Type() string {
	return "Integer"
}

type String struct {
	Value token.Token
}

func (node *String) Type() string {
	return "String"
}

type Boolean struct {
	Value token.Token
}

func (node *Boolean) Type() string {
	return "Boolean"
}

type ParenExpr struct {
	LparenToken token.Token
	RparenToken token.Token

	Value Node
}

func (node *ParenExpr) Type() string {
	return "ParenExpr"
}
