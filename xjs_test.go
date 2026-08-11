package xjs_test

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func ExampleScannerBuilder() {
	// register custom types first
	likeTyp := token.RegisterType("LIKE", "~~")
	notLikeTyp := token.RegisterType("NOT_LIKE", "~!")

	sb := xjs.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(sc); err != nil {
			return
		}
		if tok.Literal == "~" {
			// after consuming the token, the current char is
			// already in the next char to be processed
			switch sc.CurrentChar() {
			case '~':
				sc.AdvanceChar()
				tok.Type = likeTyp
				tok.Literal = "~~"
			case '!':
				sc.AdvanceChar()
				tok.Type = notLikeTyp
				tok.Literal = "~!"
			}
		}
		return
	})

	input := `
	name ~~ "john*"
	name ~! "admin*"`
	sc := sb.Build([]byte(input))
	for tok := sc.NextToken(); tok.Type != token.EOF; tok = sc.NextToken() {
		fmt.Printf("[%q: %s]", tok.Literal, tok.Type.Name())
	}

	// Output:
	// ["name": IDENT]["~~": LIKE]["\"john*\"": STRING]["name": IDENT]["~!": NOT_LIKE]["\"admin*\"": STRING]
}
