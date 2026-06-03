package printer

func Fork(p *Printer) *Printer {
	p1 := &Printer{printer: p.printer}
	p1.Init(
		WithLineComments(p.lineComments),
		WithBlockComments(p.blockComments),
		WithNewLines(p.newLines),
		WithIndent(p.indent),
	)
	return p1
}

func isNewLine(r rune) bool {
	return r == eol || r == '\r' || r == '\n'
}

func isWhitespace(r rune) bool {
	return isNewLine(r) || r == ' ' || r == '\t'
}
