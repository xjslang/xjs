// Package ast defines the Abstract Syntax Tree nodes for the xjslang language.
// It provides interfaces and concrete types representing different language constructs.
package ast

import "github.com/xjslang/xjs/token"

// Node represents any node in the AST
type Node interface {
	String() string
}

// Statement represents all statement nodes
type Statement interface {
	Node
}

// Expression represents all expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out string
	for _, stmt := range p.Statements {
		out += stmt.String()
	}
	return out
}

// Statements
type LetStatement struct {
	Token token.Token // the LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	out := "let " + ls.Name.String()
	if ls.Value != nil {
		out += " = " + ls.Value.String()
	}
	return out
}

type ReturnStatement struct {
	Token       token.Token // the RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	out := "return"
	if rs.ReturnValue != nil {
		out += " " + rs.ReturnValue.String()
	}
	return out
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type FunctionDeclaration struct {
	Token      token.Token // the FUNCTION token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) String() string {
	out := "function " + fd.Name.String() + "("
	for i, param := range fd.Parameters {
		if i > 0 {
			out += ", "
		}
		out += param.String()
	}
	out += ") " + fd.Body.String()
	return out
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	out := "{"
	for i, stmt := range bs.Statements {
		out += stmt.String()
		if i < len(bs.Statements)-1 {
			out += ";"
		}
	}
	out += "}"
	return out
}

type IfStatement struct {
	Token      token.Token // the IF token
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement // can be nil
}

func (ifs *IfStatement) String() string {
	out := "if (" + ifs.Condition.String() + ") " + ifs.ThenBranch.String()
	if ifs.ElseBranch != nil {
		out += " else " + ifs.ElseBranch.String()
	}
	return out
}

type WhileStatement struct {
	Token     token.Token // the WHILE token
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) String() string {
	return "while (" + ws.Condition.String() + ") " + ws.Body.String()
}

type ForStatement struct {
	Token     token.Token // the FOR token
	Init      Statement   // can be nil
	Condition Expression  // can be nil
	Update    Expression  // can be nil
	Body      Statement
}

func (fs *ForStatement) String() string {
	out := "for ("
	if fs.Init != nil {
		out += fs.Init.String()
	}
	out += "; "
	if fs.Condition != nil {
		out += fs.Condition.String()
	}
	out += "; "
	if fs.Update != nil {
		out += fs.Update.String()
	}
	out += ") " + fs.Body.String()
	return out
}

// Expressions
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type IntegerLiteral struct {
	Token token.Token // the INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token // the FLOAT token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string  { return fl.Token.Literal }

type StringLiteral struct {
	Token token.Token // the STRING token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return "\"" + sl.Value + "\"" }

type BooleanLiteral struct {
	Token token.Token // the TRUE or FALSE token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string  { return bl.Token.Literal }

type NullLiteral struct {
	Token token.Token // the NULL token
}

func (nl *NullLiteral) expressionNode() {}
func (nl *NullLiteral) String() string  { return "null" }

type BinaryExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return "(" + be.Left.String() + " " + be.Operator + " " + be.Right.String() + ")"
}

type UnaryExpression struct {
	Token    token.Token // the operator token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	return "(" + ue.Operator + ue.Right.String() + ")"
}

type CallExpression struct {
	Token     token.Token // the ( token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	out := ce.Function.String() + "("
	for i, arg := range ce.Arguments {
		if i > 0 {
			out += ", "
		}
		out += arg.String()
	}
	out += ")"
	return out
}

type MemberExpression struct {
	Token    token.Token // the . or [ token
	Object   Expression
	Property Expression
	Computed bool // true for obj[prop], false for obj.prop
}

func (me *MemberExpression) expressionNode() {}
func (me *MemberExpression) String() string {
	if me.Computed {
		return me.Object.String() + "[" + me.Property.String() + "]"
	}
	return me.Object.String() + "." + me.Property.String()
}

type AssignmentExpression struct {
	Token token.Token // the = token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) String() string {
	return ae.Left.String() + " = " + ae.Value.String()
}

type FunctionExpression struct {
	Token      token.Token // the FUNCTION token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fe *FunctionExpression) expressionNode() {}
func (fe *FunctionExpression) String() string {
	out := "function("
	for i, param := range fe.Parameters {
		if i > 0 {
			out += ", "
		}
		out += param.String()
	}
	out += ") " + fe.Body.String()
	return out
}

type ArrayLiteral struct {
	Token    token.Token // the [ token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) String() string {
	out := "["
	for i, elem := range al.Elements {
		if i > 0 {
			out += ", "
		}
		out += elem.String()
	}
	out += "]"
	return out
}

type ObjectLiteral struct {
	Token      token.Token // the { token
	Properties map[Expression]Expression
}

func (ol *ObjectLiteral) expressionNode() {}
func (ol *ObjectLiteral) String() string {
	out := "{"
	i := 0
	for key, value := range ol.Properties {
		if i > 0 {
			out += ", "
		}
		out += key.String() + ": " + value.String()
		i++
	}
	out += "}"
	return out
}
