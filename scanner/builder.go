package scanner

import "github.com/xjslang/xjs/token"

type Builder struct {
	scanners []func(*Scanner, func() (token.Token, error)) (token.Token, error)
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) UseScanner(scanner func(sc *Scanner, next func() (token.Token, error)) (token.Token, error)) *Builder {
	b.scanners = append(b.scanners, scanner)
	return b
}

func (b *Builder) Build(input []byte) *Scanner {
	sc := &Scanner{}
	for _, scanner := range b.scanners {
		sc.useScanner(scanner)
	}
	sc.init(input)
	return sc
}
