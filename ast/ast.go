// Package ast defines the Abstract Syntax Tree nodes for the xjslang language.
// It provides interfaces and concrete types representing different language constructs.
package ast

import (
	"sort"

	"github.com/xjslang/xjs/token"
)

type CompileOptions struct {
	GenerateSourceMap bool
}

type Node interface {
	WriteTo(cw *CodeWriter)
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

func (p *Program) WriteTo(cw *CodeWriter) {
	for _, stmt := range p.Statements {
		stmt.WriteTo(cw)
	}
}

// CommentBlock represents one or more consecutive line comments.
// It is treated as a statement in the AST.
type CommentBlock struct {
	Comments []token.Token // each token is a COMMENT token
}

func (cb *CommentBlock) WriteTo(cw *CodeWriter) {
	if !cw.PrettyPrint {
		return
	}
	for _, comment := range cb.Comments {
		cw.AddMapping(comment.Start)
		cw.WriteString("//")
		cw.WriteString(comment.Literal)
		cw.WriteNewline()
	}
}

// Statements
type LetStatement struct {
	Token token.Token // the LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ls.Token.Start)
	cw.WriteString("let ")
	ls.Name.WriteTo(cw)
	if ls.Value != nil {
		cw.WriteSpace()
		cw.WriteRune('=')
		cw.WriteSpace()
		ls.Value.WriteTo(cw)
	}
	cw.WriteRune(';')
	cw.WriteNewline()
}

type ReturnStatement struct {
	Token       token.Token // the RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) WriteTo(cw *CodeWriter) {
	cw.AddMapping(rs.Token.Start)
	cw.WriteString("return")
	if rs.ReturnValue != nil {
		cw.WriteRune(' ')
		rs.ReturnValue.WriteTo(cw)
	}
	cw.WriteRune(';')
	cw.WriteNewline()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) WriteTo(cw *CodeWriter) {
	if es.Expression == nil {
		return
	}
	es.Expression.WriteTo(cw)
	cw.WriteRune(';')
	cw.WriteNewline()
}

type FunctionDeclaration struct {
	Token      token.Token // the FUNCTION token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) WriteTo(cw *CodeWriter) {
	cw.AddMapping(fd.Token.Start)
	cw.WriteString("function ")
	fd.Name.WriteTo(cw)
	cw.WriteRune('(')
	for i, param := range fd.Parameters {
		if i > 0 {
			cw.WriteRune(',')
			cw.WriteSpace()
		}
		param.WriteTo(cw)
	}
	cw.WriteRune(')')
	cw.WriteSpace()
	fd.Body.WriteTo(cw)
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) WriteTo(cw *CodeWriter) {
	cw.AddMapping(bs.Token.Start)
	cw.WriteRune('{')
	cw.IncreaseIndent()
	cw.WriteNewline()
	for _, stmt := range bs.Statements {
		cw.WriteIndent()
		stmt.WriteTo(cw)
	}
	cw.DecreaseIndent()
	cw.WriteIndent()
	cw.WriteRune('}')
	cw.WriteNewline()
}

type IfStatement struct {
	Token      token.Token // the IF token
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement // can be nil
}

func (ifs *IfStatement) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ifs.Token.Start)
	cw.WriteString("if")
	cw.WriteSpace()
	cw.WriteRune('(')
	ifs.Condition.WriteTo(cw)
	cw.WriteRune(')')
	cw.WriteSpace()
	ifs.ThenBranch.WriteTo(cw)
	if ifs.ElseBranch != nil {
		// TODO: do we need to use `PrettyPrint` here?
		if !cw.PrettyPrint {
			cw.WriteRune(' ')
		} else {
			cw.WriteSpace()
		}
		cw.WriteString("else")
		if !cw.PrettyPrint {
			cw.WriteRune(' ')
		} else {
			cw.WriteSpace()
		}
		ifs.ElseBranch.WriteTo(cw)
	}
}

