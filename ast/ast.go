// Package ast defines the Abstract Syntax Tree nodes for the xjslang language.
// It provides interfaces and concrete types representing different language constructs.
package ast

import (
	"github.com/xjslang/xjs/token"
)

// Operator precedence levels (must match parser precedences)
const (
	_ int = iota
	PrecedenceLowest
	PrecedenceAssignment // =
	PrecedenceLogicalOr  // ||
	PrecedenceLogicalAnd // &&
	PrecedenceEquality   // == !=
	PrecedenceComparison // > < >= <=
	PrecedenceSum        // +
	PrecedenceProduct    // * / %
	PrecedenceUnary      // -x !x ++x --x
	PrecedencePostfix    // x++ x--
	PrecedenceCall       // myFunction(X)
	PrecedenceMember     // obj.prop obj[prop]
	PrecedenceAtomic     // literals, identifiers, grouped expressions
)

func operatorPrecedence(tokenType token.Type) int {
	switch tokenType {
	case token.ASSIGN, token.PLUS_ASSIGN, token.MINUS_ASSIGN:
		return PrecedenceAssignment
	case token.OR:
		return PrecedenceLogicalOr
	case token.AND:
		return PrecedenceLogicalAnd
	case token.EQ, token.NOT_EQ:
		return PrecedenceEquality
	case token.LT, token.GT, token.LTE, token.GTE:
		return PrecedenceComparison
	case token.PLUS, token.MINUS:
		return PrecedenceSum
	case token.MULTIPLY, token.DIVIDE, token.MODULO:
		return PrecedenceProduct
	case token.INCREMENT, token.DECREMENT:
		return PrecedencePostfix
	case token.LPAREN:
		return PrecedenceCall
	case token.DOT, token.LBRACKET:
		return PrecedenceMember
	default:
		return PrecedenceLowest
	}
}

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
	// Precedence returns the operator precedence level for this expression.
	// Higher values indicate higher precedence.
	Precedence() int
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

// BlankLine represents one or more consecutive blank lines.
// Multiple blank lines are condensed to a single blank line in output.
type BlankLine struct {
	Token token.Token
}

func (bl *BlankLine) WriteTo(cw *CodeWriter) {
	if !cw.PrettyPrint {
		return
	}
	cw.WriteRune('\n')
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
	cw.WriteSemi()
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
	cw.WriteSemi()
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
	cw.WriteSemi()
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
	cw.WriteNewline()
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
		cw.WriteString(" else ")
		ifs.ElseBranch.WriteTo(cw)
	}
	cw.WriteNewline()
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
	Init      Expression  // can be nil, typically LetExpression or AssignmentExpression
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
	}
	cw.WriteRune(';')
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

type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) WriteTo(cw *CodeWriter) {
	cw.AddNamedMapping(i.Token.Start.Line, i.Token.Start.Column, i.Value)
	cw.WriteString(i.Value)
}

func (i *Identifier) Precedence() int {
	return PrecedenceAtomic
}

type IntegerLiteral struct {
	Token token.Token // the INT token
}

func (il *IntegerLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(il.Token.Start)
	cw.WriteString(il.Token.Literal)
}

func (il *IntegerLiteral) Precedence() int {
	return PrecedenceAtomic
}

type FloatLiteral struct {
	Token token.Token // the FLOAT token
}

func (fl *FloatLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(fl.Token.Start)
	cw.WriteString(fl.Token.Literal)
}

