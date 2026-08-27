package parser

import "github.com/xjslang/xjs/token"

type ScannerBuilder struct {
	scanners []func(*Scanner, func(*Scanner) token.Token) token.Token
}

// TODO: The middleware function type is repeated inline in both the Builder field and UseScanner signature, which makes the exported API harder to read and leads to long, duplicated types across the repo. Consider introducing named function types (e.g., NextScannerFunc / ScannerMiddleware) and using them in Builder/UseScanner (and in Scanner.useScanner) to keep signatures concise and consistent.
func (b *ScannerBuilder) UseScanner(scanner func(sc *Scanner, next func(*Scanner) token.Token) token.Token) *ScannerBuilder {
	b.scanners = append(b.scanners, scanner)
	return b
}

// TODO: Scanner.Builder.Build now accepts a string and the scanner populates Token.Literal/LeadingTrivia using string slices. In Go, substrings retain a reference to the entire original string, so keeping tokens/AST nodes alive will also keep the full source input alive. If this is intentional for perf, it should be documented on the public Build API (and/or Scanner) so callers understand the memory-ownership tradeoff.
func (b *ScannerBuilder) Build(input string) *Scanner {
	sc := &Scanner{}
	for _, scanner := range b.scanners {
		sc.useScanner(scanner)
	}
	sc.init(input)
	return sc
}
