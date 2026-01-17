package ast

func (cw *CodeWriter) WriteLeadingComments(comments []string) {
	if !cw.PrettyPrint || len(comments) == 0 {
		return
	}

	for i, comment := range comments {
		isComment := len(comment) > 0
		if i == 0 {
			// Add a space between the current line and the comment
			if cw.pendingNewline && isComment {
				cw.Builder.WriteRune(' ')
			}
		} else {
			cw.writeNewline()
			cw.writeIndent()
		}
		if isComment {
			cw.Builder.WriteString("//")
		}
		cw.Builder.WriteString(comment)
	}
	cw.WriteNewline()
	cw.WriteIndent()
}
