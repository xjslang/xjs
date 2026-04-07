package ast

type Statement interface {
	statementNode()
}
type Expression interface {
	expressionNode()
}

// Identifier Represents a variable name or reference.
// Note: Identifier is not a Statement or Expression.
type Identifier struct {
	Value string
}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

type FunctionStatement struct {
	Body *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}

// Implements Expression
type IntegerLiteral struct {
	Value string
}

func (il *IntegerLiteral) expressionNode() {}

// Implements Expression
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
