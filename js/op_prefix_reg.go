package js

import (
	"errors"
	"slices"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type PrefixRegOp struct {
	ast.BaseExpr
	Value token.Token
}

func ParsePrefixRegOp(p *parser.Parser) (node *PrefixRegOp, err error) {
	node = &PrefixRegOp{}
	if node.Value, err = p.ExpectWith(scanRegex); err != nil {
		return
	}
	p.AdvanceToken()
	return
}

func PrintPrefixRegOp(pr *printer.Printer, node *PrefixRegOp) error {
	pr.Print(node.Value)
	return nil
}

func scanRegex(s *scanner.Scanner) (string, error) {
	sb := strings.Builder{}
	if s.CurrentChar() != '/' {
		return "", errors.New("/ expected")
	}
	sb.WriteRune(s.CurrentChar())
	s.AdvanceChar() // consume /
	inClass := false
	for {
		ch := s.CurrentChar()
		if ch == '\\' {
			sb.WriteRune(ch)
			s.AdvanceChar()
			ch = s.CurrentChar()
			if ch == scanner.EOF || ch == '\n' || ch == '\r' {
				return sb.String(), errors.New("unexpected end of line")
			}
			sb.WriteRune(ch)
			s.AdvanceChar()
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
			s.AdvanceChar()
			break
		}
		sb.WriteRune(ch)
		s.AdvanceChar()
	}
	flags := []rune("dgimsuvy")
	for slices.Contains(flags, s.CurrentChar()) {
		sb.WriteRune(s.CurrentChar())
		s.AdvanceChar()
	}
	if c := s.CurrentChar(); (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return "", errors.New("unknown flag " + string(c))
	}
	return sb.String(), nil
}
