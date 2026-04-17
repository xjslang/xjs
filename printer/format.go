package printer

type formatConfig struct {
	indent    string
	semicolon bool
}

type FormatOption func(*formatConfig)

func WithIndent(value string) FormatOption {
	return func(cfg *formatConfig) {
		cfg.indent = value
	}
}

func WithSemicolon() FormatOption {
	return func(cfg *formatConfig) {
		cfg.semicolon = true
	}
}

func (p *Printer) WithFormat(opts ...FormatOption) *Printer {
	cfg := &formatConfig{
		indent:    p.indentString,
		semicolon: p.semicolon,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	p.formatted = true
	p.indentString = cfg.indent
	p.semicolon = cfg.semicolon
	return p
}

func (p *Printer) PrintNewline() {
	if !p.formatted {
		return
	}
	p.doc.WriteRune('\n')
}

func (p *Printer) PrintSemicolon() {
	if p.formatted && !p.semicolon {
		return
	}
	p.doc.WriteRune(';')
}

func (p *Printer) IncreaseIndent() {
	p.indentLevel++
}

func (p *Printer) DecreaseIndent() {
	p.indentLevel--
}

func (p *Printer) PrintIndent() {
	if !p.formatted {
		return
	}
	for range p.indentLevel {
		p.doc.WriteString(p.indentString)
	}
}

func (p *Printer) PrintWhitespace() {
	if !p.formatted {
		return
	}
	p.doc.WriteRune(' ')
}
