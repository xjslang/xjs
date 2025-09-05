// Package ast defines the Abstract Syntax Tree nodes for the xjslang language.
// It provides interfaces and concrete types representing different language constructs.
package ast

import (
	"sort"
	"strings"

	"github.com/xjslang/xjs/token"
)

type Node interface {
	WriteTo(b *strings.Builder)
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var b strings.Builder
	p.WriteTo(&b)
	return b.String()
}

func (p *Program) WriteTo(b *strings.Builder) {
	for i, stmt := range p.Statements {
		if i > 0 {
			b.WriteRune(';')
		}
		stmt.WriteTo(b)
	}
}

// Statements
type LetStatement struct {
	Token token.Token // the LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) WriteTo(b *strings.Builder) {
	b.WriteString("let ")
	ls.Name.WriteTo(b)
	if ls.Value != nil {
		b.WriteRune('=')
		ls.Value.WriteTo(b)
	}
}

type ReturnStatement struct {
	Token       token.Token // the RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) WriteTo(b *strings.Builder) {
	b.WriteString("return")
	if rs.ReturnValue != nil {
		b.WriteRune(' ')
		rs.ReturnValue.WriteTo(b)
	}
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) WriteTo(b *strings.Builder) {
	if es.Expression != nil {
		es.Expression.WriteTo(b)
	}
}

type FunctionDeclaration struct {
	Token      token.Token // the FUNCTION token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) WriteTo(b *strings.Builder) {
	b.WriteString("function ")
	fd.Name.WriteTo(b)
	b.WriteRune('(')
	for i, param := range fd.Parameters {
		if i > 0 {
			b.WriteRune(',')
		}
		param.WriteTo(b)
	}
	b.WriteRune(')')
	fd.Body.WriteTo(b)
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) WriteTo(b *strings.Builder) {
	b.WriteRune('{')
	for i, stmt := range bs.Statements {
		if i > 0 {
			b.WriteRune(';')
		}
		stmt.WriteTo(b)
	}
	b.WriteRune('}')
}

type IfStatement struct {
	Token      token.Token // the IF token
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement // can be nil
}

func (ifs *IfStatement) WriteTo(b *strings.Builder) {
	b.WriteString("if (")
	ifs.Condition.WriteTo(b)
	b.WriteRune(')')
	ifs.ThenBranch.WriteTo(b)
	if ifs.ElseBranch != nil {
		b.WriteString(" else ")
		ifs.ElseBranch.WriteTo(b)
	}
}

type WhileStatement struct {
	Token     token.Token // the WHILE token
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) WriteTo(b *strings.Builder) {
	b.WriteString("while (")
	ws.Condition.WriteTo(b)
	b.WriteRune(')')
	ws.Body.WriteTo(b)
}

type ForStatement struct {
	Token     token.Token // the FOR token
	Init      Statement   // can be nil
	Condition Expression  // can be nil
	Update    Expression  // can be nil
	Body      Statement
}

func (fs *ForStatement) WriteTo(b *strings.Builder) {
	b.WriteString("for (")
	if fs.Init != nil {
		fs.Init.WriteTo(b)
	}
	b.WriteRune(';')
	if fs.Condition != nil {
		fs.Condition.WriteTo(b)
	}
	b.WriteRune(';')
	if fs.Update != nil {
		fs.Update.WriteTo(b)
	}
	b.WriteRune(')')
	fs.Body.WriteTo(b)
}

// Expressions
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) WriteTo(b *strings.Builder) {
	b.WriteString(i.Value)
}

type IntegerLiteral struct {
	Token token.Token // the INT token
}

func (il *IntegerLiteral) WriteTo(b *strings.Builder) {
	b.WriteString(il.Token.Literal)
}

type FloatLiteral struct {
	Token token.Token // the FLOAT token
}

func (fl *FloatLiteral) WriteTo(b *strings.Builder) {
	b.WriteString(fl.Token.Literal)
}

type StringLiteral struct {
	Token token.Token // the STRING token
	Value string
}

func (sl *StringLiteral) WriteTo(b *strings.Builder) {
	b.WriteRune('"')
	b.WriteString(sl.Value)
	b.WriteRune('"')
}

type MultiStringLiteral struct {
	Token token.Token // the MULTI_STRING token
	Value string
}

func (sl *MultiStringLiteral) WriteTo(b *strings.Builder) {
	b.WriteRune('`')
	b.WriteString(sl.Value)
	b.WriteRune('`')
}

