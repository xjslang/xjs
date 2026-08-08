package token

import (
	"strconv"
	"sync"
)

type ForkableScanner interface {
	ForkFrom(Position) Scanner
	Fork() Scanner
	Apply(Scanner)
}

type Scanner interface {
	CurrentChar() rune
	AdvanceChar()
	NextToken() Token
}

type Type int

func (tt Type) Name() string {
	registerMu.RLock()
	defer registerMu.RUnlock()
	info, ok := tokenInfo[tt]
	if !ok {
		return "unknown(" + strconv.Itoa(int(tt)) + ")"
	}
	return info.name
}

func (tt Type) Literal() string {
	registerMu.RLock()
	defer registerMu.RUnlock()
	info, ok := tokenInfo[tt]
	if !ok {
		return "unknown(" + strconv.Itoa(int(tt)) + ")"
	}
	return info.literal
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
	Offset int `json:"offset"`
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

var tokenInfo = map[Type]struct {
	name, literal string
}{
	// special keywords
	EOF:     {"EOF", ""},
	IDENT:   {"IDENT", "identifier"},
	ILLEGAL: {"ILLEGAL", ""},
	UNKNOWN: {"UNKNOWN", ""},
	// operators
	ASSIGN:   {"ASSIGN", "="},
	PLUS:     {"PLUS", "+"},
	MINUS:    {"MINUS", "-"},
	MULTIPLY: {"MULTIPLY", "*"},
	DIVIDE:   {"DIVIDE", "/"},
	MODULO:   {"MODULO", "%"},
	// incremental operators
	INCREMENT: {"INCREMENT", "++"},
	DECREMENT: {"DECREMENT", "--"},
	// comparison operators
	EQ:     {"EQ", "=="},
	NOT_EQ: {"NOT_EQ", "!="},
	LT:     {"LT", "<"},
	LTE:    {"LTE", "<="},
	GT:     {"GT", ">"},
	GTE:    {"GTE", ">="},
	// logical operators
	AND: {"AND", "&&"},
	OR:  {"OR", "||"},
	NOT: {"NOT", "!"},
	// delimiters
	COMMA:     {"COMMA", ","},
	SEMICOLON: {"SEMICOLON", ";"},
	COLON:     {"COLON", ":"},
	DOT:       {"DOT", "."},
	LPAREN:    {"LPAREN", "("},
	RPAREN:    {"RPAREN", ")"},
	LBRACE:    {"LBRACE", "{"},
	RBRACE:    {"RBRACE", "}"},
	LBRACKET:  {"LBRACKET", "["},
	RBRACKET:  {"RBRACKET", "]"},
	NEWLINE:   {"NEWLINE", "new line"},
	// others
	LINE_COMMENT:  {"LINE_COMMENT", ""},
	BLOCK_COMMENT: {"BLOCK_COMMENT", ""},
	STRING:        {"STRING", ""},
	NUMBER:        {"NUMBER", ""},
}

const initCustomType Type = 1000

var (
	nextType   Type = initCustomType
	registerMu sync.RWMutex
)

func RegisterType(name, lit string) Type {
	registerMu.Lock()
	defer registerMu.Unlock()
	typ := nextType
	tokenInfo[typ] = struct{ name, literal string }{name, lit}
	nextType++
	return typ
}

var binaryOps = map[Type]int{
	COMMA: 0,
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

var prefixOpTypes = map[Type]bool{
	NOT:      true,
	PLUS:     true,
	MINUS:    true,
	LPAREN:   true,
	LBRACE:   true,
	LBRACKET: true,
}

func (typ Type) IsPrefixOp() (ok bool) {
	registerMu.RLock()
	defer registerMu.RUnlock()
	_, ok = prefixOpTypes[typ]
	return
}

// RegisterPrefixOp registers a token type as a "prefix operator".
func RegisterPrefixOp(typ Type) {
	registerMu.Lock()
	defer registerMu.Unlock()
	prefixOpTypes[typ] = true
}
