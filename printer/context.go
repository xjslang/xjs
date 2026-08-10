package printer

func (pr *Printer) PushContext() map[string]any {
	ctx := make(map[string]any)
	pr.context = append(pr.context, ctx)
	return ctx
}

func (pr *Printer) PopContext() {
	if l := len(pr.context); l > 0 {
		pr.context = pr.context[:len(pr.context)-1]
	}
}

func (pr *Printer) Context() map[string]any {
	if l := len(pr.context); l > 0 {
		return pr.context[l-1]
	}
	return nil
}