type BooleanLiteral struct {
	Token token.Token // the TRUE or FALSE token
	Value bool
}

func (bl *BooleanLiteral) WriteTo(b *strings.Builder) {
	b.WriteString(bl.Token.Literal)
}

type NullLiteral struct {
	Token token.Token // the NULL token
}

func (nl *NullLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("null")
}

type BinaryExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) WriteTo(b *strings.Builder) {
	b.WriteRune('(')
	be.Left.WriteTo(b)
	switch be.Operator {
	case "==":
		b.WriteString("===")
	case "!=":
		b.WriteString("!==")
	default:
		b.WriteString(be.Operator)
	}
	be.Right.WriteTo(b)
	b.WriteRune(')')
}

type UnaryExpression struct {
	Token    token.Token // the operator token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) WriteTo(b *strings.Builder) {
	b.WriteRune('(')
	b.WriteString(ue.Operator)
	ue.Right.WriteTo(b)
	b.WriteRune(')')
}

type PostfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) WriteTo(b *strings.Builder) {
	b.WriteRune('(')
	pe.Left.WriteTo(b)
	b.WriteString(pe.Operator)
	b.WriteRune(')')
}

type CallExpression struct {
	Token     token.Token // the ( token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) WriteTo(b *strings.Builder) {
	ce.Function.WriteTo(b)
	b.WriteRune('(')
	for i, arg := range ce.Arguments {
		if i > 0 {
			b.WriteRune(',')
		}
		arg.WriteTo(b)
	}
	b.WriteRune(')')
}

type MemberExpression struct {
	Token    token.Token // the . or [ token
	Object   Expression
	Property Expression
	Computed bool // true for obj[prop], false for obj.prop
}

func (me *MemberExpression) WriteTo(b *strings.Builder) {
	if me.Computed {
		me.Object.WriteTo(b)
		b.WriteRune('[')
		me.Property.WriteTo(b)
		b.WriteRune(']')
	} else {
		me.Object.WriteTo(b)
		b.WriteRune('.')
		me.Property.WriteTo(b)
	}
}

type AssignmentExpression struct {
	Token token.Token // the = token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) WriteTo(b *strings.Builder) {
	ae.Left.WriteTo(b)
	b.WriteRune('=')
	ae.Value.WriteTo(b)
}

type CompoundAssignmentExpression struct {
	Token    token.Token // the += or -= token
	Left     Expression
	Operator string // "+" or "-"
	Value    Expression
}

func (cae *CompoundAssignmentExpression) WriteTo(b *strings.Builder) {
	cae.Left.WriteTo(b)
	b.WriteRune(' ')
	b.WriteString(cae.Operator)
	b.WriteRune('=')
	cae.Value.WriteTo(b)
}

type FunctionExpression struct {
	Token      token.Token // the FUNCTION token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fe *FunctionExpression) WriteTo(b *strings.Builder) {
	b.WriteString("function(")
	for i, param := range fe.Parameters {
		if i > 0 {
			b.WriteRune(',')
		}
		param.WriteTo(b)
	}
	b.WriteRune(')')
	fe.Body.WriteTo(b)
}

type ArrayLiteral struct {
	Token    token.Token // the [ token
	Elements []Expression
}

func (al *ArrayLiteral) WriteTo(b *strings.Builder) {
	b.WriteRune('[')
	for i, elem := range al.Elements {
		if i > 0 {
			b.WriteRune(',')
		}
		elem.WriteTo(b)
	}
	b.WriteRune(']')
}

type ObjectLiteral struct {
	Token      token.Token // the { token
	Properties map[Expression]Expression
}

func (ol *ObjectLiteral) WriteTo(b *strings.Builder) {
	b.WriteRune('{')

	// Extract keys and sort them for deterministic output
	keys := make([]Expression, 0, len(ol.Properties))
	for key := range ol.Properties {
		keys = append(keys, key)
	}

	// Sort keys by their string representation
	sort.Slice(keys, func(i, j int) bool {
		var keyI, keyJ strings.Builder
		keys[i].WriteTo(&keyI)
		keys[j].WriteTo(&keyJ)
		return keyI.String() < keyJ.String()
	})

	// Write properties in sorted order
	for i, key := range keys {
		if i > 0 {
			b.WriteRune(',')
		}
		key.WriteTo(b)
		b.WriteRune(':')
		ol.Properties[key].WriteTo(b)
	}

	b.WriteRune('}')
}
