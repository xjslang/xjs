package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Ident struct {
	ast.ExprNode
	Name token.Token
}

func (node *Ident) Type() string {
	return "Ident"
}

type Literal struct {
	ast.ExprNode
	Value token.Token
}

func (node *Literal) Type() string {
	return "Literal"
}
