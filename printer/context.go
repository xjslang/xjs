package printer

func (pr *Printer) PushContext() map[string]string {
	ctx := make(map[string]string)
	pr.context = append(pr.context, ctx)
	return ctx
}

func (pr *Printer) PopContext() {
	if l := len(pr.context); l > 0 {
		pr.context = pr.context[:len(pr.context)-1]
	}
}

func (pr *Printer) Context() map[string]string {
	if l := len(pr.context); l > 0 {
		return pr.context[l-1]
	}
	return nil
}
