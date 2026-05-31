package js

import (
	"github.com/xjslang/xjs/token"
)

type Ident struct {
	Name token.Token
}

func (node *Ident) Type() string {
	return "Ident"
}

type Literal struct {
	Value token.Token
}

func (node *Literal) Type() string {
	return "Literal"
}
