package token

import (
	"strconv"
	"sync"
)

type ForkableScanner interface {
	Fork() Scanner
	Apply(Scanner)
}

type Scanner interface {
	NextToken() Token
}

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
	QUOTE     // ' "
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
	// others
	DIGIT // 0..9
)

var tokenLiterals = map[Type]string{
	// special keywords
	EOF:     "end of file",
	IDENT:   "identifier",
	ILLEGAL: "illegal",
	UNKNOWN: "unknown",
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
	QUOTE:     "quote",
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
	// others
	DIGIT: "digit",
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
	// =
	ASSIGN: 1,
	// ||
	OR: 2,
	// &&
	AND: 3,
	// == !=
	EQ:     4,
	NOT_EQ: 4,
	// < <= > >=
	LT:  5,
	LTE: 5,
	GT:  5,
	GTE: 5,
	// + -
	PLUS:  6,
	MINUS: 6,
	// * / %
	MULTIPLY: 7,
	DIVIDE:   7,
	MODULO:   7,
	// ( [ . ++ --
	LPAREN:    8,
	LBRACKET:  8,
	DOT:       8,
	INCREMENT: 8,
	DECREMENT: 8,
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
