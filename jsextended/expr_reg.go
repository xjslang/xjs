package jsextended

import (
	"errors"
	"slices"
	"strings"
	"text/scanner"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type RegExpr struct {
	ast.BaseExpr
	Value token.Token
}

func ParseRegExpr(p *parser.Parser) (node *RegExpr, err error) {
	node = &RegExpr{}
	if node.Value, err = p.ExpectWith(scanRegex); err != nil {
		return
	}
	p.AdvanceToken()
	return
}

func PrintRegExpr(pr *printer.Printer, node *RegExpr) error {
	pr.Print(node.Value)
	return nil
}

func scanRegex(sc token.Scanner) (string, error) {
	sb := strings.Builder{}
	if sc.CurrentChar() != '/' {
		return "", errors.New("/ expected")
	}
	sb.WriteRune(sc.CurrentChar())
	sc.AdvanceChar() // consume /
	inClass := false
	for {
		ch := sc.CurrentChar()
		if ch == '\\' {
			sb.WriteRune(ch)
			sc.AdvanceChar()
			ch = sc.CurrentChar()
			if ch == scanner.EOF || ch == '\n' || ch == '\r' {
				return sb.String(), errors.New("unexpected end of line")
			}
			sb.WriteRune(ch)
			sc.AdvanceChar()
			continue
		}
		if ch == scanner.EOF || ch == '\n' || ch == '\r' {
			return sb.String(), errors.New("unexpected end of line")
		}
		if ch == '[' && !inClass {
			inClass = true
		} else if ch == ']' && inClass {
			inClass = false
		}
		if ch == '/' && !inClass {
			sb.WriteRune(ch)
			sc.AdvanceChar()
			break
		}
		sb.WriteRune(ch)
		sc.AdvanceChar()
	}
	flags := []rune("dgimsuvy")
	for slices.Contains(flags, sc.CurrentChar()) {
		sb.WriteRune(sc.CurrentChar())
		sc.AdvanceChar()
	}
	if c := sc.CurrentChar(); (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return "", errors.New("unknown flag " + string(c))
	}
	return sb.String(), nil
}
