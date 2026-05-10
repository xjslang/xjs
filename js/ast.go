package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type InfixOperator struct {
	LeftValue  ast.Expression
	Operator   token.Token
	RightValue ast.Expression
}

func (node *InfixOperator) Type() string {
	return "InfixOperator"
}

type LetStatement struct {
	Name  token.Token
	Value ast.Expression
}

func (node *LetStatement) Type() string {
	return "LetStatement"
}

type FunctionDeclaration struct {
	Name token.Token
	Body *BlockStatement
}

func (node *FunctionDeclaration) Type() string {
	return "FunctionDeclaration"
}

type BlockStatement struct {
	Statements []ast.Statement
}

func (node *BlockStatement) Type() string {
	return "BlockStatement"
}
