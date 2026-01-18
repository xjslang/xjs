package ast

func (cw *CodeWriter) flushPending() {
	for _, ch := range cw.pendings {
		cw.Builder.WriteRune(ch)
	}
	cw.pendings = []rune{}

	if cw.pendingIndent {
		cw.writeIndent()
		cw.pendingIndent = false
	}
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
	cw.pendingIndent = true
}

// WriteSemi writes a semicolon if WriteSemicolons is true.
func (cw *CodeWriter) WriteSemi() {
	if !cw.PrettyPrint {
		cw.WriteRune(';')
		return
	}

	if cw.WriteSemicolons {
		cw.WriteRune(';')
	}
}

// WriteNewline writes a newline character if PrettyPrint is enabled
func (cw *CodeWriter) WriteNewline() {
	if !cw.PrettyPrint {
		return
	}

	if n := len(cw.pendings); n > 0 && (cw.pendings[n-1] == ' ' || cw.pendings[n-1] == '\n') {
		cw.pendings[n-1] = '\n'
	} else {
		cw.pendings = append(cw.pendings, '\n')
	}
}

// WriteSpace writes a space character if PrettyPrint is enabled
func (cw *CodeWriter) WriteSpace() {
	if !cw.PrettyPrint {
		return
	}

	if n := len(cw.pendings); n > 0 {
		lastChar := cw.pendings[n-1]
		// '\n' has higher priority than ' '
		if lastChar != '\n' {
			cw.pendings[n-1] = ' '
		}
	} else {
		cw.pendings = append(cw.pendings, ' ')
	}
}