func (fl *FloatLiteral) Precedence() int {
	return PrecedenceAtomic
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

func (sl *StringLiteral) Precedence() int {
	return PrecedenceAtomic
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

func (sl *MultiStringLiteral) Precedence() int {
	return PrecedenceAtomic
}

type BooleanLiteral struct {
	Token token.Token // the TRUE or FALSE token
	Value bool
}

func (bl *BooleanLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(bl.Token.Start)
	cw.WriteString(bl.Token.Literal)
}

func (bl *BooleanLiteral) Precedence() int {
	return PrecedenceAtomic
}

type NullLiteral struct {
	Token token.Token // the NULL token
}

func (nl *NullLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(nl.Token.Start)
	cw.WriteString("null")
}

func (nl *NullLiteral) Precedence() int {
	return PrecedenceAtomic
}

// LetExpression represents a let declaration as an expression (used in for loops).
// Unlike LetStatement, this does not add semicolon or newline when writing.
type LetExpression struct {
	Token token.Token // the LET token
	Name  *Identifier
	Value Expression
}

func (le *LetExpression) WriteTo(cw *CodeWriter) {
	cw.AddMapping(le.Token.Start)
	cw.WriteString("let ")
	le.Name.WriteTo(cw)
	if le.Value != nil {
		cw.WriteSpace()
		cw.WriteRune('=')
		cw.WriteSpace()
		le.Value.WriteTo(cw)
	}
}

func (le *LetExpression) Precedence() int {
	return PrecedenceAssignment
}

type BinaryExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) WriteTo(cw *CodeWriter) {
	myPrecedence := be.Precedence()

	// Left side needs parens if its precedence is lower than ours
	leftNeedsParens := be.Left.Precedence() < myPrecedence
	if leftNeedsParens {
		cw.WriteRune('(')
	}
	be.Left.WriteTo(cw)
	if leftNeedsParens {
		cw.WriteRune(')')
	}

	cw.WriteSpace()
	cw.AddMapping(be.Token.Start)
	cw.WriteString(be.Operator)
	cw.WriteSpace()

	// Right side needs parens if its precedence is lower OR equal (for left-associativity)
	// For example: 1-2-3 should be ((1-2)-3) not (1-(2-3))
	rightNeedsParens := be.Right.Precedence() <= myPrecedence
	if rightNeedsParens {
		cw.WriteRune('(')
	}
	be.Right.WriteTo(cw)
	if rightNeedsParens {
		cw.WriteRune(')')
	}
}

func (be *BinaryExpression) Precedence() int {
	return operatorPrecedence(be.Token.Type)
}

type UnaryExpression struct {
	Token    token.Token // the operator token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ue.Token.Start)
	cw.WriteString(ue.Operator)
	// Right side needs parens if its precedence is lower than unary
	if ue.Right.Precedence() < PrecedenceUnary {
		cw.WriteRune('(')
		ue.Right.WriteTo(cw)
		cw.WriteRune(')')
	} else {
		ue.Right.WriteTo(cw)
	}
}

func (ue *UnaryExpression) Precedence() int {
	return PrecedenceUnary
}

type PostfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) WriteTo(cw *CodeWriter) {
	// Left side needs parens if its precedence is lower than postfix
	if pe.Left.Precedence() < PrecedencePostfix {
		cw.WriteRune('(')
		pe.Left.WriteTo(cw)
		cw.WriteRune(')')
	} else {
		pe.Left.WriteTo(cw)
	}
	cw.AddMapping(pe.Token.Start)
	cw.WriteString(pe.Operator)
}

func (pe *PostfixExpression) Precedence() int {
	return PrecedencePostfix
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

func (ge *GroupedExpression) Precedence() int {
	// Grouped expressions have atomic precedence because parens are explicit
	return PrecedenceAtomic
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

func (ce *CallExpression) Precedence() int {
	return PrecedenceCall
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

func (me *MemberExpression) Precedence() int {
	return PrecedenceMember
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

func (ae *AssignmentExpression) Precedence() int {
	return PrecedenceAssignment
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

func (cae *CompoundAssignmentExpression) Precedence() int {
	return PrecedenceAssignment
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

func (fe *FunctionExpression) Precedence() int {
	return PrecedenceAtomic
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

func (al *ArrayLiteral) Precedence() int {
	return PrecedenceAtomic
}

type ObjectProperty struct {
	Key   Expression
	Value Expression
}

type ObjectLiteral struct {
	Token      token.Token // the { token
	Properties []ObjectProperty
}

func (ol *ObjectLiteral) WriteTo(cw *CodeWriter) {
	cw.AddMapping(ol.Token.Start)
	cw.WriteRune('{')

	// Write properties in order they were defined
	for i, prop := range ol.Properties {
		if i > 0 {
			cw.WriteRune(',')
			cw.WriteSpace()
		}
		prop.Key.WriteTo(cw)
		cw.WriteRune(':')
		cw.WriteSpace()
		prop.Value.WriteTo(cw)
	}

	cw.WriteRune('}')
}

func (ol *ObjectLiteral) Precedence() int {
	return PrecedenceAtomic
}
