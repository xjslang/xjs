package ast

func (cw *CodeWriter) WriteLeadingComments(comments []string) {
	if !cw.PrettyPrint || len(comments) == 0 {
		return
	}

	for i, comment := range comments {
		if i > 0 {
			cw.writeNewline()
			cw.writeIndent()
		}
		if len(comment) > 0 {
			cw.Builder.WriteString("//")
		}
		cw.Builder.WriteString(comment)
	}
	cw.WriteNewline()
}
