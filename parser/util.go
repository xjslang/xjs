package parser

import "github.com/xjslang/xjs/ast"

func Resolve[T ast.Node](p *Parser, parsers ...func(p *Parser) (T, error)) (node T, err error) {
	for _, parser := range parsers {
		f := p.Fork()
		if node, err = parser(f); err != nil {
			continue
		}
		p.Apply(f)
		break
	}
	return
}
