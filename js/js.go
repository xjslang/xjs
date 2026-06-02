package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Ident struct {
	ast.ExprNode
	Name token.Token
}

type Literal struct {
	ast.ExprNode
	Value token.Token
}
