package scanner

import "github.com/xjslang/xjs/token"

type Builder struct {
	scanners []func(*Scanner, func() (token.Token, error)) (token.Token, error)
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) UseScanner(scanner func(s *Scanner, next func() (token.Token, error)) (token.Token, error)) *Builder {
	b.scanners = append(b.scanners, scanner)
	return b
}

func (b *Builder) Build(input []byte) *Scanner {
	s := &Scanner{}
	for _, scanner := range b.scanners {
		s.useScanner(scanner)
	}
	s.init(input)
	return s
}
