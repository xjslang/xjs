package ast

func (cw *CodeWriter) clearPending() {
	cw.pendings = []rune{}
}

func (cw *CodeWriter) flushPending() {
	for _, ch := range cw.pendings {
		if ch == '\t' {
			cw.writeIndent()
		} else {
			cw.Builder.WriteRune(ch)
		}
	}
	cw.clearPending()
}

func (cw *CodeWriter) writeNewline() {
	cw.Builder.WriteRune('\n')
}

func (cw *CodeWriter) writeIndent() {
	indent := cw.IndentString
	if indent == "" {
		indent = "  " // default: 2 spaces
	}
	for i := 0; i < cw.IndentLevel; i++ {
		cw.Builder.WriteString(indent)
	}
}

// IncreaseIndent increases the indentation level
func (cw *CodeWriter) IncreaseIndent() {
	if !cw.PrettyPrint {
		return
	}
	cw.IndentLevel++
}

// DecreaseIndent decreases the indentation level
func (cw *CodeWriter) DecreaseIndent() {
	if !cw.PrettyPrint {
		return
	}
	if cw.IndentLevel > 0 {
		cw.IndentLevel--
	}
}

// WriteIndent writes the current indentation level
func (cw *CodeWriter) WriteIndent() {
	if !cw.PrettyPrint {
		return
	}
	if n := len(cw.pendings); n == 0 || cw.pendings[n-1] != '\t' {
		cw.pendings = append(cw.pendings, '\t')
	}
}

// WriteNewline writes a newline character if PrettyPrint is enabled
func (cw *CodeWriter) WriteNewline() {
	if !cw.PrettyPrint {
		return
	}
	cw.clearPending()
	cw.pendings = append(cw.pendings, '\n')
}

// WriteSpace writes a space character if PrettyPrint is enabled
func (cw *CodeWriter) WriteSpace() {
	if !cw.PrettyPrint {
		return
	}
	if n := len(cw.pendings); n == 0 || cw.pendings[n-1] != ' ' {
		cw.pendings = append(cw.pendings, ' ')
	}
}
