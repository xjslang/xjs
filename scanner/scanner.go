package scanner

import "github.com/xjslang/xjs/parser"

const EOF = parser.EOF

type (
	Builder = parser.ScannerBuilder
	Scanner = parser.Scanner
)

func IsLetter(r rune) bool {
	return parser.IsLetter(r)
}
