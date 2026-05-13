package scanner

import (
	"strconv"
	"sync"
)

type Kind int

func (tt Kind) String() string {
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
	Type          Kind
	Literal       string
	LeadingTrivia []Token
	AfterNewline  bool
}

const (
	// special keywords
	EOF Kind = iota
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

var tokenLiterals = map[Kind]string{
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

const initCustomType Kind = 1000

var (
	nextType   Kind = initCustomType
	registerMu sync.RWMutex
)

func RegisterKind(lit string) Kind {
	registerMu.Lock()
	defer registerMu.Unlock()
	typ := nextType
	tokenLiterals[typ] = lit
	nextType++
	return typ
}
