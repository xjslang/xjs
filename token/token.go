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
	// logical operators
	NOT    // !
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	LTE    // <=
	GT     // >
	GTE    // >=
	// delimiters
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	NEWLINE   // \r, \n, \r\n
	// comments
	LINE_COMMENT  // //
	BLOCK_COMMENT // /* ... */
	// keywords
	LET
	FUNCTION
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
	// logical operators
	NOT:    "!",
	EQ:     "==",
	NOT_EQ: "!=",
	LT:     "<",
	LTE:    "<=",
	GT:     ">",
	GTE:    ">=",
	// delimiters
	COMMA:     ",",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	NEWLINE:   "new line",
	// comments
	LINE_COMMENT:  "line comment",
	BLOCK_COMMENT: "block comment",
	// keywords
	LET:      "let",
	FUNCTION: "function",
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

var binOperators = map[Type]int{
	PLUS:     1,
	MINUS:    1,
	MULTIPLY: 2,
	DIVIDE:   2,
	MODULO:   2,
	LPAREN:   3,
}

func (typ Type) IsBinaryOperator() (ok bool) {
	registerMu.RLock()
	defer registerMu.RUnlock()
	_, ok = binOperators[typ]
	return
}

func (typ Type) Precedence() int {
	registerMu.RLock()
	defer registerMu.RUnlock()
	return binOperators[typ]
}

func (typ Type) Register(precedence int) {
	registerMu.Lock()
	defer registerMu.Unlock()
	binOperators[typ] = precedence
}

func RegisterBinaryOperator(lit string, precedence int) Type {
	typ := RegisterType(lit)
	typ.Register(precedence)
	return typ
}
