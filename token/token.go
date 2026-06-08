package token

import (
	"strconv"
	"sync"
)

type Type int

func (tt Type) String() string {
	registerMu.RLock()
	defer registerMu.RUnlock()
	lit, ok := tokenLiterals[tt]
	if !ok {
		return "unknown(" + strconv.Itoa(int(tt)) + ")"
	}
	return lit
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Token struct {
	Position
	Type          Type
	Literal       string
	LeadingTrivia []Token
	AfterNewline  bool
}

const (
	// special keywords
	EOF Type = iota
	IDENT
	ILLEGAL
	UNKNOWN
	// literals
	STRING
	NUMBER
	BOOLEAN
	// operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %
	// incremental operators
	INCREMENT // ++
	DECREMENT // --
	// comparison operators
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	LTE    // <=
	GT     // >
	GTE    // >=
	// Logical operators
	AND // &&
	OR  // ||
	NOT // !
	// delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	NEWLINE   // \r, \n, \r\n
	// comments
	LINE_COMMENT  // //
	BLOCK_COMMENT // /* ... */
)

var tokenLiterals = map[Type]string{
	// special keywords
	EOF:     "end of file",
	IDENT:   "identifier",
	ILLEGAL: "illegal",
	UNKNOWN: "unknown",
	// literals
	STRING:  "string",
	NUMBER:  "number",
	BOOLEAN: "boolean",
	// operators
	ASSIGN:   "=",
	PLUS:     "+",
	MINUS:    "-",
	MULTIPLY: "*",
	DIVIDE:   "/",
	MODULO:   "%",
	// incremental operators
	INCREMENT: "++",
	DECREMENT: "--",
	// comparison operators
	EQ:     "==",
	NOT_EQ: "!=",
	LT:     "<",
	LTE:    "<=",
	GT:     ">",
	GTE:    ">=",
	// logical operators
	AND: "&&",
	OR:  "||",
	NOT: "!",
	// delimiters
	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",
	DOT:       ".",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	NEWLINE:   "new line",
	// comments
	LINE_COMMENT:  "line comment",
	BLOCK_COMMENT: "block comment",
}

const initCustomType Type = 1000

var (
	nextType   Type = initCustomType
	registerMu sync.RWMutex
)

func RegisterType(lit string) Type {
	registerMu.Lock()
	defer registerMu.Unlock()
	typ := nextType
	tokenLiterals[typ] = lit
	nextType++
	return typ
}

var binaryOps = map[Type]int{
	// || (lowest precedence for operators)
	OR: 1,
	// &&
	AND: 2,
	// == !=
	EQ:     3,
	NOT_EQ: 3,
	// < <= > >=
	LT:  4,
	LTE: 4,
	GT:  4,
	GTE: 4,
	// + -
	PLUS:  5,
	MINUS: 5,
	// * / %
	MULTIPLY: 6,
	DIVIDE:   6,
	MODULO:   6,
	// ( [ .
	LPAREN:   7,
	LBRACKET: 7,
	DOT:      7,
}

func (typ Type) IsBinaryOp() (ok bool) {
	registerMu.RLock()
	defer registerMu.RUnlock()
	_, ok = binaryOps[typ]
	return
}

func (typ Type) Precedence() int {
	registerMu.RLock()
	defer registerMu.RUnlock()
	return binaryOps[typ]
}

// RegisterBinaryType registers a token type as a "binary operator".
func RegisterBinaryType(typ Type, precedence int) {
	registerMu.Lock()
	defer registerMu.Unlock()
	binaryOps[typ] = precedence
}

// RegisterBinaryOp registers a "binary operator".
func RegisterBinaryOp(lit string, precedence int) Type {
	typ := RegisterType(lit)
	RegisterBinaryType(typ, precedence)
	return typ
}

var unaryTypes = map[Type]bool{
	NOT:      true,
	PLUS:     true,
	MINUS:    true,
	LPAREN:   true,
	LBRACE:   true,
	LBRACKET: true,
}

func (typ Type) IsUnaryOp() (ok bool) {
	registerMu.RLock()
	defer registerMu.RUnlock()
	_, ok = unaryTypes[typ]
	return
}

// RegisterUnaryType registers a token type as a "unary operator".
func RegisterUnaryType(typ Type) {
	registerMu.Lock()
	defer registerMu.Unlock()
	unaryTypes[typ] = true
}

// RegisterUnaryOp registers a "unary operator".
func RegisterUnaryOp(lit string) Type {
	typ := RegisterType(lit)
	RegisterUnaryType(typ)
	return typ
}
