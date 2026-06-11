package printer

func (p *Printer) PushContext() map[string]string {
	ctx := make(map[string]string)
	p.context = append(p.context, ctx)
	return ctx
}

func (p *Printer) PopContext() {
	if l := len(p.context); l > 0 {
		p.context = p.context[:len(p.context)-1]
	}
}

func (p *Printer) Context() map[string]string {
	if l := len(p.context); l > 0 {
		return p.context[l-1]
	}
	return nil
}
