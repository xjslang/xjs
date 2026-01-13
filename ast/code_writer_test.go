package ast

import (
	"fmt"
	"testing"
)

func ExampleCodeWriter() {
	cw := CodeWriter{
		PrettyPrint: true,
	}
	cw.WriteString("function test(")
	for i := range 3 {
		if i > 0 {
			cw.WriteString(", ")
		}
		cw.WriteString(fmt.Sprintf("param_%d", i))
	}
	cw.WriteRune(')')
	cw.WriteSpace()
	cw.WriteRune('{')
	cw.writeNewline()
	cw.IncreaseIndent()
	for i := range 3 {
		cw.WriteIndent()
		cw.WriteString(fmt.Sprintf("command_%d()", i))
		cw.WriteNewline()
	}
	cw.DecreaseIndent()
	cw.WriteRune('}')

	fmt.Println(cw.String())
	// Output:
	// function test(param_0, param_1, param_2) {
	//   command_0()
	//   command_1()
	//   command_2()
	// }
}

// TestWriteNewline verifies that WriteNewline doesn't add a new line at the end.
func TestWriteNewline(t *testing.T) {
	cw := CodeWriter{PrettyPrint: true}
	for i := range 2 {
		cw.WriteString(fmt.Sprintf("command_%d()", i))
		cw.WriteNewline()
	}

	expected := "command_0()\ncommand_1()"
	output := cw.String()
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

// TestWriteIndent verifies that WriteIndent doesn't add a new indent at the end.
func TestWriteIndent(t *testing.T) {
	cw := CodeWriter{
		PrettyPrint:  true,
		IndentString: ", ",
	}
	cw.IncreaseIndent()
	for i := range 2 {
		cw.WriteString(fmt.Sprintf("param_%d()", i))
		cw.WriteIndent()
	}
	cw.DecreaseIndent()

	expected := "param_0(), param_1()"
	output := cw.String()
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}