type WhileStatement struct {
	Token     token.Token // the WHILE token
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ws.Token.Start)
	cw.WriteString("while")
	cw.WriteSpace()
	cw.WriteRune('(')
	ws.Condition.WriteTo(cw)
	cw.WriteRune(')')
	cw.WriteSpace()
	ws.Body.WriteTo(cw)
}

type ForStatement struct {
	Token     token.Token // the FOR token
	Init      Statement   // can be nil
	Condition Expression  // can be nil
	Update    Expression  // can be nil
	Body      Statement
}

func (fs *ForStatement) WriteTo(cw *CodeWriter) {
	cw.AddMapping(fs.Token.Start)
	cw.WriteString("for")
	cw.WriteSpace()
	cw.WriteRune('(')
	if fs.Init != nil {
		fs.Init.WriteTo(cw)
	} else {
		cw.WriteRune(';')
	}
	cw.WriteSpace()
	if fs.Condition != nil {
		fs.Condition.WriteTo(cw)
	}
	cw.WriteRune(';')
	cw.WriteSpace()
	if fs.Update != nil {
		fs.Update.WriteTo(cw)
	}
	cw.WriteRune(')')
	cw.WriteSpace()
	fs.Body.WriteTo(cw)
}

// Expressions
type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) WriteTo(cw *CodeWriter) {
	cw.AddNamedMapping(i.Token.Start.Line, i.Token.Start.Column, i.Value)
	cw.WriteString(i.Value)
}

type IntegerLiteral struct {
	Token token.Token // the INT token
}

func (il *IntegerLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(il.Token.Start)
	cw.WriteString(il.Token.Literal)
}

type FloatLiteral struct {
	Token token.Token // the FLOAT token
}

func (fl *FloatLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(fl.Token.Start)
	cw.WriteString(fl.Token.Literal)
}

type StringLiteral struct {
	Token token.Token // the STRING token
	Value string
}

func (sl *StringLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(sl.Token.Start)
	cw.WriteRune('"')
	cw.WriteString(sl.Value)
	cw.WriteRune('"')
}

type MultiStringLiteral struct {
	Token token.Token // the MULTI_STRING token
	Value string
}

func (sl *MultiStringLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(sl.Token.Start)
	cw.WriteRune('`')
	cw.WriteString(sl.Value)
	cw.WriteRune('`')
}

type BooleanLiteral struct {
	Token token.Token // the TRUE or FALSE token
	Value bool
}

func (bl *BooleanLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(bl.Token.Start)
	cw.WriteString(bl.Token.Literal)
}

type NullLiteral struct {
	Token token.Token // the NULL token
}

func (nl *NullLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(nl.Token.Start)
	cw.WriteString("null")
}

type BinaryExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) WriteTo(cw *CodeWriter) {
	cw.WriteRune('(')
	be.Left.WriteTo(cw)
	cw.WriteSpace()
	cw.AddMapping(be.Token.Start)
	cw.WriteString(be.Operator)
	cw.WriteSpace()
	be.Right.WriteTo(cw)
	cw.WriteRune(')')
}

type UnaryExpression struct {
	Token    token.Token // the operator token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ue.Token.Start)
	cw.WriteRune('(')
	cw.WriteString(ue.Operator)
	ue.Right.WriteTo(cw)
	cw.WriteRune(')')
}

type PostfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) WriteTo(cw *CodeWriter) {
	cw.WriteRune('(')
	pe.Left.WriteTo(cw)
	cw.AddMapping(pe.Token.Start)
	cw.WriteString(pe.Operator)
	cw.WriteRune(')')
}

type GroupedExpression struct {
	Token      token.Token // the opening ( token
	Expression Expression
}

func (ge *GroupedExpression) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ge.Token.Start)
	cw.WriteRune('(')
	ge.Expression.WriteTo(cw)
	cw.WriteRune(')')
}

