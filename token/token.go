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
	NUMBER        // 3.1415
	LINE_COMMENT  // // ..
	BLOCK_COMMENT // /* .. */
	STRING        // '..' or ".."
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
	LINE_COMMENT:  "line comment",
	BLOCK_COMMENT: "block comment",
	STRING:        "string",
	NUMBER:        "number",
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
	ASSIGN: 1000,
	// ||
	OR: 2000,
	// &&
	AND: 3000,
	// == !=
	EQ:     4000,
	NOT_EQ: 4000,
	// < <= > >=
	LT:  5000,
	LTE: 5000,
	GT:  5000,
	GTE: 5000,
	// + -
	PLUS:  6000,
	MINUS: 6000,
	// * / %
	MULTIPLY: 7000,
	DIVIDE:   7000,
	MODULO:   7000,
	// ( [ . ++ --
	LPAREN:    8000,
	LBRACKET:  8000,
	DOT:       8000,
	INCREMENT: 8000,
	DECREMENT: 8000,
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
