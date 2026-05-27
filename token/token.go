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

var infixOperators = map[Type]int{
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
	// specials operators
	LPAREN: 7,
}

func (typ Type) IsInfixOperator() (ok bool) {
	registerMu.RLock()
	defer registerMu.RUnlock()
	_, ok = infixOperators[typ]
	return
}

func (typ Type) Precedence() int {
	registerMu.RLock()
	defer registerMu.RUnlock()
	return infixOperators[typ]
}

func (typ Type) RegisterInfixOperator(precedence int) {
	registerMu.Lock()
	defer registerMu.Unlock()
	infixOperators[typ] = precedence
}

func RegisterInfixOperator(lit string, precedence int) Type {
	typ := RegisterType(lit)
	typ.RegisterInfixOperator(precedence)
	return typ
}

var prefixOperators = map[Type]bool{
	NOT: true,
}

func (typ Type) IsPrefixOperator() (ok bool) {
	registerMu.RLock()
	defer registerMu.RUnlock()
	_, ok = prefixOperators[typ]
	return
}

func (typ Type) RegisterPrefixOperator() {
	registerMu.Lock()
	defer registerMu.Unlock()
	prefixOperators[typ] = true
}

func RegisterPrefixOperator(lit string) Type {
	typ := RegisterType(lit)
	typ.RegisterPrefixOperator()
	return typ
}
