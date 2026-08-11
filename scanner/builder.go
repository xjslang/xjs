package scanner

import "github.com/xjslang/xjs/token"

type Builder struct {
	scanners []func(*Scanner, func(*Scanner) (token.Token, error)) (token.Token, error)
}

// TODO: The middleware function type is repeated inline in both the Builder field and UseScanner signature, which makes the exported API harder to read and leads to long, duplicated types across the repo. Consider introducing named function types (e.g., NextScannerFunc / ScannerMiddleware) and using them in Builder/UseScanner (and in Scanner.useScanner) to keep signatures concise and consistent.
func (b *Builder) UseScanner(scanner func(sc *Scanner, next func(*Scanner) (token.Token, error)) (token.Token, error)) *Builder {
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