type CallExpression struct {
	Token     token.Token // the ( token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) WriteTo(cw *CodeWriter) {
	ce.Function.WriteTo(cw)
	cw.AddMapping(ce.Token.Start)
	cw.WriteRune('(')
	for i, arg := range ce.Arguments {
		if i > 0 {
			cw.WriteRune(',')
			cw.WriteSpace()
		}
		arg.WriteTo(cw)
	}
	cw.WriteRune(')')
}

type MemberExpression struct {
	Token    token.Token // the . or [ token
	Object   Expression
	Property Expression
	Computed bool // true for obj[prop], false for obj.prop
}

func (me *MemberExpression) WriteTo(cw *CodeWriter) {
	me.Object.WriteTo(cw)
	if me.Computed {
		cw.AddMapping(me.Token.Start)
		cw.WriteRune('[')
		me.Property.WriteTo(cw)
		cw.WriteRune(']')
	} else {
		cw.AddMapping(me.Token.Start)
		cw.WriteRune('.')
		me.Property.WriteTo(cw)
	}
}

type AssignmentExpression struct {
	Token token.Token // the = token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) WriteTo(cw *CodeWriter) {
	ae.Left.WriteTo(cw)
	cw.WriteSpace()
	cw.AddMapping(ae.Token.Start)
	cw.WriteRune('=')
	cw.WriteSpace()
	ae.Value.WriteTo(cw)
}

type CompoundAssignmentExpression struct {
	Token    token.Token // the += or -= token
	Left     Expression
	Operator string // "+" or "-"
	Value    Expression
}

func (cae *CompoundAssignmentExpression) WriteTo(cw *CodeWriter) {
	cae.Left.WriteTo(cw)
	cw.AddMapping(cae.Token.Start)
	cw.WriteRune(' ')
	cw.WriteString(cae.Operator)
	cw.WriteRune('=')
	cae.Value.WriteTo(cw)
}

type FunctionExpression struct {
	Token      token.Token // the FUNCTION token
	Name       *Identifier // optional name
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fe *FunctionExpression) WriteTo(cw *CodeWriter) {
	cw.AddMapping(fe.Token.Start)
	cw.WriteString("function")
	if fe.Name != nil {
		cw.WriteRune(' ')
		fe.Name.WriteTo(cw)
	}
	cw.WriteRune('(')
	for i, param := range fe.Parameters {
		if i > 0 {
			cw.WriteRune(',')
			cw.WriteSpace()
		}
		param.WriteTo(cw)
	}
	cw.WriteRune(')')
	cw.WriteSpace()
	fe.Body.WriteTo(cw)
}

type ArrayLiteral struct {
	Token    token.Token // the [ token
	Elements []Expression
}

func (al *ArrayLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(al.Token.Start)
	cw.WriteRune('[')
	for i, elem := range al.Elements {
		if i > 0 {
			cw.WriteRune(',')
			cw.WriteSpace()
		}
		elem.WriteTo(cw)
	}
	cw.WriteRune(']')
}

type ObjectLiteral struct {
	Token      token.Token // the { token
	Properties map[Expression]Expression
}

func (ol *ObjectLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ol.Token.Start)
	cw.WriteRune('{')

	// Extract keys and sort them for deterministic output
	keys := make([]Expression, 0, len(ol.Properties))
	for key := range ol.Properties {
		keys = append(keys, key)
	}

	// Sort keys by their string representation
	sort.Slice(keys, func(i, j int) bool {
		var keyI, keyJ CodeWriter
		keys[i].WriteTo(&keyI)
		keys[j].WriteTo(&keyJ)
		return keyI.String() < keyJ.String()
	})

	// Write properties in sorted order
	for i, key := range keys {
		if i > 0 {
			cw.WriteRune(',')
			cw.WriteSpace()
		}
		key.WriteTo(cw)
		cw.WriteRune(':')
		cw.WriteSpace()
		ol.Properties[key].WriteTo(cw)
	}

	cw.WriteRune('}')
}
