package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type ForStmt struct {
	ForToken    token.Token
	LparenToken token.Token
	RparenToken token.Token

	Condition ast.Node
}
